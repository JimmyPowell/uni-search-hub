package handler

import (
	"sort"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/provider"
	"uni-search-hub/internal/server"

	"github.com/gin-gonic/gin"
)

type serviceEndpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Auth   string `json:"auth"` // token|session
	Notes  string `json:"notes,omitempty"`

	ExampleHeaders map[string]string `json:"example_headers,omitempty"`
	ExampleBody    any               `json:"example_body,omitempty"`
}

type serviceItem struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Provider             string            `json:"provider"`
	EnabledChannelsCount int64             `json:"enabled_channels_count"`
	Endpoints            []serviceEndpoint `json:"endpoints"`
}

type serviceCatalog struct {
	Services []serviceItem `json:"services"`
}

// GetServiceCatalog 服务目录（阶段 1）：
// - 展示系统已注册的 provider
// - 展示统一聚合与透传端点 + 示例请求
// - 统计各 provider 启用渠道数量（不泄露 api_key）
func GetServiceCatalog(c *gin.Context) {
	providers := provider.ListProviders()
	sort.Strings(providers)

	counts, err := model.CountEnabledChannelsByProvider(providers)
	if err != nil {
		server.ApiError(c, err)
		return
	}

	services := make([]serviceItem, 0, len(providers))
	for _, p := range providers {
		pt := passthroughPath(p)
		services = append(services, serviceItem{
			ID:                   "search:" + p,
			Name:                 "搜索服务（" + p + "）",
			Provider:             p,
			EnabledChannelsCount: counts[p],
			Endpoints: []serviceEndpoint{
				{
					Method: "POST",
					Path:   "/api/search?provider=" + p,
					Auth:   "token",
					Notes:  "统一聚合接口（统一请求结构）",
					ExampleHeaders: map[string]string{
						"Authorization": "Bearer <token_key>",
						"Content-Type":  "application/json",
					},
					ExampleBody: map[string]any{
						"query":           "北京大学 校训",
						"search_depth":    "advanced",
						"max_results":     5,
						"time_range":      "week",
						"include_content": "full",
						"auto_parameters": false,
					},
				},
				{
					Method: "POST",
					Path:   pt,
					Auth:   "token",
					Notes:  "原生透传接口（请求/响应体保持上游格式，仅替换上游 Authorization）",
					ExampleHeaders: map[string]string{
						"Authorization": "Bearer <token_key>",
						"Content-Type":  "application/json",
					},
					ExampleBody: passthroughExampleBody(p),
				},
			},
		})
	}

	server.ApiSuccess(c, serviceCatalog{Services: services})
}

func passthroughPath(p string) string {
	switch p {
	case "tavily":
		return "/api/tavily/search"
	case "metaso":
		return "/api/metaso/search"
	case "zhipu_web_search":
		return "/api/zhipu_web_search/search"
	default:
		// 兜底：约定 /api/{provider}/search（方便未来扩展）
		return "/api/" + p + "/search"
	}
}

func passthroughExampleBody(p string) any {
	switch p {
	case "tavily":
		return map[string]any{
			"query":       "北京大学 校训",
			"max_results": 5,
		}
	case "metaso":
		return map[string]any{
			"q":              "北京大学 校训",
			"scope":          "webpage",
			"includeSummary": true,
			"page":           1,
		}
	case "zhipu_web_search":
		return map[string]any{
			"search_query":  "北京大学 校训",
			"search_engine": "search_pro",
			"search_intent": false,
			"count":         5,
		}
	default:
		return map[string]any{}
	}
}
