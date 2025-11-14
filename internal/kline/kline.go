package kline

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"trade/db"
	"trade/model"
	"trade/utils"

	"gorm.io/gorm/clause"
)

var TimeMap = map[time.Time]time.Time{
	time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
}

func GetAllKlines(symbol string, interval string) {
	for startTime, endTime := range TimeMap {
		GetKline(symbol, interval, startTime, endTime)
	}
}

func GetKline(symbol string, interval string, startTime, endTime time.Time) {
	klines, err := db.BinanceClient.NewContinuousKlinesService().
		ContractType("PERPETUAL").
		Pair(symbol).
		Interval(interval).
		StartTime(startTime.UnixMilli()).
		EndTime(endTime.UnixMilli()).
		Limit(1000).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	klineModels := make([]model.KlineWs, 0, len(klines))
	for _, k := range klines {
		openTime := time.Unix(k.OpenTime/1000, 0).UTC()
		closeTime := time.Unix(k.CloseTime/1000, 0).UTC()
		year := openTime.Year()

		klineModel := model.KlineWs{
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
			fmt.Printf("批量插入 %d 条数据\n", result.RowsAffected)
		}
	}
}
