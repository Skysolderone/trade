package main

import (
	"fmt"

	"trade/db"
	"trade/kline"
)

func main() {
	// 初始化数据库连接
	db.InitPostgreSql()
	// 更新工具  更新所有时间段k线数据
	// 获取K线数据
	kline.GetKline("BTCUSDT", "1d")
	fmt.Println("Hello, World!")
}
