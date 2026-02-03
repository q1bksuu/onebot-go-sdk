//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-server-actions-register-output=./http_server_actions_register.gen.go
//go:generate go run ../cmd/event-bindings-gen -config=../cmd/event-bindings-gen/config.yaml -http-server-events-register-output=./http_server_events_register.gen.go
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
)

// HTTPConfig HTTP 服务配置.
var (
	errInvalidFormData = errors.New("invalid form data")
	errInvalidJSON     = errors.New("invalid json")
	errUnsupportedCT   = errors.New("unsupported content type")
)

type HTTPConfig struct {
	Addr              string // 监听地址，例 ":5700"
	APIPathPrefix     string // api接口路由前缀，可为空或"/"
	EventPath         string // 事件接口路由，可为空或"/"
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	AccessToken       string // 可选鉴权，若为空则不校验
}

// HTTPServer 实现 OneBot HTTP 传输层.
type HTTPServer struct {
	*BaseServer

	mux           *http.ServeMux
	cfg           HTTPConfig
	actionHandler dispatcher.ActionRequestHandler
	eventHandler  EventRequestHandler // 可选的事件处理器
}

// HTTPServerOption 用于配置 HTTPServer 的选项函数类型.
type HTTPServerOption func(*HTTPServer)

// WithHTTPConfig 设置 HTTP 服务配置（会覆盖之前的配置）.
func WithHTTPConfig(cfg HTTPConfig) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg = cfg
	}
}

// WithAddr 设置监听地址.
func WithAddr(addr string) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.Addr = addr
	}
}

// WithAPIPathPrefix 设置 API 路由前缀.
func WithAPIPathPrefix(prefix string) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.APIPathPrefix = prefix
	}
}

// WithEventPath 设置事件路由.
func WithEventPath(path string) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.EventPath = path
	}
}

// WithReadHeaderTimeout 设置 ReadHeaderTimeout.
func WithReadHeaderTimeout(timeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.ReadHeaderTimeout = timeout
	}
}

// WithReadTimeout 设置 ReadTimeout.
func WithReadTimeout(timeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.ReadTimeout = timeout
	}
}

// WithWriteTimeout 设置 WriteTimeout.
func WithWriteTimeout(timeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.WriteTimeout = timeout
	}
}

// WithIdleTimeout 设置 IdleTimeout.
func WithIdleTimeout(timeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.IdleTimeout = timeout
	}
}

// WithAccessToken 设置访问令牌.
func WithAccessToken(token string) HTTPServerOption {
	return func(s *HTTPServer) {
		s.cfg.AccessToken = token
	}
}

// WithActionHandler 设置动作请求处理器选项.
func WithActionHandler(actionHandler dispatcher.ActionRequestHandler) HTTPServerOption {
	return func(s *HTTPServer) {
		s.actionHandler = actionHandler
	}
}

// WithEventHandler 设置事件处理器选项.
func WithEventHandler(eventHandler EventRequestHandler) HTTPServerOption {
	return func(s *HTTPServer) {
		s.eventHandler = eventHandler
	}
}

// NewHTTPServer 创建 HTTPServer，配置由 opts 提供.
func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	mux := http.NewServeMux()

	server := &HTTPServer{cfg: HTTPConfig{}, mux: mux}

	// 应用选项（顺序生效，后者覆盖前者）
	for _, opt := range opts {
		opt(server)
	}

	trimmedPrefix := strings.Trim(server.cfg.APIPathPrefix, "/")
	if trimmedPrefix == "" {
		server.cfg.APIPathPrefix = "/"
	} else {
		server.cfg.APIPathPrefix = "/" + trimmedPrefix + "/"
	}

	mux.HandleFunc("/", server.handleRoot)

	// 如果配置了 EventPath，注册事件路由
	if server.cfg.EventPath != "" {
		eventPath := util.NormalizePath(server.cfg.EventPath)
		mux.HandleFunc(eventPath, server.handleEvent)
	}

	baseCfg := ServerConfig{
		Addr:              server.cfg.Addr,
		ReadHeaderTimeout: server.cfg.ReadHeaderTimeout,
		ReadTimeout:       server.cfg.ReadTimeout,
		WriteTimeout:      server.cfg.WriteTimeout,
		IdleTimeout:       server.cfg.IdleTimeout,
	}
	server.BaseServer = NewBaseServer(baseCfg, mux)

	return server
}

