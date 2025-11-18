package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sc-player-service/model/api/response"
)

const (
	httpTimeout = 5 * time.Second
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func CallAPI[T any](c *Client, ctx context.Context,
	method, path string,
	body any, headers map[string]string,
) (*response.Resp, *T, error) {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal_request_body -> %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
		headers["Content-Type"] = "application/json; charset=utf-8"
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("create_http_request -> %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do_http_request -> %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read_response_body -> %w", err)
	}

	var res response.Resp
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, nil, fmt.Errorf("unmarshal_response -> %w", err)
	}

	var data T
	if res.Data != nil {
		dataBytes, err := json.Marshal(res.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal_data -> %w", err)
		}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return nil, nil, fmt.Errorf("unmarshal_data -> %w", err)
		}
	}

	return &res, &data, nil
}
