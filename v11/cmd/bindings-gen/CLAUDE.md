[根目录](../../../CLAUDE.md) > [v11](../../) > [cmd](../) > **bindings-gen**

---

# bindings-gen - API 绑定代码生成工具

基于 YAML 配置的代码生成工具，自动生成 HTTP 客户端 API 方法和服务端注册代码。

---

## 变更记录 (Changelog)

### 2025-12-21 15:53:08

- **初始化**: 生成工具文档
- **覆盖**: 扫描了配置文件和核心模型定义

---

## 模块职责

bindings-gen 工具负责：

1. **客户端方法生成**: 为 `HTTPClient` 生成所有 OneBot API 的封装方法
2. **服务端注册代码生成**: 生成服务端 action 注册辅助代码
3. **配置驱动**: 通过 YAML 配置文件声明式定义所有 API
4. **分组管理**: 支持按功能分组（消息、好友、群管理等）
5. **类型安全**: 生成的代码使用泛型，确保类型安全

---

## 入口与启动

### 命令行调用

```bash
# 生成客户端方法
go run . -config=config.yaml -http-client-actions-output=../../client/http_client_actions.go

# 生成服务端注册代码
go run . -config=config.yaml -http-server-actions-register-output=../../server/http_server_actions_register.go
```

### go generate 集成

在 `client/http_client.go` 和 `server/http_server.go` 顶部：

```go
//go:generate go run ../cmd/bindings-gen -config=../cmd/bindings-gen/config.yaml -http-client-actions-output=./http_client_actions.go
```

运行：

```bash
go generate ./...
```

---

## 对外接口

### 命令行参数

| 参数                                  | 说明                              | 示例                                                       |
| ------------------------------------- | --------------------------------- | ---------------------------------------------------------- |
| `-config`                             | YAML 配置文件路径                 | `-config=config.yaml`                                      |
| `-http-client-actions-output`         | 客户端方法输出文件路径            | `-http-client-actions-output=../../client/http_client_actions.go` |
| `-http-server-actions-register-output`| 服务端注册代码输出文件路径        | `-http-server-actions-register-output=../../server/http_server_actions_register.go` |

### 配置文件结构 (config.yaml)

```yaml
combined_service:
  name: OneBotService              # 总服务名称
  desc: OneBot 11 协议服务          # 总服务描述

groups:
  - name: message                  # 分组名称
    service_name: MessageService   # 分组服务名称
    service_desc: 消息服务          # 分组服务描述
    actions:
      - method: SendPrivateMsg     # 方法名（PascalCase）
        action: send_private_msg   # action 名称（snake_case）
        desc: 发送私聊消息          # 方法描述
        request: entity.SendPrivateMsgRequest   # 请求类型（完整路径）
        response: entity.SendPrivateMsgResponse # 响应类型（完整路径）
        http_method: POST          # HTTP 方法（可选，默认 POST）
        path: /send_private_msg    # URL 路径（可选，默认 /{action}）

  - name: friend
    service_name: FriendService
    # ... 更多 actions
```

**完整配置示例**: [config.yaml](./config.yaml)

---

## 关键依赖与配置

### 外部依赖

- `gopkg.in/yaml.v3`: YAML 解析

### 配置文件

**config.yaml 包含 8 个功能分组**:

1. **message** (6 个 action): 消息发送、撤回、获取
2. **friend** (4 个 action): 好友点赞、请求处理、信息查询
3. **group_admin** (11 个 action): 群管理（踢人、禁言、设置管理员等）
4. **group_info** (5 个 action): 群信息查询
5. **account** (4 个 action): 账号凭证（登录信息、Cookies、Token）
6. **media** (2 个 action): 媒体获取（语音、图片）
7. **capability** (2 个 action): 能力检查（图片、语音）
8. **system** (4 个 action): 系统操作（状态、版本、重启、清理缓存）

**共计**: 38 个 API action

---

## 数据模型

### 核心类型 (models.go)

```go
type Config struct {
    CombinedService CombinedService `yaml:"combined_service"`
    Groups          []Group         `yaml:"groups"`
}

type CombinedService struct {
    Name string `yaml:"name"`  // 如 "OneBotService"
    Desc string `yaml:"desc"`  // 如 "OneBot 11 协议服务"
}

type Group struct {
    Name        string   `yaml:"name"`         // 如 "message"
    ServiceName string   `yaml:"service_name"` // 如 "MessageService"
    ServiceDesc string   `yaml:"service_desc"` // 如 "消息服务"
    Actions     []Action `yaml:"actions"`
}

type Action struct {
    Method     string `yaml:"method"`      // 如 "SendPrivateMsg"
    Action     string `yaml:"action"`      // 如 "send_private_msg"
    Desc       string `yaml:"desc"`        // 如 "发送私聊消息"
    Request    string `yaml:"request"`     // 如 "entity.SendPrivateMsgRequest"
    Response   string `yaml:"response"`    // 如 "entity.SendPrivateMsgResponse"
    HTTPMethod string `yaml:"http_method"` // 如 "POST" (可选)
    Path       string `yaml:"path"`        // 如 "/send_private_msg" (可选)
}
```

### 生成流程

```
1. 读取配置文件 (config.yaml)
   ↓
2. 解析 YAML 到 Config 结构体
   ↓
3. 遍历所有 groups 和 actions
   ↓
4. 生成客户端方法
   ├─ 方法签名：func (c *HTTPClient) {Method}(ctx context.Context, req *{Request}) (*{Response}, error)
   ├─ 方法体：调用 c.do() 执行 HTTP 请求
   └─ 文档注释：{Desc}
   ↓
5. 生成服务端注册代码
   ├─ 注册辅助函数：RegisterXxxActions(dispatcher *Dispatcher, handler XxxHandler)
   └─ 接口定义：XxxHandler interface { Handle{Method}(...) }
   ↓
6. 写入输出文件
```

