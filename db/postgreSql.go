package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreSql() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=trade port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
