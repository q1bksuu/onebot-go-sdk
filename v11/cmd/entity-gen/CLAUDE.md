[根目录](../../../CLAUDE.md) > [v11](../../) > [cmd](../) > **entity-gen**

---

# entity-gen - 实体 Getter/Setter 代码生成工具

基于 Go AST 分析的代码生成工具，自动为结构体生成 Getter 和 Setter 方法。

---

## 变更记录 (Changelog)

### 2026-01-05

- **更新**: 同步文档与实际代码实现
- **修正**: Generator 结构体字段名
- **修正**: generator.go 行数 (612 行)
- **添加**: 新增 `templateField`, `templateStruct`, `structDef` 结构体说明
- **添加**: 泛型支持说明 (`receiverTypeParams`)

### 2025-12-21 15:53:08

- **初始化**: 生成工具文档
- **覆盖**: 扫描了完整的 README 和主要实现文件

---

## 模块职责

entity-gen 工具负责：

1. **自动生成样板代码**: 为结构体字段生成 Getter/Setter 方法
2. **智能类型处理**: 区分基本类型和复杂类型，采用不同策略
3. **空指针安全**: 所有 Getter 方法包含空指针检查
4. **链式调用支持**: Setter 方法返回接收者指针
5. **类型别名识别**: 自动扫描和识别自定义类型别名
6. **泛型支持**: 支持带类型参数的泛型结构体

---

## 入口与启动

### 命令行调用

```bash
# 基本用法
entity-gen -file=model.go

# 指定类型
entity-gen -file=model.go -type=User

# 多个类型（逗号分隔）
entity-gen -file=model.go -type=User,Product,Order

# 自定义输出文件
entity-gen -file=model.go -output=custom_output.go

# 指定额外的常量文件
entity-gen -file=model.go -consts=types_consts.go,status_consts.go

# 禁用自动扫描 *_consts.go
entity-gen -file=model.go -no-auto-scan
```

### go generate 集成（推荐）

在源文件顶部添加：

```go
//go:generate go run ../cmd/entity-gen/main.go -type=StatusMeta
package entity

type StatusMeta struct {
    Online bool `json:"online"`
    Good   bool `json:"good"`
}
```

运行：

```bash
go generate ./...
```

---

## 对外接口

### 命令行参数

| 参数              | 说明                   | 默认值                           | 示例                     |
| ----------------- | ---------------------- | -------------------------------- | ------------------------ |
| `-file`           | 要处理的 Go 源文件     | `$GOFILE` 环境变量               | `-file=model.go`         |
| `-type`           | 要处理的结构体类型名   | 空（处理所有）                   | `-type=User,Product`     |
| `-output`         | 输出文件路径           | `{filename}_setter_getter.go`    | `-output=generated.go`   |
| `-consts`         | 额外的常量文件         | 空                               | `-consts=types.go`       |
| `-no-auto-scan`   | 禁用自动扫描 `*_consts.go` | `false`                      | `-no-auto-scan`          |

### 生成的方法签名

**Getter**:

```go
// 基本类型字段 (int64, string, bool 等)
func (r *User) GetID() int64

// 基本类型指针字段 (*int64)
func (r *User) GetAge() int64  // 返回值类型，自动解引用

// 复杂类型字段 (struct, slice, map 等)
func (r *User) GetProfile() *Profile  // 返回指针

// 复杂类型指针字段
func (r *User) GetTags() []*Tag  // 返回原始指针
```

**Setter**:

```go
// 基本类型字段
func (r *User) SetID(v int64) *User

// 基本类型指针字段
func (r *User) SetAge(v int64) *User  // 参数是值类型，内部自动取地址

// 复杂类型字段
func (r *User) SetProfile(v *Profile) *User

// 复杂类型指针字段
func (r *User) SetTags(v []*Tag) *User
```

---

## 关键依赖与配置

### 外部依赖

无（仅使用标准库 `go/ast`, `go/parser`, `go/format`, `go/token`, `text/template`）

### 配置

无配置文件，所有选项通过命令行参数传递。

---

## 数据模型

