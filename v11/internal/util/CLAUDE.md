[根目录](../../../CLAUDE.md) > [v11](../../) > [internal](../) > **util**

---

# util - 内部工具函数模块

提供 JSON 映射、类型转换、Radix Tree 数据结构等内部工具函数。

---

## 变更记录 (Changelog)

### 2026-01-05

- **新增**: `RadixTree` 泛型数据结构 (`radix_tree.go`)
- **新增**: `NormalizePath` 路径规范化函数
- **新增**: `radix_tree_test.go` (12 个单元测试)
- **依赖**: 新增 `github.com/armon/go-radix` 依赖

### 2025-12-21 15:53:08

- **初始化**: 生成模块级文档
- **覆盖**: 扫描了工具函数实现和完整的测试套件（含性能基准）

---

## 模块职责

util 模块负责：

1. **JSON 同步映射**: 将 JSON 数据同时解析到 map 和 struct
2. **类型安全转换**: 基于 `json` 标签的 map 到 struct 映射
3. **动态字段保留**: 支持保留未在 struct 定义的额外字段
4. **路径规范化**: URL 路径格式统一处理
5. **Radix Tree**: 提供泛型前缀树数据结构，用于高效前缀匹配

---

## 入口与启动

util 是纯函数库，无需初始化，直接调用即可。

---

## 对外接口

### 1. JsonUnmarshalToMapAndStruct (util.go:14-30)

**签名**:

```go
func JsonUnmarshalToMapAndStruct(data []byte, dest any, destMap *map[string]any) error
```

**功能**: 将 JSON 数据同时解析到 struct 和 map。

**参数**:

- `data`: JSON 字节数组
- `dest`: 目标结构体指针（如 `*StatusMeta`）
- `destMap`: 目标 map 指针（用于保存所有原始字段）

**返回**:

- `error`: 解析失败或映射失败时返回错误

**使用场景**:

在 `entity.StatusMeta` 和 `entity.GetVersionInfoResponse` 的 `UnmarshalJSON` 方法中使用，保留未知扩展字段。

**示例**:

```go
type User struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

var user User
var origin map[string]any
data := []byte(`{"id": 123, "name": "Alice", "extra_field": "value"}`)

err := util.JsonUnmarshalToMapAndStruct(data, &user, &origin)
// user.ID = 123, user.Name = "Alice"
// origin["extra_field"] = "value"
```

### 2. JsonTagMapping (util.go:32-47)

**签名**:

```go
func JsonTagMapping(source, dest any) error
```

**功能**: 将 map 或 struct 映射到另一个 struct，基于 `json` 标签。

**参数**:

- `source`: 源数据（通常是 `map[string]any`）
- `dest`: 目标结构体指针

**返回**:

- `error`: 创建 decoder 失败或解码失败时返回错误

**内部实现**:

使用 `github.com/go-viper/mapstructure/v2` 库，配置 `TagName: "json"`。

**使用场景**:

- `client.HTTPClient.encodeToParams`: 将 Request 结构体转为 map
- `server.Binder.Handler`: 将 map 参数绑定到 Request 结构体

**示例**:

```go
params := map[string]any{
    "user_id": 123456,
    "message": "Hello",
}

var req entity.SendPrivateMsgRequest
err := util.JsonTagMapping(params, &req)
// req.UserId = 123456, req.Message = "Hello"
```

### 3. NormalizePath (util.go:49-56)

**签名**:

```go
func NormalizePath(path string) string
```

**功能**: 规范化 URL 路径，确保以 `/` 开头且不以 `/` 结尾。

**参数**:

- `path`: 原始路径字符串

**返回**:

- `string`: 规范化后的路径，空路径返回空字符串

**处理规则**:

- 移除首尾的 `/`
- 如果结果非空，在开头添加 `/`
- 空输入返回空字符串

**使用场景**:

- `server.NewHTTPServer`: 规范化 `EventPath` 和 `ActionPath`
- `server.NewWebSocketServer`: 规范化 `PathPrefix`
- `server.WebSocketServer.matchPath`: 路径匹配时的规范化

**示例**:

```go
util.NormalizePath("/api/")     // "/api"
util.NormalizePath("api")       // "/api"
util.NormalizePath("/api/v1/")  // "/api/v1"
util.NormalizePath("")          // ""
util.NormalizePath("/")         // ""
```

### 4. RadixTree (radix_tree.go)

**类型定义**:

```go
// RadixTree 是对 radix.Tree 的泛型封装，提供类型安全的操作接口
// K 是键类型，约束为底层类型为 string 的类型
// T 是值类型
type RadixTree[K ~string, T any] struct {
    tree *radix.Tree
}

// RadixTreeStrKey 是字符串键类型的别名
type RadixTreeStrKey[T any] = RadixTree[string, T]
```

**功能**: 泛型 Radix Tree（前缀树/基数树）数据结构，支持高效的前缀匹配和查找。

**构造函数**:

```go
// 创建空树
func NewRadixTree[K ~string, T any]() *RadixTree[K, T]

// 从 map 创建
func NewRadixTreeFromMap[K ~string, T any](m map[K]T) *RadixTree[K, T]
```

**方法列表**:

| 方法 | 签名 | 功能 |
|------|------|------|
| `Insert` | `(key K, value T)` | 插入或更新键值对 |
| `Get` | `(key K) (T, bool)` | 根据键获取值 |
| `Delete` | `(key K) (T, bool)` | 删除指定键 |
| `LongestPrefix` | `(prefix K) (K, T, bool)` | 查找最长前缀匹配 |
| `Minimum` | `() (K, T, bool)` | 返回最小键值对 |
| `Maximum` | `() (K, T, bool)` | 返回最大键值对 |
| `Walk` | `(fn func(key K, value T) bool)` | 遍历所有键值对 |
| `WalkPrefix` | `(prefix K, fn func(key K, value T) bool)` | 遍历指定前缀的键值对 |
| `Len` | `() int` | 返回键值对数量 |
| `ToMap` | `() map[K]T` | 转换为 map |

**使用场景**:

- `server.HTTPServer`: 使用 `eventRadixTree` 进行事件类型的前缀匹配路由

**示例**:

```go
tree := util.NewRadixTree[string, int]()
tree.Insert("/api/users", 1)
tree.Insert("/api/groups", 2)
tree.Insert("/api/users/profile", 3)

// 精确查找
value, ok := tree.Get("/api/users")  // 1, true

// 最长前缀匹配
key, value, ok := tree.LongestPrefix("/api/users/123")  // "/api/users", 1, true

// 遍历前缀
tree.WalkPrefix("/api/", func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)
    return true  // 继续遍历
})
```

---

## 关键依赖与配置

### 外部依赖

- `github.com/go-viper/mapstructure/v2` (v2.4.0): 类型映射库
- `github.com/armon/go-radix`: Radix Tree 底层实现

### 无配置

纯函数库，无需配置。

---

## 数据模型

### 核心类型

| 类型 | 文件 | 描述 |
|------|------|------|
| `RadixTree[K, T]` | radix_tree.go | 泛型前缀树 |
| `RadixTreeStrKey[T]` | radix_tree.go | 字符串键前缀树别名 |

### 处理流程

**JsonUnmarshalToMapAndStruct 流程**:

```
1. 验证 destMap 指针非 nil
   ↓
2. 使用 json.Unmarshal 解析到 map
   ↓
3. 调用 JsonTagMapping 映射 map 到 struct
   ↓
4. 返回结果
   (struct 包含已知字段，map 包含所有字段)
```

**JsonTagMapping 流程**:

```
1. 创建 mapstructure.Decoder
   ├─ TagName: "json"
   └─ Result: dest
   ↓
2. 调用 decoder.Decode(source)
   ├─ 根据 json 标签匹配字段
   └─ 自动类型转换（如 float64 -> int64）
   ↓
3. 返回结果
```

**NormalizePath 流程**:

```
1. strings.Trim(path, "/") 移除首尾斜杠
   ↓
2. 如果结果为空，返回空字符串
   ↓
3. 否则返回 "/" + path
```

---

## 测试与质量

### 测试文件

- `util_test.go`: **23 个单元测试 + 4 个性能基准测试** (~579 行)
- `radix_tree_test.go`: **12 个单元测试** (~237 行)

### util_test.go 单元测试场景

| 测试函数                                      | 测试内容                              |
| --------------------------------------------- | ------------------------------------- |
| `TestJsonUnmarshalToMapAndStructBasic`        | 基本字段解析                          |
| `TestUnmarshalToMapAndStructNested`           | 嵌套结构体                            |
| `TestUnmarshalToMapAndStructExtraFields`      | 额外字段保留到 map                    |
| `TestUnmarshalToMapAndStructArray`            | 数组字段                              |
| `TestUnmarshalToMapAndStructEmptyJSON`        | 空 JSON 对象                          |
| `TestUnmarshalToMapAndStructNullValues`       | null 值处理                           |
| `TestUnmarshalToMapAndStructInvalidJSON`      | 无效 JSON 错误处理                    |
| `TestUnmarshalToMapAndStructTypeConversion`   | 类型转换（float64 -> int64）          |
| `TestUnmarshalToMapAndStructUnicodeCharacters`| Unicode 字符处理                      |
| `TestUnmarshalToMapAndStructLargeNumbers`     | 大数字处理                            |
| `TestUnmarshalToMapAndStructMapPointerNil`    | map 指针为 nil 错误处理               |
| `TestUnmarshalToMapAndStructComplexStructure` | 复杂结构（嵌套 + 数组 + 额外字段）    |
| `TestUnmarshalToMapAndStructWithSpecialCharacters` | 特殊字符处理                     |
| `TestUnmarshalToMapAndStructBooleanValues`    | 布尔值处理                            |

### radix_tree_test.go 单元测试场景

