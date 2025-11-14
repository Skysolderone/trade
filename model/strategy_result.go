package model

import (
	"time"
)

// Strategy1Result 策略一分析结果表
type Strategy1Result struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Symbol     string    `json:"symbol" gorm:"index:idx_strategy1_unique,unique"`      // 交易对
	Interval   string    `json:"interval" gorm:"index:idx_strategy1_unique,unique"`    // 时间周期
	AnalyzeDay string    `json:"analyze_day" gorm:"index:idx_strategy1_unique,unique"` // 分析日期(MM-DD)
	Month      int       `json:"month"`                                                 // 月份
	Day        int       `json:"day"`                                                   // 日
	TotalCount int       `json:"total_count"`                                           // 总样本数
	UpCount    int       `json:"up_count"`                                              // 上涨次数
	DownCount  int       `json:"down_count"`                                            // 下跌次数
	FlatCount  int       `json:"flat_count"`                                            // 平盘次数
	UpRate     float64   `json:"up_rate"`                                               // 上涨概率
	BestMonth  int       `json:"best_month"`                                            // 最佳月份
	BestUpRate float64   `json:"best_up_rate"`                                          // 最佳月份上涨率
	WorstMonth int       `json:"worst_month"`                                           // 最差月份
	WorstUpRate float64   `json:"worst_up_rate"`                                        // 最差月份上涨率
	CreatedAt  time.Time `json:"created_at"`                                            // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                            // 更新时间
}

// Strategy1DetailRecord 策略一详细记录表
type Strategy1DetailRecord struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ResultID  int       `json:"result_id" gorm:"index"`  // 关联Strategy1Result的ID
	Year      string    `json:"year"`                    // 年份
	OpenPrice float64   `json:"open_price"`              // 开盘价
	ClosePrice float64  `json:"close_price"`             // 收盘价
	PriceDiff float64   `json:"price_diff"`              // 价差
	IsUp      bool      `json:"is_up"`                   // 是否上涨
	CloseTime time.Time `json:"close_time"`              // 收盘时间
	CreatedAt time.Time `json:"created_at"`              // 创建时间
}

// Strategy2Result 策略二分析结果表(小时级别)
type Strategy2Result struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Symbol     string    `json:"symbol" gorm:"index:idx_strategy2_unique,unique"`   // 交易对
	Interval   string    `json:"interval" gorm:"index:idx_strategy2_unique,unique"` // 时间周期
	Hour       int       `json:"hour" gorm:"index:idx_strategy2_unique,unique"`     // 小时(0-23)
	TotalCount int       `json:"total_count"`                                        // 总样本数
	UpCount    int       `json:"up_count"`                                           // 上涨次数
	DownCount  int       `json:"down_count"`                                         // 下跌次数
	FlatCount  int       `json:"flat_count"`                                         // 平盘次数
	UpRate     float64   `json:"up_rate"`                                            // 上涨概率
	CreatedAt  time.Time `json:"created_at"`                                         // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                         // 更新时间
}

// Strategy2DetailRecord 策略二详细记录表
type Strategy2DetailRecord struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	ResultID   int       `json:"result_id" gorm:"index"` // 关联Strategy2Result的ID
	Date       string    `json:"date"`                   // 日期
	OpenPrice  float64   `json:"open_price"`             // 开盘价
	ClosePrice float64   `json:"close_price"`            // 收盘价
	PriceDiff  float64   `json:"price_diff"`             // 价差
	IsUp       bool      `json:"is_up"`                  // 是否上涨
	CloseTime  time.Time `json:"close_time"`             // 收盘时间
	CreatedAt  time.Time `json:"created_at"`             // 创建时间
}
