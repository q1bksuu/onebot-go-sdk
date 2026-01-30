[根目录](../../CLAUDE.md) > [v11](../) > **server**

---

# server - HTTP/WebSocket 服务端模块

OneBot 11 协议的服务端实现，支持 HTTP 和 WebSocket 双协议，用于接收 OneBot 实现发送的动作请求和事件上报。

---

## 变更记录 (Changelog)

### 2026-01-05

- **重大重构**: 架构升级，支持 HTTP/WebSocket 双协议
- **新增**: UnifiedServer - 统一服务器，HTTP 和 WS 共用端口
- **新增**: WebSocketServer - 正向 WebSocket 服务支持
- **新增**: EventDispatcher - 事件分发器
- **新增**: BaseServer - 通用服务器基础结构
- **变更**: ActionHandler/ActionRequestHandler 移至 `v11/dispatcher` 包
- **变更**: HTTPServer 使用 Option 模式，支持事件接收

### 2025-12-21 15:53:08

- **初始化**: 生成模块级文档

---

## 模块职责

server 模块负责：

1. **HTTP 服务端实现**: 接收和路由 OneBot 动作请求，支持事件上报接收
2. **WebSocket 服务端实现**: 支持 /api、/event、/ 三种端点
3. **统一服务器**: HTTP 和 WebSocket 共用同一端口，通过 Upgrade 头区分
4. **事件分发器**: 基于事件类型的多级路由匹配
5. **鉴权支持**: 可选的访问令牌验证（Header 或 Query）
6. **多种参数格式**: 支持 Query、Form、JSON 参数

---

## 入口与启动

### 方式一：使用统一服务器（推荐）

```go
import (
    "github.com/q1bksuu/onebot-go-sdk/v11/server"
    "github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
)

// 1. 创建 Action 分发器
actionDispatcher := dispatcher.NewDispatcher()
actionDispatcher.Register("send_private_msg", handleSendPrivateMsg)

// 2. 创建事件处理器（可选）
eventHandler := server.EventRequestHandlerFunc(func(ctx context.Context, event entity.Event) (map[string]any, error) {
    // 处理事件，返回快速操作
    return nil, nil
})

// 3. 配置统一服务器
cfg := server.UnifiedConfig{
    ServerConfig: server.ServerConfig{
        Addr: ":5700",
    },
    HTTP: server.UnifiedHTTPConfig{
        APIPathPrefix: "/",
        EventPath:     "/event",
        AccessToken:   "your-secret",
        ActionHandler: actionDispatcher,
        EventHandler:  eventHandler,
    },
    WS: server.UnifiedWSConfig{
        PathPrefix:  "/",
        AccessToken: "your-secret",
        ActionHandler: actionDispatcher,
    },
}

// 4. 创建并启动
srv := server.NewUnifiedServer(cfg)

ctx := context.Background()
if err := srv.Start(ctx); err != nil {
    log.Fatal(err)
}
```

### 方式二：仅使用 HTTP 服务器

```go
cfg := server.HTTPConfig{
    Addr:          ":5700",
    APIPathPrefix: "/",
    EventPath:     "/event",
    AccessToken:   "your-secret",
}

httpSrv := server.NewHTTPServer(
    server.WithHTTPConfig(cfg),
    server.WithActionHandler(actionDispatcher),
    server.WithEventHandler(eventHandler),
)

if err := httpSrv.Start(context.Background()); err != nil {
    log.Fatal(err)
}
```

### 方式三：仅使用 WebSocket 服务器

```go
cfg := server.WSConfig{
    Addr:        ":6700",
    PathPrefix:  "/",
    AccessToken: "your-secret",
}

wsSrv := server.NewWebSocketServer(
    server.WithWSConfig(cfg),
    server.WithWSActionHandler(actionDispatcher),
)

if err := wsSrv.Start(context.Background()); err != nil {
    log.Fatal(err)
}
```

---

## 对外接口

### 1. UnifiedServer (unified_server.go)

统一服务器，HTTP 和 WebSocket 共用同一端口。

**NewUnifiedServer** (unified_server.go:36-72)

```go
func NewUnifiedServer(
    cfg UnifiedConfig,
) *UnifiedServer
```

**UnifiedConfig 配置**:

