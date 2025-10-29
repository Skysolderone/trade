package strategy

import (
	"fmt"
	"time"

	"trade/db"
	"trade/model"
)

// 从一年的第一天开始对比其他几年的开盘收盘价，计算涨跌概率

func Strategy1() {
	// 获取今天的时间
	now := time.Now()
	// 解析当前是几月几号
	month := now.Month()
	day := now.Day()
	date := fmt.Sprintf("%d-%d", int(month), day)
	fmt.Printf("查询日期: %s\n", date)

	// 从数据库查找对应的open_time的日期为当前日期的时间段k线数据
	var klines []model.Kline
	err := db.Pog.Where("day = ?", date).Find(&klines).Error
	if err != nil {
		fmt.Printf("查询数据库失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 条K线数据\n", len(klines))
	// 根据不同date的k线数据，计算涨跌概率
	upCount := 0
	downCount := 0
	for _, kline := range klines {
		if kline.Close > kline.Open {
			upCount++
		} else {
			downCount++
		}
	}
	upRate := float64(upCount) / float64(len(klines))
	downRate := float64(downCount) / float64(len(klines))
	// 还需要显示哪年哪天涨了多少跟跌了多少
	fmt.Printf("涨概率: %f, 跌概率: %f\n", upRate, downRate)
	for _, kline := range klines {
		if kline.Close > kline.Open {
			fmt.Printf("涨了: %f, 年份-日期: %s 开盘价: %f 收盘价: %f\n", kline.Close-kline.Open, kline.Date+"-"+kline.Day, kline.Open, kline.Close)
		} else {
			fmt.Printf("跌了: %f, 年份-日期: %s 开盘价: %f 收盘价: %f\n", kline.Open-kline.Close, kline.Date+"-"+kline.Day, kline.Open, kline.Close)
		}
	}
}
