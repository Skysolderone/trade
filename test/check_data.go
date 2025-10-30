package main

import (
	"fmt"

	"trade/db"
	"trade/model"
)

func main() {
	// 初始化数据库
	db.InitPostgreSql()

	fmt.Println("========== 数据库诊断 ==========\n")

	// 1. 检查总数据量
	var totalCount int64
	db.Pog.Model(&model.Kline{}).Count(&totalCount)
	fmt.Printf("1. 数据库总记录数: %d\n\n", totalCount)

	if totalCount == 0 {
		fmt.Println("❌ 数据库中没有任何数据！")
		fmt.Println("请先运行 main.go 采集数据")
		return
	}

	// 2. 检查不同交易对
	var symbols []string
	db.Pog.Model(&model.Kline{}).Distinct("symbol").Pluck("symbol", &symbols)
	fmt.Printf("2. 交易对列表: %v\n\n", symbols)

	// 3. 检查不同时间周期
	var intervals []string
	db.Pog.Model(&model.Kline{}).Distinct("interval").Pluck("interval", &intervals)
	fmt.Printf("3. 时间周期列表: %v\n", intervals)
	if len(intervals) == 0 {
		fmt.Println("   ⚠️  interval 字段为空！这是问题所在！")
	}
	fmt.Println()

	// 4. 查看前10条数据的详细信息
	var samples []model.Kline
	db.Pog.Limit(10).Find(&samples)
	fmt.Println("4. 前10条数据示例：")
	fmt.Println("   Symbol | Interval | Day     | OpenTime")
	fmt.Println("   -------|----------|---------|---------------------------")
	for _, s := range samples {
		fmt.Printf("   %-7s| %-9s| %-8s| %s\n",
			s.Symbol, s.Interval, s.Day, s.OpenTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// 5. 检查 day 字段格式
	var days []string
	db.Pog.Model(&model.Kline{}).Distinct("day").Order("day").Limit(20).Pluck("day", &days)
	fmt.Printf("5. Day 字段样本（前20个）: %v\n", days)
	fmt.Println()

	// 6. 测试查询 29 号数据
	var count29 int64
	db.Pog.Model(&model.Kline{}).Where("day LIKE ?", "%-29").Count(&count29)
	fmt.Printf("6. 包含 '29' 的记录数: %d\n", count29)

	if count29 > 0 {
		var sample29 []model.Kline
		db.Pog.Where("day LIKE ?", "%-29").Limit(5).Find(&sample29)
		fmt.Println("   示例数据:")
		for _, s := range sample29 {
			fmt.Printf("   - Day: %s, Symbol: %s, Interval: %s, OpenTime: %s\n",
				s.Day, s.Symbol, s.Interval, s.OpenTime.Format("2006-01-02"))
		}
	}
	fmt.Println()

	// 7. 测试策略查询条件
	symbol := "BTCUSDT"
	interval := "1d"
	day := "10-29"

	var testCount int64
	db.Pog.Model(&model.Kline{}).
		Where("symbol = ? AND interval = ? AND day = ?", symbol, interval, day).
		Count(&testCount)

	fmt.Printf("7. 策略查询测试:\n")
	fmt.Printf("   查询条件: symbol='%s' AND interval='%s' AND day='%s'\n", symbol, interval, day)
	fmt.Printf("   结果数量: %d\n", testCount)

	if testCount == 0 {
		fmt.Println("\n❌ 查询失败！可能的原因:")

		// 测试没有 interval 的情况
		var countNoInterval int64
		db.Pog.Model(&model.Kline{}).
			Where("symbol = ? AND day = ?", symbol, day).
			Count(&countNoInterval)
		fmt.Printf("   - 不加 interval 条件的结果: %d 条\n", countNoInterval)

		// 测试 interval 为空的情况
		var countNullInterval int64
		db.Pog.Model(&model.Kline{}).
			Where("symbol = ? AND (interval IS NULL OR interval = '') AND day = ?", symbol, day).
			Count(&countNullInterval)
		fmt.Printf("   - interval 为空的记录: %d 条\n", countNullInterval)

		// 检查实际的 interval 值
		var actualIntervals []string
		db.Pog.Model(&model.Kline{}).
			Where("symbol = ?", symbol).
			Distinct("interval").
			Pluck("interval", &actualIntervals)
		fmt.Printf("   - BTCUSDT 实际的 interval 值: %v\n", actualIntervals)

		// 检查实际的 day 值（29号相关）
		var actualDays []string
		db.Pog.Model(&model.Kline{}).
			Where("symbol = ? AND day LIKE ?", symbol, "%-29").
			Distinct("day").
			Order("day").
			Pluck("day", &actualDays)
		fmt.Printf("   - BTCUSDT 包含'29'的 day 值: %v\n", actualDays)
	} else {
		fmt.Println("\n✅ 查询成功！")
	}

	fmt.Println("\n========== 诊断完成 ==========")
}
