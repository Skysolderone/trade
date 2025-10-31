package main

import (
	"fmt"

	"trade/api/router"
	"trade/db"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 初始化数据库连接
	fmt.Println("正在初始化数据库连接...")
	db.InitPostgreSql()
	fmt.Println("✓ 数据库连接成功")

	// 初始化币安客户端(如果需要API密钥，请在这里填写)
	db.InitBinance("", "")

	// 创建Hertz服务器
	h := server.Default(
		server.WithHostPorts(":8080"),
		server.WithMaxRequestBodySize(4*1024*1024), // 4MB
	)

	// 注册路由
	router.Register(h)

	// 打印启动信息
	fmt.Println("\n╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           Trade Strategy API Server                          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println("\n服务已启动，监听端口: 8080")
	fmt.Println("\n可用接口:")
	fmt.Println("  - GET  http://localhost:8080/")
	fmt.Println("  - GET  http://localhost:8080/health")
	fmt.Println("  - GET  http://localhost:8080/api/v1/strategy/analyze")
	fmt.Println("  - POST http://localhost:8080/api/v1/strategy/analyze")
	fmt.Println("\n示例请求:")
	fmt.Println("  # 策略一 - 历史同期涨跌分析")
	fmt.Println("  curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d&date=2024-10-30'")
	fmt.Println("\n  # 策略二 - 小时级别涨跌分析")
	fmt.Println("  curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h&hour=14'")
	fmt.Println("\n按 Ctrl+C 停止服务")
	fmt.Println("═══════════════════════════════════════════════════════════════\n")

	// 启动服务器
	h.Spin()
}
