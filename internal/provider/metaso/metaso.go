package metaso

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

const ProviderName = "metaso"

// Client 秘塔（MetaSo）搜索客户端。
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://metaso.cn",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Name() string { return ProviderName }

type metaSearchRequest struct {
	Q string `json:"q"`

	Scope string `json:"scope,omitempty"` // webpage...

	IncludeSummary    bool `json:"includeSummary,omitempty"`
	IncludeRawContent bool `json:"includeRawContent,omitempty"`
	ConciseSnippet    bool `json:"conciseSnippet,omitempty"`

	// MetaSo 的分页和 size 是两套独立模式：
	// - 使用 page：固定每页 10 条
	// - 使用 size：不可分页
	Page any `json:"page,omitempty"`
	Size any `json:"size,omitempty"`
}

type metaWebpage struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Score   string `json:"score"` // high/medium/low（示例为 high）
	Snippet string `json:"snippet"`
	Date    string `json:"date"`
}

type metaResponse struct {
	Credits          int           `json:"credits"`
	SearchParameters map[string]any `json:"searchParameters"`
	Webpages         []metaWebpage `json:"webpages"`
	Total            int           `json:"total"`
}

func (c *Client) Search(ctx context.Context, ch *model.Channel, req *dto.UnifiedSearchRequest) (*dto.UnifiedSearchResponse, error) {
	if ch == nil || ch.ApiKey == "" {
		return nil, fmt.Errorf("metaso channel api key 为空")
	}
	if req == nil || req.Query == "" {
		return nil, fmt.Errorf("query 不能为空")
	}

	mReq := metaSearchRequest{
		Q:      req.Query,
		Scope:  "webpage",
		Page:   nil,
		Size:   nil,
	}

	// include_content -> includeSummary/includeRawContent
	switch req.IncludeContent {
	case "summary":
		mReq.IncludeSummary = true
	case "full":
		mReq.IncludeRawContent = true
	default:
		// none/空：两者均 false
	}

	// 两套独立模式：优先使用 max_results(size)，否则使用 offset(page)
	if req.MaxResults > 0 {
		mReq.Size = req.MaxResults
		mReq.Page = nil
	} else {
		// 统一接口的 offset 这里按“页偏移”理解：offset=0 -> page=1
		page := req.Offset + 1
		if page < 1 {
			page = 1
		}
		mReq.Page = page
	}

	bodyBytes, err := json.Marshal(&mReq)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/api/v1/search"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+ch.ApiKey)
	httpReq.Header.Set("Accept", "application/json")
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
		return nil, fmt.Errorf("metaso 请求失败，status=%d", resp.StatusCode)
	}

	var mResp metaResponse
	if err := json.Unmarshal(raw, &mResp); err != nil {
		return nil, err
	}

	out := &dto.UnifiedSearchResponse{
		Query: req.Query,
		Usage: &dto.UnifiedSearchUsage{Credits: mResp.Credits},
	}

	for _, w := range mResp.Webpages {
		out.Results = append(out.Results, dto.UnifiedSearchResult{
			Title:       w.Title,
			URL:         w.Link,
			Snippet:     w.Snippet,
			Content:     w.Snippet,
			Score:       mapScoreToFloat(w.Score),
			PublishedAt: w.Date,
		})
	}

	return out, nil
}

func mapScoreToFloat(score string) float64 {
	switch score {
	case "high":
		return 0.9
	case "medium":
		return 0.6
	case "low":
		return 0.3
	default:
		// 尝试直接解析成数字
		if f, err := strconv.ParseFloat(score, 64); err == nil {
			return f
		}
		return 0
	}
}
