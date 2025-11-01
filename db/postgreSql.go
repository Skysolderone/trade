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
		log.Println("已创建唯一索引: idx_unique_kline (symbol + interval + open_time)")
		log.Println("已创建策略结果表")
	}

	log.Println("数据库连接成功")
	return Pog
}
