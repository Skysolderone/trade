package kline

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"trade/db"
	"trade/model"
	"trade/utils"

	"github.com/adshao/go-binance/v2/futures"
	"gorm.io/gorm/clause"
)

func UpdateKline(symbol string, interval string) {
	// 查询数据库中该交易对的最新记录
	var latestKline model.Kline
	result := db.Pog.Where("symbol = ?", symbol).
		Order("close_time DESC").
		Limit(1).
		Find(&latestKline)

	var startTime time.Time
	if result.Error != nil || result.RowsAffected == 0 {
		// 如果没有找到记录,从默认时间开始
		fmt.Printf("未找到 %s 的历史数据,将从 2018-08-01 开始获取\n", symbol)
		startTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		// 删除最新记录(因为它可能是不完整的)
		deleteResult := db.Pog.Delete(&latestKline)
		if deleteResult.Error != nil {
			fmt.Printf("删除最新记录失败: %v\n", deleteResult.Error)
			return
		}
		fmt.Printf("已删除 %s 的最新记录: %s\n",
			symbol,
			latestKline.CloseTime.Format("2006-01-02 15:04:05"))

		// 从被删除记录的开始时间重新获取
		startTime = latestKline.OpenTime
		fmt.Printf("将从 %s 开始重新获取数据\n", startTime.Format("2006-01-02 15:04:05"))
	}

	// 结束时间设置为昨日最后一刻
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	endTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.UTC)

	// 如果开始时间已经超过结束时间,说明数据已经是最新的
	if startTime.After(endTime) {
		fmt.Printf("%s 的数据已经是最新的,无需更新\n", symbol)
		return
	}

	// 调用更新函数
	updateKlineData(db.BinanceClient, symbol, interval, startTime, endTime)
}

func GetKline(symbol string, interval string) {
	// 使用全局币安客户端
	getAllKlines(db.BinanceClient, symbol, interval)
}

// updateKlineData 更新K线数据(支持自定义时间范围)
func updateKlineData(api *futures.Client, symbol string, interval string, startTime, endTime time.Time) {
	totalCount := 0
	var lastKlineTime time.Time

	for {
		// 计算当前批次的结束时间
		batchEndTime := startTime.Add(time.Duration(1000) * getIntervalDuration(interval))
		if batchEndTime.After(endTime) {
			batchEndTime = endTime
		}

		fmt.Printf("正在获取 %s 从 %s 到 %s 的K线数据...\n",
			symbol, startTime.Format("2006-01-02"), batchEndTime.Format("2006-01-02"))

		// 请求当前批次的K线数据
		kline, err := api.NewContinuousKlinesService().
			ContractType("PERPETUAL").
			Pair(symbol).
			Interval(interval).
			StartTime(startTime.UnixMilli()).
			EndTime(batchEndTime.UnixMilli()).
			Limit(1000).
			Do(context.Background())
		if err != nil {
			fmt.Printf("获取K线数据失败: %v\n", err)
			return
		}
		fmt.Println(kline)
		if len(kline) == 0 {
			fmt.Println("没有更多数据了")
			break
		}

		// 批量构建K线数据（提高性能）
		klineModels := make([]model.Kline, 0, len(kline))
		for _, k := range kline {
			openTime := time.Unix(k.OpenTime/1000, 0).UTC()
			closeTime := time.Unix(k.CloseTime/1000, 0).UTC()
			year := openTime.Year()

			klineModel := model.Kline{
				Symbol:    symbol,
				Interval:  interval, // 添加时间周期字段
				Open:      utils.StringToFloat64(k.Open),
				High:      utils.StringToFloat64(k.High),
				Low:       utils.StringToFloat64(k.Low),
				Close:     utils.StringToFloat64(k.Close),
				OpenTime:  openTime,
				CloseTime: closeTime,
				Date:      strconv.Itoa(year),
				Day:       fmt.Sprintf("%02d-%02d", int(openTime.Month()), openTime.Day()), // 修复：使用两位数格式
				Hour:      strconv.Itoa(openTime.Hour()),
				Week:      strconv.Itoa(int(openTime.Weekday())%7 + 1),
				Min:       strconv.Itoa(openTime.Minute()),
			}
			klineModels = append(klineModels, klineModel)
			lastKlineTime = closeTime
		}

		// 批量插入，遇到重复则跳过（ON CONFLICT DO NOTHING）
		// 性能提升：1000条数据从 ~2000ms 降到 ~50ms
		if len(klineModels) > 0 {
			result := db.Pog.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "symbol"}, {Name: "interval"}, {Name: "open_time"}},
				DoNothing: true, // 遇到冲突直接跳过，不报错
			}).Create(&klineModels)

			if result.Error != nil {
				fmt.Printf("批量插入失败: %v\n", result.Error)
			} else {
				insertedCount := result.RowsAffected
				totalCount += int(insertedCount)
				fmt.Printf("批量插入 %d 条数据（跳过 %d 条重复数据）\n",
					insertedCount, len(klineModels)-int(insertedCount))
			}
		}

		// 更新开始时间为最后一条K线的收盘时间
		if len(kline) > 0 {
			lastKline := kline[len(kline)-1]
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

// getAllKlines 分页获取所有K线数据
func getAllKlines(api *futures.Client, symbol string, interval string) {
	// 设置开始时间（可以根据需要调整）
	startTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	// 结束时间设置为今日零点
	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 调用统一的更新函数
	updateKlineData(api, symbol, interval, startTime, endTime)
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
		return time.Hour // 默认1小时
	}
}
