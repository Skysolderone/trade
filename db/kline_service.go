package db

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// BinanceKline 币安API返回的K线数据格式
type BinanceKline []interface{}

// FetchKlinesFromBinance 从币安获取K线数据
// symbol: 交易对，如 BTCUSDT
// interval: K线周期，如 1m, 5m, 1h, 1d
// limit: 获取数量，默认500，最大1500
func FetchKlinesFromBinance(symbol string, interval string, limit int) ([]BinanceKline, error) {
	if limit == 0 {
		limit = 1440 // 默认获取24小时的1分钟数据
	}

	// 币安永续合约K线接口
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=%d",
		symbol, interval, limit)

	log.Printf("正在获取K线数据: %s", url)

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON
	var klines []BinanceKline
	err = json.Unmarshal(body, &klines)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	log.Printf("成功获取 %d 条K线数据", len(klines))
	return klines, nil
}

// SaveKlinesToDB 保存K线数据到数据库
func SaveKlinesToDB(symbol string, interval string, binanceKlines []BinanceKline) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 计算当前时间（UTC+8，即北京时间）
	now := time.Now().UTC().Add(8 * time.Hour)

	// 计算当天的开始时间（UTC+8的00:00:00）
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStartTimestamp := todayStart.UnixMilli()

	log.Printf("当前时间: %s, 今天开始时间戳: %d", now.Format("2006-01-02 15:04:05"), todayStartTimestamp)

	successCount := 0
	skipCount := 0
	filteredCount := 0 // 被过滤的未结束K线数量

	for _, bk := range binanceKlines {
		if len(bk) < 6 {
			log.Printf("数据格式错误，跳过: %v", bk)
			continue
		}

		// 解析币安返回的数据
		// [timestamp, open, high, low, close, volume, ...]
		openTime := int64(bk[0].(float64))

		// 过滤未结束的K线：如果K线的开盘时间 >= 今天的开始时间，说明是今天或未来的K线，跳过
		if openTime >= todayStartTimestamp {
			filteredCount++
			continue
		}

		openPrice := parseFloat(bk[1])
		highPrice := parseFloat(bk[2])
		lowPrice := parseFloat(bk[3])
		closePrice := parseFloat(bk[4])
		volume := parseFloat(bk[5])

		// 创建Kline对象
		kline := Kline{
			Symbol:     symbol,
			Timeframe:  interval,
			OpenTime:   openTime,
			OpenTimeDt: time.UnixMilli(openTime),
			OpenPrice:  openPrice,
			HighPrice:  highPrice,
			LowPrice:   lowPrice,
			ClosePrice: closePrice,
			Volume:     volume,
		}

		// 使用 GORM 的 Clauses 来实现 ON CONFLICT DO UPDATE
		// 如果记录已存在（根据唯一索引），则更新数据
		result := DB.Where("symbol = ? AND timeframe = ? AND open_time = ?",
			symbol, interval, openTime).
			Assign(map[string]interface{}{
				"open_price":  openPrice,
				"high_price":  highPrice,
				"low_price":   lowPrice,
				"close_price": closePrice,
				"volume":      volume,
			}).
			FirstOrCreate(&kline)

		if result.Error != nil {
			log.Printf("保存K线数据失败: %v", result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			successCount++
		} else {
			skipCount++
		}
	}

	log.Printf("保存完成: 成功 %d 条, 跳过 %d 条（重复）, 过滤 %d 条（未结束）", successCount, skipCount, filteredCount)
	return nil
}

// parseFloat 辅助函数：将interface{}转换为float64
func parseFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case string:
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	default:
		return 0
	}
}

// FetchAndSaveKlines 获取并保存K线数据（组合函数）
// symbol: 交易对，如 BTCUSDT (注意：币安API不需要斜杠)
// interval: K线周期，如 1m, 5m, 1h, 1d
// limit: 获取数量
func FetchAndSaveKlines(symbol string, interval string, limit int) error {
	// 1. 从币安获取数据
	klines, err := FetchKlinesFromBinance(symbol, interval, limit)
	if err != nil {
		return err
	}

	// 2. 保存到数据库
	err = SaveKlinesToDB(symbol, interval, klines)
	if err != nil {
		return err
	}

	return nil
}

// GetKlinesFromDB 从数据库查询K线数据
func GetKlinesFromDB(symbol string, interval string, limit int) ([]Kline, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	var klines []Kline
	result := DB.Where("symbol = ? AND timeframe = ?", symbol, interval).
		Order("open_time DESC").
		Limit(limit).
		Find(&klines)

	if result.Error != nil {
		return nil, result.Error
	}

	return klines, nil
}
