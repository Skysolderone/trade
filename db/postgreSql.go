package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库连接
var DB *gorm.DB

// InitPostgreSql 初始化PostgreSQL数据库连接
func InitPostgreSql() *gorm.DB {
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=trade port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 测试连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(&Kline{})
	if err != nil {
		log.Printf("自动迁移失败: %v", err)
	}

	log.Println("数据库连接成功")
	return DB
}
