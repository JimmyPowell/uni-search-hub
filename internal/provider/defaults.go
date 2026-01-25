package provider

import "uni-search-hub/internal/provider/tavily"

func init() {
	// 默认注册 Tavily（MVP）。
	RegisterProvider(tavily.NewClient())
}

