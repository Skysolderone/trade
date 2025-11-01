package strategy

import (
	"fmt"
	"sort"
	"time"

	"trade/db"
	"trade/model"
)

// KlineRecord 单条K线记录
type KlineRecord struct {
	Year       string  // 年份
	OpenPrice  float64 // 开盘价
	ClosePrice float64 // 收盘价
	PriceDiff  float64 // 价差 (收盘价 - 开盘价)
	IsUp       bool    // 是否上涨
}

// DayStats 每个日期的统计数据
type DayStats struct {
	Day        string        // 日期格式：MM-DD
	Month      int           // 月份
	TotalCount int           // 总样本数
	UpCount    int           // 上涨次数
	DownCount  int           // 下跌次数
	FlatCount  int           // 平盘次数
	UpRate     float64       // 上涨概率
	Records    []KlineRecord // 每年的K线记录
}

// Strategy1 根据配置文件分析所有交易对
func Strategy1(config *model.Config) {
	if config == nil || len(config.Symbols) == 0 {
		fmt.Println("⚠️  配置文件为空，无法执行策略分析")
		return
	}

	// 获取今天的时间
	now := time.Now()
	month := now.Month()
	day := now.Day()

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║          策略一：历史同期涨跌分析（跨月对比）                  ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("当前日期: %02d月%02d日\n", int(month), day)
	fmt.Printf("将分析 %d 个交易对的历史数据\n\n", len(config.Symbols))

	// 遍历配置文件中的所有交易对
	for i, symbolConfig := range config.Symbols {
		fmt.Printf("\n")
		fmt.Printf("═══════════════════════════════════════════════════════════════\n")
		fmt.Printf("  交易对 [%d/%d]: %s\n", i+1, len(config.Symbols), symbolConfig.Symbol)
		fmt.Printf("═══════════════════════════════════════════════════════════════\n")

		// 遍历该交易对的所有时间周期
		for _, interval := range symbolConfig.Intervals {
			analyzeSymbolInterval(symbolConfig.Symbol, interval, int(month), day)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    所有策略分析完成                             ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

// analyzeSymbolInterval 分析单个交易对的单个时间周期
func analyzeSymbolInterval(symbol, interval string, month, day int) {
	fmt.Printf("\n【时间周期: %s】\n", interval)

	// 1. 分析当前月当前日（例如：10-30）
	currentDayStats := analyzeSingleDay(symbol, interval, month, day)

	// 2. 分析其他月相同日期（例如：01-30, 02-30, ..., 12-30）- 跨月对比
	allMonthStats := analyzeAllMonthsSameDay(symbol, interval, day)

	// 3. 分析所有年份同一日期（例如：2018-10-30, 2019-10-30...）- 跨年对比
	allYearStats := analyzeAllYearsSameDate(symbol, interval, month, day)

	// 4. 保存结果到数据库
	saveStrategy1Result(symbol, interval, month, day, currentDayStats, allMonthStats, allYearStats)

	// 5. 输出对比结果
	printComparisonResults(currentDayStats, allMonthStats, allYearStats, month, day)
}

// analyzeSingleDay 分析单个日期的统计数据
func analyzeSingleDay(symbol, interval string, month, day int) *DayStats {
	dateStr := fmt.Sprintf("%02d-%02d", month, day)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND day = ?", symbol, interval, dateStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return &DayStats{
			Day:        dateStr,
			Month:      month,
			TotalCount: 0,
			Records:    []KlineRecord{},
		}
	}

	return calculateStats(dateStr, month, klines)
}

// analyzeAllMonthsSameDay 分析所有月份相同日期的数据（跨月对比）
func analyzeAllMonthsSameDay(symbol, interval string, day int) []*DayStats {
	stats := make([]*DayStats, 0, 12)

	for month := 1; month <= 12; month++ {
		// 检查该月是否有这一天（例如 2月没有30日）
		if !isValidDate(month, day) {
			continue
		}

		stat := analyzeSingleDay(symbol, interval, month, day)
		if stat.TotalCount > 0 {
			stats = append(stats, stat)
		}
	}

	return stats
}

// analyzeAllYearsSameDate 分析所有年份同一日期的数据（跨年对比）
func analyzeAllYearsSameDate(symbol, interval string, month, day int) []KlineRecord {
	dateStr := fmt.Sprintf("%02d-%02d", month, day)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND day = ?", symbol, interval, dateStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return []KlineRecord{}
	}

	records := make([]KlineRecord, 0, len(klines))
	for _, kline := range klines {
		priceDiff := kline.Close - kline.Open
		isUp := kline.Close > kline.Open

		record := KlineRecord{
			Year:       kline.Date,
			OpenPrice:  kline.Open,
			ClosePrice: kline.Close,
			PriceDiff:  priceDiff,
			IsUp:       isUp,
		}
		records = append(records, record)
	}

	// 按年份排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Year < records[j].Year
	})

	return records
}

