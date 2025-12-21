[根目录](../../CLAUDE.md) > [v11](../) > **server**

---

# server - HTTP 服务端模块

OneBot 11 协议的 HTTP 服务端实现，用于接收 OneBot 实现发送的动作请求。

---

## 变更记录 (Changelog)

### 2025-12-21 15:53:08

- **初始化**: 生成模块级文档
- **覆盖**: 扫描了 HTTP 服务端、分发器、绑定器及自动生成的注册代码

---

## 模块职责

server 模块负责：

1. **HTTP 服务端实现**: 接收和路由 OneBot 动作请求
2. **分发器机制**: 基于 action 名称路由到对应处理器
3. **绑定器模式**: 类型安全的参数绑定和业务逻辑封装
4. **鉴权支持**: 可选的访问令牌验证（Header 或 Query）
5. **多种参数格式**: 支持 Query、Form、JSON 参数

---

## 入口与启动

### 创建服务器

```go
import "github.com/q1bksuu/onebot-go-sdk/v11/server"

// 1. 创建分发器
dispatcher := server.NewDispatcher()

// 2. 注册处理器（方式一：直接注册）
dispatcher.Register("send_private_msg", func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error) {
    // 处理逻辑
    return &entity.ActionRawResponse{
        Status:  entity.StatusOK,
        Retcode: 0,
        Data:    json.RawMessage(`{"message_id": 12345}`),
    }, nil
})

// 3. 注册处理器（方式二：使用 Binder）
binder := server.NewBinder("send_group_msg", handleSendGroupMsg)
dispatcher.Register(binder.Action(), binder.Handler())

// 4. 创建 HTTP 服务器
cfg := server.HTTPConfig{
    Addr:        ":5700",
    PathPrefix:  "/",
    AccessToken: "your-secret",
    ReadTimeout: 30 * time.Second,
}
httpServer := server.NewHTTPServer(cfg, dispatcher)

// 5. 启动服务器
ctx := context.Background()
if err := httpServer.Start(ctx); err != nil {
    log.Fatal(err)
}
```

### 处理器函数示例

```go
func handleSendGroupMsg(ctx context.Context, req *entity.SendGroupMsgRequest) (*entity.ActionResponse[entity.SendGroupMsgResponse], error) {
    // 业务逻辑
    msgID := sendMessageToGroup(req.GroupId, req.Message)

    return &entity.ActionResponse[entity.SendGroupMsgResponse]{
        Status:  entity.StatusOK,
        Retcode: 0,
        Data: &entity.SendGroupMsgResponse{
            MessageId: msgID,
        },
    }, nil
}
```

---

## 对外接口

### 1. HTTPServer (http_server.go)

**NewHTTPServer** (http_server.go:42-58)

```go
func NewHTTPServer(cfg HTTPConfig, handler ActionRequestHandler) *HTTPServer
```

**参数**:

- `cfg`: HTTP 服务器配置
- `handler`: 动作请求处理器（通常是 `Dispatcher`）

**方法**:

- `Start(ctx context.Context) error`: 启动服务器（阻塞）
- `Shutdown(ctx context.Context) error`: 优雅关闭
- `Handler() http.Handler`: 返回 HTTP 处理器（可挂载到外部路由）

**HTTPConfig 配置**:

```go
type HTTPConfig struct {
    Addr              string        // 监听地址，如 ":5700"
    PathPrefix        string        // 路由前缀，如 "/"
    ReadHeaderTimeout time.Duration // 读取头部超时
    ReadTimeout       time.Duration // 读取超时
    WriteTimeout      time.Duration // 写入超时
    IdleTimeout       time.Duration // 空闲超时
    AccessToken       string        // 访问令牌（可选）
}
```

### 2. Dispatcher (dispatcher.go)

**NewDispatcher** (dispatcher.go:17-19)

```go
func NewDispatcher() *Dispatcher
```

**方法**:

- `Register(action string, h ActionHandler)`: 注册 action 处理器
- `HandleActionRequest(ctx, req) (*entity.ActionRawResponse, error)`: 处理动作请求（实现 `ActionRequestHandler` 接口）

### 3. Binder (binder.go)

**NewBinder** (binder.go:18-23)

```go
func NewBinder[Req any, Resp any](
    action string,
    fn func(context.Context, *Req) (*entity.ActionResponse[Resp], error),
) *Binder[Req, Resp]
```

**类型参数**:

- `Req`: 请求类型（如 `entity.SendPrivateMsgRequest`）
- `Resp`: 响应类型（如 `entity.SendPrivateMsgResponse`）

**方法**:

- `Action() string`: 返回绑定的 action 名称
- `Handler() ActionHandler`: 返回类型安全的处理器

**工作原理**:

1. 接收 `map[string]any` 参数
2. 使用 `util.JsonTagMapping` 绑定到 `Req` 类型
3. 调用业务函数 `fn`
4. 将 `ActionResponse[Resp]` 转换为 `ActionRawResponse`

### 4. 处理器接口 (handler.go)

```go
// 处理动作请求
type ActionRequestHandler interface {
    HandleActionRequest(ctx context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error)
}

// 适配函数类型
type ActionRequestHandlerFunc func(ctx context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error)

// 处理具体 action
type ActionHandler func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error)
```

### 5. 错误定义 (errors.go)

```go
var (
    ErrActionNotFound = errors.New("action not found")
    ErrBadRequest     = errors.New("bad request")
)
```

### 6. 自动生成的注册辅助 (http_server_actions_register.go)

通过 `//go:generate go run ../cmd/bindings-gen` 生成，提供便捷的批量注册方法（具体实现取决于生成器配置）。

---

## 关键依赖与配置

### 内部依赖

- `github.com/q1bksuu/onebot-go-sdk/v11/entity`: 协议实体定义
- `github.com/q1bksuu/onebot-go-sdk/v11/internal/util`: JSON 映射工具

### 外部依赖

无（仅使用标准库）

### 代码生成配置

**触发方式**:

在 `http_server.go` 文件顶部：

```go
//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-server-actions-register-output=./http_server_actions_register.go
```

---

## 数据模型

### 核心类型

```go
type HTTPServer struct {
    srv     *http.Server           // 标准 HTTP 服务器
    mux     *http.ServeMux         // 路由多路复用器
    cfg     HTTPConfig             // 配置
    handler ActionRequestHandler   // 请求处理器
}

type Dispatcher struct {
    handlers map[string]ActionHandler // action -> 处理器映射
}

type Binder[Req any, Resp any] struct {
    action string                                                                        // action 名称
    fn     func(ctx context.Context, req *Req) (*entity.ActionResponse[Resp], error)   // 业务函数
}
```

### 请求处理流程

```
1. HTTP 请求到达 (如 POST /send_private_msg)
   ↓
2. handleRoot (http_server.go:97-132)
   ├─ extractAction: 从 URL 路径提取 action
   ├─ checkAccess: 验证访问令牌（如果配置了）
   └─ parseParams: 解析参数（Query + Form + JSON Body）
   ↓
3. 构造 ActionRequest
   {
     "action": "send_private_msg",
     "params": {"user_id": 123456, "message": "Hello"}
   }
   ↓
4. 调用 handler.HandleActionRequest (通常是 Dispatcher)
   ↓
5. Dispatcher.HandleActionRequest (dispatcher.go:27-37)
   ├─ 查找 handlers[action]
   └─ 调用对应的 ActionHandler
   ↓
6. ActionHandler (可能来自 Binder)
   ├─ Binder.Handler: 绑定参数到类型化结构体
   ├─ 调用业务函数 fn(ctx, req)
   └─ 转换响应为 ActionRawResponse
   ↓
7. 返回 JSON 响应
   {
     "status": "ok",
     "retcode": 0,
     "data": {"message_id": 12345}
   }
```

### 鉴权机制

支持两种传递方式（优先级：Header > Query）：

1. **HTTP Header**: `Authorization: Bearer {token}`
2. **URL Query**: `?access_token={token}`

验证逻辑 (http_server.go:152-173):

- 如果 `cfg.AccessToken` 为空，跳过验证
- 否则检查请求中的 token 是否匹配
- 不匹配返回 401 (Unauthorized) 或 403 (Forbidden)

### 参数解析

(http_server.go:175-219)

1. **解析 Form 参数**: `r.ParseForm()`，支持 Query 和 `application/x-www-form-urlencoded`
2. **解析 JSON Body** (仅 POST):
   - 检查 `Content-Type: application/json`
   - 使用 `json.Decoder` 解码到 `map[string]any`
   - 合并到 Form 参数
3. **返回合并后的参数 map**

---

## 测试与质量

### 测试文件

- `http_server_test.go`: 10 个单元测试
- `dispatcher_test.go`: 3 个单元测试
- `binder_test.go`: 2 个单元测试

### 测试场景

#### http_server_test.go

