[根目录](../../CLAUDE.md) > [v11](../) > **entity**

---

# entity - OneBot 协议实体模块

OneBot 11 协议的完整类型定义，包括消息、事件、API 请求/响应、通信层实体。

---

## 变更记录 (Changelog)

### 2026-01-05

- **新增**: `Event` 接口，统一所有事件类型的公共方法
- **新增**: `ActionRequestEnvelope` 和 `ActionResponseEnvelope`，支持 Echo 字段的通信封装
- **重构**: 事件常量类型简化，从每事件独立类型改为统一枚举（`EventPostType`、`EventNoticeType` 等）
- **更新**: 行数和行号引用

### 2025-12-21 15:53:08

- **初始化**: 生成模块级文档
- **覆盖**: 扫描了所有实体定义和常量文件

---

## 模块职责

entity 模块是整个 SDK 的**类型基石**，负责：

1. **协议实体定义**: 完整实现 OneBot 11 标准的所有数据结构
2. **类型安全保证**: 通过强类型避免运行时字段错误
3. **动态字段支持**: 提供 `origin` 字段机制保留未知扩展字段
4. **代码自动生成**: 通过 `entity-gen` 工具生成 Getter/Setter 方法
5. **常量定义**: 集中管理所有枚举类型和常量值

---

## 入口与启动

**主要入口文件**:

- `base.go`: 基础类型（StatusMeta、GroupAnonymousUser）
- `message.go`: 消息值类型（MessageValue、Segment）
- `event.go`: 所有事件类型（私聊/群聊/通知/请求/元事件）
- `api.go`: 所有 API 请求/响应类型
- `communication.go`: 通信层类型（ActionRequest、ActionResponse、ActionError）

**代码生成标记**:

每个主要文件的第一行都有 `//go:generate` 指令：

```go
//go:generate go run ../cmd/entity-gen
```

运行 `go generate ./...` 会自动生成 `*_setter_getter.go` 文件。

---

## 对外接口

### 1. 消息实体

**MessageValue** (message.go:10-72)

支持 OneBot 的混合消息格式：

```go
type MessageValue struct {
    Type        MessageValueType // "string" | "array"
    StringValue string           // CQ 码格式
    ArrayValue  []*Segment       // 消息段数组
}
```

特性：

- 实现了 `json.Unmarshaler` 和 `json.Marshaler`
- 自动识别 JSON 中的字符串或数组格式
- 支持双向转换

### 2. 事件类型

**Event 接口** (event.go:3-7)

所有事件类型的公共接口：

```go
type Event interface {
    GetTime() int64
    GetSelfId() int64
    GetPostType() EventPostType
}
```

所有具体事件类型（如 `PrivateMessageEvent`、`GroupMessageEvent` 等）都实现了此接口。

**消息事件** (event.go:9-102)

- `PrivateMessageEvent`: 私聊消息（好友/群临时/其他）
- `GroupMessageEvent`: 群消息（普通/匿名/系统提示）

**通知事件** (event.go:104-350)

- 群文件上传、管理员变动、成员增减、禁言
- 好友添加、消息撤回、戳一戳、红包运气王、荣誉变更

**请求事件** (event.go:352-396)

- `FriendRequestEvent`: 加好友请求
- `GroupRequestEvent`: 加群请求/邀请

**元事件** (event.go:398-432)

- `LifecycleEvent`: 生命周期（启用/停用/连接）
- `HeartbeatEvent`: 心跳

### 3. API 定义

**消息 API** (api.go:10-127)

- `SendPrivateMsgRequest/Response`: 发送私聊消息
- `SendGroupMsgRequest/Response`: 发送群消息
- `SendMsgRequest/Response`: 通用发送消息
- `DeleteMsgRequest/Response`: 撤回消息
- `GetMsgRequest/Response`: 获取消息
- `GetForwardMsgRequest/Response`: 获取合并转发消息

**群管理 API** (api.go:141-305)

- 踢人、禁言、设置管理员、群名片、群名、退群、专属头衔等

**信息查询 API** (api.go:307-487)

- 登录信息、陌生人信息、好友列表
- 群信息、群列表、群成员信息、群荣誉

**凭证 API** (api.go:489-525)

- Cookies、CSRF Token、Credentials

**媒体 API** (api.go:527-553)

- 获取语音、获取图片

**系统 API** (api.go:555-647)

- 能力检查、运行状态、版本信息、重启、清理缓存

### 4. 通信层实体

**ActionRequest** (communication.go:9-14)

传输层的动作请求：

```go
type ActionRequest struct {
    Action string         `json:"action"`
    Params map[string]any `json:"params,omitempty"`
}
```

**ActionRawResponse** (communication.go:17-22)

