package main

import (
	"flag"
	"fmt"

	"trade/db"
	"trade/web"
)

func main() {
	// 定义命令行参数
	port := flag.Int("port", 8080, "Web服务器端口")
	flag.Parse()

	// 初始化数据库连接
	db.InitPostgreSql()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          交易策略分析 - Web服务器                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("\n")
	fmt.Printf("🌐 访问地址: http://localhost:%d\n", *port)
	fmt.Printf("📊 查看策略一: http://localhost:%d/api/strategy1\n", *port)
	fmt.Printf("🕐 查看策略二: http://localhost:%d/api/strategy2\n", *port)
	fmt.Printf("\n")
	fmt.Println("按 Ctrl+C 停止服务器")
	fmt.Println("════════════════════════════════════════════════════════════════")

	// 启动Web服务器
	web.StartServer(*port)
}
