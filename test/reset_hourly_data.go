package main

import (
	"fmt"

	"trade/db"
	"trade/model"
)

func main() {
	// 初始化数据库
	db.InitPostgreSql()

	intervals := []string{"8h", "4h", "2h", "1h"}

	fmt.Println("========== 清理小时级K线数据 ==========\n")

	for _, interval := range intervals {
		var count int64
		db.Pog.Model(&model.Kline{}).
			Where("interval = ?", interval).
			Count(&count)

		if count > 0 {
			result := db.Pog.Where("interval = ?", interval).Delete(&model.Kline{})
			if result.Error != nil {
				fmt.Printf("❌ 删除 %s 数据失败: %v\n", interval, result.Error)
			} else {
				fmt.Printf("✅ 已删除 %s 的 %d 条数据\n", interval, result.RowsAffected)
			}
		} else {
			fmt.Printf("⚠️  %s 没有数据\n", interval)
		}
	}

	fmt.Println("\n========== 清理完成 ==========")
}