// calculateStats 计算统计数据
func calculateStats(dateStr string, month int, klines []model.Kline) *DayStats {
	stats := &DayStats{
		Day:        dateStr,
		Month:      month,
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
	stats.UpRate = float64(stats.UpCount) / float64(stats.TotalCount) * 100

	return stats
}

// printComparisonResults 打印对比结果
func printComparisonResults(currentStats *DayStats, allMonthStats []*DayStats, allYearRecords []KlineRecord, currentMonth, currentDay int) {
	// 1. 跨年对比：所有年份同一日期（例如：2018-10-30, 2019-10-30...）
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📅 跨年对比：历年 %02d月%02d日 的数据\n", currentMonth, currentDay)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if len(allYearRecords) == 0 {
		fmt.Printf("⚠️  没有找到 %02d月%02d日 的历史数据\n\n", currentMonth, currentDay)
	} else {
		printYearlyRecords(allYearRecords, currentMonth, currentDay)
	}

	// 2. 跨月对比：所有月份相同日期（例如：01-30, 02-30...）
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 跨月对比：所有月份的 %02d号（详细数据）\n", currentDay)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if len(allMonthStats) == 0 {
		fmt.Printf("⚠️  没有找到任何数据\n")
		return
	}

	// 按月份排序（1月到12月）
	sort.Slice(allMonthStats, func(i, j int) bool {
		return allMonthStats[i].Month < allMonthStats[j].Month
	})

	// 打印每个月份的详细记录
	for i, stat := range allMonthStats {
		marker := ""
		if stat.Month == currentMonth {
			marker = " 👉 当前月份"
		}

		fmt.Printf("\n【%02d月%02d日】%s\n", stat.Month, currentDay, marker)
		fmt.Printf("样本数: %d 条 | 上涨率: %.2f%% (%d涨/%d跌/%d平)\n",
			stat.TotalCount, stat.UpRate, stat.UpCount, stat.DownCount, stat.FlatCount)

		// 打印该月份各年份的详细记录
		if len(stat.Records) > 0 {
			fmt.Printf("\n%-8s %-15s %-15s %-15s %-8s\n",
				"年份", "开盘价", "收盘价", "价差", "涨跌")
			fmt.Printf("%-8s %-15s %-15s %-15s %-8s\n",
				"────", "─────────────", "─────────────", "─────────────", "──────")

			for _, record := range stat.Records {
				direction := "📉 跌"
				if record.IsUp {
					direction = "📈 涨"
				} else if record.PriceDiff == 0 {
					direction = "➡️ 平"
				}

				fmt.Printf("%-8s %15.2f %15.2f %+15.2f %s\n",
					record.Year,
					record.OpenPrice,
					record.ClosePrice,
					record.PriceDiff,
					direction,
				)
			}
		}

		// 如果不是最后一个，添加分隔线
		if i < len(allMonthStats)-1 {
			fmt.Printf("\n%s\n", "─────────────────────────────────────────────────────────────")
		}
	}

	// 3. 统计摘要
	printMonthComparisonSummary(allMonthStats, currentMonth, currentDay)
}

