[根目录](../../../CLAUDE.md) > [v11](../../) > [cmd](../) > **bindings-gen**

---

# bindings-gen - API 绑定代码生成工具

基于 YAML 配置的代码生成工具，自动生成 HTTP 客户端 API 方法和服务端注册代码。

---

## 变更记录 (Changelog)

### 2026-01-05

- **更新**: 同步文档与实际代码实现
- **修正**: 更新生成代码示例（客户端返回 `ActionResponse[T]`、服务端使用 `dispatcher.APIFuncToActionHandler`）
- **修正**: 更新行数统计

### 2025-12-21 15:53:08

- **初始化**: 生成工具文档
- **覆盖**: 扫描了配置文件和核心模型定义

---

## 模块职责

bindings-gen 工具负责：

1. **客户端方法生成**: 为 `HTTPClient` 生成所有 OneBot API 的封装方法
2. **服务端注册代码生成**: 生成服务接口定义、聚合接口、Unimplemented 实现和批量注册函数
3. **配置驱动**: 通过 YAML 配置文件声明式定义所有 API
4. **分组管理**: 支持按功能分组（消息、好友、群管理等）
5. **类型安全**: 生成的代码使用泛型 `ActionResponse[T]`，确保类型安全

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
  name: OneBotService              # 聚合服务接口名称
  desc: OneBot 11 协议服务          # 聚合服务描述

groups:
  - name: message                  # 分组名称
    service_name: MessageService   # 分组服务接口名称
    service_desc: 消息服务          # 分组服务描述
    actions:
      - method: SendPrivateMsg     # 方法名（PascalCase）
        action: send_private_msg   # action 名称（snake_case）
        desc: 发送私聊消息          # 方法描述
        request: entity.SendPrivateMsgRequest   # 请求类型（完整路径）
        response: entity.SendPrivateMsgResponse # 响应类型（完整路径）
        http_method: POST          # HTTP 方法（可选，默认 POST）
        path: /send_private_msg    # URL 路径（可选，默认空）

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
// Config 配置根.
type Config struct {
    Groups          []Group         `yaml:"groups"`
    CombinedService CombinedService `yaml:"combined_service"`
}

type CombinedService struct {
    Name string `yaml:"name"`  // 如 "OneBotService"
    Desc string `yaml:"desc"`  // 如 "OneBot 11 协议服务"
}

// Group 表示一组业务接口，可生成独立 Service.
type Group struct {
    Name        string   `yaml:"name"`         // 如 "message"
    ServiceName string   `yaml:"service_name"` // 如 "MessageService"
    ServiceDesc string   `yaml:"service_desc"` // 如 "消息服务"
    Actions     []Action `yaml:"actions"`
}

// Action 定义单个 action 的生成规则.
type Action struct {
    // Method 生成的 Service 方法名（必填）.
    Method string `yaml:"method"`
    // Action 协议动作名（必填）.
    Action string `yaml:"action"`
    // Desc 方法描述（必填）.
    Desc string `yaml:"desc"`
    // Request / Response 类型（必填），需包含包名，例如 entity.SendPrivateMsgRequest.
    Request  string `yaml:"request"`
    Response string `yaml:"response"`
    // HTTPMethod 可选，默认 POST，可指定 GET/POST.
    HTTPMethod string `yaml:"http_method"`
    // Path 可选，默认 "/{action}".
    Path string `yaml:"path"`
}
```

### 生成流程

```
1. 读取配置文件 (config.yaml)
   ↓
2. 解析 YAML 到 Config 结构体
   ↓
3. 验证配置
   ├─ 检查 groups 不为空
   ├─ 检查每个 group 的 name 不为空
   ├─ 检查每个 action 的 method、request、response 不为空
   └─ 检查 http_method 为空或 GET/POST
   ↓
4. 遍历所有 groups 和 actions
   ↓
5. 渲染模板并格式化
   ↓
6. 写入输出文件
```

### 生成的客户端代码示例

```go
// 生成到 client/http_client_actions.go