```go
type UnifiedConfig struct {
    ServerConfig

    HTTP UnifiedHTTPConfig
    WS   UnifiedWSConfig
}

type UnifiedHTTPConfig struct {
    APIPathPrefix string
    EventPath     string
    AccessToken   string
    ActionHandler dispatcher.ActionRequestHandler
    EventHandler  EventRequestHandler
}

type UnifiedWSConfig struct {
    PathPrefix   string
    AccessToken  string
    CheckOrigin  func(r *http.Request) bool
    ActionHandler dispatcher.ActionRequestHandler
}
```

**方法**:

- `Start(ctx context.Context) error`: 启动服务器（阻塞，直到 context 取消）
- `Shutdown(ctx context.Context) error`: 优雅关闭

### 2. HTTPServer (http_server.go)

**NewHTTPServer** (http_server.go)

```go
func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer
```

**HTTPServerOption 选项**:

- `WithHTTPConfig(cfg HTTPConfig)`: 设置 HTTP 配置（覆盖）
- `WithAddr(addr string)`: 设置监听地址
- `WithAPIPathPrefix(prefix string)`: 设置 API 路由前缀
- `WithEventPath(path string)`: 设置事件路由
- `WithReadHeaderTimeout(timeout time.Duration)`: 设置 ReadHeaderTimeout
- `WithReadTimeout(timeout time.Duration)`: 设置 ReadTimeout
- `WithWriteTimeout(timeout time.Duration)`: 设置 WriteTimeout
- `WithIdleTimeout(timeout time.Duration)`: 设置 IdleTimeout
- `WithAccessToken(token string)`: 设置访问令牌
- `WithActionHandler(h dispatcher.ActionRequestHandler)`: 设置动作处理器
- `WithEventHandler(h EventRequestHandler)`: 设置事件处理器

**HTTPConfig 配置**:

```go
type HTTPConfig struct {
    Addr              string        // 监听地址，例 ":5700"
    APIPathPrefix     string        // API 接口路由前缀，可为空或"/"
    EventPath         string        // 事件接口路由，可为空或"/"
    ReadHeaderTimeout time.Duration
    ReadTimeout       time.Duration
    WriteTimeout      time.Duration
    IdleTimeout       time.Duration
    AccessToken       string        // 可选鉴权，若为空则不校验
}
```

**方法**:

- `Start(ctx context.Context) error`: 启动服务器
- `Handler() http.Handler`: 返回 HTTP 处理器（可挂载到外部路由）

### 3. WebSocketServer (websocket.go)

**NewWebSocketServer** (websocket.go)

```go
func NewWebSocketServer(opts ...WebSocketServerOption) *WebSocketServer
```

**WebSocketServerOption 选项**:

- `WithWSConfig(cfg WSConfig)`: 设置 WebSocket 配置（覆盖）
- `WithWSAddr(addr string)`: 设置监听地址
- `WithWSPathPrefix(prefix string)`: 设置路径前缀
- `WithWSAccessToken(token string)`: 设置访问令牌
- `WithWSCheckOrigin(fn func(*http.Request) bool)`: 设置跨域检查函数
- `WithWSReadTimeout(timeout time.Duration)`: 设置 ReadTimeout
- `WithWSWriteTimeout(timeout time.Duration)`: 设置 WriteTimeout
- `WithWSIdleTimeout(timeout time.Duration)`: 设置 IdleTimeout
- `WithWSActionHandler(h dispatcher.ActionRequestHandler)`: 设置动作处理器

**WSConfig 配置**:

```go
type WSConfig struct {
    Addr         string                     // 监听地址，例 ":6700"
    PathPrefix   string                     // 路径前缀，用于 /api、/event、/ 路由
    AccessToken  string                     // 可选鉴权，若为空则不校验
    CheckOrigin  func(r *http.Request) bool // 可选跨域校验，默认全放行
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}
```

**方法**:

- `Start(ctx context.Context) error`: 启动服务器
- `Shutdown(ctx context.Context) error`: 优雅关闭
- `Handler() http.Handler`: 返回 HTTP 处理器
- `BroadcastEvent(event entity.Event) error`: 向所有事件连接广播事件

**WebSocket 端点**:

- `{prefix}/api`: API 端点，用于接收动作请求
- `{prefix}/event`: 事件端点，用于推送事件
- `{prefix}` 或 `/`: 通用端点，同时支持 API 和事件

### 4. EventDispatcher (event_dispatcher.go)

事件分发器，基于事件类型进行多级路由匹配。

**NewEventDispatcher** (event_dispatcher.go:18-20)

```go
func NewEventDispatcher() *EventDispatcher
```