传输层的原始响应：

```go
type ActionRawResponse struct {
    Status  ActionResponseStatus  `json:"status"`
    Retcode ActionResponseRetcode `json:"retcode"`
    Data    json.RawMessage       `json:"data,omitempty"`
    Message string                `json:"message,omitempty"`
}
```

**ActionRequestEnvelope** (communication.go:25-29)

包含 Echo 字段的请求封装，用于 WebSocket 等需要请求-响应关联的场景：

```go
type ActionRequestEnvelope struct {
    ActionRequest
    Echo json.RawMessage `json:"echo,omitempty"`
}
```

**ActionResponseEnvelope** (communication.go:32-36)

包含 Echo 字段的响应封装：

```go
type ActionResponseEnvelope struct {
    ActionRawResponse
    Echo json.RawMessage `json:"echo,omitempty"`
}
```

**ActionResponse[T]** (communication.go:39-69)

泛型响应类型，Data 字段已解码：

```go
type ActionResponse[T any] struct {
    Status  ActionResponseStatus
    Retcode ActionResponseRetcode
    Data    *T
    Message string
}
```

**ActionError** (communication.go:71-85)

错误类型，实现了 `error` 接口。

---

## 关键依赖与配置

### 内部依赖

- `github.com/q1bksuu/onebot-go-sdk/v11/internal/util`: JSON 映射工具

### 外部依赖

无直接外部依赖（仅使用标准库）

### 代码生成配置

**entity-gen 触发**:

在 `base.go`, `message.go`, `event.go`, `api.go`, `communication.go` 文件顶部：

```go
//go:generate go run ../cmd/entity-gen
```

**bindings-gen 配置**:

`../cmd/bindings-gen/config.yaml` 定义了所有 API 的绑定关系：

- 8 个功能分组（message、friend、group_admin、group_info、account、media、capability、system）
- 每个 action 映射到对应的 Request/Response 类型

---

## 数据模型

### 类型层次

```
entity/
├── 基础类型
│   ├── StatusMeta (带 origin 字段的状态元信息)
│   ├── GroupAnonymousUser (匿名用户)
│   └── BaseUser (基础用户信息)
│
├── 消息类型
│   ├── MessageValue (字符串或数组)
│   └── Segment (消息段，含 SegmentData)
│
├── 事件类型
│   ├── Event (公共接口：GetTime, GetSelfId, GetPostType)
│   ├── 消息事件 (PrivateMessageEvent, GroupMessageEvent)
│   ├── 通知事件 (10+ 种)
│   ├── 请求事件 (FriendRequestEvent, GroupRequestEvent)
│   └── 元事件 (LifecycleEvent, HeartbeatEvent)
│
├── API 类型
│   ├── 消息 API (Send*, Delete*, Get*)
│   ├── 群管理 API (SetGroup*, SetFriendAddRequest)
│   ├── 查询 API (Get*Info, Get*List)
│   ├── 凭证 API (GetCookies, GetCsrfToken, GetCredentials)
│   ├── 媒体 API (GetRecord, GetImage)
│   └── 系统 API (GetStatus, GetVersionInfo, SetRestart, CleanCache)
│
├── 通信类型
│   ├── ActionRequest (动作请求)
│   ├── ActionRequestEnvelope (带 Echo 的请求封装)
│   ├── ActionRawResponse (原始响应)
│   ├── ActionResponseEnvelope (带 Echo 的响应封装)
│   ├── ActionResponse[T] (泛型响应)
│   └── ActionError (错误)
│
└── 常量定义
    ├── base_consts.go (性别、群角色等)
    ├── message_consts.go (消息类型、值类型)
    ├── event_consts.go (统一事件类型枚举)
    ├── api_consts.go (群荣誉类型、录音格式等)
    ├── communication_consts.go (响应状态、返回码)
    └── segment_data_consts.go (消息段类型、录音格式等)
```

### 动态字段机制

部分实体（如 `StatusMeta`, `GetVersionInfoResponse`）支持动态字段：

1. **字段定义**:
   ```go
   type StatusMeta struct {
       Online bool `json:"online"`
       Good   bool `json:"good"`
       origin map[string]any // 私有字段
   }
   ```

2. **自定义反序列化**:
   ```go
   func (r *StatusMeta) UnmarshalJSON(data []byte) error {
       return util.JsonUnmarshalToMapAndStruct(data, r, &r.origin)
   }
   ```

3. **访问方法**:
   ```go
   func (r *StatusMeta) GetOrigin(key string) any
   func (r *StatusMeta) SetOrigin(key string, value any) *StatusMeta
   ```

**优势**:

