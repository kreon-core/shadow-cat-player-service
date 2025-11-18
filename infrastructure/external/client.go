package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kreon-core/shadow-cat-common/logc"

	"sc-player-service/infrastructure/config"
	"sc-player-service/model/api/response"
)

const (
	httpTimeout = 5 * time.Second
)

type HTTPClient struct {
	BaseURL string
	Paths   map[string]string
	Client  *http.Client
}

func NewClient(clientConfig *config.Client) *HTTPClient {
	return &HTTPClient{
		BaseURL: clientConfig.BaseURL,
		Paths:   clientConfig.Paths,
		Client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func CallAPI[T any](ctx context.Context, c *HTTPClient,
	method, path string,
	body any, headers map[string]string,
) (int, *response.Resp, *T, error) {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return 0, nil, nil, fmt.Errorf("marshal_request_body -> %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
		headers["Content-Type"] = "application/json; charset=utf-8"
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("create_http_request -> %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logc.Info().Str("method", method).Str("url", url).
		Any("body", body).Any("headers", headers).
		Msg("Making HTTP request to external service")

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("do_http_request -> %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("read_response_body -> %w", err)
	}

	var res response.Resp
	if err := json.Unmarshal(resBody, &res); err != nil {
		return 0, nil, nil, fmt.Errorf("unmarshal_response -> %w", err)
	}

	var data T
	if res.Data != nil {
		dataBytes, err := json.Marshal(res.Data)
		if err != nil {
			return 0, nil, nil, fmt.Errorf("marshal_data -> %w", err)
		}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return 0, nil, nil, fmt.Errorf("unmarshal_data -> %w", err)
		}
	}

	return resp.StatusCode, &res, &data, nil
}
