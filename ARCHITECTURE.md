# 架构设计文档

## 概览

OneBot 11 Go SDK 代码生成器采用模块化设计，通过将 Markdown 文档作为"真实数据源"，自动生成强类型的 Go 代码，避免手工维护。

## 核心模块

### 1. Schema 模块 (`schema.py`)

定义所有数据结构的类型系统。

**主要类:**
- `FieldType` - 字段类型枚举（int64, string, bool, message 等）
- `MessageTypeVariant` - Message 字段的变体类型（string 或 array）
- `Field` - 单个字段的完整描述
- `APIModel` - API 的请求或响应模型
- `APIDefinition` - 完整的 API 定义
- `EventModel` - 事件模型定义
- `MessageSegment` - 消息段定义

### 2. Type Mapper 模块 (`type_mapper.py`)

处理类型系统的映射和转换。

**主要功能:**
- `parse_type()` - 将 Markdown 类型字符串映射到 Go 类型
- `snake_to_camel()` - 命名规则转换（蛇形 → 驼峰）
- `snake_to_pascal()` - 命名规则转换（蛇形 → 帕斯卡）
- `determine_go_type_with_omitempty()` - 确定字段是否需要 omitempty tag

**类型映射表:**

```
Markdown          →  Go 类型
number (int64)   →  int64
number (int32)   →  int32
string           →  string
boolean          →  bool
object           →  *map[string]interface{}
array            →  *[]interface{}
message          →  MessageValue
```

### 3. Markdown Parser 模块 (`markdown_parser.py`)

从 Markdown 文档中提取 API 定义。

**处理流程:**

1. **文件加载** - 读取 Markdown 文件
2. **API 定义查找** - 使用正则查找 `## \`api_name\`` 格式的 API
3. **表格提取** - 解析 `### 参数` 和 `### 响应数据` 下的表格
4. **字段解析** - 从表格行中提取字段信息
5. **数据清理** - 去除 Markdown 格式的特殊字符

**关键方法:**

```python
parse_api_file(file_path)              # 解析整个 API 文件
extract_apis_from_content(content)     # 从内容中提取 API 列表
_extract_model_from_lines()            # 从行列表中提取模型
_parse_table_row()                     # 解析单行表格
```

### 4. Go Code Generator 模块 (`go_generator.py`)

根据 API 定义生成 Go 代码。

**主要方法:**

```python
generate_model_file_header()           # 生成文件头和导入
generate_api_models(api_def)           # 为单个 API 生成模型
_generate_request_model()              # 生成请求模型
_generate_response_model()             # 生成响应模型
_generate_field_code()                 # 生成单个字段代码
generate_message_value_type()          # 生成 MessageValue 类型
generate_all_apis(apis)                # 生成所有 API 的代码
```

### 5. Main 模块 (`main.py`)

命令行入口和编排逻辑。

**流程:**

1. 解析命令行参数
2. 初始化 Parser 和 Generator
3. 读取 Markdown 文件
4. 生成 Go 代码
5. 输出到文件

## 数据流

```
┌─────────────────┐
│ Markdown 文件   │ (api/public.md)
└────────┬────────┘
         │
         ↓
┌──────────────────────┐
│ MarkdownParser       │
│ - 解析 API 定义      │
│ - 提取表格          │
│ - 清理字段名        │
└────────┬─────────────┘
         │
         ↓
┌──────────────────────┐
│ APIDefinition        │ (内存中的 API 模型)
│ - request_model      │
│ - response_model     │
└────────┬─────────────┘
         │
         ↓
┌──────────────────────┐
│ GoCodeGenerator      │
│ - 类型映射          │
│ - 生成 struct       │
│ - 生成 JSON tag     │
└────────┬─────────────┘
         │
         ↓
┌──────────────────────┐
│ Go 源代码文件        │ (output/models.go)
│ - MessageValue       │
│ - *Request struct   │
│ - *Response struct  │
└──────────────────────┘
```

## MessageValue 类型设计

由于 OneBot 中的 `message` 字段可以是两种类型（字符串或数组），我们设计了 `MessageValue` 类型来精确表示，参考了 Protocol Buffer 的 oneof 思想。

```go
type MessageValue struct {
    Type        string              // 标志类型：string 或 array
    StringValue string              // 纯文本字符串
    ArrayValue  []MessageSegment    // 消息段数组
}
```

**实现 JSON 序列化/反序列化:**

- `UnmarshalJSON()` - 在反序列化时自动检测类型
- `MarshalJSON()` - 在序列化时输出正确的格式

## 设计原则

### KISS (Keep It Simple, Stupid)

- 代码清晰易懂，没有过度的抽象
- 直接处理问题，避免不必要的中间层
- 函数职责单一，易于测试和维护

### YAGNI (You Aren't Gonna Need It)

- 只生成当前需要的代码
- 不预留未来可能的扩展点
- 不包含未使用的配置选项

### DRY (Don't Repeat Yourself)

- 提取公共逻辑到工具函数
- 使用强类型避免重复的类型转换
- 利用代码生成避免手工维护的重复代码

## 扩展性

### 添加新的 API

无需修改生成器，只需在 Markdown 中添加新的 API 定义，重新运行生成器即可。

### 支持新的数据类型

修改 `type_mapper.py` 中的 `TYPE_MAPPING` 即可：

```python
TYPE_MAPPING = {
    r"custom_type": (FieldType.CUSTOM, "CustomGoType"),
    # ...
}
```

### 自定义生成格式

修改 `go_generator.py` 中的模板方法，如 `_generate_field_code()`。

## 限制和已知问题

1. **嵌套对象** - 不支持复杂的嵌套对象结构，会被简化为 `*map[string]interface{}`
2. **Array 元素类型** - 复杂的 array 元素类型不支持，默认为 `interface{}`
3. **字段验证** - 不支持从 Markdown 自动生成字段验证规则
4. **事件和消息段** - 目前只支持 API 模型生成，未来需要支持事件和消息段

## 性能考虑

- Markdown 解析使用正则表达式，对于大文件性能足够
- 代码生成是一次性操作，性能不是关键因素
- 内存占用极小（<10MB 用于生成 38 个 API）

## 测试策略

由于生成的代码需要符合 Go 的语法要求，我们使用：

1. **语法验证** - 运行 `go fmt` 检查生成的代码
2. **集成测试** - 验证生成的代码是否能正确反序列化 JSON
3. **手工审查** - 检查生成代码的质量和准确性

## 未来改进

- [ ] 支持 `go:generate` 集成
- [ ] 添加事件模型生成
- [ ] 添加消息段类型生成
- [ ] 生成 API 客户端方法
- [ ] 支持自定义代码模板
- [ ] 添加单元测试框架
- [ ] 性能基准测试
