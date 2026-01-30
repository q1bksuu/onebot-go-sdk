[根目录](../../CLAUDE.md) > [v11](../) > **client**

---

# client - 客户端模块

OneBot 11 协议的客户端实现，支持 HTTP 和 WebSocket 两种通信方式。

---

## 变更记录 (Changelog)

### 2026-01-05

- **新增 WebSocket 客户端**: 添加 `websocket.go`，支持 WebSocket 正向连接
- **更新 API 方法数量**: 38 个自动生成的 API 方法
- **配置选项调整**: `WithHeader`/`WithQuery` 改为接受单个键值对参数

### 2025-12-21 15:53:08

- **初始化**: 生成模块级文档
- **覆盖**: 扫描了 HTTP 客户端核心实现和自动生成的 API 方法

---

## 模块职责

client 模块负责：

1. **HTTP 客户端**: 统一的请求构建、鉴权、错误处理
2. **WebSocket 客户端**: 正向 WebSocket 连接，支持自动重连和 Action 请求处理
3. **API 方法生成**: 通过 `bindings-gen` 工具自动生成所有 OneBot API 的封装方法
4. **灵活配置**: 支持自定义 HTTP Client、超时、访问令牌、路径前缀

---

## 入口与启动

### HTTP 客户端

```go
import "github.com/q1bksuu/onebot-go-sdk/v11/client"

// 基础创建
c, err := client.NewHTTPClient("http://localhost:5700")

// 带访问令牌
c, err := client.NewHTTPClient("http://localhost:5700",
    client.WithAccessToken("your-secret-token"))

// 自定义超时
c, err := client.NewHTTPClient("http://localhost:5700",
    client.WithTimeout(60 * time.Second))

// 自定义 HTTP 客户端
httpClient := &http.Client{Transport: customTransport}
c, err := client.NewHTTPClient("http://localhost:5700",
    client.WithHTTPClient(httpClient))

// 设置路径前缀
c, err := client.NewHTTPClient("http://localhost:5700",
    client.WithPathPrefix("/bot"))
```

### WebSocket 客户端

```go
import "github.com/q1bksuu/onebot-go-sdk/v11/client"

// 创建 WebSocket 客户端
wsClient := client.NewWebSocketClient(
    client.WithWSURL("ws://localhost:6700"),
    client.WithWSReconnectInterval(5*time.Second), // 断线重连间隔
    client.WithWSSelfID(123456789),                // 机器人 QQ 号
    client.WithWSAccessToken("your-token"),        // 可选鉴权令牌
    client.WithWSReadTimeout(30*time.Second),      // 可选读取超时
    client.WithWSWriteTimeout(30*time.Second),     // 可选写入超时
)

// 启动客户端（阻塞直到 context 取消）
ctx, cancel := context.WithCancel(context.Background())
go func() {
    if err := wsClient.Start(ctx); err != nil {
        log.Printf("WebSocket client error: %v", err)
    }
}()

// 关闭客户端
cancel()
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer shutdownCancel()
wsClient.Shutdown(shutdownCtx)
```

### 调用 API

```go
ctx := context.Background()

// 发送私聊消息
resp, err := c.SendPrivateMsg(ctx, &entity.SendPrivateMsgRequest{
    UserId: 123456,
    Message: &entity.MessageValue{
        Type: entity.MessageValueTypeString,
        StringValue: "Hello!",
    },
})

// 获取群列表
groups, err := c.GetGroupList(ctx, &entity.GetGroupListRequest{})
```

---

## 对外接口

### 1. HTTPClient 类型

**NewHTTPClient** (http_client.go:43-63)

```go
func NewHTTPClient(baseURL string, opts ...Option) (*HTTPClient, error)
```

**参数**:

- `baseURL`: OneBot 实现的 HTTP 地址（必填，如 `http://localhost:5700`）
- `opts`: 可选配置项

**错误**:

- `errBaseURLEmpty`: baseURL 为空
- `errMissingSchemeOrHost`: baseURL 缺少 scheme 或 host

### 2. WebSocketClient 类型

**NewWebSocketClient** (websocket.go)

```go
func NewWebSocketClient(opts ...WSClientOption) *WebSocketClient
```

**WSClientConfig 字段**（可通过 `WithWSConfig` 一次性设置）:

