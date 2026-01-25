package provider

import (
	"context"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

// SearchProvider 对外部搜索供应商的最小抽象。
// 后续接入 Brave/Jina 等，只需新增实现并在 registry 中注册即可。
type SearchProvider interface {
	Name() string
	Search(ctx context.Context, ch *model.Channel, req *dto.UnifiedSearchRequest) (*dto.UnifiedSearchResponse, error)
}

