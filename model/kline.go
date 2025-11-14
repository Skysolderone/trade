package model

import (
	"time"
)

// Kline 表示K线数据
type Kline struct {
	ID        int       `json:"id" db:"id"`                                                    // 主键ID
	Symbol    string    `json:"symbol" db:"symbol" gorm:"index:idx_unique_kline,unique"`       // 交易对符号
	Interval  string    `json:"interval" db:"interval" gorm:"index:idx_unique_kline,unique"`   // 时间周期(1m,1h,1d等)
	Open      float64   `json:"open" db:"open"`                                                // 开盘价
	Close     float64   `json:"close" db:"close"`                                              // 收盘价
	High      float64   `json:"high" db:"high"`                                                // 最高价
	Low       float64   `json:"low" db:"low"`                                                  // 最低价
	OpenTime  time.Time `json:"open_time" db:"open_time" gorm:"index:idx_unique_kline,unique"` // 开盘时间
	CloseTime time.Time `json:"close_time" db:"close_time"`                                    // 收盘时间
	Date      string    `json:"date" db:"date"`                                                // 日期字符串
	Day       string    `json:"day" db:"day"`                                                  // 日
	Hour      string    `json:"hour" db:"hour"`                                                // 小时
	Week      string    `json:"week" db:"week"`                                                // 周
	Min       string    `json:"min" db:"min"`                                                  // 分钟
}

type KlineWs struct {
	ID        int       `json:"id" db:"id"`                                                    // 主键ID
	Symbol    string    `json:"symbol" db:"symbol" gorm:"index:idx_unique_kline,unique"`       // 交易对符号
	Interval  string    `json:"interval" db:"interval" gorm:"index:idx_unique_kline,unique"`   // 时间周期(1m,1h,1d等)
	Open      float64   `json:"open" db:"open"`                                                // 开盘价
	Close     float64   `json:"close" db:"close"`                                              // 收盘价
	High      float64   `json:"high" db:"high"`                                                // 最高价
	Low       float64   `json:"low" db:"low"`                                                  // 最低价
	OpenTime  time.Time `json:"open_time" db:"open_time" gorm:"index:idx_unique_kline,unique"` // 开盘时间
	CloseTime time.Time `json:"close_time" db:"close_time"`                                    // 收盘时间
	Date      string    `json:"date" db:"date"`                                                // 日期字符串
	Day       string    `json:"day" db:"day"`                                                  // 日
	Hour      string    `json:"hour" db:"hour"`                                                // 小时
	Week      string    `json:"week" db:"week"`                                                // 周
	Min       string    `json:"min" db:"min"`                                                  // 分钟
}
