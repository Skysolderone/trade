package strategy

import (
	"fmt"
	"sort"
	"time"

	"trade/db"
	"trade/model"
)

// HourStats 每个小时的统计数据
type HourStats struct {
	Hour       int           // 小时 (0-23)
	TotalCount int           // 总样本数
	UpCount    int           // 上涨次数
	DownCount  int           // 下跌次数
	FlatCount  int           // 平盘次数
	UpRate     float64       // 上涨概率
	Records    []KlineRecord // K线记录
}

// Strategy2 小时级别分析策略
func Strategy2(config *model.Config) {
	if config == nil || len(config.Symbols) == 0 {
		fmt.Println("⚠️  配置文件为空，无法执行策略分析")
		return
	}

	// 获取当前时间
	now := time.Now()
	currentHour := now.Hour()

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║          策略二：小时级别涨跌分析（日内时段对比）              ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("当前时间: %02d:00\n", currentHour)
	fmt.Printf("将分析 %d 个交易对的小时级别数据\n\n", len(config.Symbols))

	// 遍历配置文件中的所有交易对
	for i, symbolConfig := range config.Symbols {
		fmt.Printf("\n")
		fmt.Printf("═══════════════════════════════════════════════════════════════\n")
		fmt.Printf("  交易对 [%d/%d]: %s\n", i+1, len(config.Symbols), symbolConfig.Symbol)
		fmt.Printf("═══════════════════════════════════════════════════════════════\n")

		// 只分析小时级别的时间周期
		hourlyIntervals := []string{"8h", "4h", "2h", "1h"}
		for _, interval := range symbolConfig.Intervals {
			// 检查是否是小时级别
			isHourly := false
			for _, hourInterval := range hourlyIntervals {
				if interval == hourInterval {
					isHourly = true
					break
				}
			}

			if isHourly {
				analyzeHourlyPattern(symbolConfig.Symbol, interval, currentHour)
			}
		}
	}

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    小时级别分析完成                             ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

// analyzeHourlyPattern 分析单个交易对的小时规律
func analyzeHourlyPattern(symbol, interval string, currentHour int) {
	fmt.Printf("\n【时间周期: %s】\n", interval)

	// 1. 分析当前小时的历史表现
	currentHourStats := analyzeSpecificHour(symbol, interval, currentHour)

	// 2. 分析24小时的整体表现
	allHourStats := analyzeAll24Hours(symbol, interval)

	// 3. 输出分析结果
	printHourlyAnalysis(currentHourStats, allHourStats, currentHour, interval)
}

// analyzeSpecificHour 分析特定小时的统计数据
func analyzeSpecificHour(symbol, interval string, hour int) *HourStats {
	hourStr := fmt.Sprintf("%d", hour)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND hour = ?", symbol, interval, hourStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return &HourStats{
			Hour:       hour,
			TotalCount: 0,
			Records:    []KlineRecord{},
		}
	}

	return calculateHourStats(hour, klines)
}

// analyzeAll24Hours 分析所有24小时的数据
func analyzeAll24Hours(symbol, interval string) []*HourStats {
	stats := make([]*HourStats, 0, 24)

	for hour := 0; hour < 24; hour++ {
		stat := analyzeSpecificHour(symbol, interval, hour)
		if stat.TotalCount > 0 {
			stats = append(stats, stat)
		}
	}

	return stats
}

// calculateHourStats 计算小时统计数据
func calculateHourStats(hour int, klines []model.Kline) *HourStats {
	stats := &HourStats{
		Hour:       hour,
		TotalCount: len(klines),
		Records:    make([]KlineRecord, 0, len(klines)),
	}

	for _, kline := range klines {
		priceDiff := kline.Close - kline.Open
		isUp := kline.Close > kline.Open

		// 记录每条K线数据
		record := KlineRecord{
			Year:       kline.Date,
			OpenPrice:  kline.Open,
			ClosePrice: kline.Close,
			PriceDiff:  priceDiff,
			IsUp:       isUp,
		}
		stats.Records = append(stats.Records, record)

		// 统计涨跌次数
		if kline.Close > kline.Open {
			stats.UpCount++
		} else if kline.Close < kline.Open {
			stats.DownCount++
		} else {
			stats.FlatCount++
		}
	}

	// 计算上涨概率
	if stats.TotalCount > 0 {
		stats.UpRate = float64(stats.UpCount) / float64(stats.TotalCount) * 100
	}

	return stats
}

// printHourlyAnalysis 打印小时级别分析结果
func printHourlyAnalysis(currentHourStats *HourStats, allHourStats []*HourStats, currentHour int, interval string) {
	// 1. 打印当前小时的历史表现
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🕐 当前时段 %02d:00 的历史表现\n", currentHour)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if currentHourStats.TotalCount == 0 {
		fmt.Printf("⚠️  没有找到 %02d:00 的历史数据\n\n", currentHour)
	} else {
		fmt.Printf("样本数量: %d 条\n", currentHourStats.TotalCount)
		if currentHourStats.TotalCount < 10 {
			fmt.Printf("⚠️  样本量较少（少于10条），统计结果可能不可靠\n")
		}
		fmt.Printf("\n基础统计：\n")
		fmt.Printf("  上涨次数: %d (%.2f%%)\n", currentHourStats.UpCount, currentHourStats.UpRate)
		fmt.Printf("  下跌次数: %d (%.2f%%)\n", currentHourStats.DownCount,
			float64(currentHourStats.DownCount)/float64(currentHourStats.TotalCount)*100)
		fmt.Printf("  平盘次数: %d (%.2f%%)\n\n", currentHourStats.FlatCount,
			float64(currentHourStats.FlatCount)/float64(currentHourStats.TotalCount)*100)
	}

	// 2. 打印24小时对比分析
	if len(allHourStats) > 0 {
		print24HourComparison(allHourStats, currentHour, interval)
	}

	// 3. 打印最佳和最差时段
	printBestAndWorstHours(allHourStats, currentHour)
}

