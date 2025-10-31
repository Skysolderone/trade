package handler

import (
	"context"
	"fmt"
	"sort"
	"time"

	"trade/api/response"
	"trade/db"
	"trade/model"
	"trade/strategy"

	"github.com/cloudwego/hertz/pkg/app"
)

// AnalyzeRequest 分析请求参数
type AnalyzeRequest struct {
	StrategyType string `json:"strategy_type" query:"strategy_type"` // 策略类型: strategy_1, strategy_2
	Symbol       string `json:"symbol" query:"symbol"`                // 交易对
	Interval     string `json:"interval" query:"interval"`            // K线周期
	Date         string `json:"date,omitempty" query:"date"`          // 日期(策略一使用，格式：2024-10-30)
	Hour         *int   `json:"hour,omitempty" query:"hour"`          // 小时(策略二使用，0-23)
}

// AnalyzeStrategy 策略分析接口
func AnalyzeStrategy(ctx context.Context, c *app.RequestContext) {
	var req AnalyzeRequest

	// 绑定参数
	if err := c.Bind(&req); err != nil {
		response.ParamError(c, fmt.Sprintf("参数错误：%v", err))
		return
	}

	// 参数校验
	if req.Symbol == "" {
		response.ParamError(c, "参数错误：缺少symbol参数")
		return
	}
	if req.Interval == "" {
		response.ParamError(c, "参数错误：缺少interval参数")
		return
	}
	if req.StrategyType == "" {
		response.ParamError(c, "参数错误：缺少strategy_type参数")
		return
	}

	// 校验interval参数
	validIntervals := []string{"1m", "5m", "15m", "30m", "1h", "2h", "4h", "8h", "1d", "1w"}
	isValidInterval := false
	for _, interval := range validIntervals {
		if req.Interval == interval {
			isValidInterval = true
			break
		}
	}
	if !isValidInterval {
		response.ParamError(c, "参数错误：interval只支持1m,5m,15m,30m,1h,2h,4h,8h,1d,1w")
		return
	}

	// 根据策略类型调用不同的处理函数
	switch req.StrategyType {
	case "strategy_1":
		handleStrategy1(ctx, c, &req)
	case "strategy_2":
		handleStrategy2(ctx, c, &req)
	default:
		response.ParamError(c, "参数错误：strategy_type只支持strategy_1或strategy_2")
	}
}

// handleStrategy1 处理策略一
func handleStrategy1(ctx context.Context, c *app.RequestContext, req *AnalyzeRequest) {
	// 解析日期
	var targetDate time.Time
	var month, day int

	if req.Date != "" {
		var err error
		targetDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.ParamError(c, fmt.Sprintf("参数错误：日期格式错误，应为YYYY-MM-DD，例如：2024-10-30"))
			return
		}
		month = int(targetDate.Month())
		day = targetDate.Day()
	} else {
		// 默认使用今天
		now := time.Now()
		month = int(now.Month())
		day = now.Day()
		targetDate = now
	}

	// 检查日期是否有效
	if !isValidDate(month, day) {
		response.ParamError(c, fmt.Sprintf("参数错误：日期不存在(%d月没有%d日)", month, day))
		return
	}

	// 分析当前日期
	dateStr := fmt.Sprintf("%02d-%02d", month, day)
	currentDayStats := analyzeSingleDay(req.Symbol, req.Interval, month, day)

	// 如果没有数据
	if currentDayStats.TotalCount == 0 {
		response.DataNotFound(c, fmt.Sprintf("未找到%s在%s的历史数据", req.Symbol, dateStr))
		return
	}

	// 分析所有年份同一日期
	allYearRecords := analyzeAllYearsSameDate(req.Symbol, req.Interval, month, day)

	// 分析所有月份相同日期
	allMonthStats := analyzeAllMonthsSameDay(req.Symbol, req.Interval, day)

	// 构建响应数据
	resp := buildStrategy1Response(req, currentDayStats, allYearRecords, allMonthStats, month, day)

	// 检查样本量
	if currentDayStats.TotalCount < 5 {
		response.SampleTooLow(c, "样本量不足，统计结果可能不可靠", resp)
		return
	}

	response.Success(c, resp)
}