- 向前兼容：OneBot 实现可以添加自定义字段
- 不丢失数据：所有原始 JSON 字段都保留在 `origin`
- 类型安全：已知字段仍然是强类型

---

## 测试与质量

### 测试文件

- `base_test.go`: 测试 `StatusMeta` 的 JSON 反序列化和原始字段访问
- `api_test.go`: 测试 `GetVersionInfoResponse` 的 JSON 反序列化和原始字段访问

### 测试场景

**base_test.go:10-54** - `TestStatusMetaUnmarshalJSON`

- 验证标准字段正确解析
- 验证扩展字段存储到 `origin`
- 验证 `GetOrigin`/`SetOrigin` 方法

**api_test.go:10-66** - `TestGetVersionInfoResponseUnmarshalJSON`

- 验证版本信息字段正确解析
- 验证自定义字段（如 `extra_field`）存储到 `origin`
- 验证空指针安全

### 质量保证

- **代码生成**: Getter/Setter 由工具生成，避免手写错误
- **常量枚举**: 所有枚举类型都定义为 `type XxxType string` + 常量
- **JSON 标签**: 强制 `snake_case` 风格（由 golangci-lint 检查）
- **文档注释**: 所有公开类型和字段都有中文注释

---

## 常见问题 (FAQ)

**Q: 为什么有些类型有 Getter/Setter，有些没有？**

- 带有 `//go:generate go run ../cmd/entity-gen` 的文件会生成 Getter/Setter
- 自动生成的文件：`*_setter_getter.go`
- 手动维护的核心逻辑：`message.go` 的 `UnmarshalJSON` 方法

**Q: 如何添加新的事件类型？**

1. 在 `event.go` 定义新的结构体（参考现有事件）
2. 在 `event_consts.go` 添加对应的常量定义
3. 运行 `go generate ./...` 生成 Getter/Setter
4. 在 `../server/dispatcher.go` 注册事件处理器（如需要）

**Q: MessageValue 的 Type 字段如何自动设置？**

`UnmarshalJSON` 方法会根据 JSON 数据的实际类型自动设置：

- 如果是 `"string"`，则 `Type = MessageValueTypeString`
- 如果是 `[...]`，则 `Type = MessageValueTypeArray`

**Q: 为什么要用 `map[string]any` 而不是 `interface{}`？**

Go 1.18+ 推荐使用 `any` 代替 `interface{}`，更简洁。

**Q: origin 字段为什么是私有的？**

设计理念：

- 强制通过 `GetOrigin`/`SetOrigin` 访问，保证一致性
- 避免直接修改 map 导致的并发问题（虽然当前未实现并发保护）
- 未来可以在 Getter/Setter 中添加验证逻辑

---

## 相关文件清单

### 主要源文件 (14 个 .go 文件)

| 文件                               | 行数估算 | 职责                              |
| ---------------------------------- | -------- | --------------------------------- |
| `base.go`                          | ~60      | 基础类型定义                      |
| `message.go`                       | ~73      | 消息值类型与自定义序列化          |
| `event.go`                         | ~432     | Event 接口与所有事件类型定义      |
| `api.go`                           | ~648     | 所有 API 请求/响应定义            |
| `communication.go`                 | ~85      | 通信层类型（含 Envelope 封装）    |
| `segment_data.go`                  | (未展开) | 消息段数据定义                    |
| `base_consts.go`                   | (常量)   | 性别、群角色等枚举                |
| `message_consts.go`                | (常量)   | 消息类型、值类型枚举              |
| `event_consts.go`                  | ~121     | 统一事件类型枚举                  |
| `api_consts.go`                    | (常量)   | 群荣誉类型、录音格式等枚举        |
| `communication_consts.go`          | (常量)   | 响应状态、返回码枚举              |
| `segment_data_consts.go`           | (常量)   | 消息段类型枚举                    |

### 自动生成文件 (6 个 *_setter_getter.go)

| 文件                                    | 生成源       | 说明                     |
| --------------------------------------- | ------------ | ------------------------ |
| `api_setter_getter.go`                  | `api.go`     | API 类型的 Getter/Setter |
| `base_setter_getter.go`                 | `base.go`    | 基础类型的 Getter/Setter |
| `communication_setter_getter.go`        | `communication.go` | 通信类型的 Getter/Setter |
| `event_setter_getter.go`                | `event.go`   | 事件类型的 Getter/Setter |
| `message_setter_getter.go`              | `message.go` | 消息类型的 Getter/Setter |
| `segment_data_setter_getter.go`         | `segment_data.go` | 消息段的 Getter/Setter |

### 测试文件 (2 个)

- `base_test.go`: 基础类型测试
- `api_test.go`: API 类型测试

---

*模块文档更新时间: 2026-01-05*
