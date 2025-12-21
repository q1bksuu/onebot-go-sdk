//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-client-actions-output=./http_client_actions.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
)

var (
	errBaseURLEmpty          = errors.New("baseURL is empty")
	errUnsupportedHTTPMethod = errors.New("unsupported http method")
	errHTTPStatus            = errors.New("unexpected http status")
	errMissingSchemeOrHost   = errors.New("missing scheme or host")
)

const maxErrorBodyBytes = 1024

type HTTPClient struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

type clientOptions struct {
	httpClient  *http.Client
	accessToken string
	pathPrefix  string
	timeout     time.Duration
}

// NewHTTPClient 创建 HTTP 客户端封装.
func NewHTTPClient(baseURL string, opts ...Option) (*HTTPClient, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("%w", errBaseURLEmpty)
	}

	options := clientOptions{timeout: 30 * time.Second}
	for _, opt := range opts {
		opt(&options)
	}

	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: options.timeout}
	}

	return &HTTPClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		accessToken: options.accessToken,
		httpClient:  httpClient,
	}, nil
}

func (c *HTTPClient) do(
	ctx context.Context,
	urlPath string,
	defaultMethod string,
	req any,
	opts ...CallOption,
) (*entity.ActionRawResponse, error) {
	options := callOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	method := resolveMethod(defaultMethod, options.methodOverride)

	err := validateMethod(method)
	if err != nil {
		return nil, err
	}

	u, err := buildTargetURL(c.baseURL, urlPath)
	if err != nil {
		return nil, err
	}

	params, err := encodeToParams(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	bodyReader, err := prepareRequestBody(method, u, params, options.query, req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	applyRequestHeaders(httpReq, method, c.accessToken, options.headers)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return parseActionResponse(resp, urlPath)
}

func resolveMethod(defaultMethod, override string) string {
	method := defaultMethod
	if override != "" {
		method = override
	}

	if method == "" {
		method = http.MethodPost
	}

	return method
}

func validateMethod(method string) error {
	if method != http.MethodGet && method != http.MethodPost {
		return fmt.Errorf("%w: %s", errUnsupportedHTTPMethod, method)
	}

	return nil
}

func prepareRequestBody(
	method string, u *url.URL, params map[string]any, query url.Values, req any,
) (io.Reader, error) {
	switch method {
	case http.MethodPost:
		return buildPostBody(u, query, req)
	case http.MethodGet:
		applyGetParamsWithQuery(u, params, query)

		return nil, nil //nolint:nilnil
	default:
		return nil, fmt.Errorf("%w: unsupported method %s", errUnsupportedHTTPMethod, method)
	}
}

func applyRequestHeaders(httpReq *http.Request, method, accessToken string, customHeaders map[string][]string) {
	if method == http.MethodPost {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	if accessToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	}

	for k, vs := range customHeaders {
		for _, v := range vs {
			httpReq.Header.Add(k, v)
		}
	}
}

func buildTargetURL(baseURL, urlPath string) (*url.URL, error) {
	base := strings.TrimRight(baseURL, "/")
	pathPart := strings.TrimLeft(urlPath, "/")

	target := base
	if pathPart != "" {
		target += "/" + pathPart
	}

	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("build url: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("build url: %w", &url.Error{Op: "parse", URL: target, Err: errMissingSchemeOrHost})
	}

	return u, nil
}

func applyGetParamsWithQuery(u *url.URL, params map[string]any, query url.Values) {
	qs := u.Query()
	mergeParams(qs, params)
	mergeValues(qs, query)
	u.RawQuery = qs.Encode()
}

func buildPostBody(u *url.URL, query url.Values, req any) (io.Reader, error) {
	qs := u.Query()
	mergeValues(qs, query)
	u.RawQuery = qs.Encode()

	bs, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return bytes.NewReader(bs), nil
}

func parseActionResponse(resp *http.Response, urlPath string) (*entity.ActionRawResponse, error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))

		return nil, fmt.Errorf("%w: %d %s", errHTTPStatus, resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var rawResponse entity.ActionRawResponse

	err := json.NewDecoder(resp.Body).Decode(&rawResponse)
	if err != nil {
		return nil, fmt.Errorf("decode action response: %w", err)
	}

	if rawResponse.Status == entity.StatusFailed || rawResponse.Retcode != 0 {
		return nil, &entity.ActionError{
			UrlPath: urlPath,
			Status:  rawResponse.Status,
			Retcode: rawResponse.Retcode,
			Message: rawResponse.Message,
		}
	}

	return &rawResponse, nil
}

func encodeToParams(req any) (map[string]any, error) {
	var m map[string]any

	err := util.JsonTagMapping(req, &m)
	if err != nil {
		return nil, fmt.Errorf("mapstructure mapping failed: %w", err)
	}

	return m, nil
}

func mergeParams(values url.Values, params map[string]any) {
	for k, val := range params {
		switch valType := val.(type) {
		case []any:
			for _, vv := range valType {
				values.Add(k, fmt.Sprint(vv))
			}
		case []string:
			for _, vv := range valType {
				values.Add(k, vv)
			}
		default:
			values.Add(k, fmt.Sprint(valType))
		}
	}
}

func mergeValues(dst url.Values, src url.Values) {
	if src == nil {
		return
	}

	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}
