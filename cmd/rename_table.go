package main

import (
	"database/sql"
	"log"

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

	// 重命名表
	log.Println("正在将表 klines 重命名为 klines_day...")
	_, err = db.Exec("ALTER TABLE klines RENAME TO klines_day;")
	if err != nil {
		log.Fatalf("重命名表失败: %v", err)
	}

	log.Println("✅ 表重命名成功！")

	// 验证新表名
	var tableName string
	err = db.QueryRow("SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename = 'klines_day'").Scan(&tableName)
	if err != nil {
		log.Printf("警告：验证表失败: %v", err)
	} else {
		log.Printf("✅ 表 '%s' 已成功重命名并验证", tableName)
	}

	// 查询表中的数据统计
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM klines_day").Scan(&count)
	if err != nil {
		log.Printf("警告：查询数据统计失败: %v", err)
	} else {
		log.Printf("📊 表中共有 %d 条数据", count)
	}
}
