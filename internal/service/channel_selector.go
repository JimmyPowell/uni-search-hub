package service

import (
	"errors"
	"sync"
	"sync/atomic"
	"uni-search-hub/internal/model"
)

// RoundRobinSelector 简单轮询选路（MVP）。
// 注意：这是进程内轮询；多实例部署时需要更强的一致性策略（后续再做）。
type RoundRobinSelector struct {
	counters sync.Map // map[provider]*uint64
}

func (s *RoundRobinSelector) counterPtr(provider string) *uint64 {
	v, _ := s.counters.LoadOrStore(provider, new(uint64))
	return v.(*uint64)
}

func (s *RoundRobinSelector) Select(provider string, channels []*model.Channel) (*model.Channel, error) {
	if provider == "" {
		return nil, errors.New("provider 为空")
	}
	if len(channels) == 0 {
		return nil, errors.New("没有可用渠道")
	}
	idx := atomic.AddUint64(s.counterPtr(provider), 1) - 1
	return channels[int(idx)%len(channels)], nil
}

