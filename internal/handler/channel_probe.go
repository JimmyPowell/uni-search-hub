package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/provider/metaso"
	"uni-search-hub/internal/provider/tavily"
	"uni-search-hub/internal/provider/zhipu_web_search"
	"uni-search-hub/internal/server"
	"uni-search-hub/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TestChannel 使用渠道的 ApiKey 对上游做一次最小请求测试，验证该渠道是否可用。
// 仅 Root 可调用（由路由层保证）。
func TestChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的渠道ID"})
		return
	}
	ch, err := model.GetChannelById(id)
	if err != nil {
		server.ApiError(c, err)
		return
	}
	if ch.ApiKey == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "渠道 API Key 为空"})
		return
	}

	query := c.Query("q")
	if query == "" {
		query = "hello"
	}

	req := &dto.UnifiedSearchRequest{
		Query:          query,
		SearchDepth:    "basic",
		MaxResults:     1,
		IncludeAnswer:  "none",
		IncludeContent: "none",
	}

	start := time.Now()
	status := http.StatusOK
	errMsg := ""
	var result any
	switch ch.Provider {
	case tavily.ProviderName:
		cli := tavily.NewClient()
		resp, err := cli.Search(c.Request.Context(), ch, req)
		if err != nil {
			status = http.StatusBadGateway
			errMsg = err.Error()
		} else {
			result = previewUnified(resp)
		}
	case metaso.ProviderName:
		cli := metaso.NewClient()
		resp, err := cli.Search(c.Request.Context(), ch, req)
		if err != nil {
			status = http.StatusBadGateway
			errMsg = err.Error()
		} else {
			result = previewUnified(resp)
		}
	case zhipu_web_search.ProviderName:
		cli := zhipu_web_search.NewClient()
		// 智谱的映射：advanced->search_pro 更贴近官方默认示例；测试走 pro，降低空结果概率。
		reqCopy := *req
		reqCopy.SearchDepth = "advanced"
		resp, err := cli.Search(c.Request.Context(), ch, &reqCopy)
		if err != nil {
			status = http.StatusBadGateway
			errMsg = err.Error()
		} else {
			result = previewUnified(resp)
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "不支持的 provider: " + ch.Provider})
		return
	}

	latencyMs := time.Since(start).Milliseconds()

	// 记录审计：类型=channel_test，不扣费。
	reqId := uuid.NewString()
	c.Header("X-Uni-Request-Id", reqId)
	log := &model.RequestLog{
		RequestId:  reqId,
		UserId:     c.GetInt("id"),
		TokenId:    0,
		ChannelId:  ch.Id,
		Provider:   ch.Provider,
		Action:     model.RequestActionChannelTest,
		Endpoint:   c.Request.URL.Path,
		StatusCode: status,
		LatencyMs:  latencyMs,
		Cost:       0,
	}
	meta := gin.H{"query": query}
	if result != nil {
		meta["preview"] = result
	}
	if b, err := json.Marshal(meta); err == nil {
		log.Meta = string(b)
	}
	if errMsg != "" {
		log.Error = utils.TruncateString(errMsg, 255)
	}
	_ = log.Insert()

	if errMsg != "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": errMsg})
		return
	}

	server.ApiSuccess(c, gin.H{
		"ok":         true,
		"provider":   ch.Provider,
		"channel_id": ch.Id,
		"channel":    ch.Name,
		"latency_ms": latencyMs,
		"query":      query,
		"preview":    result,
	})
}

func previewUnified(resp *dto.UnifiedSearchResponse) any {
	if resp == nil {
		return nil
	}
	out := gin.H{
		"request_id":    resp.RequestID,
		"results_count": len(resp.Results),
		"response_time": resp.ResponseTime,
	}
	if len(resp.Results) > 0 {
		out["top_result"] = gin.H{
			"title": resp.Results[0].Title,
			"url":   resp.Results[0].URL,
		}
	}
	return out
}