// SendPrivateMsg calls action "send_private_msg".
func (c *HTTPClient) SendPrivateMsg(
    ctx context.Context,
    req *entity.SendPrivateMsgRequest,
    opts ...CallOption,
) (*entity.ActionResponse[entity.SendPrivateMsgResponse], error) {
    rawResponse, err := c.do(ctx, "", "", req, opts...)
    if err != nil {
        return nil, err
    }

    out := entity.ActionResponse[entity.SendPrivateMsgResponse]{
        Status:  rawResponse.Status,
        Retcode: rawResponse.Retcode,
        Message: rawResponse.Message,
    }
    err = json.Unmarshal(rawResponse.GetData(), &out.Data)
    if err != nil {
        return nil, err
    }
    return &out, nil
}
```

### 生成的服务端代码示例

```go
// 生成到 server/http_server_actions_register.go

// MessageService 消息服务
type MessageService interface {
    // SendPrivateMsg 发送私聊消息.
    SendPrivateMsg(ctx context.Context, req *entity.SendPrivateMsgRequest) (*entity.ActionResponse[entity.SendPrivateMsgResponse], error)
    // SendGroupMsg 发送群消息.
    SendGroupMsg(ctx context.Context, req *entity.SendGroupMsgRequest) (*entity.ActionResponse[entity.SendGroupMsgResponse], error)
    // ... 其他方法
}

// OneBotService aggregates all service groups.
type OneBotService interface {
    // 消息服务.
    MessageService
    // 好友与陌生人管理.
    FriendService
    // ... 其他分组服务
}

// RegisterGenerated registers actions to dispatcher.
func RegisterGenerated(d *dispatcher.Dispatcher, svc OneBotService) {
    // Group: message
    d.Register("send_private_msg", dispatcher.APIFuncToActionHandler(svc.SendPrivateMsg))
    d.Register("send_group_msg", dispatcher.APIFuncToActionHandler(svc.SendGroupMsg))
    // ... 其他 action
}

// UnimplementedOneBotService aggregates unimplemented group services.
type UnimplementedOneBotService struct {
    UnimplementedMessageService
    UnimplementedFriendService
    // ... 其他 Unimplemented
}

// UnimplementedMessageService provides default empty implementations.
type UnimplementedMessageService struct{}

// SendPrivateMsg 发送私聊消息 (unimplemented).
func (*UnimplementedMessageService) SendPrivateMsg(
    ctx context.Context,
    req *entity.SendPrivateMsgRequest,
) (*entity.ActionResponse[entity.SendPrivateMsgResponse], error) {
    panic("unimplemented")
}
// ... 其他 unimplemented 方法
```

---

## 测试与质量

### 测试策略

- **配置验证**: 确保 config.yaml 语法正确（通过 YAML 解析器验证）
- **生成代码质量**: 生成的代码通过 golangci-lint 检查
- **集成测试**: 在 `client` 和 `server` 模块的测试中使用生成的代码

### 质量保证

- **类型安全**: 使用泛型 `ActionResponse[T]` 确保请求/响应类型匹配
- **一致性**: 所有 API 方法遵循统一的命名和结构
- **可维护性**: 通过配置文件管理，无需手写大量样板代码
- **代码格式化**: 使用 `go/format` 自动格式化生成的代码

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

**Q: Unimplemented 类型有什么用？**

`UnimplementedXxxService` 类型提供了默认的空实现，方便：

- 快速搭建服务骨架
- 只实现需要的方法，其他方法返回 panic（开发阶段提醒未实现）
- 嵌入到自定义服务中，逐步实现各个方法

---

## 相关文件清单

### 主要源文件

| 文件           | 行数  | 职责                              |
| -------------- | ----- | --------------------------------- |
| `main.go`      | ~315  | 命令行参数解析、配置验证、模板渲染、代码生成 |
| `models.go`    | ~38   | 配置文件数据结构定义              |
| `config.yaml`  | ~239  | 所有 API 的声明式配置             |

### 输出文件

| 文件                                  | 生成位置     | 内容                              |
| ------------------------------------- | ------------ | --------------------------------- |
| `http_client_actions.go`              | `v11/client` | 38 个客户端 API 方法              |
| `http_server_actions_register.go`     | `v11/server` | 服务接口、聚合接口、Unimplemented 实现、批量注册函数 |

---

## 扩展阅读

完整的配置文件示例：

- **配置文件**: [config.yaml](./config.yaml)
- **生成的客户端代码**: [v11/client/http_client_actions.go](../../client/http_client_actions.go)
- **生成的服务端代码**: [v11/server/http_server_actions_register.go](../../server/http_server_actions_register.go)

---

*工具文档更新时间: 2026-01-05*
