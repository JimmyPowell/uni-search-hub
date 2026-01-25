package provider

import (
	"uni-search-hub/internal/provider/metaso"
	"uni-search-hub/internal/provider/tavily"
	"uni-search-hub/internal/provider/zhipu_web_search"
)

func init() {
	// 默认注册 Tavily（MVP）。
	RegisterProvider(tavily.NewClient())
	// 默认注册 秘塔（MetaSo）搜索。
	RegisterProvider(metaso.NewClient())
	// 默认注册 智谱 Web Search。
	RegisterProvider(zhipu_web_search.NewClient())
}