### 核心类型

```go
// Generator 代码生成器主结构
type Generator struct {
    filename          string           // 输入文件路径
    fileSet           *token.FileSet   // 文件集
    astFile           *ast.File        // 主文件 AST
    customTypeAliases map[string]bool  // 自定义类型别名缓存（如 string 的别名）
    basicTypesCache   map[string]bool  // Go 基本类型缓存，提升性能
}

// fieldInfo 字段分析信息
type fieldInfo struct {
    name               string  // 字段名
    typeName           string  // 类型名称
    isPointer          bool    // 是否为指针
    isBasic            bool    // 是否为基本类型
    comments           string  // 注释内容
    receiverType       string  // 接收者类型
    receiverTypeParams string  // 接收者类型参数（泛型支持）
}

// templateField 模板渲染字段
type templateField struct {
    HasComments        bool    // 是否有注释
    Comments           string  // 注释内容
    ReceiverWithParams string  // 带类型参数的接收者
    GetterName         string  // Getter 方法名
    SetterName         string  // Setter 方法名
    ReturnType         string  // 返回类型
    ParamType          string  // 参数类型
    FieldName          string  // 字段名
    IsBasicPointer     bool    // 是否为基本类型指针
}

// templateStruct 模板渲染结构体
type templateStruct struct {
    Fields []templateField
}

// structDef 结构体定义信息
type structDef struct {
    name       string          // 结构体名称
    typeParams string          // 类型参数（泛型）
    structType *ast.StructType // AST 结构体类型
}
```

### 生成流程

```
1. 文件解析阶段
   ├─ 读取源文件和常量文件
   ├─ 使用 go/parser 构建 AST
   └─ 扫描自定义类型别名（如 type MyString string）

2. 结构体发现阶段
   ├─ 遍历 AST 寻找目标结构体
   ├─ 如果指定了 -type 参数，仅处理匹配的类型
   ├─ 提取类型参数（泛型支持）
   └─ 提取所有导出字段（首字母大写）

3. 字段分析阶段
   ├─ 判断是否为指针类型（*Type）
   ├─ 判断是否为基本类型或别名
   │  ├─ Go 内置：int, string, bool 等
   │  └─ 自定义别名：type Status string
   ├─ 提取注释文档
   └─ 记录接收者类型和类型参数

4. 代码生成阶段
   ├─ 为每个字段生成 Getter
   │  ├─ 添加空指针检查
   │  ├─ 基本类型指针：解引用
   │  └─ 复杂类型：直接返回
   ├─ 为每个字段生成 Setter
   │  ├─ 基本类型指针：自动取地址
   │  ├─ 复杂类型：直接赋值
   │  └─ 返回接收者（支持链式调用）
   └─ 使用 text/template 渲染代码

5. 文件写入阶段
   ├─ 使用 go/format 格式化输出
   └─ 写入到输出文件（默认 {filename}_setter_getter.go）
```

### 类型识别机制

**基本类型判断逻辑**:

```
是否为基本类型？
├─ 是 Go 内置基本类型？（int, string, bool 等）
│  └─ 是 → 基本类型
└─ 是自定义类型别名？（type Status string）
   └─ 别名的基础类型是基本类型？
      └─ 是 → 基本类型
```

**自定义类型别名扫描**:

工具会自动扫描：

1. **主文件**：当前处理的 Go 源文件
2. **同目录 `*_consts.go` 文件**：自动发现（除非 `-no-auto-scan`）
3. **手动指定文件**：通过 `-consts` 参数指定

识别形如 `type MyString string` 的定义，并记录到 `customTypeAliases`。

---

## 代码模板

生成器使用 `text/template` 渲染代码，核心模板结构：

```go
// Getter 模板
func (r *{{.ReceiverWithParams}}) {{.GetterName}}() {{.ReturnType}} {
    if r == nil {
        var zero {{.ReturnType}}
        return zero
    }
    // 基本类型指针：额外 nil 检查 + 解引用
    // 其他类型：直接返回
}

// Setter 模板
func (r *{{.ReceiverWithParams}}) {{.SetterName}}(v {{.ParamType}}) *{{.ReceiverWithParams}} {
    // 基本类型指针：val := v; r.Field = &val
    // 其他类型：r.Field = v
    return r
}
```