**方法**:

- `Register(key string, h EventHandler)`: 注册事件处理器
- `HandleEvent(ctx context.Context, event entity.Event) (map[string]any, error)`: 处理事件

**路由键格式**:

事件分发器支持多级匹配，从最具体到最通用：

```go
// 消息事件
"message"                           // 所有消息
"message.private"                   // 私聊消息
"message.private.friend"            // 好友私聊

// 通知事件
"notice"                            // 所有通知
"notice.group_increase"             // 群成员增加
"notice.group_increase.approve"     // 管理员同意入群

// 请求事件
"request"                           // 所有请求
"request.friend"                    // 好友请求

// 元事件
"meta_event"                        // 所有元事件
"meta_event.heartbeat"              // 心跳
```

### 5. BaseServer (base.go)

通用服务器基础结构，封装 http.Server 的启动与关闭逻辑。

```go
type BaseServer struct {
    Srv *http.Server
}

func NewBaseServer(cfg ServerConfig, handler http.Handler) *BaseServer
func (s *BaseServer) Start(ctx context.Context, onShutdown func(context.Context) error) error
func (s *BaseServer) Shutdown(ctx context.Context) error
```

### 6. 处理器接口 (handler.go)

```go
// 处理事件，返回快速操作响应（可选）
type EventHandler func(ctx context.Context, event entity.Event) (map[string]any, error)

// 事件请求处理器接口
type EventRequestHandler interface {
    HandleEvent(ctx context.Context, event entity.Event) (map[string]any, error)
}

// 适配函数类型
type EventRequestHandlerFunc func(ctx context.Context, event entity.Event) (map[string]any, error)
```

**注意**: `ActionHandler` 和 `ActionRequestHandler` 已移至 `v11/dispatcher` 包。

### 7. 错误定义 (errors.go)

```go
var (
    ErrBadRequest                = errors.New("bad request")
    ErrUniversalClientURLEmpty   = errors.New("universal client URL is empty")
    ErrMissingTypeField          = errors.New("missing type field")
    ErrUnknownEventType          = errors.New("unknown event type")
    ErrInvalidEventTreeStructure = errors.New("invalid event tree structure")
    ErrMissingOrInvalidPostType  = errors.New("missing or invalid post_type field")
    ErrUnknownPostType           = errors.New("unknown post_type")
    ErrNoEventHandler            = errors.New("no event handler")
)
```

### 8. 自动生成的代码

**http_server_actions_register.go** (~550行)

通过 `//go:generate go run ../cmd/bindings-gen` 生成，提供动作处理服务接口。

**http_server_events_register.go** (~205行)

通过 `//go:generate go run ../cmd/event-bindings-gen` 生成，提供事件处理服务接口：

```go
type MessageEventService interface {
    HandlePrivateMessage(ctx context.Context, ev *entity.PrivateMessageEvent) (map[string]any, error)
    HandleGroupMessage(ctx context.Context, ev *entity.GroupMessageEvent) (map[string]any, error)
}

type NoticeEventService interface {
    HandleGroupFileUpload(ctx context.Context, ev *entity.GroupFileUploadEvent) (map[string]any, error)
    // ... 更多处理方法
}

type RequestEventService interface { ... }
type MetaEventService interface { ... }
```

---

## 关键依赖与配置

### 内部依赖

- `github.com/q1bksuu/onebot-go-sdk/v11/entity`: 协议实体定义
- `github.com/q1bksuu/onebot-go-sdk/v11/dispatcher`: 动作分发器
- `github.com/q1bksuu/onebot-go-sdk/v11/internal/util`: 工具函数

### 外部依赖

- `github.com/gorilla/websocket`: WebSocket 协议支持

### 代码生成配置

**触发方式**:

在 `http_server.go` 文件顶部：

```go
//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-server-actions-register-output=./http_server_actions_register.go
//go:generate go run ../cmd/event-bindings-gen -config=../cmd/event-bindings-gen/config.yaml -output=./http_server_events_register.go
```

---

## 数据模型

### 核心类型

