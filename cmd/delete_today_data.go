package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// 数据库连接字符串
	connStr := "host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=trade port=5432 sslmode=disable"

	// 连接数据库
	log.Println("正在连接数据库...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping数据库失败: %v", err)
	}
	log.Println("数据库连接成功！")

	// 计算今天的开始时间戳（UTC+8）
	now := time.Now().UTC().Add(8 * time.Hour)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStartTimestamp := todayStart.UnixMilli()

	log.Printf("当前时间: %s", now.Format("2006-01-02 15:04:05"))
	log.Printf("今天开始时间: %s (时间戳: %d)", todayStart.Format("2006-01-02 15:04:05"), todayStartTimestamp)

	// 删除今天的数据
	log.Println("正在删除今天的未结束K线数据...")
	result, err := db.Exec("DELETE FROM klines_day WHERE open_time >= $1", todayStartTimestamp)
	if err != nil {
		log.Fatalf("删除数据失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ 已删除 %d 条未结束的K线数据", rowsAffected)

	// 查询剩余数据统计
	var count int
	var maxDate time.Time
	err = db.QueryRow("SELECT COUNT(*), MAX(open_time_dt) FROM klines_day").Scan(&count, &maxDate)
	if err != nil {
		log.Printf("警告：查询数据统计失败: %v", err)
	} else {
		log.Printf("📊 表中剩余 %d 条数据", count)
		log.Printf("📅 最新数据日期: %s", maxDate.Format("2006-01-02"))
	}
}
