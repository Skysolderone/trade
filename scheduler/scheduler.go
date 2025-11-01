package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"trade/kline"
	"trade/model"
	"trade/strategy"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron   *cron.Cron
	config *model.Config
}

// NewScheduler 创建新的调度器
func NewScheduler(config *model.Config) *Scheduler {
	// 创建带秒级精度的cron调度器
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:   c,
		config: config,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	// 每天零点执行策略更新任务
	// cron表达式: "0 0 0 * * *" 表示每天的00:00:00执行
	_, err := s.cron.AddFunc("0 0 0 * * *", s.runDailyUpdate)
	if err != nil {
		log.Fatalf("添加定时任务失败: %v", err)
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              定时任务调度器已启动                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("⏰ 每日更新时间: 00:00:00\n")
	fmt.Printf("📊 监控交易对数量: %d\n", len(s.config.Symbols))
	fmt.Printf("🕐 当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	nextRun := s.getNextRunTime()
	fmt.Printf("⏭️  下次执行时间: %s\n", nextRun.Format("2006-01-02 15:04:05"))
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()

	// 启动调度器
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	fmt.Println("\n正在停止定时任务调度器...")
	s.cron.Stop()
	fmt.Println("定时任务调度器已停止")
}

// RunNow 立即执行一次更新任务（用于测试）
func (s *Scheduler) RunNow() {
	fmt.Println("\n手动触发策略更新任务...")
	s.runDailyUpdate()
}

// runDailyUpdate 执行每日更新任务
func (s *Scheduler) runDailyUpdate() {
	startTime := time.Now()

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              开始执行每日策略更新任务                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("🕐 执行时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()

	// 1. 更新K线数据
	fmt.Printf("开始更新 %d 个交易对的K线数据...\n", len(s.config.Symbols))
	for i, symbolConfig := range s.config.Symbols {
		fmt.Printf("\n[%d/%d] 处理交易对: %s\n", i+1, len(s.config.Symbols), symbolConfig.Symbol)

		for _, interval := range symbolConfig.Intervals {
			fmt.Printf("  - 更新时间周期: %s\n", interval)
			kline.UpdateKline(symbolConfig.Symbol, interval)
		}
	}
	fmt.Println("\n✅ K线数据更新完成")

	// 2. 运行策略一
	fmt.Println("\n========== 开始运行策略一 ==========")
	strategy.Strategy1(s.config)

	// 3. 运行策略二
	fmt.Println("\n========== 开始运行策略二 ==========")
	strategy.Strategy2(s.config)

	// 计算耗时
	duration := time.Since(startTime)

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              每日策略更新任务执行完成                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("⏱️  总耗时: %s\n", duration)
	fmt.Printf("🕐 完成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	nextRun := s.getNextRunTime()
	fmt.Printf("⏭️  下次执行时间: %s\n", nextRun.Format("2006-01-02 15:04:05"))
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()
}

// getNextRunTime 获取下次执行时间
func (s *Scheduler) getNextRunTime() time.Time {
	now := time.Now()

	// 计算明天零点
	tomorrow := now.AddDate(0, 0, 1)
	nextRun := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		0, 0, 0, 0, now.Location())

	// 如果当前时间还没到今天的零点，则下次执行时间是今天零点
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location())
	if now.Before(todayMidnight) {
		return todayMidnight
	}

	return nextRun
}

// GetSchedulerInfo 获取调度器信息
func (s *Scheduler) GetSchedulerInfo() {
	fmt.Println("\n定时任务调度器信息：")
	fmt.Printf("  当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("  下次执行时间: %s\n", s.getNextRunTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("  监控交易对数量: %d\n", len(s.config.Symbols))

	fmt.Println("\n  交易对列表:")
	for i, symbolConfig := range s.config.Symbols {
		fmt.Printf("    %d. %s (时间周期: %v)\n", i+1, symbolConfig.Symbol, symbolConfig.Intervals)
	}
}