| 字段                | 类型            | 说明                     |
| ------------------- | --------------- | ------------------------ |
| `URL`               | `string`        | WebSocket 连接地址（必填）|
| `ReconnectInterval` | `time.Duration` | 断线重连间隔             |
| `SelfID`            | `int64`         | 机器人 QQ 号（X-Self-ID）|
| `AccessToken`       | `string`        | 可选鉴权令牌             |
| `ReadTimeout`       | `time.Duration` | 读取超时（默认 0 无限）  |
| `WriteTimeout`      | `time.Duration` | 写入超时（默认 0 无限）  |

**WSClientOption 选项**（顺序生效，后者覆盖前者）:

| Option                              | 说明                             | 示例                                   |
| ----------------------------------- | -------------------------------- | -------------------------------------- |
| `WithWSConfig(cfg WSClientConfig)`  | 批量设置客户端配置               | `WithWSConfig(client.WSClientConfig{...})` |
| `WithWSURL(url string)`             | 设置 WebSocket 连接地址          | `WithWSURL("ws://localhost:6700")`     |
| `WithWSReconnectInterval(d)`        | 设置断线重连间隔                 | `WithWSReconnectInterval(3*time.Second)` |
| `WithWSSelfID(id int64)`            | 设置 SelfID                      | `WithWSSelfID(123456789)`              |
| `WithWSAccessToken(token string)`   | 设置访问令牌                     | `WithWSAccessToken("secret")`          |
| `WithWSReadTimeout(d)`              | 设置读取超时                     | `WithWSReadTimeout(30*time.Second)`    |
| `WithWSWriteTimeout(d)`             | 设置写入超时                     | `WithWSWriteTimeout(30*time.Second)`   |
| `WithWSActionHandler(handler)`      | 设置 Action 请求处理器            | `WithWSActionHandler(handler)`         |

**主要方法**:

| 方法                              | 说明                              |
| --------------------------------- | --------------------------------- |
| `Start(ctx context.Context)`      | 启动客户端，阻塞直到 ctx 取消     |
| `Shutdown(ctx context.Context)`   | 优雅关闭，等待所有 goroutine 完成 |
| `BroadcastEvent(event entity.Event)` | 广播事件到所有连接             |

### 3. HTTP 客户端配置选项 (http_client_options.go)

| Option                              | 说明                             | 示例                              |
| ----------------------------------- | -------------------------------- | --------------------------------- |
| `WithAccessToken(token string)`     | 设置访问令牌（Bearer 认证）      | `WithAccessToken("secret")`       |
| `WithTimeout(timeout time.Duration)`| 设置请求超时（默认 30s）         | `WithTimeout(60 * time.Second)`   |
| `WithHTTPClient(client *http.Client)`| 使用自定义 HTTP 客户端          | `WithHTTPClient(customClient)`    |
| `WithPathPrefix(prefix string)`     | 设置路径前缀                     | `WithPathPrefix("/bot")`          |

### 4. 调用选项 (CallOption)

在调用 API 方法时可传入：

| CallOption                      | 说明                             | 示例                              |
| ------------------------------- | -------------------------------- | --------------------------------- |
| `WithMethod(method string)`     | 覆盖默认 HTTP 方法（GET/POST）   | `WithMethod(http.MethodGet)`      |
| `WithQuery(key, value string)`  | 添加单个 URL 查询参数            | `WithQuery("key", "val")`         |
| `WithHeader(key, value string)` | 添加单个自定义 HTTP 头部         | `WithHeader("X-Custom", "value")` |

### 5. 自动生成的 API 方法 (http_client_actions.go)

通过 `//go:generate go run ../cmd/bindings-gen` 自动生成，包含 **38 个方法**：

#### 消息 API

- `SendPrivateMsg(ctx, req) (*entity.SendPrivateMsgResponse, error)`
- `SendGroupMsg(ctx, req) (*entity.SendGroupMsgResponse, error)`
- `SendMsg(ctx, req) (*entity.SendMsgResponse, error)`
- `DeleteMsg(ctx, req) (*entity.DeleteMsgResponse, error)`
- `GetMsg(ctx, req) (*entity.GetMsgResponse, error)`
- `GetForwardMsg(ctx, req) (*entity.GetForwardMsgResponse, error)`

#### 好友 API

- `SendLike(ctx, req) (*entity.SendLikeResponse, error)`
- `SetFriendAddRequest(ctx, req) (*entity.SetFriendAddRequestResponse, error)`
- `GetStrangerInfo(ctx, req) (*entity.GetStrangerInfoResponse, error)`
- `GetFriendList(ctx, req) (*entity.GetFriendListResponse, error)`

#### 群管理 API

