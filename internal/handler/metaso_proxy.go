package handler

import (
	"io"
	"net/http"
	"time"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/provider/metaso"
	"uni-search-hub/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var metasoClient = metaso.NewClient()

// MetaSoSearchProxy 秘塔（MetaSo）原生接口透传：
// - 请求体/响应体保持秘塔（MetaSo）原生格式
// - 渠道轮询、鉴权、审计与扣费（MVP：固定 1）
func MetaSoSearchProxy(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无法读取请求体"})
		return
	}

	providerName := metaso.ProviderName
	channels, err := model.GetEnabledChannelsByProvider(providerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取渠道失败"})
		return
	}
	ch, err := rrSelector.Select(providerName, channels)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "没有可用渠道"})
		return
	}

	contentType := c.GetHeader("Content-Type")
	start := time.Now()
	status, respCT, respBody, err := metasoClient.ProxySearch(c.Request.Context(), ch.ApiKey, rawBody, contentType)
	latencyMs := time.Since(start).Milliseconds()

	reqId := uuid.NewString()
	c.Header("X-Uni-Request-Id", reqId)

	userId := c.GetInt("id")
	tokenIdAny, _ := c.Get("token_id")
	tokenId, _ := tokenIdAny.(int)

	log := &model.RequestLog{
		RequestId:  reqId,
		UserId:     userId,
		TokenId:    tokenId,
		ChannelId:  ch.Id,
		Provider:   providerName,
		Action:     model.RequestActionProxyRequest,
		Endpoint:   "/api/metaso/search",
		StatusCode: status,
		LatencyMs:  latencyMs,
		Cost:       0,
	}

	if err != nil {
		log.StatusCode = http.StatusBadGateway
		log.Error = utils.TruncateString(err.Error(), 255)
		_ = log.Insert()
		c.JSON(http.StatusBadGateway, gin.H{"message": "上游请求失败"})
		return
	}

	// 审计 + 固定扣费（只要上游有返回，就计费一次）
	cost := 1
	log.Cost = cost
	_ = log.Insert()

	if tokenAny, ok := c.Get("token"); ok {
		if t, ok := tokenAny.(*model.Token); ok && t != nil && !t.UnlimitedQuota {
			_ = model.DecreaseTokenQuota(t.Id, t.Key, cost)
		}
	}
	if userId != 0 {
		model.UpdateUserUsedQuotaAndRequestCount(userId, cost)
	}

	if respCT == "" {
		respCT = "application/json; charset=utf-8"
	}
	c.Data(status, respCT, respBody)
}