// printYearlyRecords 打印跨年对比的详细记录
func printYearlyRecords(records []KlineRecord, month, day int) {
	if len(records) == 0 {
		return
	}

	// 统计涨跌次数
	upCount := 0
	downCount := 0
	flatCount := 0
	for _, record := range records {
		if record.IsUp {
			upCount++
		} else if record.PriceDiff < 0 {
			downCount++
		} else {
			flatCount++
		}
	}

	upRate := float64(upCount) / float64(len(records)) * 100

	fmt.Printf("样本数量: %d 条\n", len(records))
	if len(records) < 5 {
		fmt.Printf("⚠️  样本量不足（少于5条），统计结果可能不可靠\n\n")
	}

	fmt.Printf("\n基础统计：\n")
	fmt.Printf("  上涨次数: %d (%.2f%%)\n", upCount, upRate)
	fmt.Printf("  下跌次数: %d (%.2f%%)\n", downCount, float64(downCount)/float64(len(records))*100)
	fmt.Printf("  平盘次数: %d (%.2f%%)\n\n", flatCount, float64(flatCount)/float64(len(records))*100)

	fmt.Printf("历年 %02d月%02d日 记录：\n", month, day)
	fmt.Printf("%-8s %-15s %-15s %-15s %-8s\n",
		"年份", "开盘价", "收盘价", "价差", "涨跌")
	fmt.Printf("%-8s %-15s %-15s %-15s %-8s\n",
		"────", "─────────────", "─────────────", "─────────────", "──────")

	for _, record := range records {
		direction := "📉 跌"
		if record.IsUp {
			direction = "📈 涨"
		} else if record.PriceDiff == 0 {
			direction = "➡️ 平"
		}

		fmt.Printf("%-8s %15.2f %15.2f %+15.2f %s\n",
			record.Year,
			record.OpenPrice,
			record.ClosePrice,
			record.PriceDiff,
			direction,
		)
	}
}

// printDetailedRecords 打印详细的K线记录
func printDetailedRecords(stats *DayStats) {
	fmt.Printf("样本数量: %d 条\n", stats.TotalCount)

	if stats.TotalCount < 5 {
		fmt.Printf("⚠️  样本量不足（少于5条），统计结果可能不可靠\n\n")
	}

	fmt.Printf("\n基础统计：\n")
	fmt.Printf("  上涨次数: %d (%.2f%%)\n", stats.UpCount, stats.UpRate)
	fmt.Printf("  下跌次数: %d (%.2f%%)\n", stats.DownCount, float64(stats.DownCount)/float64(stats.TotalCount)*100)
	fmt.Printf("  平盘次数: %d (%.2f%%)\n\n", stats.FlatCount, float64(stats.FlatCount)/float64(stats.TotalCount)*100)

	// 按年份排序
	sort.Slice(stats.Records, func(i, j int) bool {
		return stats.Records[i].Year < stats.Records[j].Year
	})

	fmt.Printf("历史记录：\n")
	fmt.Printf("%-8s %-15s %-15s %-15s %-8s\n",
		"年份", "开盘价", "收盘价", "价差", "涨跌")
	fmt.Printf("%-8s %-15s %-15s %-15s %-8s\n",
		"────", "─────────────", "─────────────", "─────────────", "──────")

	for _, record := range stats.Records {
		direction := "📉 跌"
		if record.IsUp {
			direction = "📈 涨"
		} else if record.PriceDiff == 0 {
			direction = "➡️ 平"
		}

		fmt.Printf("%-8s %15.2f %15.2f %+15.2f %s\n",
			record.Year,
			record.OpenPrice,
			record.ClosePrice,
			record.PriceDiff,
			direction,
		)
	}
}