// handleStrategy2 处理策略二
func handleStrategy2(ctx context.Context, c *app.RequestContext, req *AnalyzeRequest) {
	// 解析小时
	var targetHour int
	if req.Hour != nil {
		targetHour = *req.Hour
		if targetHour < 0 || targetHour > 23 {
			response.ParamError(c, "参数错误：hour参数必须在0-23之间")
			return
		}
	} else {
		// 默认使用当前小时
		targetHour = time.Now().Hour()
	}

	// 分析当前小时
	currentHourStats := analyzeSpecificHour(req.Symbol, req.Interval, targetHour)

	// 如果没有数据
	if currentHourStats.TotalCount == 0 {
		response.DataNotFound(c, fmt.Sprintf("未找到%s在%02d:00的历史数据", req.Symbol, targetHour))
		return
	}

	// 分析24小时
	allHourStats := analyzeAll24Hours(req.Symbol, req.Interval)

	// 构建响应数据
	resp := buildStrategy2Response(req, currentHourStats, allHourStats, targetHour)

	// 检查样本量
	if currentHourStats.TotalCount < 10 {
		response.SampleTooLow(c, "样本量不足，统计结果可能不可靠", resp)
		return
	}

	response.Success(c, resp)
}

// analyzeSingleDay 分析单个日期的统计数据
func analyzeSingleDay(symbol, interval string, month, day int) *strategy.DayStats {
	dateStr := fmt.Sprintf("%02d-%02d", month, day)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND day = ?", symbol, interval, dateStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return &strategy.DayStats{
			Day:        dateStr,
			Month:      month,
			TotalCount: 0,
			Records:    []strategy.KlineRecord{},
		}
	}

	return calculateStats(dateStr, month, klines)
}

