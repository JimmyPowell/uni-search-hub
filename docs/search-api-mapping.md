# 搜索 API 对齐说明（三家 + 统一标准）

本文聚焦三家搜索 API（Tavily、Brave、Jina）的请求格式与参数说明，并给出一份参考 Tavily 的统一接口定义与字段映射。

## 1. 三家搜索 API 的标准请求格式与参数说明

### 1.1 Tavily Search API

- Endpoint: `POST https://api.tavily.com/search`
- 认证: `Authorization: Bearer <TAVILY_API_KEY>`
- Content-Type: `application/json`

请求体示例:

```json
{
  "query": "who is Leo Messi?",
  "search_depth": "basic",
  "max_results": 5,
  "topic": "general",
  "time_range": "week",
  "include_answer": false,
  "include_raw_content": false,
  "include_images": false,
  "include_image_descriptions": false,
  "include_favicon": false,
  "include_domains": [],
  "exclude_domains": [],
  "country": null,
  "auto_parameters": false,
  "include_usage": false
}
```

参数列表（与搜索相关的核心字段）:

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| query | string | 是 | 搜索查询 |
| search_depth | string | 否 | `advanced|basic|fast|ultra-fast` |
| max_results | int | 否 | 最大结果数（0-20） |
| topic | string | 否 | `general|news|finance` |
| time_range | string | 否 | `day|week|month|year` 或短写 `d|w|m|y` |
| start_date | string | 否 | 起始日期 `YYYY-MM-DD` |
| end_date | string | 否 | 结束日期 `YYYY-MM-DD` |
| include_answer | bool/string | 否 | `true|basic|advanced` |
| include_raw_content | bool/string | 否 | `true|markdown|text` |
| include_images | bool | 否 | 是否返回图片结果 |
| include_image_descriptions | bool | 否 | 图片描述 |
| include_favicon | bool | 否 | 返回 favicon |
| include_domains | string[] | 否 | 仅包含指定域名 |
| exclude_domains | string[] | 否 | 排除指定域名 |
| country | string | 否 | 仅 `topic=general` 生效 |
| auto_parameters | bool | 否 | 自动推断参数（可能提高成本） |
| include_usage | bool | 否 | 返回用量信息 |

### 1.2 Brave Web Search API

- Endpoint: `GET https://api.search.brave.com/res/v1/web/search`
- 认证: `X-Subscription-Token: <BRAVE_API_KEY>`
- Accept: `application/json`

请求示例:

```bash
curl -s "https://api.search.brave.com/res/v1/web/search?q=brave+search&count=10&offset=0" \
  -H "Accept: application/json" \
  -H "X-Subscription-Token: <BRAVE_API_KEY>"
```

参数列表（Query Params）:

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| q | string | 是 | 搜索查询（<= 400 chars, 50 words） |
| count | int | 否 | 结果数（上限 20） |
| offset | int | 否 | 分页偏移（0-9） |
| country | string | 否 | 两位国家码 |
| search_lang | string | 否 | 语言偏好 |
| ui_lang | string | 否 | UI 语言 |
| safesearch | string | 否 | `off|moderate|strict` |
| freshness | string | 否 | `pd|pw|pm|py` 或日期区间 |
| result_filter | string | 否 | `web,news,videos,locations,...` |
| extra_snippets | bool | 否 | 返回额外片段 |
| summary | bool | 否 | 返回 summarizer key |

可选定位 Header（地理相关）:

- `X-Loc-Lat` / `X-Loc-Long`
- `X-Loc-City` / `X-Loc-State` / `X-Loc-Country`

### 1.3 Jina Search API

- Endpoint: `POST https://s.jina.ai/`
- 认证: `Authorization: Bearer $JINA_API_KEY`
- 必须 Header: `Accept: application/json`
- Content-Type: `application/json`

请求体示例:

```json
{
  "q": "When was Jina AI founded?",
  "gl": "US",
  "hl": "en",
  "location": "San Francisco",
  "num": 10,
  "page": 1
}
```

参数列表（Body）:

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| q | string | 是 | 搜索查询 |
| gl | string | 否 | 国家（两位） |
| hl | string | 否 | 语言（两位） |
| location | string | 否 | 位置（推荐城市级） |
| num | number | 否 | 最大结果数 |
| page | number | 否 | 结果偏移（分页） |

常用控制 Header:

