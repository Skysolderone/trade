# 策略分析API响应格式说明

## 概述

本文档说明币安合约交易策略分析系统的API响应格式，提供给移动端客户端调用。

## 文件说明

- `api_response_format.json` - 策略一(历史同期涨跌分析)的响应格式
- `api_response_format_strategy2.json` - 策略二(小时级别涨跌分析)的响应格式

## 通用响应结构

```json
{
  "code": 0,              // 状态码：0-成功，非0-失败
  "message": "string",    // 状态信息
  "timestamp": 1698764800, // Unix时间戳(秒)
  "data": {}              // 具体业务数据
}
```

### 状态码说明

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 数据不存在 |
| 1003 | 样本量不足 |
| 5000 | 服务器内部错误 |

## 策略一：历史同期涨跌分析

### 分析原理

通过查询历史K线数据，统计特定日期（如10月30日）在历年的涨跌表现，并进行跨年、跨月对比分析。

### 数据统计方法

1. **数据来源**：PostgreSQL数据库中的Kline表
2. **查询条件**：
   - `symbol`：交易对（如BTCUSDT）
   - `interval`：K线周期（如1d）
   - `day`：日期字段，格式为"MM-DD"（如"10-30"）
3. **统计维度**：
   - 跨年对比：查询所有年份的同一日期（2018-10-30, 2019-10-30...）
   - 跨月对比：查询所有月份的相同日号（01-30, 02-30, ..., 12-30）

### 关键字段说明

#### strategy_info (策略信息)
- `strategy_type`: 策略类型标识
- `strategy_name`: 策略名称
- `description`: 策略描述
- `analysis_method`: 数据分析方法说明

#### data_statistics (数据统计说明)
- `data_source`: 数据来源
- `date_range`: 数据时间范围
- `total_records_used`: 使用的K线记录总数
- `query_method`: 数据库查询方法

#### current_period_result (当前周期结果)
- `sample_count`: 样本数量（历年该日期的K线数量）
- `up_count`: 上涨次数（收盘价>开盘价）
- `down_count`: 下跌次数（收盘价<开盘价）
- `flat_count`: 平盘次数（收盘价=开盘价）
- `up_rate`: 上涨概率（百分比）
- `reliability`: 可靠性等级
  - `high`: 样本数≥10
  - `medium`: 5≤样本数<10
  - `low`: 样本数<5

#### cross_year_analysis (跨年对比分析)
统计历年同一日期的整体表现，找出表现最好和最差的年份

#### cross_month_analysis (跨月对比分析)
对比所有月份的相同日号，分析月份间的差异

#### trading_recommendation (交易建议)
- `signal`: 交易信号
  - `bullish`: 看涨（上涨率>55%）
  - `bearish`: 看跌（上涨率<45%）
  - `neutral`: 中性（45%≤上涨率≤55%）
- `confidence_level`: 置信度等级
  - `high`: 样本量充足且趋势明显
  - `medium`: 样本量适中或趋势一般
  - `low`: 样本量不足或趋势不明显

## 策略二：小时级别涨跌分析

### 分析原理

分析特定小时（0-23）的历史涨跌表现，通过24小时对比找出最佳和最差交易时段。

### 数据统计方法

1. **数据来源**：PostgreSQL数据库中的Kline表
2. **查询条件**：
   - `symbol`：交易对（如BTCUSDT）
   - `interval`：K线周期（建议1h, 2h, 4h, 8h）
   - `hour`：小时字段，值为0-23
3. **统计维度**：
   - 当前小时分析：统计特定小时的历史表现
   - 24小时对比：对比所有小时时段的表现
   - 时区特征分析：按亚洲、欧洲、美洲交易时段分析

### 关键字段说明

#### current_hour_result (当前小时结果)
- `hour_label`: 小时标签（如"14:00"）
- `sample_count`: 该小时的历史K线样本数
- `up_rate`: 该小时的历史上涨概率

#### hourly_comparison (小时对比)
- `best_hour`: 24小时中表现最好的时段
- `worst_hour`: 24小时中表现最差的时段
- `current_hour_ranking`: 当前小时的排名情况

#### high_win_hours (高胜率时段)
上涨率≥60%且样本数≥10的时段列表，建议优先交易

#### low_win_hours (低胜率时段)
上涨率<45%的时段列表，建议谨慎交易或避免交易

#### time_zone_analysis (时区特征分析)
按全球主要交易时段划分：
- 亚洲时段：00:00-08:00 UTC
- 欧洲时段：08:00-16:00 UTC
- 美洲时段：16:00-24:00 UTC

## 使用建议

### 移动端展示建议

1. **概览页面**
   - 显示当前信号（看涨/看跌/中性）
   - 显示置信度和主要原因
   - 显示风险警告

2. **详细分析页面**
   - 当前周期统计数据（饼图或柱状图）
   - 跨年/跨月对比（折线图）
   - 最佳/最差时段对比

3. **推荐建议页面**
   - 交易信号和置信度
   - 支持因素列表
   - 风险因素列表
   - 优化交易时段（策略二）

### 数据缓存建议

- 策略一：每日更新一次（每天0点UTC）
- 策略二：每小时更新一次
- 建议客户端缓存数据，减少API调用

### 风险提示

所有策略分析结果仅基于历史数据统计，不构成投资建议。使用时请注意：

1. 历史表现不代表未来走势
2. 样本量不足时统计结果可能不可靠
3. 需结合实时行情、技术指标、基本面等综合判断
4. 加密货币市场波动大，请控制风险和仓位

## 示例用法

### 请求示例

```bash
# 策略一：查询10月30日的历史表现
GET /api/v1/strategy/analyze
{
  "strategy_type": "strategy_1",
  "symbol": "BTCUSDT",
  "interval": "1d",
  "date": "2024-10-30"
}

# 策略二：查询14点的历史表现
GET /api/v1/strategy/analyze
{
  "strategy_type": "strategy_2",
  "symbol": "BTCUSDT",
  "interval": "1h",
  "hour": 14
}
```

### 响应示例

参见 `api_response_format.json` 和 `api_response_format_strategy2.json`

## 技术实现参考

### Go代码实现位置

- 策略一实现：`strategy/strategy_1.go`
- 策略二实现：`strategy/strategy_2.go`
- 数据模型：`model/kline.go`
- 数据库层：`db/postgreSql.go`

### 数据库表结构

K线表关键字段：
- `symbol`: 交易对
- `interval`: K线周期
- `open_time`: 开盘时间
- `open`: 开盘价
- `close`: 收盘价
- `date`: 日期字符串（YYYY-MM-DD）
- `day`: 日期字符串（MM-DD）
- `hour`: 小时字符串（0-23）
- `week`: 周几（1-7）

## 更新日志

- 2024-10-31: 初始版本，包含策略一和策略二的响应格式
