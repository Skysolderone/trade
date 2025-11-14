package main

import (
	"trade/internal/config"
	"trade/internal/kline"
)

func main() {
	// 初始化配置
	config.InitConfig()
	kline.GetAllKlines("BTCUSDT", "1h")
}