// printMonthComparisonSummary 打印跨月对比统计摘要
func printMonthComparisonSummary(allStats []*DayStats, currentMonth, currentDay int) {
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("💡 跨月对比摘要\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// 找出上涨率最高和最低的月份
	var best, worst *DayStats
	for _, stat := range allStats {
		if best == nil || stat.UpRate > best.UpRate {
			best = stat
		}
		if worst == nil || stat.UpRate < worst.UpRate {
			worst = stat
		}
	}

	if best != nil && worst != nil {
		fmt.Printf("📈 最佳月份: %02d月%02d日 - 上涨率 %.2f%% (%d涨/%d跌, 样本%d条)\n",
			best.Month, currentDay, best.UpRate, best.UpCount, best.DownCount, best.TotalCount)
		fmt.Printf("📉 最差月份: %02d月%02d日 - 上涨率 %.2f%% (%d涨/%d跌, 样本%d条)\n\n",
			worst.Month, currentDay, worst.UpRate, worst.UpCount, worst.DownCount, worst.TotalCount)
	}

	// 找出当前月份并显示信息
	var currentStat *DayStats
	for _, stat := range allStats {
		if stat.Month == currentMonth {
			currentStat = stat
			break
		}
	}

	if currentStat != nil {
		// 计算排名（按上涨率）
		rank := 1
		for _, stat := range allStats {
			if stat.UpRate > currentStat.UpRate {
				rank++
			}
		}

		fmt.Printf("🎯 当前月份 (%02d月): 上涨率 %.2f%%, 排名 %d/%d\n",
			currentMonth, currentStat.UpRate, rank, len(allStats))

		if rank <= len(allStats)/3 {
			fmt.Printf("✅ 当前月份表现优秀，历史上涨概率较高\n")
		} else if rank >= len(allStats)*2/3 {
			fmt.Printf("⚠️  当前月份表现较差，建议谨慎操作\n")
		} else {
			fmt.Printf("ℹ️  当前月份表现中等\n")
		}
	}

	fmt.Printf("\n⚠️  风险提示：历史数据不代表未来表现，请结合其他技术指标综合判断！\n")
	fmt.Printf("════════════════════════════════════════════════════════════════\n\n")
}

// isValidDate 检查日期是否有效
func isValidDate(month, day int) bool {
	daysInMonth := []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month < 1 || month > 12 {
		return false
	}
	return day >= 1 && day <= daysInMonth[month]
}

// saveStrategy1Result 保存策略一结果到数据库
func saveStrategy1Result(symbol, interval string, month, day int, currentStats *DayStats, allMonthStats []*DayStats, allYearRecords []KlineRecord) {
	analyzeDay := fmt.Sprintf("%02d-%02d", month, day)

	// 找出最佳和最差月份
	var bestMonth, worstMonth int
	var bestUpRate, worstUpRate float64
	if len(allMonthStats) > 0 {
		best := allMonthStats[0]
		worst := allMonthStats[0]
		for _, stat := range allMonthStats {
			if stat.UpRate > best.UpRate {
				best = stat
			}
			if stat.UpRate < worst.UpRate {
				worst = stat
			}
		}
		bestMonth = best.Month
		bestUpRate = best.UpRate
		worstMonth = worst.Month
		worstUpRate = worst.UpRate
	}

	// 创建或更新策略一结果
	result := &model.Strategy1Result{
		Symbol:      symbol,
		Interval:    interval,
		AnalyzeDay:  analyzeDay,
		Month:       month,
		Day:         day,
		TotalCount:  currentStats.TotalCount,
		UpCount:     currentStats.UpCount,
		DownCount:   currentStats.DownCount,
		FlatCount:   currentStats.FlatCount,
		UpRate:      currentStats.UpRate,
		BestMonth:   bestMonth,
		BestUpRate:  bestUpRate,
		WorstMonth:  worstMonth,
		WorstUpRate: worstUpRate,
	}

	// 使用upsert保存结果
	err := db.Pog.Where("symbol = ? AND interval = ? AND analyze_day = ?", symbol, interval, analyzeDay).
		Assign(result).
		FirstOrCreate(result).Error

	if err != nil {
		fmt.Printf("⚠️ 保存策略一结果失败: %v\n", err)
		return
	}

	// 删除旧的详细记录
	db.Pog.Where("result_id = ?", result.ID).Delete(&model.Strategy1DetailRecord{})

	// 保存详细记录
	for _, record := range allYearRecords {
		detailRecord := &model.Strategy1DetailRecord{
			ResultID:   result.ID,
			Year:       record.Year,
			OpenPrice:  record.OpenPrice,
			ClosePrice: record.ClosePrice,
			PriceDiff:  record.PriceDiff,
			IsUp:       record.IsUp,
		}
		if err := db.Pog.Create(detailRecord).Error; err != nil {
			fmt.Printf("⚠️ 保存详细记录失败: %v\n", err)
		}
	}
}
