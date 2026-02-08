# onebot-go-sdk

OneBot 11 协议的 Go 语言 SDK 实现。

[![codecov](https://codecov.io/github/q1bksuu/onebot-go-sdk/graph/badge.svg?token=F51V2M4L6G)](https://codecov.io/github/q1bksuu/onebot-go-sdk)

> **注意**: 本项目目前处于早期开发阶段，个人使用为主，API 尚未稳定，可能随时发生破坏性变更。不建议在生产环境中使用。

## 功能特性

- 完整的 OneBot 11 协议支持（消息、事件、API）
- 强类型实体定义，基于代码生成保证一致性
- 多种通信方式：
  - HTTP 客户端/服务端
  - WebSocket 客户端/服务端（正向/反向）
  - UnifiedServer 统一服务器（同时支持 HTTP 和 WebSocket）
- 灵活的事件分发机制

## 安装

```bash
go get github.com/q1bksuu/onebot-go-sdk
```

要求 Go 1.24 或更高版本。

## 快速开始

### HTTP 客户端

```go
package main

import (
    "context"
    "fmt"

    "github.com/q1bksuu/onebot-go-sdk/v11/client"
    "github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

func main() {
    c, err := client.NewHTTPClient(
        "http://127.0.0.1:5700",
        client.WithAccessToken("your-access-token"),
    )
    if err != nil {
        panic(err)
    }

    resp, err := c.SendPrivateMsg(context.Background(), &entity.SendPrivateMsgRequest{
        UserId: 123456789,
        Message: &entity.MessageValue{
            Type:        entity.MessageValueTypeString,
            StringValue: "Hello!",
        },
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("消息已发送，ID: %d\n", resp.MessageID)
}
```

### WebSocket 客户端

```go
package main

import (
    "context"
    "time"

    "github.com/q1bksuu/onebot-go-sdk/v11/client"
)

func main() {
    wsClient := client.NewWebSocketClient(
        client.WithWSURL("ws://127.0.0.1:6700"),
        client.WithWSAccessToken("your-access-token"),
        client.WithWSSelfID(123456789),
    )

    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        if err := wsClient.Start(ctx); err != nil {
            panic(err)
        }
    }()

    // ... your logic ...

    cancel()
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer shutdownCancel()
    _ = wsClient.Shutdown(shutdownCtx)
}
```

### HTTP 服务端（接收事件上报）

推荐使用“覆盖默认实现类 + 自动注册”的方式，写法更清晰，事件类型也更明确。

```go
package main

import (
    "context"
    "log"

    "github.com/q1bksuu/onebot-go-sdk/v11/entity"
    "github.com/q1bksuu/onebot-go-sdk/v11/server"
)

// 1) 覆盖默认实现类，只实现你关心的事件
type MyEventService struct {
    server.UnimplementedOneBotEventService
}

func (s *MyEventService) HandlePrivateMessage(ctx context.Context, ev *entity.PrivateMessageEvent) (map[string]any, error) {
    log.Printf("收到私聊消息: %s\n", ev.RawMessage)
    return map[string]any{"reply": "收到"}, nil
}

func main() {
    // 2) HTTP 服务配置
    cfg := server.HTTPConfig{
        Addr:        ":8080",
        EventPath:   "/event",
        AccessToken: "your-access-token",
    }

    // 3) 一步创建事件服务并启动
    srv := server.NewHTTPEventServer(cfg, &MyEventService{})

    if err := srv.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

注意：未覆盖的方法默认会 panic。请在 OneBot 客户端只启用你需要的事件，
或自行实现相关方法返回 `nil, nil` 来忽略事件。
`NewHTTPEventServer` 在 `EventPath` 为空时默认使用 `/event`，且 action 请求统一返回 404（如需处理 action，请使用 `NewHTTPServer` + `ActionHandler`）。

## 项目结构

```
onebot-go-sdk/
├── v11/
│   ├── entity/          # OneBot 协议实体定义
│   ├── client/          # HTTP/WebSocket 客户端
│   ├── server/          # HTTP/WebSocket 服务端
│   ├── dispatcher/      # Action 请求分发器
│   └── internal/util/   # 内部工具函数
└── go.mod
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/q1bksuu/onebot-go-sdk.git
cd onebot-go-sdk

# 运行测试
go test ./...

# 代码检查
golangci-lint run

# 重新生成代码
go generate ./...
```
