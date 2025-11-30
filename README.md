# OneBot 11 Go SDK 代码生成器

这是一个 OneBot 11 标准 Go SDK 代码生成器，可以将 Markdown 格式的 API 文档自动转换为强类型的 Go 代码。

## 核心特性 ✨

- **强类型支持** - 为每个 API 的请求和响应生成强类型的 Go 结构体
- **自动类型映射** - 将 Markdown 中的数据类型自动映射到对应的 Go 类型
- **Message 类型精确表示** - 使用 `MessageValue` 自定义类型（类似 Proto3 oneof）精确表示 message 字段，而不是使用 `interface{}`
- **自动 JSON Tag** - 根据字段是否有默认值自动生成 `json` tag 和 `omitempty`
- **文档注释生成** - 为每个字段和结构体自动生成中文注释

## 项目结构

```
onebot-11-go-sdk/
├── generator/                  # 代码生成器模块
│   ├── __init__.py            # 模块入口
│   ├── main.py                # 命令行入口
│   ├── schema.py              # 数据结构定义
│   ├── type_mapper.py         # 类型映射系统
│   ├── markdown_parser.py     # Markdown 文档解析器
│   └── go_generator.py        # Go 代码生成器
├── output/                     # 生成的 Go 代码输出目录
│   └── models.go              # 生成的模型代码
├── pyproject.toml             # Python 项目配置
└── README.md                  # 本文件
```

## 快速开始

### 1. 运行生成器

```bash
cd generator
python3 main.py
```

### 2. 查看生成结果

生成的 Go 代码会输出到 `output/models.go`：

```bash
cat output/models.go
```

### 3. 自定义输出

```bash
python3 main.py --input-dir /path/to/api/docs --output-dir /path/to/output --package mypackage
```

## 使用示例

### 简单的 API 调用

```go
package main

import (
    "encoding/json"
    "onebot"
)

func main() {
    // 创建发送私聊消息请求
    req := &onebot.SendPrivateMsgRequest{
        UserId: 123456789,
        Message: onebot.MessageValue{
            Type:        "string",
            StringValue: "Hello, World!",
        },
    }

    // 序列化为 JSON
    data, _ := json.Marshal(req)
    // 输出: {"user_id":123456789,"message":"Hello, World!"}
}
```

### Message 字段的两种用法

```go
// 方式 1: 纯文本字符串 (CQ 码格式)
msg1 := &onebot.SendPrivateMsgRequest{
    UserId: 123456789,
    Message: onebot.MessageValue{
        Type:        "string",
        StringValue: "[CQ:face,id=123]",
    },
}

// 方式 2: 消息段数组
msg2 := &onebot.SendPrivateMsgRequest{
    UserId: 123456789,
    Message: onebot.MessageValue{
        Type: "array",
        ArrayValue: []onebot.MessageSegment{
            {
                Type: "text",
                Data: map[string]interface{}{"text": "Hello"},
            },
            {
                Type: "face",
                Data: map[string]interface{}{"id": "123"},
            },
        },
    },
}

// 序列化时会自动选择正确的格式
```

## 类型系统设计

### Message 类型（Proto3 oneof 思想）

OneBot 中的 `message` 字段可以是两种类型：
1. **字符串** - CQ 码格式的消息字符串
2. **数组** - 消息段对象数组

使用 `MessageValue` 类型精确表示这两种情况：

```go
type MessageValue struct {
    Type        string              // "string" 或 "array"
    StringValue string              // 纯文本字符串
    ArrayValue  []MessageSegment    // 消息段数组
}

type MessageSegment struct {
    Type string                 `json:"type"`      // 消息段类型
    Data map[string]interface{} `json:"data"`      // 消息段参数
}
```

### 类型映射表

| Markdown 类型 | Go 类型 |
|-------------|--------|
| `number (int64)` | `int64` |
| `number (int32)` | `int32` |
| `string` | `string` |
| `boolean` | `bool` |
| `object` | `*map[string]interface{}` |
| `array` | `*[]interface{}` |
| `message` | `MessageValue` |

## 代码生成流程

1. **解析 Markdown** - 使用正则表达式从 Markdown 中提取 API 定义
2. **提取表格** - 解析参数表和响应表，提取字段信息
3. **类型映射** - 将 Markdown 类型字符串映射到 Go 类型
4. **生成代码** - 根据字段信息生成 Go 结构体代码
5. **添加注释** - 为每个字段添加对应的中文注释

## 生成示例

### 输入 (Markdown)

```markdown
## `send_private_msg` 发送私聊消息

### 参数

| 字段名 | 数据类型 | 默认值 | 说明 |
| ----- | ------- | ----- | --- |
| `user_id` | number (int64) | - | 对方 QQ 号 |
| `message` | message | - | 要发送的内容 |
| `auto_escape` | boolean | `false` | 消息内容是否作为纯文本发送 |

### 响应数据

| 字段名 | 数据类型 | 说明 |
| ----- | ------- | --- |
| `message_id` | number (int32) | 消息 ID |
```

### 输出 (Go Code)

```go
// SendPrivateMsgRequest 表示 send_private_msg API 的请求参数
// 对应文档: 发送私聊消息
type SendPrivateMsgRequest struct {
    // 对方 QQ 号
    UserId int64 `json:"user_id,omitempty"`
    // 要发送的内容
    // 可以是字符串 (CQ 码格式) 或消息段数组
    Message MessageValue `json:"message,omitempty"`
    // 消息内容是否作为纯文本发送
    AutoEscape bool `json:"auto_escape,omitempty"`
}

// SendPrivateMsgResponse 表示 send_private_msg API 的响应数据
type SendPrivateMsgResponse struct {
    // 消息 ID
    MessageId int32 `json:"message_id"`
}
```

## 设计原则

### KISS (Keep It Simple, Stupid)
- 代码简洁直观
- 没有过度的模板和复杂性
- 直接生成易于理解的 Go 代码

### YAGNI (You Aren't Gonna Need It)
- 只生成当前需要的代码
- 不预留未来可能的扩展点
- 不包含未使用的依赖

### DRY (Don't Repeat Yourself)
- 避免重复代码
- 公共逻辑提取到工具函数
- 使用强类型避免重复的类型检查

## 开发指南

### 添加新的 API

1. 在 OneBot 文档中添加新的 API 定义
2. 重新运行生成器
3. 新的 API 对应的代码会自动生成

### 修改现有的类型映射

编辑 `type_mapper.py` 中的 `TYPE_MAPPING` 字典来修改类型映射规则。

### 自定义生成的代码格式

编辑 `go_generator.py` 中的 `_generate_*` 方法来自定义生成的代码格式。

## 限制和已知问题

- 目前不支持嵌套的对象类型（只能识别为 `*map[string]interface{}`）
- 对于复杂的 array 类型，元素类型默认为 `interface{}`
- 不支持自定义的字段验证规则

## 下一步

- [ ] 实现事件模型生成
- [ ] 实现消息段类型生成
- [ ] 支持嵌套对象类型
- [ ] 添加单元测试
- [ ] 生成 API 客户端方法

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
