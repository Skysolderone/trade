package main

import (
	"fmt"
	"time"

	"trade/db"
	"trade/model"
)

func main() {
	// 初始化数据库
	db.InitPostgreSql()

	fmt.Println("\n╔═══════════════════════════════════════════╗")
	fmt.Println("║       旧数据修复工具                      ║")
	fmt.Println("╚═══════════════════════════════════════════╝\n")

	// 1. 统计需要修复的数据
	var needFixCount int64
	db.Pog.Model(&model.Kline{}).
		Where("interval IS NULL OR interval = ''").
		Count(&needFixCount)

	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   需要修复的记录数: %d 条\n\n", needFixCount)

	if needFixCount == 0 {
		fmt.Println("✅ 没有需要修复的数据！")
		return
	}

	fmt.Printf("⚠️  即将修复以下问题:\n")
	fmt.Printf("   1. 将空的 interval 字段设置为 '1d'\n")
	fmt.Printf("   2. 统一 day 字段格式为两位数 (例如: 9-8 → 09-08)\n\n")

	fmt.Print("确认执行？(输入 'yes' 继续): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		fmt.Println("❌ 操作已取消")
		return
	}

	fmt.Println("\n开始修复...\n")

	// 2. 分批修复数据
	batchSize := 100
	offset := 0
	totalFixed := 0

	for {
		var klines []model.Kline
		result := db.Pog.Where("interval IS NULL OR interval = ''").
			Offset(offset).
			Limit(batchSize).
			Find(&klines)

		if result.Error != nil {
			fmt.Printf("❌ 查询失败: %v\n", result.Error)
			return
		}

		if len(klines) == 0 {
			break
		}

		// 修复每条记录
		for _, kline := range klines {
			// 修复 interval
			if kline.Interval == "" {
				kline.Interval = "1d"
			}

			// 修复 day 格式
			if kline.OpenTime.IsZero() {
				continue
			}
			correctDay := fmt.Sprintf("%02d-%02d",
				int(kline.OpenTime.Month()),
				kline.OpenTime.Day())
			kline.Day = correctDay

			// 保存修改
			err := db.Pog.Save(&kline).Error
			if err != nil {
				fmt.Printf("⚠️  修复失败 ID=%d: %v\n", kline.ID, err)
				continue
			}

			totalFixed++
			if totalFixed%100 == 0 {
				fmt.Printf("   已修复 %d 条...\n", totalFixed)
			}
		}

		offset += batchSize
		time.Sleep(10 * time.Millisecond) // 避免数据库压力过大
	}

	fmt.Printf("\n✅ 修复完成！共修复 %d 条记录\n\n", totalFixed)

	// 3. 验证修复结果
	var stillNeedFix int64
	db.Pog.Model(&model.Kline{}).
		Where("interval IS NULL OR interval = ''").
		Count(&stillNeedFix)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 验证结果:")
	fmt.Printf("   修复前: %d 条\n", needFixCount)
	fmt.Printf("   修复后: %d 条\n", stillNeedFix)

	if stillNeedFix == 0 {
		fmt.Println("   状态: ✅ 所有数据已修复")
	} else {
		fmt.Printf("   状态: ⚠️  还有 %d 条数据未修复\n", stillNeedFix)
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 4. 显示修复后的样本数据
	fmt.Println("📋 修复后的数据样本:")
	var samples []model.Kline
	db.Pog.Where("interval = '1d'").Limit(5).Find(&samples)
	for _, s := range samples {
		fmt.Printf("   Symbol: %s | Interval: %s | Day: %s | OpenTime: %s\n",
			s.Symbol, s.Interval, s.Day, s.OpenTime.Format("2006-01-02"))
	}

	fmt.Println("\n🎉 数据修复完成！现在可以正常使用策略了。")
}