- `SetGroupKick`, `SetGroupBan`, `SetGroupAnonymousBan`, `SetGroupWholeBan`
- `SetGroupAdmin`, `SetGroupAnonymous`, `SetGroupCard`, `SetGroupName`
- `SetGroupLeave`, `SetGroupSpecialTitle`, `SetGroupAddRequest`

#### 群信息 API

- `GetGroupInfo`, `GetGroupList`, `GetGroupMemberInfo`, `GetGroupMemberList`, `GetGroupHonorInfo`

#### 账号凭证 API

- `GetLoginInfo`, `GetCookies`, `GetCsrfToken`, `GetCredentials`

#### 媒体 API

- `GetRecord`, `GetImage`

#### 能力检查 API

- `CanSendImage`, `CanSendRecord`

#### 系统 API

- `GetStatus`, `GetVersionInfo`, `SetRestart`, `CleanCache`

---

## 关键依赖与配置

### 内部依赖

- `github.com/q1bksuu/onebot-go-sdk/v11/entity`: 协议实体定义
- `github.com/q1bksuu/onebot-go-sdk/v11/dispatcher`: 分发器接口定义（WebSocket 客户端使用 ActionRequestHandler）
- `github.com/q1bksuu/onebot-go-sdk/v11/server`: 服务端组件
- `github.com/q1bksuu/onebot-go-sdk/v11/internal/util`: JSON 映射工具

### 外部依赖

- `github.com/gorilla/websocket`: WebSocket 协议支持

### 代码生成配置

**触发方式**:

在 `http_client.go` 文件顶部：

```go
//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-client-actions-output=./http_client_actions.go
```

运行 `go generate` 后，会根据 `config.yaml` 生成所有 API 方法到 `http_client_actions.go`。

---

## 数据模型

### HTTP 客户端核心类型

```go
type HTTPClient struct {
    baseURL     string        // OneBot 实现的地址
    accessToken string        // 访问令牌（可选）
    httpClient  *http.Client  // HTTP 客户端实例
}

type clientOptions struct {
    httpClient  *http.Client
    accessToken string
    pathPrefix  string
    timeout     time.Duration
}

type callOptions struct {
    headers        http.Header     // 自定义请求头
    query          url.Values      // 查询参数
    methodOverride string          // HTTP 方法覆盖
}
```

### WebSocket 客户端核心类型

```go
type WSClientConfig struct {
    URL               string
    ReconnectInterval time.Duration // 断线重连间隔
    SelfID            int64         // 机器人 QQ 号（用于 X-Self-ID 请求头）
    AccessToken       string        // 可选鉴权令牌
    ReadTimeout       time.Duration // 读取超时（可选），默认 0
    WriteTimeout      time.Duration // 写入超时（可选），默认 0
}

type WebSocketClient struct {
    cfg           WSClientConfig
    actionHandler dispatcher.ActionRequestHandler

    // 连接管理
    mu            sync.Mutex
    conn          *websocket.Conn

    // 控制
    cancel        context.CancelFunc
    wg            sync.WaitGroup
}
```

### HTTP 请求流程

```
1. 调用 API 方法（如 SendPrivateMsg）
   ↓
2. 进入 do() 方法
   ↓
3. 解析方法、URL、参数
   ├─ resolveMethod: 确定 HTTP 方法（GET/POST）
   ├─ buildTargetURL: 构建完整 URL
   └─ encodeToParams: 将 Request 结构体转为 map[string]any
   ↓
4. 准备请求体
   ├─ GET: 参数合并到 URL Query
   └─ POST: 参数序列化为 JSON Body
   ↓
5. 设置请求头
   ├─ Content-Type: application/json (POST)
   └─ Authorization: Bearer {token} (如果有)
   ↓
6. 执行 HTTP 请求
   ↓
7. 解析响应
   ├─ 检查 HTTP 状态码 (2xx)
   ├─ 解码 JSON 为 ActionRawResponse
   └─ 检查 retcode (非 0 返回 ActionError)
   ↓
8. 返回结果
```

### 错误处理

| 错误类型                  | 触发条件                        | 示例                             |
| ------------------------- | ------------------------------- | -------------------------------- |
| `errBaseURLEmpty`         | baseURL 为空                    | `NewHTTPClient("")`              |
| `errUnsupportedHTTPMethod`| 使用了非 GET/POST 方法          | `WithMethod("PUT")`              |
| `errHTTPStatus`           | HTTP 状态码非 2xx               | OneBot 实现返回 500              |
| `errMissingSchemeOrHost`  | URL 缺少协议或主机              | `NewHTTPClient("localhost:5700")`|
| `entity.ActionError`      | retcode 非 0 或 status 为 failed| OneBot 返回错误码                |