// Start 启动 HTTP 服务器（异步监听）.
func (s *HTTPServer) Start(ctx context.Context) error {
	return s.BaseServer.Start(ctx, nil)
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

	resp, err := s.actionHandler.HandleActionRequest(r.Context(), req)
	if err != nil {
		s.writeError(w, err)

		return
	}

	if resp == nil {
		resp = &entity.ActionRawResponse{Status: entity.StatusFailed, Retcode: -1, Message: "empty response"}
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *HTTPServer) handleEvent(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// 检查访问权限
	if err := s.checkAccess(r); err != nil {
		http.Error(w, err.message, err.code)

		return
	}

	// 解析事件 JSON
	event, err := s.parseEvent(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// 如果没有事件处理器，返回 204
	if s.eventHandler == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	// 调用事件处理器
	quickOp, err := s.eventHandler.HandleEvent(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// 如果没有快速操作，返回 204
	if len(quickOp) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	// 返回快速操作 JSON
	s.writeJSON(w, http.StatusOK, quickOp)
}

// rawEvent 原始事件结构体，包含所有可能的事件字段（冗余字段）.
type rawEvent struct {
	Time          int64                   `json:"time"`
	SelfId        int64                   `json:"self_id"`
	PostType      entity.EventPostType    `json:"post_type"`
	MessageType   entity.EventMessageType `json:"message_type,omitempty"`
	NoticeType    entity.EventNoticeType  `json:"notice_type,omitempty"`
	RequestType   entity.EventRequestType `json:"request_type,omitempty"`
	MetaEventType entity.EventMetaType    `json:"meta_event_type,omitempty"`
	SubType       string                  `json:"sub_type,omitempty"`
}

// eventTreeValue 存储第二层树的值，可以是直接的事件构造函数或第三层树.
type eventTreeValue struct {
	constructor     func() entity.Event
	nextLevelGetter func(event *rawEvent) string
	subTree         *util.RadixTreeStrKey[*eventTreeValue]
}

var (
	// 使用多层嵌套 Radix Tree 实现前缀匹配，完全避免字符串拼接和分配
	// 第一层：按 post_type 分类，第二层：按 type_field 分类，第三层（可选）：按 sub_type 分类.
	// 这些全局变量用于延迟初始化，确保只初始化一次
	//nolint:gochecknoglobals // 延迟初始化的全局变量是合理的
	eventRadixTree *util.RadixTree[entity.EventPostType, *eventTreeValue]
	//nolint:gochecknoglobals // 延迟初始化的全局变量是合理的
	eventRadixOnce sync.Once
)

// initEventTypeMap 初始化事件类型多层 Radix Tree，使用高性能前缀匹配，完全避免字符串拼接.
func initEventTypeMap() {
	// 第一层：按 post_type 分类
	eventRadixTree = util.NewRadixTreeFromMap(map[entity.EventPostType]*eventTreeValue{
		// message 类型：message.private, message.group
		entity.EventPostTypeMessage: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.MessageType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventMessageTypePrivate): {
					constructor: func() entity.Event { return &entity.PrivateMessageEvent{} }},
				string(entity.EventMessageTypeGroup): {
					constructor: func() entity.Event { return &entity.GroupMessageEvent{} }},
			}),
		},
		// notice 类型
		entity.EventPostTypeNotice: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.NoticeType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventNoticeTypeGroupUpload): {
					constructor: func() entity.Event { return &entity.GroupFileUploadEvent{} }},
				string(entity.EventNoticeTypeGroupAdmin): {
					constructor: func() entity.Event { return &entity.GroupAdminChangeEvent{} }},
				string(entity.EventNoticeTypeGroupDecrease): {
					constructor: func() entity.Event { return &entity.GroupMemberDecreaseEvent{} }},
				string(entity.EventNoticeTypeGroupIncrease): {
					constructor: func() entity.Event { return &entity.GroupMemberIncreaseEvent{} }},
				string(entity.EventNoticeTypeGroupBan): {
					constructor: func() entity.Event { return &entity.GroupBanEvent{} }},
				string(entity.EventNoticeTypeFriendAdd): {
					constructor: func() entity.Event { return &entity.FriendAddEvent{} }},
				string(entity.EventNoticeTypeGroupRecall): {
					constructor: func() entity.Event { return &entity.GroupRecallEvent{} }},
				string(entity.EventNoticeTypeFriendRecall): {
					constructor: func() entity.Event { return &entity.FriendRecallEvent{} }},
				string(entity.EventNoticeTypeNotify): {
					subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
						string(entity.EventNoticeSubTypeGroupPoke): {
							constructor: func() entity.Event { return &entity.GroupPokeEvent{} }},
						string(entity.EventNoticeSubTypeGroupLuckyKing): {
							constructor: func() entity.Event { return &entity.GroupLuckyKingEvent{} }},
						string(entity.EventNoticeSubTypeGroupHonor): {
							constructor: func() entity.Event { return &entity.GroupHonorChangeEvent{} }},
					})},
			}),
		},
		// request 类型：request.friend, request.group
		entity.EventPostTypeRequest: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.RequestType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventRequestTypeFriend): {
					constructor: func() entity.Event { return &entity.FriendRequestEvent{} }},
				string(entity.EventRequestTypeGroup): {
					constructor: func() entity.Event { return &entity.GroupRequestEvent{} }},
			}),
		},
		// meta_event 类型：meta_event.lifecycle, meta_event.heartbeat
		entity.EventPostTypeMetaEvent: {
			nextLevelGetter: func(event *rawEvent) string {
				return string(event.MetaEventType)
			},
			subTree: util.NewRadixTreeFromMap(map[string]*eventTreeValue{
				string(entity.EventMetaTypeLifecycle): {
					constructor: func() entity.Event { return &entity.LifecycleEvent{} }},
				string(entity.EventMetaTypeHeartbeat): {
					constructor: func() entity.Event { return &entity.HeartbeatEvent{} }},
			}),
		},
	})
}

