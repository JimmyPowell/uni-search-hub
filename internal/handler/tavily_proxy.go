package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/provider/tavily"
	"uni-search-hub/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var tavilyClient = tavily.NewClient()

// TavilySearchProxy Tavily 透传接口：
// - 请求体/响应体保持 Tavily 原生格式
// - 仅做渠道轮询、鉴权、审计与扣费（MVP：固定 1）
func TavilySearchProxy(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无法读取请求体"})
		return
	}
	// 兼容部分 SDK：请求体可能携带 api_key（本系统 token）。转发上游前必须剥离，避免泄露。
	rawBody = stripJSONField(rawBody, "api_key")

	providerName := tavily.ProviderName
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
	status, respCT, respBody, err := tavilyClient.ProxySearch(c.Request.Context(), ch.ApiKey, rawBody, contentType)
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
		Endpoint:   "/api/tavily/search",
		StatusCode: status,
		LatencyMs:  latencyMs,
		Cost:       0,
	}

	// 发生网络/请求构造错误：未命中上游响应，按系统错误返回
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

func stripJSONField(raw []byte, field string) []byte {
	if len(raw) == 0 || field == "" {
		return raw
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return raw
	}
	if _, ok := m[field]; !ok {
		return raw
	}
	delete(m, field)
	out, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return bytes.TrimSpace(out)
}
