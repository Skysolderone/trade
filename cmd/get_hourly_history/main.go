package main

import (
	"context"
	"fmt"
	"time"

	"trade/db"
	"trade/model"
	"trade/utils"

	"github.com/adshao/go-binance/v2/futures"
)

// getEarliestKlineTime 获取交易对最早的K线时间（通过API查询）
func getEarliestKlineTime(api *futures.Client, symbol string) time.Time {
	// 从2019年初开始尝试（币安永续合约大约从这个时间开始）
	testTime := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	// 尝试获取一条数据
	klines, err := api.NewContinuousKlinesService().
		ContractType("PERPETUAL").
		Pair(symbol).
		Interval("1d").
		StartTime(testTime.UnixMilli()).
		Limit(1).
		Do(context.Background())

	if err != nil || len(klines) == 0 {
		fmt.Printf("⚠️  无法获取 %s 的最早时间，使用 2019-09-01\n", symbol)
		return time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)
	}

	earliestTime := time.Unix(klines[0].OpenTime/1000, 0).UTC()
	fmt.Printf("📅 %s 的最早数据时间: %s\n", symbol, earliestTime.Format("2006-01-02"))
	return earliestTime
}

func main() {
	// 初始化币安客户端
	db.InitBinance("", "")

	// 初始化数据库连接
	db.InitPostgreSql()

	// 加载配置文件
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("❌ 加载配置文件失败: %v\n", err)
		return
	}

	// 只获取小时级别的数据
	hourlyIntervals := []string{"8h", "4h", "2h", "1h"}

	fmt.Println("========== 开始获取小时级K线历史数据 ==========\n")

	// 遍历配置文件中的所有交易对
	for _, symbolConfig := range config.Symbols {
		fmt.Printf("\n========== 处理交易对: %s ==========\n", symbolConfig.Symbol)

		for _, interval := range hourlyIntervals {
			fmt.Printf("\n--- 时间区间: %s ---\n", interval)

			// 先查询该交易对该时间周期的最新记录
			var latestKline model.Kline
			result := db.Pog.Where("symbol = ? AND interval = ?", symbolConfig.Symbol, interval).
				Order("close_time DESC").
				Limit(1).
				Find(&latestKline)

			var startTime time.Time
			var updateMode string

			if result.Error != nil || result.RowsAffected == 0 {
				// 没有该时间周期的数据，需要全量获取
				updateMode = "全量获取"
				fmt.Printf("📊 数据库中没有 %s %s 的数据，准备全量获取...\n", symbolConfig.Symbol, interval)

				// 先从数据库查询该交易对1d数据的最早时间
				var earliestKline model.Kline
				result1d := db.Pog.Where("symbol = ? AND interval = ?", symbolConfig.Symbol, "1d").
					Order("open_time ASC").
					Limit(1).
					Find(&earliestKline)

				if result1d.Error != nil || result1d.RowsAffected == 0 {
					// 如果数据库没有1d数据，通过API查询最早时间
					fmt.Printf("📊 数据库中没有 %s 的1d数据，通过API查询最早时间...\n", symbolConfig.Symbol)
					startTime = getEarliestKlineTime(db.BinanceClient, symbolConfig.Symbol)
				} else {
					startTime = earliestKline.OpenTime
					fmt.Printf("📅 从1d数据获取最早时间: %s\n", startTime.Format("2006-01-02"))
				}
			} else {
				// 有数据，增量更新
				updateMode = "增量更新"
				fmt.Printf("✅ 找到最新记录: %s\n", latestKline.CloseTime.Format("2006-01-02 15:04:05"))

				// 删除最新记录(因为它可能是不完整的)
				deleteResult := db.Pog.Delete(&latestKline)
				if deleteResult.Error != nil {
					fmt.Printf("❌ 删除最新记录失败: %v\n", deleteResult.Error)
					continue
				}
				fmt.Printf("🗑️  已删除最新记录，将从 %s 重新获取\n", latestKline.OpenTime.Format("2006-01-02 15:04:05"))

				// 从被删除记录的开始时间重新获取
				startTime = latestKline.OpenTime
			}

			// 计算结束时间（昨天）
			now := time.Now().UTC()
			yesterday := now.AddDate(0, 0, -1)
			endTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.UTC)

			// 如果开始时间已经超过结束时间，说明数据已经是最新的
			if startTime.After(endTime) {
				fmt.Printf("✅ %s %s 的数据已经是最新的，无需更新\n", symbolConfig.Symbol, interval)
				continue
			}

			fmt.Printf("🚀 开始%s: 从 %s 到 %s\n", updateMode, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

			// 手动调用updateKlineData获取指定时间范围的数据
			updateKlineData(db.BinanceClient, symbolConfig.Symbol, interval, startTime, endTime)
		}
	}

	fmt.Println("\n========== 所有小时级历史数据获取完成 ==========")
}

