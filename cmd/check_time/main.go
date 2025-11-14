package main

import (
	"fmt"
	"time"

	"trade/db"
	"trade/model"
)

func main() {
	db.InitPostgreSql()

	// 查询策略一的详细记录，看看 close_time
	var records []model.Strategy1DetailRecord
	db.Pog.Limit(10).Order("id ASC").Find(&records)

	fmt.Println("策略一详细记录样本（前10条）：")
	fmt.Println("ID\tYear\t\tCloseTime")
	for _, r := range records {
		fmt.Printf("%d\t%s\t%v\n", r.ID, r.Year, r.CloseTime)
	}

	// 统计零值时间的数量
	var count int64
	zeroTime := time.Time{}
	db.Pog.Model(&model.Strategy1DetailRecord{}).Where("close_time = ?", zeroTime).Count(&count)
	fmt.Printf("\n策略一零值时间记录数: %d\n", count)

	// 统计早于1900年的时间数量
	earlyTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	db.Pog.Model(&model.Strategy1DetailRecord{}).Where("close_time < ?", earlyTime).Count(&count)
	fmt.Printf("策略一早于1900年的记录数: %d\n", count)

	// 查询策略二的详细记录
	var records2 []model.Strategy2DetailRecord
	db.Pog.Limit(10).Order("id ASC").Find(&records2)

	fmt.Println("\n策略二详细记录样本（前10条）：")
	fmt.Println("ID\tDate\t\tCloseTime")
	for _, r := range records2 {
		fmt.Printf("%d\t%s\t%v\n", r.ID, r.Date, r.CloseTime)
	}

	// 统计零值时间的数量
	db.Pog.Model(&model.Strategy2DetailRecord{}).Where("close_time = ?", zeroTime).Count(&count)
	fmt.Printf("\n策略二零值时间记录数: %d\n", count)

	// 统计早于1900年的时间数量
	db.Pog.Model(&model.Strategy2DetailRecord{}).Where("close_time < ?", earlyTime).Count(&count)
	fmt.Printf("策略二早于1900年的记录数: %d\n", count)
}
