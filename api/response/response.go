package response

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 状态码定义
const (
	CodeSuccess       = 0    // 成功
	CodeParamError    = 1001 // 参数错误
	CodeDataNotFound  = 1002 // 数据不存在
	CodeSampleTooLow  = 1003 // 样本量不足
	CodeInternalError = 5000 // 服务器内部错误
)

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code      int         `json:"code"`      // 状态码
	Message   string      `json:"message"`   // 状态信息
	Timestamp int64       `json:"timestamp"` // Unix时间戳(秒)
	Data      interface{} `json:"data"`      // 业务数据
}

// StrategyInfo 策略信息
type StrategyInfo struct {
	StrategyType   string `json:"strategy_type"`   // 策略类型
	StrategyName   string `json:"strategy_name"`   // 策略名称
	Description    string `json:"description"`     // 策略描述
	AnalysisMethod string `json:"analysis_method"` // 分析方法
}

// DataStatistics 数据统计说明
type DataStatistics struct {
	DataSource       string    `json:"data_source"`       // 数据来源
	DateRange        DateRange `json:"date_range"`        // 数据时间范围
	TotalRecordsUsed int       `json:"total_records_used"` // 使用的K线记录总数
	QueryMethod      string    `json:"query_method"`      // 数据库查询方法
}

// DateRange 日期范围
type DateRange struct {
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
}

// AnalysisTarget 分析目标
type AnalysisTarget struct {
	Symbol           string `json:"symbol"`             // 交易对
	Interval         string `json:"interval"`           // K线周期
	AnalysisDate     string `json:"analysis_date,omitempty"`     // 分析日期(策略一)
	AnalysisDatetime string `json:"analysis_datetime,omitempty"` // 分析时间(策略二)
	TargetPeriod     string `json:"target_period,omitempty"`     // 目标周期(策略一)
	TargetHour       int    `json:"target_hour,omitempty"`       // 目标小时(策略二)
}

// PeriodResult 周期结果
type PeriodResult struct {
	PeriodLabel      string  `json:"period_label,omitempty"`      // 周期标签(策略一)
	HourLabel        string  `json:"hour_label,omitempty"`        // 小时标签(策略二)
	SampleCount      int     `json:"sample_count"`      // 样本数量
	UpCount          int     `json:"up_count"`          // 上涨次数
	DownCount        int     `json:"down_count"`        // 下跌次数
	FlatCount        int     `json:"flat_count"`        // 平盘次数
	UpRate           float64 `json:"up_rate"`           // 上涨概率
	DownRate         float64 `json:"down_rate"`         // 下跌概率
	Reliability      string  `json:"reliability"`       // 可靠性等级
	ReliabilityNote  string  `json:"reliability_note"`  // 可靠性说明
}

// Performance 表现数据
type Performance struct {
	Year             string  `json:"year,omitempty"`        // 年份(策略一)
	Month            int     `json:"month,omitempty"`       // 月份(策略一)
	Hour             int     `json:"hour,omitempty"`        // 小时(策略二)
	MonthLabel       string  `json:"month_label,omitempty"` // 月份标签(策略一)
	HourLabel        string  `json:"hour_label,omitempty"`  // 小时标签(策略二)
	UpRate           float64 `json:"up_rate"`           // 上涨率
	SampleCount      int     `json:"sample_count"`      // 样本数量
	UpCount          int     `json:"up_count"`          // 上涨次数
	DownCount        int     `json:"down_count"`        // 下跌次数
	Result           string  `json:"result,omitempty"`      // 结果(策略一)
	Rate             float64 `json:"rate,omitempty"`        // 概率(策略一)
	PerformanceNote  string  `json:"performance_note,omitempty"` // 表现说明
}

// CrossYearAnalysis 跨年对比分析(策略一)
type CrossYearAnalysis struct {
	Title            string       `json:"title"`             // 标题
	Description      string       `json:"description"`       // 描述
	YearsAnalyzed    int          `json:"years_analyzed"`    // 分析的年数
	OverallUpRate    float64      `json:"overall_up_rate"`   // 总体上涨率
	Trend            string       `json:"trend"`             // 趋势
	BestPerformance  *Performance `json:"best_performance"`  // 最佳表现
	WorstPerformance *Performance `json:"worst_performance"` // 最差表现
}

// CrossMonthAnalysis 跨月对比分析(策略一)
type CrossMonthAnalysis struct {
	Title              string       `json:"title"`               // 标题
	Description        string       `json:"description"`         // 描述
	MonthsAnalyzed     int          `json:"months_analyzed"`     // 分析的月数
	BestMonth          *Performance `json:"best_month"`          // 最佳月份
	WorstMonth         *Performance `json:"worst_month"`         // 最差月份
	CurrentMonthRanking *Ranking    `json:"current_month_ranking"` // 当前月份排名
}

// HourlyComparison 小时对比(策略二)
type HourlyComparison struct {
	Title              string       `json:"title"`               // 标题
	Description        string       `json:"description"`         // 描述
	HoursAnalyzed      int          `json:"hours_analyzed"`      // 分析的小时数
	AverageUpRate      float64      `json:"average_up_rate"`     // 平均上涨率
	BestHour           *Performance `json:"best_hour"`           // 最佳时段
	WorstHour          *Performance `json:"worst_hour"`          // 最差时段
	CurrentHourRanking *Ranking     `json:"current_hour_ranking"` // 当前小时排名
}

