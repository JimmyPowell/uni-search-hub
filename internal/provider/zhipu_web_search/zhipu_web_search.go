package zhipu_web_search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

const ProviderName = "zhipu_web_search"

// Client 智谱 Web Search API 客户端（可注入 BaseURL / HTTPClient 便于测试）。
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://open.bigmodel.cn/api",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Name() string { return ProviderName }

type webSearchRequest struct {
	SearchQuery         string `json:"search_query"`
	SearchEngine        string `json:"search_engine"`
	SearchIntent        bool   `json:"search_intent"`
	Count               int    `json:"count,omitempty"`
	SearchDomainFilter  string `json:"search_domain_filter,omitempty"`
	SearchRecencyFilter string `json:"search_recency_filter,omitempty"`
	ContentSize         string `json:"content_size,omitempty"`
}

type webSearchIntentItem struct {
	Query    string `json:"query"`
	Intent   string `json:"intent"`
	Keywords string `json:"keywords"`
}

type webSearchResultItem struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	Media       string `json:"media"`
	Icon        string `json:"icon"`
	Refer       string `json:"refer"`
	PublishDate string `json:"publish_date"`
}

type webSearchResponse struct {
	ID           string                `json:"id"`
	Created      int64                 `json:"created"`
	RequestID    string                `json:"request_id"`
	SearchIntent []webSearchIntentItem `json:"search_intent"`
	SearchResult []webSearchResultItem `json:"search_result"`
}

func (c *Client) Search(ctx context.Context, ch *model.Channel, req *dto.UnifiedSearchRequest) (*dto.UnifiedSearchResponse, error) {
	if ch == nil || ch.ApiKey == "" {
		return nil, fmt.Errorf("zhipu_web_search channel api key 为空")
	}
	if req == nil || req.Query == "" {
		return nil, fmt.Errorf("query 不能为空")
	}

	zReq := webSearchRequest{
		SearchQuery:  req.Query,
		SearchEngine: mapSearchEngine(req.SearchDepth),
		SearchIntent: req.AutoParameters,
		Count:        clampCount(req.MaxResults),
		ContentSize:  mapContentSize(req.IncludeContent),
	}

	// 智谱仅支持单域名白名单；统一结构里 site 优先，其次 include_domains[0]
	if req.Site != "" {
		zReq.SearchDomainFilter = req.Site
	} else if len(req.IncludeDomains) > 0 {
		zReq.SearchDomainFilter = req.IncludeDomains[0]
	}

	if rf := mapRecencyFilter(req.TimeRange); rf != "" {
		zReq.SearchRecencyFilter = rf
	}

	bodyBytes, err := json.Marshal(&zReq)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/paas/v4/web_search"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+ch.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("zhipu_web_search 请求失败，status=%d", resp.StatusCode)
	}

	var zResp webSearchResponse
	if err := json.Unmarshal(raw, &zResp); err != nil {
		return nil, err
	}

	out := &dto.UnifiedSearchResponse{
		Query:     req.Query,
		RequestID: zResp.RequestID,
	}
	for _, r := range zResp.SearchResult {
		out.Results = append(out.Results, dto.UnifiedSearchResult{
			Title:       r.Title,
			URL:         r.Link,
			Snippet:     r.Content,
			Content:     r.Content,
			Favicon:     r.Icon,
			PublishedAt: r.PublishDate,
		})
	}

	return out, nil
}

func clampCount(n int) int {
	// 智谱：count 范围 1-50，默认 10。
	if n <= 0 {
		return 10
	}
	if n > 50 {
		return 50
	}
	return n
}

func mapContentSize(includeContent string) string {
	// 智谱：content_size=medium|high
	if includeContent == "full" {
		return "high"
	}
	return "medium"
}

func mapSearchEngine(searchDepth string) string {
	// 你选择“映射”策略：advanced -> search_pro；basic/fast/ultra-fast -> search_std；默认 search_pro。
	switch searchDepth {
	case "advanced":
		return "search_pro"
	case "basic", "fast", "ultra-fast":
		return "search_std"
	case "":
		return "search_pro"
	default:
		// 未知值：按默认 search_pro，避免无效 enum 导致上游报错
		return "search_pro"
	}
}

func mapRecencyFilter(timeRange string) string {
	// 统一接口：day|week|month|year 或短写 d|w|m|y
	switch timeRange {
	case "day", "d":
		return "oneDay"
	case "week", "w":
		return "oneWeek"
	case "month", "m":
		return "oneMonth"
	case "year", "y":
		return "oneYear"
	default:
		return ""
	}
}
