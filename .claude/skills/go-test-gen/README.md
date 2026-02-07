# go-test-gen Skill 使用指南

这是 onebot-go-sdk 项目的单元测试生成 skill，可以智能分析 Go 源代码并生成符合项目规范的测试代码。

## 快速开始

### 1. 为单个文件生成测试

```bash
# 基础用法
/go-test-gen v11/entity/base.go

# 生成结果
✅ 已生成: v11/entity/base_test.go
📊 新增测试函数: 8 个
🎯 预期覆盖率: 0% → 85.3%
```

### 2. 覆盖率驱动生成

```bash
# 分析覆盖率并补充缺失的测试
/go-test-gen --coverage v11/client/websocket.go

# 输出
当前覆盖率: 62.3%
目标覆盖率: 80.0%
缺口: -17.7%

未覆盖的函数:
  ❌ handlePing    0.0%
  ❌ handlePong    0.0%
  ⚠️  reconnect   45.2%

🔧 正在生成缺失的测试...
✅ 完成！预期覆盖率: 82.5%
```

### 3. 生成基准测试

```bash
# 为性能敏感的函数生成基准测试
/go-test-gen --benchmark v11/internal/util/radix_tree.go

# 生成内容
✅ BenchmarkRadixTree_Insert
✅ BenchmarkRadixTree_Insert_Parallel
✅ BenchmarkRadixTree_Search
✅ BenchmarkRadixTree_Delete
```

### 4. 生成 Mock 接口

```bash
# 为接口生成 Mock 实现
/go-test-gen --mock MessageSender v11/client/http_client.go

# 生成文件
✅ v11/client/mocks/message_sender_mock.go
```

### 5. 批量生成（整个包）

```bash
# 为整个包生成测试
/go-test-gen --package v11/entity

# 输出
正在分析包: v11/entity
  ✅ base.go → base_test.go (8 tests)
  ✅ message.go → message_test.go (12 tests)
  ✅ event.go → event_test.go (15 tests)
  ✅ api.go → api_test.go (38 tests)

📊 总计: 73 个测试函数
🎯 包覆盖率: 45.2% → 83.7%
```

## 功能特性

### ✅ 智能分析

- 🔍 解析函数签名和返回值
- 🎯 识别错误处理路径
- 📏 检测边界条件（nil、空值、零值）
- 🔗 分析依赖注入点
- ⚡ 识别并发敏感函数

### ✅ 规范遵循

- 📐 AAA 模式（Arrange-Act-Assert）
- ✨ testify 断言（assert/require）
- 🚀 并行测试（t.Parallel()）
- 🛠️ 测试辅助函数（t.Helper()）
- 🧹 资源清理（t.Cleanup()）

### ✅ 多种测试类型

- 📋 **表驱动测试**: 多场景测试
- 🎯 **简单测试**: 单一场景测试
- ⚡ **并发测试**: 竞态检测
- 📊 **基准测试**: 性能测试
- 🎭 **Mock 测试**: 接口模拟

## 生成的测试结构

### 表驱动测试示例

```go
func TestValidateUser(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        user    *User
        wantErr bool
        errMsg  string
    }{
        {
            name:    "nil_user_returns_error",
            user:    nil,
            wantErr: true,
            errMsg:  "user cannot be nil",
        },
        // ... 更多场景
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            // Act
            err := ValidateUser(tt.user)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                return
            }

            assert.NoError(t, err)
        })
    }
}
```

### 基准测试示例

```go
func BenchmarkRadixTree_Insert(b *testing.B) {
    tree := util.NewRadixTree[string]()

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        tree.Insert("key", "value")
    }
}
```

### Mock 生成示例

```go
type MockMessageSender struct {
    mock.Mock
}

func (m *MockMessageSender) SendPrivateMsg(ctx context.Context, req *SendPrivateMsgRequest) (*SendPrivateMsgResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*SendPrivateMsgResponse), args.Error(1)
}
```

## 测试场景生成策略

| 参数类型 | 自动生成的场景 |
|---------|---------------|
| 指针 | nil、有效值 |
| 切片 | nil、空切片、单元素、多元素 |
| 字符串 | 空字符串、普通字符串、特殊字符 |
| 数值 | 零值、负数、正数、边界值 |
| error 返回 | 成功场景、错误场景 |

## 最佳实践

### 1. 增量生成

```bash
# ✅ 推荐：新增函数后立即生成测试
git add websocket.go
/go-test-gen websocket.go
git add websocket_test.go

# ❌ 不推荐：等到功能完成后再生成
```

### 2. 覆盖率驱动

```bash
# 定期检查覆盖率
/go-test-gen --coverage ./...

# 优先补充关键路径
/go-test-gen --coverage v11/client
```

### 3. CI 集成

在 `.github/workflows/test.yml` 中：

```yaml
- name: Generate tests and check coverage
  run: |
    /go-test-gen --coverage --check ./...
    # 覆盖率低于目标则失败
```

### 4. Review 生成的代码

⚠️ 生成的测试是起点，不是终点：

- ✅ Review 测试逻辑是否正确
- ✅ 补充特定业务场景
- ✅ 调整断言和错误消息
- ✅ 添加文档注释

## 项目集成

### 作为 Claude Code Skill

将此目录放在 `.claude/skills/go-test-gen/`：

```
onebot-go-sdk/
├── .claude/
│   └── skills/
│       └── go-test-gen/
│           ├── skill.md           # Skill 定义
│           ├── README.md          # 本文件
│           └── examples/          # 示例代码
└── ...
```

### 在 IDE 中使用

如果你使用支持 Claude Code Skills 的 IDE：

1. 打开命令面板
2. 输入 `/go-test-gen`
3. 选择源文件
4. 自动生成测试

## 示例文件

- [source_example.go](examples/source_example.go) - 示例源代码
- [generated_test_example.go](examples/generated_test_example.go) - 生成的测试示例

## 限制和注意事项

### 当前限制

1. **泛型支持有限** - 泛型函数测试需手动补充类型实例化
2. **复杂 Mock** - 复杂接口的 Mock 可能需要调整
3. **黑盒测试优先** - 仅生成导出函数的测试
4. **并发测试需调优** - 自动生成的并发测试可能需要调整

### 使用建议

- ✅ 快速生成测试骨架
- ✅ 覆盖率驱动开发
- ✅ 生成 Mock 样板代码
- ⚠️ 生成后需 Review
- ⚠️ 复杂逻辑需手动补充
- ❌ 不能完全替代手工测试

## 常见问题

### Q: 生成的测试文件在哪里？

A: 与源文件同目录，文件名为 `<source>_test.go`

### Q: 如何跳过某些函数？

A: 在配置文件中设置 `exclude.functions` 正则表达式

### Q: 支持哪些测试框架？

A: 目前仅支持 testify，计划支持原生 testing

### Q: 如何生成并发测试？

A: Skill 会自动识别并发敏感函数（使用 sync、channel），或使用 `--concurrent` 标志

### Q: 覆盖率目标如何设置？

A: 使用 `--target` 标志指定目标覆盖率，例如 `/go-test-gen --coverage --target 85 <file>`

## 相关文档

- [项目测试策略](../../CLAUDE.md#单元测试编写流程)
- [testify 文档](https://github.com/stretchr/testify)
- [Go 测试最佳实践](https://go.dev/doc/tutorial/add-a-test)

---

**版本**: 1.0.0
**创建时间**: 2026-02-04
**维护者**: onebot-go-sdk 团队
