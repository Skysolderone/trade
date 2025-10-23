# 币安永续合约K线数据保存功能

## 功能概述

本项目实现了从币安永续合约API获取K线数据并保存到PostgreSQL数据库的功能，支持：

- ✅ 多交易对支持（BTC/USDT、ETH/USDT等）
- ✅ 多时间周期（1m、5m、15m、1h、1d等）
- ✅ 自动去重（同一时间的K线只保存一次）
- ✅ 定时更新机制
- ✅ 数据查询接口

## 数据库表结构

### 表名：`klines`

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGSERIAL | 主键，自增 |
| symbol | VARCHAR(20) | 交易对，如 BTCUSDT |
| timeframe | VARCHAR(10) | K线周期，如 1m, 5m, 1h |
| open_time | BIGINT | 开盘时间戳（毫秒） |
| open_time_dt | TIMESTAMP | 开盘时间（可读格式） |
| open_price | DECIMAL(20,8) | 开盘价 |
| high_price | DECIMAL(20,8) | 最高价 |
| low_price | DECIMAL(20,8) | 最低价 |
| close_price | DECIMAL(20,8) | 收盘价 |
| volume | DECIMAL(20,8) | 成交量 |
| created_at | TIMESTAMP | 记录创建时间 |
| updated_at | TIMESTAMP | 记录更新时间 |

### 索引

- 唯一索引：`(symbol, timeframe, open_time)` - 防止重复数据
- 普通索引：`symbol`, `timeframe`, `open_time` - 提高查询性能
- 组合索引：`(symbol, timeframe, open_time DESC)` - 优化常用查询

## 核心函数说明

### 1. 初始化数据库

```go
import "trade/db"

func main() {
    db.InitPostgreSql()
}
```

### 2. 获取并保存K线数据

```go
// FetchAndSaveKlines(symbol, interval, limit)
// symbol: 交易对，如 "BTCUSDT" (注意：币安API不需要斜杠)
// interval: K线周期，如 "1m", "5m", "1h", "1d"
// limit: 获取数量，默认500，最大1500

err := db.FetchAndSaveKlines("BTCUSDT", "1m", 1440)
if err != nil {
    log.Fatal(err)
}
```

### 3. 从数据库查询K线数据

```go
// GetKlinesFromDB(symbol, interval, limit)
klines, err := db.GetKlinesFromDB("BTCUSDT", "1m", 100)
if err != nil {
    log.Fatal(err)
}

// 遍历K线数据
for _, kline := range klines {
    fmt.Printf("时间: %s, 开盘: %.2f, 收盘: %.2f\n",
        kline.OpenTimeDt.Format("2006-01-02 15:04:05"),
        kline.OpenPrice,
        kline.ClosePrice)
}
```

## 使用示例

### 示例1：基本使用

运行示例代码：

```bash
cd examples
go run kline_example.go
```

该示例会：
1. 初始化数据库连接
2. 获取BTC和ETH的最近24小时1分钟K线数据
3. 保存到数据库
4. 查询并打印前5条数据

### 示例2：定时任务

运行定时更新程序：

```bash
cd examples
go run kline_scheduler.go
```

该程序会：
1. 每分钟自动更新BTC、ETH、BNB的K线数据
2. 持续运行，直到按 Ctrl+C 停止

## 支持的K线周期

| 周期 | 说明 |
|------|------|
| 1m | 1分钟 |
| 3m | 3分钟 |
| 5m | 5分钟 |
| 15m | 15分钟 |
| 30m | 30分钟 |
| 1h | 1小时 |
| 2h | 2小时 |
| 4h | 4小时 |
| 6h | 6小时 |
| 8h | 8小时 |
| 12h | 12小时 |
| 1d | 1天 |
| 3d | 3天 |
| 1w | 1周 |
| 1M | 1月 |

## 常用交易对

- BTCUSDT - 比特币/USDT
- ETHUSDT - 以太坊/USDT
- BNBUSDT - 币安币/USDT
- ADAUSDT - 艾达币/USDT
- SOLUSDT - Solana/USDT
- DOGEUSDT - 狗狗币/USDT

完整交易对列表可访问：https://fapi.binance.com/fapi/v1/exchangeInfo

## API 限制

币安API有以下限制：

- **请求频率限制**：每分钟最多2400次请求
- **权重限制**：每分钟最多6000权重
- **K线数量限制**：单次请求最多1500条

建议：
- 定时任务间隔不要太短（建议≥1分钟）
- 避免同时发起大量请求
- 实现请求失败重试机制

## 数据更新策略

### 实时数据更新

对于需要实时数据的场景：

```go
// 每分钟更新最近100条1分钟K线
ticker := time.NewTicker(1 * time.Minute)
for range ticker.C {
    db.FetchAndSaveKlines("BTCUSDT", "1m", 100)
}
```

### 历史数据回填

如果需要回填历史数据：

```go
// 由于API限制单次最多1500条，需要分批获取
// 1分钟K线：1500条 ≈ 25小时
// 1小时K线：1500条 ≈ 62.5天

// 获取最近1500条（约25小时）
db.FetchAndSaveKlines("BTCUSDT", "1m", 1500)
```

## 故障排查

### 1. 数据库连接失败

检查 `db/postgreSql.go` 中的数据库配置：
- 主机地址
- 用户名密码
- 数据库名称
- 端口号

### 2. API请求失败

可能原因：
- 网络问题（检查是否能访问 fapi.binance.com）
- 请求频率超限（降低请求频率）
- 交易对名称错误（确保使用正确的symbol格式）

### 3. 数据重复

数据表有唯一索引，重复数据会自动更新而不是插入新记录。

## 项目文件结构

```
trade/
├── db/
│   ├── postgreSql.go        # 数据库连接初始化
│   ├── models.go            # K线数据模型
│   ├── kline_service.go     # K线数据获取和保存服务
│   └── migrations/
│       └── create_klines_table.sql  # 建表SQL
├── examples/
│   ├── kline_example.go     # 基本使用示例
│   └── kline_scheduler.go   # 定时任务示例
├── docs/
│   └── KLINE_README.md      # 本文档
└── main.go
```

## 下一步改进建议

1. **添加WebSocket实时推送** - 使用币安WebSocket API实时接收K线数据
2. **数据分析功能** - 添加技术指标计算（MA、MACD、RSI等）
3. **告警机制** - 价格突破、成交量异常等告警
4. **数据导出** - 支持导出CSV、Excel等格式
5. **可视化界面** - 使用Grafana或自建前端展示K线图表
6. **多交易所支持** - 扩展支持OKX、Bybit等交易所

## 许可证

MIT License