// Ranking 排名信息
type Ranking struct {
	CurrentMonth     int    `json:"current_month,omitempty"`     // 当前月份(策略一)
	Rank             int    `json:"rank"`                        // 排名
	TotalMonths      int    `json:"total_months,omitempty"`      // 总月数(策略一)
	TotalHours       int    `json:"total_hours,omitempty"`       // 总小时数(策略二)
	PerformanceLevel string `json:"performance_level"`           // 表现等级
	PerformanceNote  string `json:"performance_note"`            // 表现说明
}

// TimeZoneAnalysis 时区特征分析(策略二)
type TimeZoneAnalysis struct {
	Title          string         `json:"title"`           // 标题
	AsianSession   *SessionInfo   `json:"asian_session"`   // 亚洲时段
	EuropeanSession *SessionInfo  `json:"european_session"` // 欧洲时段
	AmericanSession *SessionInfo  `json:"american_session"` // 美洲时段
}

// SessionInfo 时段信息
type SessionInfo struct {
	Hours          string  `json:"hours"`          // 时间范围
	AverageUpRate  float64 `json:"average_up_rate"` // 平均上涨率
	Characteristic string  `json:"characteristic"`  // 特征
}

// TradingRecommendation 交易建议
type TradingRecommendation struct {
	Signal              string   `json:"signal"`                // 交易信号
	ConfidenceLevel     string   `json:"confidence_level"`      // 置信度等级
	ConfidenceScore     float64  `json:"confidence_score"`      // 置信度分数
	MainReason          string   `json:"main_reason"`           // 主要原因
	SupportingFactors   []string `json:"supporting_factors"`    // 支持因素
	RiskFactors         []string `json:"risk_factors,omitempty"` // 风险因素
	OptimalTradingHours []string `json:"optimal_trading_hours,omitempty"` // 最佳交易时段(策略二)
	AvoidTradingHours   []string `json:"avoid_trading_hours,omitempty"`   // 避免交易时段(策略二)
}

// RiskWarning 风险警告
type RiskWarning struct {
	Level    string   `json:"level"`    // 风险级别
	Warnings []string `json:"warnings"` // 警告信息列表
}

// Strategy1Response 策略一响应数据
type Strategy1Response struct {
	StrategyInfo          *StrategyInfo          `json:"strategy_info"`           // 策略信息
	AnalysisTarget        *AnalysisTarget        `json:"analysis_target"`         // 分析目标
	DataStatistics        *DataStatistics        `json:"data_statistics"`         // 数据统计
	CurrentPeriodResult   *PeriodResult          `json:"current_period_result"`   // 当前周期结果
	CrossYearAnalysis     *CrossYearAnalysis     `json:"cross_year_analysis"`     // 跨年对比分析
	CrossMonthAnalysis    *CrossMonthAnalysis    `json:"cross_month_analysis"`    // 跨月对比分析
	TradingRecommendation *TradingRecommendation `json:"trading_recommendation"`  // 交易建议
	RiskWarning           *RiskWarning           `json:"risk_warning"`            // 风险警告
}

// Strategy2Response 策略二响应数据
type Strategy2Response struct {
	StrategyInfo          *StrategyInfo          `json:"strategy_info"`           // 策略信息
	AnalysisTarget        *AnalysisTarget        `json:"analysis_target"`         // 分析目标
	DataStatistics        *DataStatistics        `json:"data_statistics"`         // 数据统计
	CurrentHourResult     *PeriodResult          `json:"current_hour_result"`     // 当前小时结果
	HourlyComparison      *HourlyComparison      `json:"hourly_comparison"`       // 小时对比
	HighWinHours          *HighWinHours          `json:"high_win_hours"`          // 高胜率时段
	LowWinHours           *LowWinHours           `json:"low_win_hours"`           // 低胜率时段
	TimeZoneAnalysis      *TimeZoneAnalysis      `json:"time_zone_analysis,omitempty"` // 时区特征分析
	TradingRecommendation *TradingRecommendation `json:"trading_recommendation"`  // 交易建议
	RiskWarning           *RiskWarning           `json:"risk_warning"`            // 风险警告
}

// HighWinHours 高胜率时段
type HighWinHours struct {
	Title       string         `json:"title"`       // 标题
	Description string         `json:"description"` // 描述
	Count       int            `json:"count"`       // 数量
	Hours       []*Performance `json:"hours"`       // 时段列表
}

// LowWinHours 低胜率时段
type LowWinHours struct {
	Title       string         `json:"title"`       // 标题
	Description string         `json:"description"` // 描述
	Count       int            `json:"count"`       // 数量
	Hours       []*Performance `json:"hours"`       // 时段列表
}

// Success 成功响应
func Success(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, &BaseResponse{
		Code:      CodeSuccess,
		Message:   "成功",
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

// Error 错误响应
func Error(c *app.RequestContext, code int, message string) {
	c.JSON(consts.StatusOK, &BaseResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Data:      nil,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *app.RequestContext, code int, message string, data interface{}) {
	c.JSON(consts.StatusOK, &BaseResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

// ParamError 参数错误
func ParamError(c *app.RequestContext, message string) {
	Error(c, CodeParamError, message)
}

// DataNotFound 数据不存在
func DataNotFound(c *app.RequestContext, message string) {
	Error(c, CodeDataNotFound, message)
}

// SampleTooLow 样本量不足
func SampleTooLow(c *app.RequestContext, message string, data interface{}) {
	ErrorWithData(c, CodeSampleTooLow, message, data)
}

// InternalError 服务器内部错误
func InternalError(c *app.RequestContext, message string) {
	Error(c, CodeInternalError, message)
}
