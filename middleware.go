package loggeriris

import (
	"github.com/kataras/iris/v12"
	"github.com/lianglong/logger"
)

// Config 中间件配置
type Config struct {
	// Logger 基础 logger 实例
	Logger logger.Logger

	// SkipPaths 跳过的路径
	SkipPaths []string

	// ExtractUserID 自定义 user_id 提取函数（可选）
	ExtractUserID func(ctx iris.Context) string

	// ExtractTraceID 自定义 trace_id 提取函数（可选）
	ExtractTraceID func(ctx iris.Context) string

	// CustomFields 自定义字段提取函数（可选）
	CustomFields func(ctx iris.Context) map[string]interface{}
}

// New 创建 Iris 日志中间件（使用默认配置）
func New(baseLogger logger.Logger) iris.Handler {
	return NewWithConfig(Config{
		Logger: baseLogger,
	})
}

// NewWithConfig 创建 Iris 日志中间件（使用自定义配置）
func NewWithConfig(cfg Config) iris.Handler {
	return func(ctx iris.Context) {
		// 跳过指定路径
		path := ctx.Path()
		for _, skip := range cfg.SkipPaths {
			if path == skip {
				ctx.Next()
				return
			}
		}

		stdCtx := ctx.Request().Context()

		// 1. 提取 request_id（兼容 Iris 的 requestid 中间件）
		requestID := ctx.GetID()
		if requestID != "" {
			stdCtx = logger.WithRequestID(stdCtx, requestID.(string))
		}

		// 2. 提取 user_id（如果配置了提取函数）
		if cfg.ExtractUserID != nil {
			if userID := cfg.ExtractUserID(ctx); userID != "" {
				stdCtx = logger.WithUserID(stdCtx, userID)
			}
		}

		// 3. 提取 trace_id（如果配置了提取函数）
		if cfg.ExtractTraceID != nil {
			if traceID := cfg.ExtractTraceID(ctx); traceID != "" {
				stdCtx = logger.WithTraceID(stdCtx, traceID)
			}
		}

		// 4. 创建带字段的 logger
		contextLogger := cfg.Logger.WithContext(stdCtx)

		// 5. 添加自定义字段（如果配置了）
		if cfg.CustomFields != nil {
			if fields := cfg.CustomFields(ctx); len(fields) > 0 {
				contextLogger = contextLogger.WithFields(fields)
			}
		}

		// 6. 将 logger 存入 context
		stdCtx = logger.WithLogger(stdCtx, contextLogger)

		// 7. 更新 request context
		ctx.ResetRequest(ctx.Request().WithContext(stdCtx))

		ctx.Next()
	}
}

// FromContext 从 Iris Context 获取 logger（便捷方法）
func FromContext(ctx iris.Context) logger.Logger {
	return logger.FromContext(ctx.Request().Context())
}

// GetLogger 从 Iris Context 获取 logger（别名）
func GetLogger(ctx iris.Context) logger.Logger {
	return FromContext(ctx)
}
