package provider

import (
	"fmt"
	"sync"
)

var (
	registryMu sync.RWMutex
	registry   = map[string]SearchProvider{}
)

// RegisterProvider 注册 provider。重复注册会 panic，便于启动时尽早暴露配置问题。
func RegisterProvider(p SearchProvider) {
	if p == nil {
		panic("provider is nil")
	}
	name := p.Name()
	if name == "" {
		panic("provider name is empty")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("provider already registered: %s", name))
	}
	registry[name] = p
}

func GetProvider(name string) (SearchProvider, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, ok := registry[name]
	return p, ok
}

// ListProviders 返回已注册 provider 的名称列表（用于 meta 展示）。
func ListProviders() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	return out
}
