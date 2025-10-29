package model

// SymbolConfig 交易对配置
type SymbolConfig struct {
	Symbol    string   `json:"symbol"`    // 交易对符号
	Intervals []string `json:"intervals"` // K线时间区间列表
}

// Config 全局配置
type Config struct {
	Symbols []SymbolConfig `json:"symbols"` // 交易对配置列表
}
