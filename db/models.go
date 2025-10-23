package db

import (
	"time"
)

// Kline K线数据模型
type Kline struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol     string    `gorm:"type:varchar(20);not null;index" json:"symbol"`              // 交易对
	Timeframe  string    `gorm:"type:varchar(10);not null;index" json:"timeframe"`           // K线周期
	OpenTime   int64     `gorm:"not null;index" json:"open_time"`                            // 开盘时间戳（毫秒）
	OpenTimeDt time.Time `gorm:"type:timestamp;not null" json:"open_time_dt"`                // 开盘时间（可读）
	OpenPrice  float64   `gorm:"type:decimal(20,8);not null" json:"open_price"`              // 开盘价
	HighPrice  float64   `gorm:"type:decimal(20,8);not null" json:"high_price"`              // 最高价
	LowPrice   float64   `gorm:"type:decimal(20,8);not null" json:"low_price"`               // 最低价
	ClosePrice float64   `gorm:"type:decimal(20,8);not null" json:"close_price"`             // 收盘价
	Volume     float64   `gorm:"type:decimal(20,8);not null" json:"volume"`                  // 成交量
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`                           // 创建时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`                           // 更新时间
}

// TableName 指定表名
func (Kline) TableName() string {
	return "klines_day"
}
