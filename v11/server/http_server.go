//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-server-actions-register-output=./http_server_actions_register.go
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// HTTPConfig HTTP 服务配置.
var (
	errInvalidFormData = errors.New("invalid form data")
	errInvalidJSON     = errors.New("invalid json")
	errUnsupportedCT   = errors.New("unsupported content type")
)

type HTTPConfig struct {
	Addr              string // 监听地址，例 ":5700"
	PathPrefix        string // 路由前缀，可为空或"/"
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	AccessToken       string // 可选鉴权，若为空则不校验
}

// HTTPServer 实现 OneBot HTTP 传输层.
type HTTPServer struct {
	srv     *http.Server
	mux     *http.ServeMux
	cfg     HTTPConfig
	handler ActionRequestHandler
}

// NewHTTPServer 创建 HTTPServer.若传入 mux 为 nil，则使用自建 ServeMux.
func NewHTTPServer(cfg HTTPConfig, handler ActionRequestHandler) *HTTPServer {
	mux := http.NewServeMux()
	cfg.PathPrefix = "/" + strings.Trim(cfg.PathPrefix, "/") + "/"
	server := &HTTPServer{cfg: cfg, mux: mux, handler: handler}
	mux.HandleFunc("/", server.handleRoot)

	server.srv = &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return server
}

// Start 启动 HTTP 服务器（异步监听）.
func (s *HTTPServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		return s.Shutdown(shutdownCtx)
	case err := <-errCh:
		return fmt.Errorf("http server listen and serve failed: %w", err)
	}
}

// Shutdown 优雅关闭.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("http server shutdown failed: %w", err)
	}

	return nil
}

// Handler 返回 http.Handler，便于挂载到外部路由.
func (s *HTTPServer) Handler() http.Handler {
	return s.mux
}

func (s *HTTPServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	action, ok := s.extractAction(r.URL.Path)
	if !ok {
		http.NotFound(w, r)

		return
	}

	if err := s.checkAccess(r); err != nil {
		http.Error(w, err.message, err.code)

		return
	}

	params, err := s.parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	req := &entity.ActionRequest{Action: action, Params: params}

	resp, err := s.handler.HandleActionRequest(r.Context(), req)
	if err != nil {
		s.writeError(w, err)

		return
	}

	if resp == nil {
		resp = &entity.ActionRawResponse{Status: entity.StatusFailed, Retcode: -1, Message: "empty response"}
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *HTTPServer) extractAction(path string) (string, bool) {
	if !strings.HasPrefix(path, s.cfg.PathPrefix) {
		return "", false
	}

	trimmed := strings.Trim(strings.TrimPrefix(path, s.cfg.PathPrefix), "/")
	if trimmed == "" {
		return "", false
	}

	return trimmed, true
}

type accessError struct {
	message string
	code    int
}

func (s *HTTPServer) checkAccess(r *http.Request) *accessError {
	if s.cfg.AccessToken == "" {
		return nil
	}

	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	} else if token == "" {
		token = r.URL.Query().Get("access_token")
	}

	if token == "" {
		return &accessError{message: "missing access token", code: http.StatusUnauthorized}
	}

	if token != s.cfg.AccessToken {
		return &accessError{message: "forbidden", code: http.StatusForbidden}
	}

	return nil
}

func (s *HTTPServer) parseParams(r *http.Request) (map[string]any, error) {
	params := make(map[string]any)

	err := r.ParseForm()
	if err != nil {
		return nil, errInvalidFormData
	}

	for k, v := range r.Form {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	if r.Method != http.MethodPost {
		return params, nil
	}

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		var m map[string]any

		dec := json.NewDecoder(r.Body)
		dec.UseNumber()

		err := dec.Decode(&m)
		if err != nil {
			return nil, errInvalidJSON
		}

		for k, v := range m {
			params[k] = v
		}

		return params, nil
	}

	if contentType != "" && !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		return nil, errUnsupportedCT
	}

	return params, nil
}

func (s *HTTPServer) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "encode response json failed", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrActionNotFound):
		http.NotFound(w, nil)
	case errors.Is(err, ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