// print24HourComparison 打印24小时对比
func print24HourComparison(allStats []*HourStats, currentHour int, interval string) {
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 24小时涨跌概率分布 (%s周期)\n", interval)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// 按小时排序
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].Hour < allStats[j].Hour
	})

	fmt.Printf("%-6s %-10s %-12s %-8s %s\n", "时段", "样本数", "上涨率", "涨/跌", "图表")
	fmt.Printf("%-6s %-10s %-12s %-8s %s\n", "────", "────────", "──────────", "──────", "────────────────────")

	for _, stat := range allStats {
		marker := "  "
		if stat.Hour == currentHour {
			marker = "👉"
		}

		// 生成可视化柱状图
		barLength := int(stat.UpRate / 5) // 每5%一个字符
		bar := ""
		for i := 0; i < barLength; i++ {
			bar += "█"
		}

		fmt.Printf("%s%02d:00 %-10d %6.2f%%    %3d/%-3d %s\n",
			marker,
			stat.Hour,
			stat.TotalCount,
			stat.UpRate,
			stat.UpCount,
			stat.DownCount,
			bar,
		)
	}
}

// printBestAndWorstHours 打印最佳和最差时段
func printBestAndWorstHours(allStats []*HourStats, currentHour int) {
	if len(allStats) == 0 {
		return
	}

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💡 时段对比摘要\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// 找出上涨率最高和最低的时段（样本数>=5）
	var best, worst *HourStats
	for _, stat := range allStats {
		if stat.TotalCount < 5 {
			continue
		}
		if best == nil || stat.UpRate > best.UpRate {
			best = stat
		}
		if worst == nil || stat.UpRate < worst.UpRate {
			worst = stat
		}
	}

	if best != nil && worst != nil {
		fmt.Printf("📈 最佳时段: %02d:00 - 上涨率 %.2f%% (%d涨/%d跌, 样本%d条)\n",
			best.Hour, best.UpRate, best.UpCount, best.DownCount, best.TotalCount)
		fmt.Printf("📉 最差时段: %02d:00 - 上涨率 %.2f%% (%d涨/%d跌, 样本%d条)\n\n",
			worst.Hour, worst.UpRate, worst.UpCount, worst.DownCount, worst.TotalCount)
	}

	// 找出当前时段并显示信息
	var currentStat *HourStats
	for _, stat := range allStats {
		if stat.Hour == currentHour {
			currentStat = stat
			break
		}
	}

	if currentStat != nil && currentStat.TotalCount >= 5 {
		// 计算排名（按上涨率）
		rank := 1
		validCount := 0
		for _, stat := range allStats {
			if stat.TotalCount >= 5 {
				validCount++
				if stat.UpRate > currentStat.UpRate {
					rank++
				}
			}
		}

		fmt.Printf("🎯 当前时段 (%02d:00): 上涨率 %.2f%%, 排名 %d/%d\n",
			currentHour, currentStat.UpRate, rank, validCount)

		if rank <= validCount/3 {
			fmt.Printf("✅ 当前时段表现优秀，历史上涨概率较高\n")
		} else if rank >= validCount*2/3 {
			fmt.Printf("⚠️  当前时段表现较差，建议谨慎操作\n")
		} else {
			fmt.Printf("ℹ️  当前时段表现中等\n")
		}
	}

	// 时段建议
	fmt.Printf("\n💡 交易时段建议：\n")

	// 找出高胜率时段（上涨率>60%且样本>=10）
	highWinRate := make([]*HourStats, 0)
	for _, stat := range allStats {
		if stat.TotalCount >= 10 && stat.UpRate >= 60 {
			highWinRate = append(highWinRate, stat)
		}
	}

	if len(highWinRate) > 0 {
		sort.Slice(highWinRate, func(i, j int) bool {
			return highWinRate[i].UpRate > highWinRate[j].UpRate
		})

		fmt.Printf("  高胜率时段（上涨率≥60%%）:\n")
		for i, stat := range highWinRate {
			if i >= 3 { // 最多显示3个
				break
			}
			fmt.Printf("    • %02d:00 (%.2f%%, 样本%d条)\n", stat.Hour, stat.UpRate, stat.TotalCount)
		}
	} else {
		fmt.Printf("  暂无高胜率时段（上涨率≥60%%且样本≥10）\n")
	}

	fmt.Printf("\n⚠️  风险提示：历史数据不代表未来表现，请结合实时行情和其他技术指标综合判断！\n")
	fmt.Printf("════════════════════════════════════════════════════════════════\n\n")
}
