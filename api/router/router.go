package router

import (
	"context"

	"trade/api/handler"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// Register 注册路由
func Register(h *server.Hertz) {
	// API v1 路由组
	v1 := h.Group("/api/v1")

	// 策略分析路由
	strategy := v1.Group("/strategy")
	{
		// POST /api/v1/strategy/analyze - 策略分析接口
		strategy.POST("/analyze", handler.AnalyzeStrategy)

		// GET /api/v1/strategy/analyze - 支持GET请求(通过query参数)
		strategy.GET("/analyze", handler.AnalyzeStrategy)
	}

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// 首页
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"name":    "Trade Strategy API",
			"version": "1.0.0",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/v1/strategy/analyze",
				"POST /api/v1/strategy/analyze",
			},
		})
	})
}
