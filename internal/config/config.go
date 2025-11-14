package config

import (
	"sync"

	"trade/db"
)

var once sync.Once

func InitConfig() {
	once.Do(func() {
		db.InitPostgreSqlWs()
		db.InitBinance("", "")
	})
}
