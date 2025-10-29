package db

import (
	"github.com/adshao/go-binance/v2/futures"
)

// BinanceClient 全局币安合约客户端
var BinanceClient *futures.Client

// InitBinance 初始化币安客户端
func InitBinance(apiKey, secretKey string) {
	BinanceClient = futures.NewClient(apiKey, secretKey)
}
