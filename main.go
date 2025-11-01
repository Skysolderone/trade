package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"trade/db"
	"trade/kline"
	"trade/model"
	"trade/scheduler"
	"trade/strategy"
	"trade/utils"
)

func main() {
	// 命令行参数
	mode := flag.String("mode", "once", "运行模式: once(单次运行) 或 daemon(定时任务)")
	runNow := flag.Bool("now", false, "daemon模式下是否立即执行一次")
	flag.Parse()

	// 初始化币安客户端(如果需要API密钥,请在这里填写)
	db.InitBinance("", "")

	// 初始化数据库连接
	db.InitPostgreSql()

	// 加载配置文件
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	switch *mode {
	case "daemon":
		// 定时任务模式
		runDaemonMode(config, *runNow)

	case "once":
		// 单次运行模式
		runOnceMode(config)

	default:
		fmt.Printf("未知的运行模式: %s\n", *mode)
		fmt.Println("支持的模式:")
		fmt.Println("  once   - 单次运行（默认）")
		fmt.Println("  daemon - 定时任务模式，每天零点自动执行")
		os.Exit(1)
	}
}

// runOnceMode 单次运行模式
func runOnceMode(config *model.Config) {
	fmt.Printf("开始更新 %d 个交易对的K线数据...\n", len(config.Symbols))

	// 遍历配置文件中的所有交易对和时间区间
	for _, symbolConfig := range config.Symbols {
		fmt.Printf("\n========== 处理交易对: %s ==========\n", symbolConfig.Symbol)

		for _, interval := range symbolConfig.Intervals {
			fmt.Printf("\n--- 时间区间: %s ---\n", interval)
			// 更新K线数据(从数据库最新记录开始更新到昨天)
			kline.UpdateKline(symbolConfig.Symbol, interval)
		}
	}

	fmt.Println("\n========== 所有数据更新完成 ==========")

	// 运行策略一: 分析历史同期涨跌概率（传入配置文件）
	fmt.Println("\n========== 开始运行策略一 ==========")
	strategy.Strategy1(config)

	// 运行策略二: 小时级别涨跌分析
	fmt.Println("\n========== 开始运行策略二 ==========")
	strategy.Strategy2(config)
}

// runDaemonMode 定时任务模式
func runDaemonMode(config *model.Config, runNow bool) {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          币安交易策略定时任务系统                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 创建调度器
	s := scheduler.NewScheduler(config)

	// 如果指定了立即运行，先执行一次
	if runNow {
		fmt.Println("⚡ 立即执行模式，先运行一次更新任务...")
		s.RunNow()
	}

	// 启动定时任务
	s.Start()

	// 显示调度器信息
	s.GetSchedulerInfo()

	fmt.Println("\n💡 提示:")
	fmt.Println("  - 程序将在每天 00:00:00 自动执行策略更新")
	fmt.Println("  - 按 Ctrl+C 退出程序")
	fmt.Println()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	s.Stop()
	fmt.Println("程序已退出")
}
