
# logger-iris

Iris 框架的日志中间件，基于 [github.com/lianglong/logger](https://github.com/lianglong/logger)。

## 安装

```bash
go get github.com/lianglong/logger-iris
go get github.com/lianglong/logger        # 核心包
go get github.com/lianglong/logger-zap    # Zap 驱动（或其他驱动）
```

## 功能特性

- 🔄 自动提取 request_id
- 👤 支持自定义 user_id 提取
- 🔍 支持 trace_id（分布式追踪）
- 📝 支持自定义字段
- 🚫 支持路径跳过
- 🎯 框架无关的设计

## 快速开始

```go
package main

import (
    "os"
    "github.com/kataras/iris/v12"
    "github.com/lianglong/logger"
    _ "github.com/lianglong/logger-zap"
    loggeriris "github.com/lianglong/logger-iris"
)

func main() {
    app := iris.New()

    // 创建 logger
    log := logger.MustNew("zap", logger.Config{
        Level:  logger.InfoLevel,
        Output: os.Stdout,
    })
    defer log.Sync()

    // 使用中间件
    app.Use(loggeriris.New(log))

    // 在 handler 中使用
    app.Get("/", func(ctx iris.Context) {
        log := loggeriris.FromContext(ctx)
        log.Info("hello world")
        ctx.JSON(iris.Map{"message": "ok"})
    })

    app.Listen(":8080")
}
```

## 高级配置

```go
app.Use(loggeriris.NewWithConfig(loggeriris.Config{
    Logger:    log,
    SkipPaths: []string{"/health", "/metrics"},
    ExtractUserID: func(ctx iris.Context) string {
        // 从 JWT 提取 user_id
        if claims := ctx.Values().Get("jwt_claims"); claims != nil {
            return claims.(map[string]interface{})["user_id"].(string)
        }
        return ""
    },
    ExtractTraceID: func(ctx iris.Context) string {
        return ctx.GetHeader("X-Trace-ID")
    },
    CustomFields: func(ctx iris.Context) map[string]interface{} {
        return map[string]interface{}{
            "ip":     ctx.RemoteAddr(),
            "method": ctx.Method(),
            "path":   ctx.Path(),
        }
    },
}))
```

## API

### New(logger.Logger) iris.Handler
创建使用默认配置的中间件

### NewWithConfig(Config) iris.Handler
创建使用自定义配置的中间件

### FromContext(iris.Context) logger.Logger
从 Iris Context 获取 logger 实例

### GetLogger(iris.Context) logger.Logger
FromContext 的别名
