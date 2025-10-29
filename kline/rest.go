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
)

func GetKline(symbol string, interval string) {
	api := futures.NewClient("", "")

	// 获取所有K线数据
	getAllKlines(api, symbol, interval)
}

// getAllKlines 分页获取所有K线数据
func getAllKlines(api *futures.Client, symbol string, interval string) {
	// 设置开始时间（可以根据需要调整）
	startTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	// 结束时间设置为今日零点
	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

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
		kline, err := api.NewContinuousKlinesService().
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

		if len(kline) == 0 {
			fmt.Println("没有更多数据了")
			break
		}

		// 处理当前批次的K线数据
		for _, k := range kline {
			openTime := time.Unix(k.OpenTime/1000, 0).UTC()
			closeTime := time.Unix(k.CloseTime/1000, 0).UTC()
			year := openTime.Year()

			klineModel := model.Kline{
				Symbol:    symbol,
				Open:      utils.StringToFloat64(k.Open),
				High:      utils.StringToFloat64(k.High),
				Low:       utils.StringToFloat64(k.Low),
				Close:     utils.StringToFloat64(k.Close),
				OpenTime:  openTime,
				CloseTime: closeTime,
				Date:      strconv.Itoa(year),
				Day:       fmt.Sprintf("%d-%d", int(openTime.Month()), openTime.Day()),
				Hour:      strconv.Itoa(openTime.Hour()),
				Week:      strconv.Itoa(int(openTime.Weekday())%7 + 1),
				Min:       strconv.Itoa(openTime.Minute()),
			}

			// 这里可以添加数据库保存逻辑
			result := db.Pog.Create(&klineModel)
			if result.Error != nil {
				fmt.Println(result.Error)
				continue
			}

			// 更新最后一条K线的时间
			lastKlineTime = closeTime

			totalCount++
			if totalCount%100 == 0 {
				fmt.Printf("已处理 %d 条K线数据\n", totalCount)
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