---

## 测试与质量

### 测试策略

工具本身没有单元测试文件，但：

- 在 `v11/entity` 中大量使用，生成的代码经过充分测试
- 生成的代码符合 golangci-lint 的 100+ linters 检查

### 质量保证

- **AST 准确性**: 基于 Go 标准库 `go/parser`，保证语法正确
- **代码格式化**: 使用 `go/format` 自动格式化生成的代码
- **注释保留**: 继承原字段的文档注释
- **错误处理**: 命令行参数验证、文件读写错误处理

---

## 常见问题 (FAQ)

**Q: 为什么基本类型指针字段的 Getter 返回值类型？**

设计理念：**对外 API 简化**

```go
// 字段定义
Age *int64

// 生成的方法
func (r *User) GetAge() int64 {  // 返回值类型，不是 *int64
    if r == nil {
        var zero int64
        return zero  // 安全的零值
    }
    if r.Age == nil {
        var zero int64
        return zero
    }
    return *r.Age  // 自动解引用
}

func (r *User) SetAge(v int64) *User {  // 参数是值类型
    val := v
    r.Age = &val  // 自动取地址
    return r
}
```

**优势**:

- 调用方无需关心指针细节
- 避免 nil 指针解引用 panic
- API 更简洁

**Q: 为什么复杂类型字段不做解引用？**

避免不必要的值拷贝：

```go
// 字段定义
Profile *UserProfile  // UserProfile 是一个大结构体

// 生成的方法
func (r *User) GetProfile() *UserProfile {  // 返回指针
    if r == nil {
        var zero *UserProfile
        return zero
    }
    return r.Profile  // 不解引用，避免拷贝
}
```

**Q: 如何判断是否为基本类型？**

内置基本类型列表（generator.go 中的 `initBasicTypesCache`）：

```
bool, string, error
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64, uintptr
float32, float64
complex64, complex128
byte, rune
```

+ 自定义类型别名（如 `type Status string`）

**Q: 链式调用如何使用？**

```go
user := &User{}
user.SetID(1001).
    SetName("alice").
    SetAge(25).
    SetStatus("active")
```

**Q: 为什么需要 -consts 参数？**

当常量定义在其他文件（如 `types.go`），且使用了 `-no-auto-scan` 时，需要手动指定：

```bash
entity-gen -file=model.go -consts=types.go -no-auto-scan
```

**Q: 生成的文件如何重新生成？**

直接覆盖：

```bash
go generate ./...
```

或删除后重新生成：

```bash
rm v11/entity/*_setter_getter.go
go generate ./v11/entity
```

**Q: 支持泛型吗？**

支持。工具会识别带类型参数的结构体：

```go
type Container[T any] struct {
    Value T
}
```

生成的方法会保留类型参数：

```go
func (r *Container[T]) GetValue() T { ... }
func (r *Container[T]) SetValue(v T) *Container[T] { ... }
```

---

## 相关文件清单

### 主要源文件

| 文件           | 行数  | 职责                              |
| -------------- | ----- | --------------------------------- |
| `main.go`      | 146   | 命令行参数解析和流程控制          |
| `generator.go` | 612   | AST 分析和代码生成逻辑            |
| `README.md`    | ~350  | 完整的工具文档                    |

### 使用示例

在 `v11/entity` 中的使用：

- `base.go` → `base_setter_getter.go`
- `api.go` → `api_setter_getter.go`
- `event.go` → `event_setter_getter.go`
- `message.go` → `message_setter_getter.go`
- `communication.go` → `communication_setter_getter.go`
- `segment_data.go` → `segment_data_setter_getter.go`

---

## 扩展阅读

详细使用说明和实现原理请参考：

- **完整文档**: [README.md](./README.md)
- **生成示例**: [v11/entity/*_setter_getter.go](../../entity/)

---

*工具文档更新时间: 2026-01-05*
