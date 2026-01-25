package handler

import (
	"net/http"
	"time"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/provider"
	"uni-search-hub/internal/server"
	"uni-search-hub/internal/service"
	"uni-search-hub/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var rrSelector = &service.RoundRobinSelector{}

// UnifiedSearch 统一搜索入口（MVP：仅接入 Tavily + 固定扣费 1）。
// 鉴权：使用 TokenAuth，外部调用通过 Authorization: Bearer <token>。
func UnifiedSearch(c *gin.Context) {
	var req dto.UnifiedSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "query 不能为空"})
		return
	}

	// MVP：默认使用 tavily；可通过 query param 切换，如 /api/search?provider=metaso
	providerName := c.Query("provider")
	if providerName == "" {
		providerName = "tavily"
	}

	p, ok := provider.GetProvider(providerName)
	if !ok {
		server.ApiErrorMsg(c, "provider 未注册")
		return
	}

	channels, err := model.GetEnabledChannelsByProvider(providerName)
	if err != nil {
		server.ApiError(c, err)
		return
	}
	ch, err := rrSelector.Select(providerName, channels)
	if err != nil {
		server.ApiError(c, err)
		return
	}

	start := time.Now()
	resp, err := p.Search(c.Request.Context(), ch, &req)
	latencyMs := time.Since(start).Milliseconds()

	// 审计日志：无论成功/失败都记录（MVP：同步写入，后续可异步/批量）。
	reqId := uuid.NewString()
	userId := c.GetInt("id")
	tokenId, _ := c.Get("token_id")
	tokenIdInt, _ := tokenId.(int)

	log := &model.RequestLog{
		RequestId:  reqId,
		UserId:     userId,
		TokenId:    tokenIdInt,
		ChannelId:  ch.Id,
		Provider:   providerName,
		Action:     model.RequestActionProxyRequest,
		Endpoint:   "/api/search",
		StatusCode: 200,
		LatencyMs:  latencyMs,
		Cost:       0,
	}

	if err != nil {
		log.StatusCode = 500
		log.Error = utils.TruncateString(err.Error(), 255)
		_ = log.Insert()
		server.ApiError(c, err)
		return
	}

	// 固定扣费 1（后续可按 provider usage 精细计费）
	cost := 1
	log.Cost = cost
	_ = log.Insert()

	// 扣费：优先扣 token 额度；同时记录 user used_quota / request_count 便于审计统计
	if tokenAny, ok := c.Get("token"); ok {
		if t, ok := tokenAny.(*model.Token); ok && t != nil && !t.UnlimitedQuota {
			_ = model.DecreaseTokenQuota(t.Id, t.Key, cost)
		}
	}
	if userId != 0 {
		model.UpdateUserUsedQuotaAndRequestCount(userId, cost)
	}

	// 统一响应附带我们自己的 request_id，便于排障/对账
	if resp != nil {
		resp.RequestID = reqId
	}
	server.ApiSuccess(c, resp)
}