```go
type BaseServer struct {
    Srv *http.Server
}

type HTTPServer struct {
    *BaseServer
    mux           *http.ServeMux
    cfg           HTTPConfig
    actionHandler dispatcher.ActionRequestHandler
    eventHandler  EventRequestHandler
}

type WebSocketServer struct {
    *BaseServer
    cfg           WSConfig
    handler       dispatcher.ActionRequestHandler
    upgrader      websocket.Upgrader
    mu            sync.Mutex
    eventConns    map[*wsConn]struct{}  // 事件连接池
    universalConn map[*wsConn]struct{}  // 通用连接池
}

type UnifiedServer struct {
    *BaseServer
    httpSrv *HTTPServer
    wsSrv   *WebSocketServer
}

type EventDispatcher struct {
    handlers map[string]EventHandler
}
```

### HTTP 请求处理流程

```
1. HTTP 请求到达 (如 POST /send_private_msg)
   ↓
2. handleRoot (http_server.go)
   ├─ extractAction: 从 URL 路径提取 action
   ├─ checkAccess: 验证访问令牌（如果配置了）
   └─ parseParams: 解析参数（Query + Form + JSON Body）
   ↓
3. 调用 actionHandler.HandleActionRequest
   ↓
4. 返回 JSON 响应
```

### HTTP 事件处理流程

```
1. HTTP POST 到达 EventPath (如 POST /event)
   ↓
2. handleEvent (http_server.go)
   ├─ checkAccess: 验证访问令牌
   ├─ parseEvent: 解析事件 JSON 为具体类型
   └─ 调用 eventHandler.HandleEvent
   ↓
3. 返回快速操作响应或 204 No Content
```

### WebSocket 处理流程

```
1. WebSocket 连接建立 (如 ws://host/api)
   ↓
2. 根据路径分发
   ├─ /api: handleAPI - 接收动作请求，返回响应
   ├─ /event: handleEvent - 接收事件推送
   └─ /: handleUniversal - 同时支持 API 和事件
   ↓
3. serveActionConn: 循环读取消息
   ├─ handleActionMessage: 解析并处理动作请求
   └─ 写入响应
```

### 统一服务器路由

UnifiedServer 使用 `combinedHandler` 在同一端口处理 HTTP 和 WebSocket：

```
请求到达
   ↓
检查 Upgrade: websocket 头
   ├─ 是: 路由到 WebSocketServer.Handler()
   └─ 否: 路由到 HTTPServer.Handler()
```

### 鉴权机制

支持两种传递方式（优先级：Header > Query）：

1. **HTTP Header**: `Authorization: Bearer {token}`
2. **URL Query**: `?access_token={token}`

---

## 测试与质量

### 测试文件

- `http_server_test.go`: 17 个单元测试
- `websocket_test.go`: 6 个单元测试
- `unified_server_test.go`: 5 个单元测试

### 测试场景

#### http_server_test.go

| 测试函数 | 测试内容 |
| --- | --- |
| `TestNewHTTPServer_PathPrefixNormalizeAndHandler` | 路径前缀标准化和 Handler 获取 |
| `TestHTTPServer_HandleRoot_PathAndNotFound` | 路由匹配和 404 处理 |
| `TestHTTPServer_AuthRequired_MissingOrWrongToken` | 鉴权失败场景（401/403） |
| `TestHTTPServer_AuthRequired_WithHeaderAndQuery` | 鉴权成功场景 |
| `TestHTTPServer_Params_QueryAndFormAndJSON` | 多种参数格式解析 |
| `TestHTTPServer_Params_InvalidForm` | 无效 Form 数据处理 |
| `TestHTTPServer_JSON_InvalidOrUnsupportedContentType` | 无效 JSON 处理 |
| `TestHTTPServer_WriteError_Mapping` | 错误映射到 HTTP 状态码 |
| `TestHTTPServer_NilResponse_DefaultFailed` | 空响应默认返回 failed |
| `TestHTTPServer_StartAndShutdown_ContextCancel` | 服务器启动和优雅关闭 |
| `TestHTTPServer_EventPath_Registration` | 事件路径注册 |
| `TestHTTPServer_EventPath_NoHandler_Returns204` | 无事件处理器返回 204 |
| `TestHTTPServer_EventPath_OnlyAcceptsPOST` | 事件端点仅接受 POST |
| `TestHTTPServer_EventPath_InvalidJSON_Returns400` | 无效 JSON 返回 400 |
| `TestHTTPServer_EventPath_EmptyQuickOp_Returns204` | 空快速操作返回 204 |
| `TestHTTPServer_EventPath_AllEventTypes` | 所有事件类型解析 |
| `TestEventDispatcher_Routing` | 事件分发器路由 |

#### websocket_test.go

