# Generator 模块文档

OneBot 11 Go SDK 代码生成器的核心模块。负责从 Markdown 文档中提取 API 定义，并生成强类型的 Go 代码。

## 快速开始

### 运行生成器

```bash
# 从项目根目录
python3 generator/main.py

# 或使用 make 命令
make generate

# 自定义参数
python3 generator/main.py \
    --input-dir ../api \
    --output-dir ../output \
    --package onebot
```

## 模块结构

### 1. schema.py - 数据结构定义

定义代码生成器使用的所有数据结构，包括：

- `FieldType` - 字段类型枚举（int64、string、bool 等）
- `MessageTypeVariant` - Message 字段的变体（string 或 array）
- `Field` - 单个字段的完整定义
- `APIModel` - API 的请求或响应模型
- `APIDefinition` - 完整的 API 定义

### 2. type_mapper.py - 类型映射系统

处理 Markdown 类型字符串到 Go 类型的转换。

**主要功能：**

- `parse_type()` - 解析类型字符串（如 "number (int64)" → int64）
- `snake_to_camel()` - 命名转换（user_id → userId）
- `snake_to_pascal()` - 命名转换（user_id → UserId）
- `determine_go_type_with_omitempty()` - 确定是否需要 omitempty

### 3. markdown_parser.py - Markdown 解析器

从 Markdown 文档中提取 API 定义。

**处理流程：**

1. 读取 Markdown 文件
2. 使用正则查找 `## \`api_name\`` 格式的 API
3. 提取 `### 参数` 和 `### 响应数据` 下的表格
4. 解析表格行提取字段信息
5. 清理 Markdown 特殊字符

### 4. go_generator.py - Go 代码生成器

根据 API 定义生成 Go 源代码。

**主要方法：**

- `generate_model_file_header()` - 生成文件头和导入
- `generate_api_models()` - 为单个 API 生成模型
- `_generate_field_code()` - 生成单个字段代码
- `generate_message_value_type()` - 生成特殊的 MessageValue 类型
- `generate_all_apis()` - 生成所有 API 的完整代码

### 5. main.py - 命令行入口

程序的主入口，处理参数和编排整个生成流程。

**命令行参数：**

```bash
python3 main.py [选项]

选项:
  --input-dir DIR      输入 Markdown 文档目录 (默认: ../api)
  --output-dir DIR     输出 Go 代码目录 (默认: ../output)
  --package NAME       Go 包名 (默认: onebot)
```

## 类型映射参考

| Markdown | Go 类型 | 说明 |
|----------|--------|------|
| number (int64) | int64 | 64-bit 整数 |
| number (int32) | int32 | 32-bit 整数 |
| string | string | 字符串 |
| boolean | bool | 布尔值 |
| object | *map[string]interface{} | 对象 |
| array | *[]interface{} | 数组 |
| message | MessageValue | OneBot 消息（特殊处理） |

## 字段可选性规则

- 有默认值 → 使用 omitempty
- 无默认值 → 不使用 omitempty
- message 字段 → 总是使用 omitempty

## 使用示例

### 解析 API

```python
from markdown_parser import MarkdownParser

parser = MarkdownParser()
apis = parser.parse_api_file("../api/public.md")

for api in apis:
    print(f"{api.api_name}: {api.description}")
    for field in api.request_model.fields:
        print(f"  - {field.go_name}: {field.go_type}")
```

### 生成代码

```python
from markdown_parser import MarkdownParser
from go_generator import GoCodeGenerator

parser = MarkdownParser()
apis = parser.parse_api_file("../api/public.md")

generator = GoCodeGenerator(package_name="onebot")
code = generator.generate_all_apis(apis)

with open("output/models.go", "w") as f:
    f.write(code)
```

## 开发指南

### 添加新的数据类型

1. 在 `schema.py` 中添加 `FieldType` 枚举值
2. 在 `type_mapper.py` 的 `TYPE_MAPPING` 中添加映射规则
3. 在 `go_generator.py` 中处理特殊生成逻辑（如需）

### 自定义生成格式

编辑 `go_generator.py` 中的：

- `generate_model_file_header()` - 修改文件头
- `_generate_field_code()` - 修改字段代码
- `_build_json_tag()` - 修改 JSON tag 格式
- `_build_field_comment()` - 修改注释生成

### 处理新的 Markdown 格式

修改 `markdown_parser.py` 中的：

- `extract_apis_from_content()` - API 定义查找逻辑
- `_extract_model_from_lines()` - 表格提取逻辑
- `_parse_table_row()` - 字段解析逻辑

## 常见问题

### Q: 为什么使用 MessageValue 而不是 interface{}？

A: MessageValue 提供精确的类型信息，避免运行时类型断言，这是强类型的优势。

### Q: 如何支持新的 OneBot API？

A: 只需在 Markdown 中添加新的 API 定义，重新运行生成器即可。

### Q: 为什么某些字段有 omitempty？

A: 有默认值的字段是可选的（加 omitempty），无默认值的字段是必需的（不加 omitempty）。

## 相关资源

- 上层文档: [../README.md](../README.md)
- 架构设计: [../ARCHITECTURE.md](../ARCHITECTURE.md)
- 实现总结: [../IMPLEMENTATION_SUMMARY.md](../IMPLEMENTATION_SUMMARY.md)

## 许可证

MIT License

---

**最后更新**: 2024-12-01
**版本**: 0.1.0
