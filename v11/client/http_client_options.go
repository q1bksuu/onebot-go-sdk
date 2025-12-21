package client

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Option func(*clientOptions)

// WithHTTPClient 注入自定义 http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(o *clientOptions) { o.httpClient = hc }
}

// WithAccessToken 设置访问令牌，发送时附加 Authorization Bearer 头.
func WithAccessToken(token string) Option {
	return func(o *clientOptions) { o.accessToken = token }
}

// WithPathPrefix 设置路径前缀，例如 "bot" 或 "/bot/".
func WithPathPrefix(prefix string) Option {
	return func(o *clientOptions) { o.pathPrefix = prefix }
}

// WithTimeout 设置默认超时（仅在未提供 http.Client 时生效）.
func WithTimeout(d time.Duration) Option {
	return func(o *clientOptions) { o.timeout = d }
}

type CallOption func(*callOptions)

type callOptions struct {
	headers        http.Header
	query          url.Values
	methodOverride string
}

// WithHeader 为单次调用追加自定义 Header.
func WithHeader(key, value string) CallOption {
	return func(co *callOptions) {
		if co.headers == nil {
			co.headers = make(http.Header)
		}

		co.headers.Add(key, value)
	}
}

// WithQuery 为单次调用追加 query 参数.
func WithQuery(key, value string) CallOption {
	return func(co *callOptions) {
		if co.query == nil {
			co.query = make(url.Values)
		}

		co.query.Add(key, value)
	}
}

// WithMethod 覆盖默认 HTTP 方法（GET/POST）.
func WithMethod(method string) CallOption {
	return func(co *callOptions) {
		co.methodOverride = strings.ToUpper(strings.TrimSpace(method))
	}
}