---

## 测试与质量

### 测试文件

- `http_client_test.go`: 7 个单元测试
- `websocket_test.go`: 1 个单元测试

### 测试场景

| 测试函数                              | 测试内容                              |
| ------------------------------------- | ------------------------------------- |
| `TestNewHTTPClient_EmptyBaseURL_Error`| 空 baseURL 返回错误                   |
| `TestHTTPClient_do_GetQueryMerge`     | GET 请求参数正确合并到 URL Query      |
| `TestHTTPClient_do_PostHeadersBody`   | POST 请求正确设置 Content-Type 和 JSON Body |
| `TestHTTPClient_do_StatusNot2xx`      | HTTP 非 2xx 状态码返回错误            |
| `TestHTTPClient_do_NonZeroRetcode`    | retcode 非 0 返回 ActionError         |
| `TestHTTPClient_do_DecodeError`       | 无效 JSON 响应返回解码错误            |
| `TestHTTPClient_SendPrivateMsg_Success` | 完整 API 调用流程测试             |
| `TestWebSocketClient_HandleActionMessage` | WebSocket Action 消息处理流程测试 |

### 质量保证

- **Mock HTTP 服务器**: 使用 `httptest.NewServer` 模拟 OneBot 实现
- **边界条件**: 覆盖空参数、错误状态码、无效 JSON 等场景
- **集成测试**: `TestHTTPClient_SendPrivateMsg_Success` 测试完整调用链
- **Mock ActionHandler**: WebSocket 测试使用模拟的 ActionRequestHandler

---

## 常见问题 (FAQ)

**Q: 如何添加新的 API 方法？**

1. 在 `../entity/api.go` 定义 `XxxRequest` 和 `XxxResponse`
2. 编辑 `../cmd/bindings-gen/config.yaml`，在对应 group 添加 action 配置
3. 运行 `go generate ./...`
4. 新方法会自动出现在 `http_client_actions.go`

**Q: 如何自定义 HTTP 传输层（如添加代理）？**

```go
transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}
httpClient := &http.Client{Transport: transport}
c, _ := client.NewHTTPClient(baseURL, client.WithHTTPClient(httpClient))
```

**Q: HTTP 客户端和 WebSocket 客户端如何选择？**

- **HTTP 客户端**: 用于主动调用 OneBot API（发消息、查询信息等）
- **WebSocket 客户端**: 用于接收 OneBot 实现推送的 Action 请求，支持断线自动重连

**Q: 为什么有些 API 使用 GET，有些使用 POST？**

- 默认使用 POST（OneBot 11 标准推荐）
- 可以通过 `WithMethod(http.MethodGet)` 覆盖
- 生成的代码默认都是 POST

**Q: AccessToken 如何传递？**

两种方式（优先级：Header > Query）：

1. **HTTP Header**: `Authorization: Bearer {token}`
2. **URL Query**: `?access_token={token}`

客户端默认使用 Header 方式。

**Q: 如何处理超时？**

```go
// 全局超时
c, _ := client.NewHTTPClient(baseURL, client.WithTimeout(60*time.Second))

// 请求级超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
resp, err := c.SendPrivateMsg(ctx, req)
```

**Q: 参数如何序列化？**

- **GET 请求**: 使用 `encodeToParams` 转为 `map[string]any`，然后合并到 URL Query
- **POST 请求**: 使用 `json.Marshal` 序列化为 JSON Body

内部使用 `util.JsonTagMapping` 确保字段名与 `json` 标签一致。

---

## 相关文件清单

### 主要源文件

| 文件                         | 行数  | 职责                              |
| ---------------------------- | ----- | --------------------------------- |
| `http_client.go`             | ~275  | HTTP 客户端核心逻辑               |
| `http_client_options.go`     | ~67   | 客户端和调用选项定义              |
| `http_client_actions.go`     | ~885  | 38 个 API 方法（自动生成）        |
| `websocket.go`               | ~414  | WebSocket 客户端实现              |

### 测试文件

- `http_client_test.go` (~167 行): 7 个单元测试，覆盖核心流程和边界条件
- `websocket_test.go` (~72 行): 1 个单元测试，覆盖 Action 消息处理逻辑

---

*模块文档更新时间: 2026-01-05*