### 生成的客户端代码示例

```go
// 生成到 client/http_client_actions.go

// SendPrivateMsg 发送私聊消息
func (c *HTTPClient) SendPrivateMsg(ctx context.Context, req *entity.SendPrivateMsgRequest) (*entity.SendPrivateMsgResponse, error) {
    rawResp, err := c.do(ctx, "send_private_msg", http.MethodPost, req)
    if err != nil {
        return nil, err
    }

    var resp entity.SendPrivateMsgResponse
    err = json.Unmarshal(rawResp.Data, &resp)
    if err != nil {
        return nil, fmt.Errorf("unmarshal response: %w", err)
    }

    return &resp, nil
}
```

### 生成的服务端代码示例

```go
// 生成到 server/http_server_actions_register.go

// MessageHandler 消息服务处理器接口
type MessageHandler interface {
    HandleSendPrivateMsg(ctx context.Context, req *entity.SendPrivateMsgRequest) (*entity.ActionResponse[entity.SendPrivateMsgResponse], error)
    HandleSendGroupMsg(ctx context.Context, req *entity.SendGroupMsgRequest) (*entity.ActionResponse[entity.SendGroupMsgResponse], error)
    // ... 其他方法
}

// RegisterMessageActions 注册消息服务的所有 action
func RegisterMessageActions(dispatcher *Dispatcher, handler MessageHandler) {
    binder := NewBinder("send_private_msg", handler.HandleSendPrivateMsg)
    dispatcher.Register(binder.Action(), binder.Handler())

    // ... 其他 action
}
```

---

## 测试与质量

### 测试策略

- **配置验证**: 确保 config.yaml 语法正确（通过 YAML 解析器验证）
- **生成代码质量**: 生成的代码通过 golangci-lint 检查
- **集成测试**: 在 `client` 和 `server` 模块的测试中使用生成的代码

### 质量保证

- **类型安全**: 使用泛型确保请求/响应类型匹配
- **一致性**: 所有 API 方法遵循统一的命名和结构
- **可维护性**: 通过配置文件管理，无需手写大量样板代码

---

## 常见问题 (FAQ)

**Q: 如何添加新的 API？**

1. 在 `config.yaml` 的对应 group 中添加 action 配置：

   ```yaml
   - method: CustomAction
     action: custom_action
     desc: 自定义操作
     request: entity.CustomActionRequest
     response: entity.CustomActionResponse
   ```

2. 在 `v11/entity/api.go` 定义 Request/Response 类型

3. 运行 `go generate ./v11/client` 和 `go generate ./v11/server`

**Q: 如何修改默认 HTTP 方法？**

在 action 配置中添加 `http_method`:

```yaml
- method: GetLoginInfo
  action: get_login_info
  desc: 获取登录信息
  request: entity.GetLoginInfoRequest
  response: entity.GetLoginInfoResponse
  http_method: GET  # 覆盖默认的 POST
```

**Q: 如何自定义 URL 路径？**

在 action 配置中添加 `path`:

```yaml
- method: SendPrivateMsg
  action: send_private_msg
  desc: 发送私聊消息
  request: entity.SendPrivateMsgRequest
  response: entity.SendPrivateMsgResponse
  path: /api/v1/send_private_msg  # 自定义路径
```

**Q: 为什么使用 YAML 而不是 JSON？**

YAML 优势：

- 支持注释，便于维护
- 更简洁的语法（无需大量引号和逗号）
- 人类可读性更好

**Q: 生成的代码可以手动修改吗？**

不建议，因为：

- 下次运行 `go generate` 会覆盖修改
- 应该修改配置文件，而不是生成的代码

如果需要自定义逻辑，可以：

- 在业务代码中包装生成的方法
- 或者直接调用底层的 `HTTPClient.do()` 方法

**Q: 如何添加新的功能分组？**

在 `config.yaml` 的 `groups` 数组中添加：

```yaml
groups:
  # ... 现有分组
  - name: custom
    service_name: CustomService
    service_desc: 自定义服务
    actions:
      - method: CustomMethod
        action: custom_action
        desc: 自定义操作
        request: entity.CustomRequest
        response: entity.CustomResponse
```

---

## 相关文件清单

### 主要源文件

| 文件           | 行数估算 | 职责                              |
| -------------- | -------- | --------------------------------- |
| `main.go`      | ~200     | 命令行参数解析、模板渲染、代码生成|
| `models.go`    | ~50      | 配置文件数据结构定义              |
| `config.yaml`  | ~240     | 所有 API 的声明式配置             |

### 输出文件

| 文件                                  | 生成位置     | 内容                              |
| ------------------------------------- | ------------ | --------------------------------- |
| `http_client_actions.go`              | `v11/client` | 38 个客户端 API 方法              |
| `http_server_actions_register.go`     | `v11/server` | 服务端注册辅助函数和接口定义      |

---

## 扩展阅读

完整的配置文件示例：

- **配置文件**: [config.yaml](./config.yaml)
- **生成的客户端代码**: [v11/client/http_client_actions.go](../../client/http_client_actions.go)
- **生成的服务端代码**: [v11/server/http_server_actions_register.go](../../server/http_server_actions_register.go)

---

*工具文档生成时间: 2025-12-21 15:53:08*
