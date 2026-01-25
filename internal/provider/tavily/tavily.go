package tavily

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

const ProviderName = "tavily"

// Client Tavily Search API 客户端（可注入 BaseURL / HTTPClient 便于测试）。
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://api.tavily.com",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Name() string { return ProviderName }

type tavilySearchRequest struct {
	Query                   string   `json:"query"`
	SearchDepth              string   `json:"search_depth,omitempty"`
	MaxResults               int      `json:"max_results,omitempty"`
	Topic                    string   `json:"topic,omitempty"`
	TimeRange                string   `json:"time_range,omitempty"`
	StartDate                string   `json:"start_date,omitempty"`
	EndDate                  string   `json:"end_date,omitempty"`
	IncludeAnswer            any      `json:"include_answer,omitempty"`
	IncludeRawContent        any      `json:"include_raw_content,omitempty"`
	IncludeImages            bool     `json:"include_images,omitempty"`
	IncludeImageDescriptions bool     `json:"include_image_descriptions,omitempty"`
	IncludeFavicon           bool     `json:"include_favicon,omitempty"`
	IncludeDomains           []string `json:"include_domains,omitempty"`
	ExcludeDomains           []string `json:"exclude_domains,omitempty"`
	Country                  any      `json:"country,omitempty"`
	AutoParameters           bool     `json:"auto_parameters,omitempty"`
	IncludeUsage             bool     `json:"include_usage,omitempty"`
}

type tavilyResult struct {
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	RawContent string  `json:"raw_content"`
}

type tavilyResponse struct {
	Query        string         `json:"query"`
	Answer       string         `json:"answer"`
	Results      []tavilyResult `json:"results"`
	ResponseTime float64        `json:"response_time"`
	RequestID    string         `json:"request_id"`
}

func (c *Client) Search(ctx context.Context, ch *model.Channel, req *dto.UnifiedSearchRequest) (*dto.UnifiedSearchResponse, error) {
	if ch == nil || ch.ApiKey == "" {
		return nil, fmt.Errorf("tavily channel api key 为空")
	}
	if req == nil || req.Query == "" {
		return nil, fmt.Errorf("query 不能为空")
	}

	tReq := tavilySearchRequest{
		Query:                   req.Query,
		SearchDepth:              req.SearchDepth,
		MaxResults:               req.MaxResults,
		Topic:                    req.Topic,
		TimeRange:                req.TimeRange,
		StartDate:                req.StartDate,
		EndDate:                  req.EndDate,
		IncludeImages:            req.IncludeImages,
		IncludeImageDescriptions: req.IncludeImageDescriptions,
		IncludeFavicon:           req.IncludeFavicon,
		IncludeDomains:           req.IncludeDomains,
		ExcludeDomains:           req.ExcludeDomains,
		AutoParameters:           req.AutoParameters,
		IncludeUsage:             req.IncludeUsage,
	}

	// include_answer: Tavily 支持 bool|string。统一结构里用 none|basic|advanced。
	switch req.IncludeAnswer {
	case "basic", "advanced":
		tReq.IncludeAnswer = req.IncludeAnswer
	default:
		// none/空：不返回 answer
		if req.IncludeAnswer != "" && req.IncludeAnswer != "none" {
			// 非法值按 none 处理
		}
		tReq.IncludeAnswer = false
	}

	// include_content/content_format: 映射到 include_raw_content 的 bool|string
	// full -> true；summary/none -> false（MVP：不细分 markdown/text/html）
	switch req.IncludeContent {
	case "full":
		tReq.IncludeRawContent = true
	default:
		tReq.IncludeRawContent = false
	}

	// country: Tavily 允许 null（topic=general 才生效）
	if req.Country == "" {
		tReq.Country = nil
	} else {
		tReq.Country = req.Country
	}

	bodyBytes, err := json.Marshal(&tReq)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/search"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+ch.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")

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
		// 不把上游返回原文直接透给前端，避免泄露敏感信息/内部细节。
		return nil, fmt.Errorf("tavily 请求失败，status=%d", resp.StatusCode)
	}

	var tResp tavilyResponse
	if err := json.Unmarshal(raw, &tResp); err != nil {
		return nil, err
	}

	out := &dto.UnifiedSearchResponse{
		Query:        tResp.Query,
		Answer:       tResp.Answer,
		ResponseTime: tResp.ResponseTime,
		RequestID:    tResp.RequestID,
	}
	for _, r := range tResp.Results {
		out.Results = append(out.Results, dto.UnifiedSearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			// content 先用 raw_content（如果有）否则退回 content
			Content: func() string {
				if r.RawContent != "" {
					return r.RawContent
				}
				return r.Content
			}(),
			Score: r.Score,
		})
	}

	// MVP：credits 暂不从 Tavily 精确解析，走固定扣费策略。
	return out, nil
}