| Header | 说明 |
| --- | --- |
| X-Site | 站内搜索（单域名） |
| X-No-Cache | 绕过缓存 |
| X-With-Links-Summary | 返回链接汇总 |
| X-With-Images-Summary | 返回图片汇总 |
| X-With-Generated-Alt | 生成图片描述 |
| X-Return-Format | `markdown|html|text|screenshot|pageshot` |
| X-With-Favicon / X-With-Favicons | favicon |
| X-Engine | `browser|direct` |
| X-Respond-With | `no-content` 仅返回元信息 |
| X-Locale | UI 语言/地区 |

## 2. 统一接口定义（主要参考 Tavily）

### 2.1 统一请求结构

```json
{
  "query": "string",
  "search_depth": "basic|fast|ultra-fast|advanced",
  "max_results": 5,
  "offset": 0,
  "topic": "general|news|finance",
  "time_range": "day|week|month|year",
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD",
  "country": "US",
  "language": "en",
  "ui_language": "en-US",
  "location": "San Francisco",
  "result_types": ["web","news","videos","locations"],
  "safe_search": "off|moderate|strict",
  "include_answer": "none|basic|advanced",
  "include_content": "none|summary|full",
  "content_format": "markdown|text|html",
  "include_images": false,
  "include_image_descriptions": false,
  "include_favicon": false,
  "include_domains": ["example.com"],
  "exclude_domains": ["spam.com"],
  "site": "example.com",
  "use_cache": true,
  "include_usage": false
}
```

说明（与 Tavily 语义对齐）:

- `include_answer`: Tavily 的 `include_answer` 直接对齐。
- `include_content`: `full` 对应 Tavily `include_raw_content=true`，`summary` 对应 `include_raw_content=false` 但保留摘要。
- `content_format`: `markdown|text` 与 Tavily `include_raw_content` 的枚举对应。
- `search_depth` / `topic` / `time_range` 的枚举值参考 Tavily。

### 2.2 统一响应结构

```json
{
  "query": "string",
  "answer": "string",
  "results": [
    {
      "title": "string",
      "url": "string",
      "snippet": "string",
      "content": "string",
      "score": 0.0,
      "favicon": "string",
      "language": "string",
      "published_at": "string"
    }
  ],
  "images": [
    {
      "url": "string",
      "description": "string"
    }
  ],
  "usage": {
    "credits": 1
  },
  "response_time": 0.0,
  "request_id": "string"
}
```

## 3. 字段映射（统一标准 -> 三家 API）

| 统一字段 | Tavily | Brave | Jina |
| --- | --- | --- | --- |
| query | `query` | `q` | `q` |
| search_depth | `search_depth` | 不支持 | 不支持 |
| max_results | `max_results` | `count` | `num` |
| offset | 不支持 | `offset` | `page` |
| topic | `topic` | 近似 `result_filter=news` | 不支持 |
| time_range | `time_range` | `freshness` | 不支持 |
| start_date | `start_date` | `freshness` (date range) | 不支持 |
| end_date | `end_date` | `freshness` (date range) | 不支持 |
| country | `country` | `country` | `gl` |
| language | 不支持 | `search_lang` | `hl` |
| ui_language | 不支持 | `ui_lang` | `X-Locale` |
| location | 不支持 | `X-Loc-*` headers | `location` |
| result_types | 不支持 | `result_filter` | 不支持 |
| safe_search | 不支持 | `safesearch` | 不支持 |
| include_answer | `include_answer` | 不支持 | 不支持 |
| include_content | `include_raw_content` | 不支持 | `X-Respond-With` (no-content) |
| content_format | `include_raw_content` | 不支持 | `X-Return-Format` |
| include_images | `include_images` | 不支持 | `X-With-Images-Summary` |
| include_image_descriptions | `include_image_descriptions` | 不支持 | `X-With-Generated-Alt` |
| include_favicon | `include_favicon` | 不支持（响应里可能带 `meta_url.favicon`） | `X-With-Favicon` |
| include_domains | `include_domains` | 查询语法 `site:` | `X-Site`（单域名） |
| exclude_domains | `exclude_domains` | 不支持 | 不支持 |
| site | 不支持 | 查询语法 `site:` | `X-Site` |
| use_cache | 不支持 | `Cache-Control: no-cache` | `X-No-Cache` |
| include_usage | `include_usage` | 响应头 `X-RateLimit-*` | `usage`（响应） |

备注:

- Brave 的 `topic` 无直接参数，常用做法是用 `result_filter` 选择 `news` 或在 `q` 中加入限定。
- Jina 的 `include_content` 属于响应控制，使用 `X-Respond-With: no-content` 可关闭正文。
- `include_domains` / `site` 在 Brave/Jina 主要通过站内搜索限制实现。
