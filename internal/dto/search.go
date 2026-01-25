package dto

// UnifiedSearchRequest 统一搜索请求结构（对齐 docs/search-api-mapping.md）。
// MVP 阶段只要求 query，其余字段按 provider 能力做 best-effort 映射。
type UnifiedSearchRequest struct {
	Query string `json:"query"`

	SearchDepth string `json:"search_depth,omitempty"` // basic|fast|ultra-fast|advanced
	MaxResults  int    `json:"max_results,omitempty"`
	Offset      int    `json:"offset,omitempty"`

	Topic     string `json:"topic,omitempty"`      // general|news|finance
	TimeRange string `json:"time_range,omitempty"` // day|week|month|year
	StartDate string `json:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   string `json:"end_date,omitempty"`   // YYYY-MM-DD

	Country    string `json:"country,omitempty"`
	Language   string `json:"language,omitempty"`
	UILanguage string `json:"ui_language,omitempty"`
	Location   string `json:"location,omitempty"`

	ResultTypes []string `json:"result_types,omitempty"`
	SafeSearch  string   `json:"safe_search,omitempty"` // off|moderate|strict

	IncludeAnswer string `json:"include_answer,omitempty"`  // none|basic|advanced
	IncludeContent string `json:"include_content,omitempty"` // none|summary|full
	ContentFormat  string `json:"content_format,omitempty"`  // markdown|text|html

	IncludeImages            bool `json:"include_images,omitempty"`
	IncludeImageDescriptions bool `json:"include_image_descriptions,omitempty"`
	IncludeFavicon           bool `json:"include_favicon,omitempty"`

	IncludeDomains []string `json:"include_domains,omitempty"`
	ExcludeDomains []string `json:"exclude_domains,omitempty"`
	Site           string   `json:"site,omitempty"`

	UseCache      bool `json:"use_cache,omitempty"`
	IncludeUsage  bool `json:"include_usage,omitempty"`
	AutoParameters bool `json:"auto_parameters,omitempty"`
}

type UnifiedSearchResult struct {
	Title       string  `json:"title,omitempty"`
	URL         string  `json:"url,omitempty"`
	Snippet     string  `json:"snippet,omitempty"`
	Content     string  `json:"content,omitempty"`
	Score       float64 `json:"score,omitempty"`
	Favicon     string  `json:"favicon,omitempty"`
	Language    string  `json:"language,omitempty"`
	PublishedAt string  `json:"published_at,omitempty"`
}

type UnifiedSearchImage struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type UnifiedSearchUsage struct {
	Credits int `json:"credits,omitempty"`
}

// UnifiedSearchResponse 统一搜索响应结构（对齐 docs/search-api-mapping.md）。
type UnifiedSearchResponse struct {
	Query        string               `json:"query,omitempty"`
	Answer       string               `json:"answer,omitempty"`
	Results      []UnifiedSearchResult `json:"results,omitempty"`
	Images       []UnifiedSearchImage  `json:"images,omitempty"`
	Usage        *UnifiedSearchUsage   `json:"usage,omitempty"`
	ResponseTime float64              `json:"response_time,omitempty"`
	RequestID    string               `json:"request_id,omitempty"`
}