// findEventConstructor 递归查找事件构造函数.
func findEventConstructor(treeValue *eventTreeValue, raw *rawEvent, path []string) (func() entity.Event, error) {
	// 如果直接有构造函数，直接返回
	if treeValue.constructor != nil {
		return treeValue.constructor, nil
	}

	// 如果有子树，继续递归查找
	if treeValue.subTree != nil && treeValue.nextLevelGetter != nil {
		nextKey := treeValue.nextLevelGetter(raw)
		if nextKey == "" {
			return nil, fmt.Errorf("%w at path: %s", ErrMissingTypeField, strings.Join(path, "."))
		}

		nextValue, ok := treeValue.subTree.Get(nextKey)
		if !ok {
			return nil, fmt.Errorf("%w: %s.%s", ErrUnknownEventType, strings.Join(path, "."), nextKey)
		}

		return findEventConstructor(nextValue, raw, append(path, nextKey))
	}

	return nil, fmt.Errorf("%w at path: %s", ErrInvalidEventTreeStructure, strings.Join(path, "."))
}

// parseEvent 解析事件 JSON 为具体的事件类型.
// 使用多层嵌套 Radix Tree 实现前缀匹配，完全避免字符串拼接，通过递归查找实现高性能.
func (s *HTTPServer) parseEvent(r *http.Request) (entity.Event, error) {
	// 读取完整的 Body 内容
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// 解析为原始事件结构体（包含冗余字段）以获取类型信息
	var raw rawEvent

	err = json.Unmarshal(bodyBytes, &raw)
	if err != nil {
		return nil, fmt.Errorf("invalid event json: %w", err)
	}

	// 验证 post_type 是否存在
	if raw.PostType == "" {
		return nil, ErrMissingOrInvalidPostType
	}

	// 确保 Radix Tree 已初始化（只初始化一次）
	eventRadixOnce.Do(initEventTypeMap)

	// 从第一层开始递归查找
	firstLayerValue, ok := eventRadixTree.Get(raw.PostType)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPostType, raw.PostType)
	}

	// 递归查找构造函数
	constructor, err := findEventConstructor(firstLayerValue, &raw, []string{string(raw.PostType)})
	if err != nil {
		return nil, err
	}

	// 创建事件实例
	event := constructor()

	// 统一解析 JSON 到具体事件类型
	err = json.Unmarshal(bodyBytes, event)
	if err != nil {
		return nil, fmt.Errorf("parse event failed: %w", err)
	}

	return event, nil
}

func (s *HTTPServer) extractAction(path string) (string, bool) {
	if !strings.HasPrefix(path, s.cfg.APIPathPrefix) {
		return "", false
	}

	trimmed := strings.Trim(strings.TrimPrefix(path, s.cfg.APIPathPrefix), "/")
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
	case errors.Is(err, dispatcher.ErrActionNotFound):
		http.NotFound(w, nil)
	case errors.Is(err, ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