// updateKlineData 更新K线数据(从kline/rest.go复制，避免import cycle)
func updateKlineData(api *futures.Client, symbol string, interval string, startTime, endTime time.Time) {
	totalCount := 0
	var lastKlineTime time.Time

	for {
		// 计算当前批次的结束时间
		batchEndTime := startTime.Add(time.Duration(1500) * getIntervalDuration(interval))
		if batchEndTime.After(endTime) {
			batchEndTime = endTime
		}

		fmt.Printf("正在获取 %s 从 %s 到 %s 的K线数据...\n",
			symbol, startTime.Format("2006-01-02"), batchEndTime.Format("2006-01-02"))

		// 请求当前批次的K线数据
		klines, err := api.NewContinuousKlinesService().
			ContractType("PERPETUAL").
			Pair(symbol).
			Interval(interval).
			StartTime(startTime.UnixMilli()).
			EndTime(batchEndTime.UnixMilli()).
			Limit(1500).
			Do(context.Background())
		if err != nil {
			fmt.Printf("获取K线数据失败: %v\n", err)
			return
		}

		if len(klines) == 0 {
			fmt.Println("没有更多数据了")
			break
		}

		// 批量构建K线数据
		klineModels := make([]model.Kline, 0, len(klines))
		for _, k := range klines {
			openTime := time.Unix(k.OpenTime/1000, 0).UTC()
			closeTime := time.Unix(k.CloseTime/1000, 0).UTC()
			year := openTime.Year()

			klineModel := model.Kline{
				Symbol:    symbol,
				Interval:  interval,
				Open:      utils.StringToFloat64(k.Open),
				High:      utils.StringToFloat64(k.High),
				Low:       utils.StringToFloat64(k.Low),
				Close:     utils.StringToFloat64(k.Close),
				OpenTime:  openTime,
				CloseTime: closeTime,
				Date:      fmt.Sprintf("%d", year),
				Day:       fmt.Sprintf("%02d-%02d", int(openTime.Month()), openTime.Day()),
				Hour:      fmt.Sprintf("%d", openTime.Hour()),
				Week:      fmt.Sprintf("%d", int(openTime.Weekday())%7+1),
				Min:       fmt.Sprintf("%d", openTime.Minute()),
			}
			klineModels = append(klineModels, klineModel)
			lastKlineTime = closeTime
		}

		// 批量插入
		if len(klineModels) > 0 {
			result := db.Pog.Create(&klineModels)
			if result.Error != nil {
				fmt.Printf("批量插入失败: %v\n", result.Error)
			} else {
				insertedCount := result.RowsAffected
				totalCount += int(insertedCount)
				fmt.Printf("批量插入 %d 条数据\n", insertedCount)
			}
		}

		// 更新开始时间为最后一条K线的收盘时间
		if len(klines) > 0 {
			lastKline := klines[len(klines)-1]
			startTime = time.Unix(lastKline.CloseTime/1000, 0)
		}

		// 如果已经到达结束时间，退出循环
		if startTime.After(endTime) || startTime.Equal(endTime) {
			break
		}

		// 添加延迟避免API限制
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("完成！总共获取了 %d 条K线数据\n", totalCount)
	if !lastKlineTime.IsZero() {
		fmt.Printf("最后一条K线时间: %s\n", lastKlineTime.Format("2006-01-02 15:04:05"))
	}
}

// getIntervalDuration 根据间隔字符串返回对应的时长
func getIntervalDuration(interval string) time.Duration {
	switch interval {
	case "1m":
		return time.Minute
	case "3m":
		return 3 * time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "2h":
		return 2 * time.Hour
	case "4h":
		return 4 * time.Hour
	case "6h":
		return 6 * time.Hour
	case "8h":
		return 8 * time.Hour
	case "12h":
		return 12 * time.Hour
	case "1d":
		return 24 * time.Hour
	case "3d":
		return 3 * 24 * time.Hour
	case "1w":
		return 7 * 24 * time.Hour
	case "1M":
		return 30 * 24 * time.Hour
	default:
		return time.Hour
	}
}
