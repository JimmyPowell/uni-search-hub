package tavily

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// ProxySearch Tavily 原生 /search 透传（仅替换 Authorization）。
func (c *Client) ProxySearch(ctx context.Context, apiKey string, rawBody []byte, contentType string) (status int, respContentType string, respBody []byte, err error) {
	if apiKey == "" {
		return 0, "", nil, fmt.Errorf("tavily api key 为空")
	}
	if c == nil {
		return 0, "", nil, fmt.Errorf("tavily client 为空")
	}
	url := c.BaseURL + "/search"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rawBody))
	if err != nil {
		return 0, "", nil, err
	}
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", nil, err
	}
	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, err
	}
	return resp.StatusCode, resp.Header.Get("Content-Type"), respBody, nil
}

