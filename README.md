# onebot-go-sdk

OneBot 11 协议的 Go 语言 SDK 实现。

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

要求 Go 1.25 或更高版本。

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

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/q1bksuu/onebot-go-sdk/v11/entity"
    "github.com/q1bksuu/onebot-go-sdk/v11/server"
)

func main() {
    srv := server.NewHTTPServer(server.HTTPServerConfig{
        AccessToken: "your-access-token",
    })

    srv.OnPrivateMessage(func(w http.ResponseWriter, r *http.Request, event *entity.PrivateMessageEvent) {
        fmt.Printf("收到私聊消息: %s\n", event.RawMessage)
    })

    http.ListenAndServe(":8080", srv)
}
```

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
