package db

import (
	"log"
	"time"

	"trade/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库连接
var Pog *gorm.DB

// InitPostgreSql 初始化PostgreSQL数据库连接
func InitPostgreSql() *gorm.DB {
	var err error
	Pog, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=trade port=5432 sslmode=disable TimeZone=UTC",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 测试连接

	sqlDB, err := Pog.DB()
	if err != nil {
		log.Fatal(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	err = Pog.AutoMigrate(
		&model.Kline{},
		&model.Strategy1Result{},
		&model.Strategy1DetailRecord{},
		&model.Strategy2Result{},
		&model.Strategy2DetailRecord{},
	)
	if err != nil {
		log.Printf("自动迁移失败: %v", err)
	} else {
		log.Println("数据库表结构迁移成功")
		log.Println("已创建策略结果表")
	}

	// 显式创建唯一索引
	createUniqueIndexes()

	log.Println("数据库连接成功")
	return Pog
}

// InitPostgreSql 初始化PostgreSQL数据库连接
func InitPostgreSqlWs() *gorm.DB {
	var err error
	Pog, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=trade port=5432 sslmode=disable TimeZone=UTC",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 测试连接

	sqlDB, err := Pog.DB()
	if err != nil {
		log.Fatal(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	err = Pog.AutoMigrate(
		&model.KlineWs{},
	)
	if err != nil {
		log.Printf("自动迁移失败: %v", err)
	} else {
		log.Println("数据库表结构迁移成功")
	}

	// 显式创建唯一索引
	createUniqueIndexesWs()

	log.Println("数据库连接成功")
	return Pog
}

// createUniqueIndexes 为 Kline 表创建唯一索引
func createUniqueIndexes() {
	// 检查索引是否已存在
	var count int64
	err := Pog.Raw("SELECT COUNT(*) FROM pg_indexes WHERE indexname = ?", "idx_unique_kline").Scan(&count).Error
	if err != nil {
		log.Printf("检查索引时出错: %v", err)
		return
	}

	if count > 0 {
		log.Println("✅ 唯一索引 idx_unique_kline 已存在")
		return
	}

	// 创建唯一索引
	sql := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_kline ON klines (symbol, interval, open_time)"
	err = Pog.Exec(sql).Error
	if err != nil {
		log.Printf("❌ 创建唯一索引失败: %v", err)
	} else {
		log.Println("✅ 成功创建唯一索引: idx_unique_kline (symbol + interval + open_time)")
	}
}

// createUniqueIndexesWs 为 KlineWs 表创建唯一索引
func createUniqueIndexesWs() {
	// 检查索引是否已存在
	var count int64
	err := Pog.Raw("SELECT COUNT(*) FROM pg_indexes WHERE indexname = ?", "idx_unique_kline_ws").Scan(&count).Error
	if err != nil {
		log.Printf("检查索引时出错: %v", err)
		return
	}

	if count > 0 {
		log.Println("✅ 唯一索引 idx_unique_kline_ws 已存在")
		return
	}

	// 创建唯一索引
	sql := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_kline_ws ON kline_ws (symbol, interval, open_time)"
	err = Pog.Exec(sql).Error
	if err != nil {
		log.Printf("❌ 创建唯一索引失败: %v", err)
	} else {
		log.Println("✅ 成功创建唯一索引: idx_unique_kline_ws (symbol + interval + open_time)")
	}
}