| 测试函数                          | 测试内容                              |
| --------------------------------- | ------------------------------------- |
| `TestNewRadixTree`                | 空树创建                              |
| `TestNewRadixTreeFromMap`         | 从 map 创建树                         |
| `TestRadixTree_Insert`            | 插入操作                              |
| `TestRadixTree_Get`               | 查找操作                              |
| `TestRadixTree_Delete`            | 删除操作                              |
| `TestRadixTree_LongestPrefix`     | 最长前缀匹配                          |
| `TestRadixTree_Minimum`           | 最小键查找                            |
| `TestRadixTree_Maximum`           | 最大键查找                            |
| `TestRadixTree_Walk`              | 全树遍历                              |
| `TestRadixTree_WalkPrefix`        | 前缀遍历                              |
| `TestRadixTree_ToMap`             | 转换为 map                            |
| `TestRadixTree_ComplexType`       | 复杂值类型支持                        |

### 性能基准测试

| Benchmark 函数                               | 测试场景                              |
| -------------------------------------------- | ------------------------------------- |
| `BenchmarkUnmarshalToMapAndStructSmallPayload` | 小负载（3 字段）                    |
| `BenchmarkUnmarshalToMapAndStructMediumPayload`| 中等负载（10+ 字段，含嵌套）        |
| `BenchmarkUnmarshalToMapAndStructLargePayload` | 大负载（20+ 字段，复杂嵌套）        |
| `BenchmarkUnmarshalToMapAndStructDeepNesting`  | 深层嵌套（4 层）                    |

**基准测试配置**: `-benchtime=10x` (每个 benchmark 运行 10 次)

### 质量保证

- **边界条件**: 空 JSON、null 值、无效 JSON、nil 指针
- **类型转换**: 测试 JSON number -> Go int64 的转换
- **Unicode 支持**: 测试中文、emoji 等字符
- **大数字**: 测试 int64 范围边界值
- **性能优化**: 通过 benchmark 确保性能可接受
- **泛型支持**: RadixTree 测试覆盖多种键值类型组合

---

## 常见问题 (FAQ)

**Q: 为什么需要 JsonUnmarshalToMapAndStruct？**

OneBot 11 协议允许实现添加自定义扩展字段，单纯使用 struct 会丢失这些字段。同时保存到 map 可以：

- 保留所有原始数据（向前兼容）
- 通过 `GetOrigin(key)` 访问扩展字段
- 不影响已知字段的类型安全

**Q: JsonTagMapping 和 json.Unmarshal 有什么区别？**

| 特性           | json.Unmarshal        | JsonTagMapping        |
| -------------- | --------------------- | --------------------- |
| 输入           | JSON 字节数组         | map 或 struct         |
| 输出           | struct                | struct                |
| 类型转换       | 严格（会报错）        | 宽松（自动转换）      |
| 标签支持       | `json` 标签           | 可配置（默认 `json`） |
| 使用场景       | HTTP 响应解析         | 内部类型转换          |

**Q: 为什么使用 mapstructure 而不是手写反射代码？**

- **成熟稳定**: mapstructure 是社区广泛使用的库
- **功能完善**: 支持复杂类型转换、自定义 decoder
- **性能优化**: 内部有缓存机制
- **维护成本低**: 避免手写反射代码的复杂性

**Q: map 指针为 nil 会发生什么？**

`JsonUnmarshalToMapAndStruct` 会返回 `errInvalidMapPointer` 错误，避免 panic。

**Q: 性能如何？**

根据 benchmark 结果（具体数值取决于硬件）：

- **小负载**: ~2-5 μs/op
- **中等负载**: ~10-20 μs/op
- **大负载**: ~30-50 μs/op

对于大多数场景（HTTP 请求处理），性能完全可接受。

**Q: 是否支持并发调用？**

是，所有函数都是无状态的，可以并发调用。RadixTree 本身不是线程安全的，需要外部同步。

**Q: 为什么需要 RadixTree？**

Radix Tree 提供 O(k) 复杂度的前缀匹配（k 为键长度），适用于：

- URL 路由匹配
- 事件类型的层次化分发
- 需要前缀查找的场景

相比 map，RadixTree 在前缀匹配场景下更高效，且支持 `LongestPrefix` 和 `WalkPrefix` 等高级操作。

**Q: NormalizePath 的作用是什么？**

确保路径格式一致，避免因 `/api` 和 `api/` 和 `/api/` 格式不同导致的匹配问题。

---

## 相关文件清单

### 主要源文件

| 文件             | 行数 | 职责                              |
| ---------------- | ---- | --------------------------------- |
| `util.go`        | ~57  | JSON 映射工具函数、路径规范化     |
| `radix_tree.go`  | ~181 | 泛型 Radix Tree 数据结构          |

### 测试文件

| 文件                 | 行数  | 职责                              |
| -------------------- | ----- | --------------------------------- |
| `util_test.go`       | ~579  | 23 个单元测试 + 4 个性能基准测试  |
| `radix_tree_test.go` | ~237  | 12 个单元测试                     |

---

*模块文档更新时间: 2026-01-05*