// analyzeAllYearsSameDate 分析所有年份同一日期的数据
func analyzeAllYearsSameDate(symbol, interval string, month, day int) []strategy.KlineRecord {
	dateStr := fmt.Sprintf("%02d-%02d", month, day)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND day = ?", symbol, interval, dateStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return []strategy.KlineRecord{}
	}

	records := make([]strategy.KlineRecord, 0, len(klines))
	for _, kline := range klines {
		priceDiff := kline.Close - kline.Open
		isUp := kline.Close > kline.Open

		record := strategy.KlineRecord{
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

// analyzeAllMonthsSameDay 分析所有月份相同日期的数据
func analyzeAllMonthsSameDay(symbol, interval string, day int) []*strategy.DayStats {
	stats := make([]*strategy.DayStats, 0, 12)

	for month := 1; month <= 12; month++ {
		// 检查该月是否有这一天
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

// analyzeSpecificHour 分析特定小时的统计数据
func analyzeSpecificHour(symbol, interval string, hour int) *strategy.HourStats {
	hourStr := fmt.Sprintf("%d", hour)

	var klines []model.Kline
	err := db.Pog.Where("symbol = ? AND interval = ? AND hour = ?", symbol, interval, hourStr).
		Order("open_time ASC").
		Find(&klines).Error

	if err != nil || len(klines) == 0 {
		return &strategy.HourStats{
			Hour:       hour,
			TotalCount: 0,
			Records:    []strategy.KlineRecord{},
		}
	}

	return calculateHourStats(hour, klines)
}

// analyzeAll24Hours 分析所有24小时的数据
func analyzeAll24Hours(symbol, interval string) []*strategy.HourStats {
	stats := make([]*strategy.HourStats, 0, 24)

	for hour := 0; hour < 24; hour++ {
		stat := analyzeSpecificHour(symbol, interval, hour)
		if stat.TotalCount > 0 {
			stats = append(stats, stat)
		}
	}

	return stats
}

// calculateStats 计算统计数据
func calculateStats(dateStr string, month int, klines []model.Kline) *strategy.DayStats {
	stats := &strategy.DayStats{
		Day:        dateStr,
		Month:      month,
		TotalCount: len(klines),
		Records:    make([]strategy.KlineRecord, 0, len(klines)),
	}

	for _, kline := range klines {
		priceDiff := kline.Close - kline.Open
		isUp := kline.Close > kline.Open

		// 记录每条K线数据
		record := strategy.KlineRecord{
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

// calculateHourStats 计算小时统计数据
func calculateHourStats(hour int, klines []model.Kline) *strategy.HourStats {
	stats := &strategy.HourStats{
		Hour:       hour,
		TotalCount: len(klines),
		Records:    make([]strategy.KlineRecord, 0, len(klines)),
	}

	for _, kline := range klines {
		priceDiff := kline.Close - kline.Open
		isUp := kline.Close > kline.Open

		// 记录每条K线数据
		record := strategy.KlineRecord{
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

// isValidDate 检查日期是否有效
func isValidDate(month, day int) bool {
	daysInMonth := []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month < 1 || month > 12 {
		return false
	}
	return day >= 1 && day <= daysInMonth[month]
}

// buildStrategy1Response 构建策略一响应
func buildStrategy1Response(req *AnalyzeRequest, currentStats *strategy.DayStats,
	allYearRecords []strategy.KlineRecord, allMonthStats []*strategy.DayStats,
	month, day int) *response.Strategy1Response {

	dateStr := fmt.Sprintf("%02d-%02d", month, day)
	analysisDate := fmt.Sprintf("2024-%02d-%02d", month, day)

	// 计算可靠性
	reliability, reliabilityNote := getReliability(currentStats.TotalCount)

	// 计算跨年分析
	crossYearAnalysis := buildCrossYearAnalysis(allYearRecords, month, day)

	// 计算跨月分析
	crossMonthAnalysis := buildCrossMonthAnalysis(allMonthStats, month, day)

	// 生成交易建议
	tradingRec := buildTradingRecommendation(currentStats, reliability)

	// 风险警告
	riskWarning := buildRiskWarning(currentStats.TotalCount)

	return &response.Strategy1Response{
		StrategyInfo: &response.StrategyInfo{
			StrategyType:   "strategy_1",
			StrategyName:   "历史同期涨跌分析",
			Description:    "基于历史K线数据，统计特定日期在历年的涨跌表现，通过跨年、跨月对比分析价格走势规律",
			AnalysisMethod: "从PostgreSQL数据库查询历史K线数据，按日期(day)字段分组统计涨跌次数和概率",
		},
		AnalysisTarget: &response.AnalysisTarget{
			Symbol:       req.Symbol,
			Interval:     req.Interval,
			AnalysisDate: analysisDate,
			TargetPeriod: dateStr,
		},
		DataStatistics: &response.DataStatistics{
			DataSource: "币安合约历史K线数据",
			DateRange: response.DateRange{
				StartDate: "2018-01-01",
				EndDate:   time.Now().Format("2006-01-02"),
			},
			TotalRecordsUsed: currentStats.TotalCount,
			QueryMethod:      "按symbol、interval和day字段查询数据库Kline表",
		},
		CurrentPeriodResult: &response.PeriodResult{
			PeriodLabel:     fmt.Sprintf("%d月%d日", month, day),
			SampleCount:     currentStats.TotalCount,
			UpCount:         currentStats.UpCount,
			DownCount:       currentStats.DownCount,
			FlatCount:       currentStats.FlatCount,
			UpRate:          currentStats.UpRate,
			DownRate:        float64(currentStats.DownCount) / float64(currentStats.TotalCount) * 100,
			Reliability:     reliability,
			ReliabilityNote: reliabilityNote,
		},
		CrossYearAnalysis:     crossYearAnalysis,
		CrossMonthAnalysis:    crossMonthAnalysis,
		TradingRecommendation: tradingRec,
		RiskWarning:           riskWarning,
	}
}

// buildStrategy2Response 构建策略二响应
func buildStrategy2Response(req *AnalyzeRequest, currentHourStats *strategy.HourStats,
	allHourStats []*strategy.HourStats, targetHour int) *response.Strategy2Response {

	analysisDatetime := time.Now().Format("2006-01-02 15:04:05")

	// 计算可靠性
	reliability, reliabilityNote := getHourReliability(currentHourStats.TotalCount)

	// 构建小时对比
	hourlyComparison := buildHourlyComparison(allHourStats, targetHour)

	// 高胜率时段
	highWinHours := buildHighWinHours(allHourStats)

	// 低胜率时段
	lowWinHours := buildLowWinHours(allHourStats)

	// 时区分析
	timeZoneAnalysis := buildTimeZoneAnalysis(allHourStats)

	// 交易建议
	tradingRec := buildHourTradingRecommendation(currentHourStats, allHourStats, targetHour, reliability)

	// 风险警告
	riskWarning := buildHourRiskWarning(currentHourStats.TotalCount)

	return &response.Strategy2Response{
		StrategyInfo: &response.StrategyInfo{
			StrategyType:   "strategy_2",
			StrategyName:   "小时级别涨跌分析",
			Description:    "分析特定小时的历史涨跌表现，通过24小时对比找出最佳交易时段",
			AnalysisMethod: "从PostgreSQL数据库查询历史K线数据，按小时(hour)字段分组统计不同时段的涨跌概率",
		},
		AnalysisTarget: &response.AnalysisTarget{
			Symbol:           req.Symbol,
			Interval:         req.Interval,
			AnalysisDatetime: analysisDatetime,
			TargetHour:       targetHour,
		},
		DataStatistics: &response.DataStatistics{
			DataSource: "币安合约历史K线数据",
			DateRange: response.DateRange{
				StartDate: "2018-01-01",
				EndDate:   time.Now().Format("2006-01-02"),
			},
			TotalRecordsUsed: currentHourStats.TotalCount,
			QueryMethod:      "按symbol、interval和hour字段查询数据库Kline表",
		},
		CurrentHourResult: &response.PeriodResult{
			HourLabel:       fmt.Sprintf("%02d:00", targetHour),
			SampleCount:     currentHourStats.TotalCount,
			UpCount:         currentHourStats.UpCount,
			DownCount:       currentHourStats.DownCount,
			FlatCount:       currentHourStats.FlatCount,
			UpRate:          currentHourStats.UpRate,
			DownRate:        float64(currentHourStats.DownCount) / float64(currentHourStats.TotalCount) * 100,
			Reliability:     reliability,
			ReliabilityNote: reliabilityNote,
		},
		HourlyComparison:      hourlyComparison,
		HighWinHours:          highWinHours,
		LowWinHours:           lowWinHours,
		TimeZoneAnalysis:      timeZoneAnalysis,
		TradingRecommendation: tradingRec,
		RiskWarning:           riskWarning,
	}
}

// getReliability 获取可靠性等级
func getReliability(sampleCount int) (string, string) {
	if sampleCount >= 10 {
		return "high", "样本数量充足，统计结果可靠性高"
	} else if sampleCount >= 5 {
		return "medium", "样本数量适中，统计结果具有一定参考价值"
	}
	return "low", fmt.Sprintf("样本数量过少(仅%d条)，统计结果不可靠", sampleCount)
}

// getHourReliability 获取小时可靠性等级
func getHourReliability(sampleCount int) (string, string) {
	if sampleCount >= 100 {
		return "high", "样本数量充足(>100条)，统计结果可靠性高"
	} else if sampleCount >= 10 {
		return "medium", "样本数量适中，统计结果具有一定参考价值"
	}
	return "low", fmt.Sprintf("样本数量较少(少于10条)，统计结果可能不可靠", sampleCount)
}

// buildCrossYearAnalysis 构建跨年分析
func buildCrossYearAnalysis(records []strategy.KlineRecord, month, day int) *response.CrossYearAnalysis {
	if len(records) == 0 {
		return nil
	}

	upCount := 0
	var bestRecord, worstRecord *strategy.KlineRecord

	for i := range records {
		if records[i].IsUp {
			upCount++
		}
		if bestRecord == nil || (records[i].IsUp && records[i].PriceDiff > bestRecord.PriceDiff) {
			bestRecord = &records[i]
		}
		if worstRecord == nil || (!records[i].IsUp && records[i].PriceDiff < worstRecord.PriceDiff) {
			worstRecord = &records[i]
		}
	}

	overallUpRate := float64(upCount) / float64(len(records)) * 100
	trend := "neutral"
	if overallUpRate > 55 {
		trend = "bullish"
	} else if overallUpRate < 45 {
		trend = "bearish"
	}

	var best, worst *response.Performance
	if bestRecord != nil {
		best = &response.Performance{
			Year:   bestRecord.Year,
			Result: "上涨",
			Rate:   100.0,
		}
	}
	if worstRecord != nil {
		worst = &response.Performance{
			Year:   worstRecord.Year,
			Result: "下跌",
			Rate:   0.0,
		}
	}

	return &response.CrossYearAnalysis{
		Title:            "跨年对比",
		Description:      fmt.Sprintf("分析历年同一日期(%02d-%02d)的涨跌情况", month, day),
		YearsAnalyzed:    len(records),
		OverallUpRate:    overallUpRate,
		Trend:            trend,
		BestPerformance:  best,
		WorstPerformance: worst,
	}
}

// buildCrossMonthAnalysis 构建跨月分析
func buildCrossMonthAnalysis(allStats []*strategy.DayStats, currentMonth, day int) *response.CrossMonthAnalysis {
	if len(allStats) == 0 {
		return nil
	}

	// 找出最佳和最差月份
	var best, worst *strategy.DayStats
	for _, stat := range allStats {
		if best == nil || stat.UpRate > best.UpRate {
			best = stat
		}
		if worst == nil || stat.UpRate < worst.UpRate {
			worst = stat
		}
	}

	// 当前月份排名
	var currentStat *strategy.DayStats
	rank := 1
	for _, stat := range allStats {
		if stat.Month == currentMonth {
			currentStat = stat
		}
		if currentStat != nil && stat.UpRate > currentStat.UpRate {
			rank++
		}
	}

	var bestPerf, worstPerf *response.Performance
	if best != nil {
		bestPerf = &response.Performance{
			Month:      best.Month,
			MonthLabel: fmt.Sprintf("%02d月%02d日", best.Month, day),
			UpRate:     best.UpRate,
			SampleCount: best.TotalCount,
			UpCount:    best.UpCount,
			DownCount:  best.DownCount,
		}
	}
	if worst != nil {
		worstPerf = &response.Performance{
			Month:      worst.Month,
			MonthLabel: fmt.Sprintf("%02d月%02d日", worst.Month, day),
			UpRate:     worst.UpRate,
			SampleCount: worst.TotalCount,
			UpCount:    worst.UpCount,
			DownCount:  worst.DownCount,
		}
	}

	var ranking *response.Ranking
	if currentStat != nil {
		perfLevel := "medium"
		perfNote := "当前月份表现中等"
		if rank <= len(allStats)/3 {
			perfLevel = "excellent"
			perfNote = "当前月份表现优秀，历史上涨概率较高"
		} else if rank >= len(allStats)*2/3 {
			perfLevel = "poor"
			perfNote = "当前月份表现较差，建议谨慎操作"
		}

		ranking = &response.Ranking{
			CurrentMonth:     currentMonth,
			Rank:             rank,
			TotalMonths:      len(allStats),
			PerformanceLevel: perfLevel,
			PerformanceNote:  perfNote,
		}
	}

	return &response.CrossMonthAnalysis{
		Title:               "跨月对比",
		Description:         fmt.Sprintf("对比所有月份的%d号，找出历史表现最好和最差的月份", day),
		MonthsAnalyzed:      len(allStats),
		BestMonth:           bestPerf,
		WorstMonth:          worstPerf,
		CurrentMonthRanking: ranking,
	}
}

// buildHourlyComparison 构建小时对比
func buildHourlyComparison(allStats []*strategy.HourStats, currentHour int) *response.HourlyComparison {
	if len(allStats) == 0 {
		return nil
	}

	// 计算平均上涨率
	totalUpRate := 0.0
	for _, stat := range allStats {
		totalUpRate += stat.UpRate
	}
	avgUpRate := totalUpRate / float64(len(allStats))

	// 找出最佳和最差时段(样本数>=5)
	var best, worst *strategy.HourStats
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

	// 当前小时排名
	var currentStat *strategy.HourStats
	rank := 1
	validCount := 0
	for _, stat := range allStats {
		if stat.Hour == currentHour {
			currentStat = stat
		}
		if stat.TotalCount >= 5 {
			validCount++
			if currentStat != nil && stat.UpRate > currentStat.UpRate {
				rank++
			}
		}
	}

	var bestPerf, worstPerf *response.Performance
	if best != nil {
		bestPerf = &response.Performance{
			Hour:            best.Hour,
			HourLabel:       fmt.Sprintf("%02d:00-15:00", best.Hour),
			UpRate:          best.UpRate,
			SampleCount:     best.TotalCount,
			UpCount:         best.UpCount,
			DownCount:       best.DownCount,
			PerformanceNote: "历史表现最佳时段",
		}
	}
	if worst != nil {
		worstPerf = &response.Performance{
			Hour:            worst.Hour,
			HourLabel:       fmt.Sprintf("%02d:00-15:00", worst.Hour),
			UpRate:          worst.UpRate,
			SampleCount:     worst.TotalCount,
			UpCount:         worst.UpCount,
			DownCount:       worst.DownCount,
			PerformanceNote: "历史表现最差时段",
		}
	}

	var ranking *response.Ranking
	if currentStat != nil && currentStat.TotalCount >= 5 {
		perfLevel := "medium"
		perfNote := "当前时段表现中等"
		if rank <= validCount/3 {
			perfLevel = "excellent"
			perfNote = "当前时段历史表现优秀，排名前10%"
		} else if rank >= validCount*2/3 {
			perfLevel = "poor"
			perfNote = "当前时段表现较差，建议谨慎操作"
		}

		ranking = &response.Ranking{
			Rank:             rank,
			TotalHours:       validCount,
			PerformanceLevel: perfLevel,
			PerformanceNote:  perfNote,
		}
	}

	return &response.HourlyComparison{
		Title:              "24小时对比",
		Description:        "对比24个小时时段，找出历史表现最好和最差的交易时段",
		HoursAnalyzed:      len(allStats),
		AverageUpRate:      avgUpRate,
		BestHour:           bestPerf,
		WorstHour:          worstPerf,
		CurrentHourRanking: ranking,
	}
}

// buildHighWinHours 构建高胜率时段
func buildHighWinHours(allStats []*strategy.HourStats) *response.HighWinHours {
	hours := make([]*response.Performance, 0)

	for _, stat := range allStats {
		if stat.TotalCount >= 10 && stat.UpRate >= 60 {
			hours = append(hours, &response.Performance{
				Hour:        stat.Hour,
				HourLabel:   fmt.Sprintf("%02d:00", stat.Hour),
				UpRate:      stat.UpRate,
				SampleCount: stat.TotalCount,
			})
		}
	}

	// 按上涨率排序
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].UpRate > hours[j].UpRate
	})

	// 最多返回3个
	if len(hours) > 3 {
		hours = hours[:3]
	}

	return &response.HighWinHours{
		Title:       "高胜率时段推荐",
		Description: "上涨率≥60%且样本数≥10的时段",
		Count:       len(hours),
		Hours:       hours,
	}
}

// buildLowWinHours 构建低胜率时段
func buildLowWinHours(allStats []*strategy.HourStats) *response.LowWinHours {
	hours := make([]*response.Performance, 0)

	for _, stat := range allStats {
		if stat.UpRate < 45 {
			hours = append(hours, &response.Performance{
				Hour:        stat.Hour,
				HourLabel:   fmt.Sprintf("%02d:00", stat.Hour),
				UpRate:      stat.UpRate,
				SampleCount: stat.TotalCount,
			})
		}
	}

	// 按上涨率排序(升序)
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].UpRate < hours[j].UpRate
	})

	// 最多返回2个
	if len(hours) > 2 {
		hours = hours[:2]
	}

	return &response.LowWinHours{
		Title:       "低胜率时段警示",
		Description: "上涨率<45%的时段，建议谨慎交易",
		Count:       len(hours),
		Hours:       hours,
	}
}

// buildTimeZoneAnalysis 构建时区分析
func buildTimeZoneAnalysis(allStats []*strategy.HourStats) *response.TimeZoneAnalysis {
	// 计算各时区平均上涨率
	asianUpRate := 0.0
	asianCount := 0
	europeanUpRate := 0.0
	europeanCount := 0
	americanUpRate := 0.0
	americanCount := 0

	for _, stat := range allStats {
		if stat.Hour >= 0 && stat.Hour < 8 {
			asianUpRate += stat.UpRate
			asianCount++
		} else if stat.Hour >= 8 && stat.Hour < 16 {
			europeanUpRate += stat.UpRate
			europeanCount++
		} else {
			americanUpRate += stat.UpRate
			americanCount++
		}
	}

	if asianCount > 0 {
		asianUpRate /= float64(asianCount)
	}
	if europeanCount > 0 {
		europeanUpRate /= float64(europeanCount)
	}
	if americanCount > 0 {
		americanUpRate /= float64(americanCount)
	}

	return &response.TimeZoneAnalysis{
		Title: "时区特征分析",
		AsianSession: &response.SessionInfo{
			Hours:          "00:00-08:00(UTC)",
			AverageUpRate:  asianUpRate,
			Characteristic: "波动较小，胜率偏低",
		},
		EuropeanSession: &response.SessionInfo{
			Hours:          "08:00-16:00(UTC)",
			AverageUpRate:  europeanUpRate,
			Characteristic: "活跃度上升，胜率较高",
		},
		AmericanSession: &response.SessionInfo{
			Hours:          "16:00-24:00(UTC)",
			AverageUpRate:  americanUpRate,
			Characteristic: "波动加大，胜率中等",
		},
	}
}

// buildTradingRecommendation 构建交易建议
func buildTradingRecommendation(stats *strategy.DayStats, reliability string) *response.TradingRecommendation {
	signal := "neutral"
	if stats.UpRate > 55 {
		signal = "bullish"
	} else if stats.UpRate < 45 {
		signal = "bearish"
	}

	confidenceLevel := reliability

	supportingFactors := []string{
		fmt.Sprintf("历史数据显示上涨概率为%.2f%%", stats.UpRate),
		fmt.Sprintf("样本数量为%d条", stats.TotalCount),
	}

	riskFactors := []string{}
	if stats.TotalCount < 10 {
		riskFactors = append(riskFactors, fmt.Sprintf("样本量不足%d条，可能存在统计偏差", stats.TotalCount))
	}
	riskFactors = append(riskFactors, "历史表现不代表未来走势")

	mainReason := fmt.Sprintf("历史数据显示上涨概率为%.2f%%", stats.UpRate)

	return &response.TradingRecommendation{
		Signal:            signal,
		ConfidenceLevel:   confidenceLevel,
		ConfidenceScore:   stats.UpRate,
		MainReason:        mainReason,
		SupportingFactors: supportingFactors,
		RiskFactors:       riskFactors,
	}
}

// buildHourTradingRecommendation 构建小时交易建议
func buildHourTradingRecommendation(currentStats *strategy.HourStats, allStats []*strategy.HourStats,
	targetHour int, reliability string) *response.TradingRecommendation {

	signal := "neutral"
	if currentStats.UpRate > 55 {
		signal = "bullish"
	} else if currentStats.UpRate < 45 {
		signal = "bearish"
	}

	// 找出最佳时段
	optimalHours := []string{}
	avoidHours := []string{}

	for _, stat := range allStats {
		if stat.TotalCount >= 10 && stat.UpRate >= 60 {
			optimalHours = append(optimalHours, fmt.Sprintf("%02d:00", stat.Hour))
		}
		if stat.UpRate < 45 {
			avoidHours = append(avoidHours, fmt.Sprintf("%02d:00", stat.Hour))
		}
	}

	// 计算排名
	rank := 1
	validCount := 0
	for _, stat := range allStats {
		if stat.TotalCount >= 5 {
			validCount++
			if stat.UpRate > currentStats.UpRate {
				rank++
			}
		}
	}

	supportingFactors := []string{
		fmt.Sprintf("样本数量充足(%d条)，统计结果可靠", currentStats.TotalCount),
		fmt.Sprintf("当前时段历史上涨概率为%.2f%%", currentStats.UpRate),
	}
	if rank <= validCount/3 {
		supportingFactors = append(supportingFactors, "当前时段属于高胜率时段")
	}

	riskFactors := []string{}
	if currentStats.TotalCount < 10 {
		riskFactors = append(riskFactors, "样本量较少，统计结果可能不稳定")
	}

	return &response.TradingRecommendation{
		Signal:              signal,
		ConfidenceLevel:     reliability,
		ConfidenceScore:     currentStats.UpRate,
		MainReason:          fmt.Sprintf("当前时段(%02d:00)历史上涨概率为%.2f%%，在24小时中排名第%d位", targetHour, currentStats.UpRate, rank),
		SupportingFactors:   supportingFactors,
		RiskFactors:         riskFactors,
		OptimalTradingHours: optimalHours,
		AvoidTradingHours:   avoidHours,
	}
}

// buildRiskWarning 构建风险警告
func buildRiskWarning(sampleCount int) *response.RiskWarning {
	level := "medium"
	warnings := []string{
		"历史数据不代表未来表现，仅供参考",
		"请结合实时行情、技术指标、资金管理等综合判断",
		"加密货币市场波动较大，请控制仓位和风险",
	}

	if sampleCount < 5 {
		level = "high"
		warnings = append([]string{"样本量严重不足，统计结果不具备参考价值"}, warnings...)
	}

	return &response.RiskWarning{
		Level:    level,
		Warnings: warnings,
	}
}

// buildHourRiskWarning 构建小时风险警告
func buildHourRiskWarning(sampleCount int) *response.RiskWarning {
	level := "medium"
	warnings := []string{
		"历史时段数据不代表未来表现，市场随时可能变化",
		"重大新闻和事件可能改变时段特征",
		"请结合实时行情、成交量、技术指标综合判断",
		"建议设置止损止盈，控制单笔交易风险",
	}

	if sampleCount < 10 {
		level = "high"
		warnings = append([]string{"样本量不足，统计结果可能不可靠"}, warnings...)
	}

	return &response.RiskWarning{
		Level:    level,
		Warnings: warnings,
	}
}