| 测试函数                                        | 测试内容                              |
| ----------------------------------------------- | ------------------------------------- |
| `TestNewHTTPServer_PathPrefixNormalizeAndHandler` | 路径前缀标准化和 Handler 获取       |
| `TestHTTPServer_HandleRoot_PathAndNotFound`     | 路由匹配和 404 处理                   |
| `TestHTTPServer_AuthRequired_MissingOrWrongToken` | 鉴权失败场景（401/403）            |
| `TestHTTPServer_AuthRequired_WithHeaderAndQuery` | 鉴权成功场景（Header 和 Query）     |
| `TestHTTPServer_Params_QueryAndFormAndJSON`     | 多种参数格式解析                      |
| `TestHTTPServer_Params_InvalidForm`             | 无效 Form 数据处理                    |
| `TestHTTPServer_JSON_InvalidOrUnsupportedContentType` | 无效 JSON 或不支持的 Content-Type |
| `TestHTTPServer_WriteError_Mapping`             | 错误映射到 HTTP 状态码                |
| `TestHTTPServer_NilResponse_DefaultFailed`      | 空响应默认返回 failed                 |
| `TestHTTPServer_StartAndShutdown_ContextCancel` | 服务器启动和优雅关闭                  |

#### dispatcher_test.go

| 测试函数                                  | 测试内容                              |
| ----------------------------------------- | ------------------------------------- |
| `TestDispatcher_RegisterAndHandle_Success` | 注册和调用处理器                      |
| `TestDispatcher_HandleActionRequest_NotFound` | 未注册的 action 返回 ErrActionNotFound |
| `TestDispatcher_Register_OverrideExisting` | 覆盖已存在的处理器                    |

#### binder_test.go

| 测试函数                              | 测试内容                              |
| ------------------------------------- | ------------------------------------- |
| `TestBinder_ActionAndHandler_Success` | Binder 绑定和调用成功                 |
| `TestBinder_Handler_PropagatesFuncError` | 业务函数错误传播                   |

### 质量保证

- **Mock 测试**: 使用 `httptest.NewRecorder` 模拟 HTTP 请求/响应
- **边界条件**: 覆盖空参数、无效 JSON、未授权等场景
- **集成测试**: `TestHTTPServer_StartAndShutdown_ContextCancel` 测试完整生命周期

---

## 常见问题 (FAQ)

**Q: 如何添加新的 action 处理器？**

**方式一：直接注册**

```go
dispatcher.Register("custom_action", func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error) {
    // 处理逻辑
    return &entity.ActionRawResponse{...}, nil
})
```

**方式二：使用 Binder（推荐）**

```go
func handleCustomAction(ctx context.Context, req *CustomRequest) (*entity.ActionResponse[CustomResponse], error) {
    // 类型安全的处理逻辑
    return &entity.ActionResponse[CustomResponse]{...}, nil
}

binder := server.NewBinder("custom_action", handleCustomAction)
dispatcher.Register(binder.Action(), binder.Handler())
```

**Q: 如何处理事件上报？**

OneBot 11 的事件上报通常通过 **反向 WebSocket** 或 **反向 HTTP** 实现，不在此模块范围内。

如果使用反向 HTTP，可以在 Dispatcher 中注册特殊 action 来接收事件。

**Q: 如何集成到现有 HTTP 服务器？**

```go
// 获取 Handler
handler := httpServer.Handler()

// 挂载到现有路由
existingMux := http.NewServeMux()
existingMux.Handle("/onebot/", http.StripPrefix("/onebot", handler))

// 启动现有服务器
http.ListenAndServe(":8080", existingMux)
```

**Q: 如何自定义错误响应？**

修改 `writeError` 方法 (http_server.go:231-240)，或在业务函数中返回自定义的 `ActionResponse`：

```go
return &entity.ActionResponse[MyResponse]{
    Status:  entity.StatusFailed,
    Retcode: 1001,
    Message: "自定义错误信息",
}, nil
```

**Q: PathPrefix 如何工作？**

假设 `PathPrefix = "/api/v1"`，则：

- URL: `http://localhost:5700/api/v1/send_private_msg`
- 提取的 action: `send_private_msg`

PathPrefix 会自动标准化为 `/api/v1/`（确保前后有斜杠）。

**Q: 如何处理并发请求？**

`http.Server` 默认为每个请求创建一个 goroutine，因此处理器需要是**并发安全**的。

- `Dispatcher` 的 `handlers` map 在注册完成后只读，并发安全
- 业务函数内的状态需要自行加锁或使用 channel

---

## 相关文件清单

### 主要源文件

| 文件                                | 行数  | 职责                              |
| ----------------------------------- | ----- | --------------------------------- |
| `http_server.go`                    | ~241  | HTTP 服务端核心实现               |
| `dispatcher.go`                     | ~38   | Action 分发器                     |
| `binder.go`                         | ~46   | 类型安全的参数绑定器              |
| `handler.go`                        | ~26   | 处理器接口定义                    |
| `errors.go`                         | ~10   | 错误定义                          |
| `http_server_actions_register.go`   | (生成) | 自动生成的注册辅助代码            |

### 测试文件

- `http_server_test.go`: 10 个单元测试
- `dispatcher_test.go`: 3 个单元测试
- `binder_test.go`: 2 个单元测试

---

*模块文档生成时间: 2025-12-21 15:53:08*