| 测试函数 | 测试内容 |
| --- | --- |
| `TestNormalizeAndMatchPath` | 路径标准化和匹配 |
| `TestCheckAccess` | WebSocket 鉴权 |
| `TestHandleActionMessageMapping` | 动作消息处理映射 |
| `TestWriteHandshakeError` | 握手错误写入 |
| `TestHandleAPIAndUniversalFlow` | API 和通用端点流程 |
| `TestBroadcastEventIntegration` | 事件广播集成测试 |

#### unified_server_test.go

| 测试函数 | 测试内容 |
| --- | --- |
| `TestUnifiedServer_Initialization` | 统一服务器初始化 |
| `TestUnifiedServer_Routing_HTTP` | HTTP 路由 |
| `TestUnifiedServer_Routing_WebSocket` | WebSocket 路由 |
| `TestUnifiedServer_Routing_UniversalPath_DistinguishByProtocol` | 协议区分 |
| `TestUnifiedServer_StartAndShutdown` | 启动和关闭 |

### 质量保证

- **Mock 测试**: 使用 `httptest.NewRecorder` 和 `httptest.NewServer`
- **WebSocket 测试**: 使用 `gorilla/websocket` 客户端
- **边界条件**: 覆盖空参数、无效 JSON、未授权等场景
- **集成测试**: 完整的服务器生命周期测试

---

## 常见问题 (FAQ)

**Q: 如何选择使用哪种服务器？**

- **UnifiedServer**: 推荐，HTTP 和 WebSocket 共用端口，简化部署
- **HTTPServer**: 仅需 HTTP API 和事件接收时使用
- **WebSocketServer**: 仅需 WebSocket 通信时使用

**Q: 如何处理事件上报？**

1. **HTTP 方式**: 配置 `HTTPConfig.EventPath` 并使用 `WithEventHandler`
2. **WebSocket 方式**: 客户端连接 `/event` 端点，服务端调用 `BroadcastEvent`

```go
// HTTP 事件处理
eventHandler := server.EventRequestHandlerFunc(func(ctx context.Context, event entity.Event) (map[string]any, error) {
    switch e := event.(type) {
    case *entity.PrivateMessageEvent:
        // 返回快速操作
        return map[string]any{"reply": "收到消息"}, nil
    }
    return nil, nil
})
```

**Q: 如何使用 EventDispatcher？**

```go
ed := server.NewEventDispatcher()

// 注册处理器（从具体到通用）
ed.Register("message.private.friend", handleFriendMessage)
ed.Register("message.private", handlePrivateMessage)
ed.Register("message", handleAllMessage)

// 作为 EventRequestHandler 使用
httpSrv := server.NewHTTPServer(server.WithHTTPConfig(cfg), server.WithEventHandler(ed))
```

**Q: 如何集成到现有 HTTP 服务器？**

```go
// 获取 Handler
handler := httpServer.Handler()

// 挂载到现有路由
existingMux := http.NewServeMux()
existingMux.Handle("/onebot/", http.StripPrefix("/onebot", handler))
```

**Q: ActionHandler 去哪了？**

`ActionHandler`、`ActionRequestHandler`、`Dispatcher`、`Binder` 已移至独立的 `v11/dispatcher` 包：

```go
import "github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"

d := dispatcher.NewDispatcher()
d.Register("send_private_msg", handler)

binder := dispatcher.NewBinder("send_group_msg", handleSendGroupMsg)
d.Register(binder.Action(), binder.Handler())
```

---

## 相关文件清单

### 主要源文件

| 文件 | 行数 | 职责 |
| --- | --- | --- |
| `http_server.go` | ~485 | HTTP 服务端核心实现 |
| `websocket.go` | ~402 | WebSocket 服务端实现 |
| `event_dispatcher.go` | ~168 | 事件分发器 |
| `unified_server.go` | ~114 | 统一服务器（HTTP+WS） |
| `base.go` | ~82 | 通用服务器基础结构 |
| `handler.go` | ~23 | 事件处理器接口定义 |
| `errors.go` | ~23 | 错误定义 |
| `http_server_actions_register.go` | ~550 | 自动生成的动作注册代码 |
| `http_server_events_register.go` | ~205 | 自动生成的事件服务接口 |

### 测试文件

| 文件 | 测试数 |
| --- | --- |
| `http_server_test.go` | 17 |
| `websocket_test.go` | 6 |
| `unified_server_test.go` | 5 |

---

*模块文档更新时间: 2026-01-05*
