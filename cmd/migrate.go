package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	// 读取SQL文件
	sqlFile := "db/migrations/create_klines_table.sql"
	log.Printf("正在读取SQL文件: %s", sqlFile)
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 执行SQL
	log.Println("正在执行建表SQL...")
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("执行SQL失败: %v", err)
	}

	log.Println("✅ 表创建成功！")

	// 验证表是否存在
	var tableName string
	err = db.QueryRow("SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename = 'klines'").Scan(&tableName)
	if err != nil {
		log.Printf("警告：验证表失败: %v", err)
	} else {
		fmt.Printf("✅ 表 '%s' 已成功创建并验证\n", tableName)
	}
}
