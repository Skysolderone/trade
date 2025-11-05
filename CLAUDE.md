# CLAUDE.md

这个文件为 Claude Code (claude.ai/code) 提供在此代码库中工作的指导。

## 项目概述

这是一个基于 Go 的加密货币策略分析系统,用于币安合约交易的 K 线数据采集、存储和策略分析。系统支持多种运行模式,包括定时任务、Web界面和API服务。

## 项目架构

### 核心组件

项目包含三个可独立运行的程序:

1. **主程序** (`main.go` / `bin/trade`)
   - 单次运行模式 (`-mode=once`): 更新K线数据并运行策略分析
   - 定时任务模式 (`-mode=daemon`): 每天00:00:00自动执行策略更新
   - 使用 `robfig/cron/v3` 实现秒级精度的定时调度

2. **Web服务器** (`cmd/web/main.go` / `bin/web_server`)
   - 提供策略分析结果的Web界面展示
   - 包含策略一(历史同期分析)和策略二(小时级别分析)的可视化
   - 默认端口: 8080

3. **API服务器** (`cmd/api/main.go`)
   - 使用 CloudWeGo Hertz 框架提供RESTful API
   - 支持实时策略分析查询
   - 端点: `/api/v1/strategy/analyze`

### 目录结构

- `api/` - HTTP API服务
  - `handler/` - 请求处理器
  - `router/` - 路由注册
  - `response/` - 响应格式化
- `cmd/` - 可执行程序入口
  - `api/` - API服务器
  - `web/` - Web服务器
  - `get_hourly_history/` - 小时级历史数据获取工具
  - `fix_old_data/` - 数据修复工具
- `db/` - 数据库层
  - `postgreSql.go` - PostgreSQL连接管理
  - `binance.go` - 币安API客户端初始化
- `kline/` - K线数据处理
  - `rest.go` - REST API批量获取历史数据
  - `ws.go` - WebSocket实时数据流
- `model/` - 数据模型
  - `kline.go` - K线数据模型
  - `strategy_result.go` - 策略结果模型
  - `config.go` - 配置文件模型
- `scheduler/` - 定时任务调度器
- `strategy/` - 策略实现
  - `strategy_1.go` - 历史同期涨跌分析(跨年/跨月对比)
  - `strategy_2.go` - 小时级别涨跌分析(日内时段对比)
- `web/` - Web服务器实现
- `scripts/` - 构建和部署脚本
- `test/` - 测试文件

### 数据流程

1. **数据更新**: `kline.UpdateKline()` 查询数据库最新记录,增量更新至昨日
2. **时间维度解析**: K线数据自动解析为 `date`(YYYY-MM-DD), `day`(MM-DD), `hour`, `week`, `min` 字段
3. **策略执行**:
   - 策略一: 按 `day` 字段分析历史同期数据,支持跨年和跨月对比
   - 策略二: 按 `hour` 字段分析24小时涨跌规律
4. **结果存储**: 策略结果保存到 `strategy_1_results` 和 `strategy_2_results` 表,包含详细记录

### 关键技术依赖

- `github.com/adshao/go-binance/v2` - 币安API客户端
- `gorm.io/gorm` - ORM框架,自动迁移表结构
- `github.com/robfig/cron/v3` - 定时任务调度(秒级精度)
- `github.com/cloudwego/hertz` - 高性能HTTP框架

## 常用命令

### 开发和测试

```bash
# 运行主程序(单次模式)
go run main.go -mode=once

# 运行主程序(定时任务模式)
go run main.go -mode=daemon

# 定时任务模式并立即执行一次
go run main.go -mode=daemon -now

# 运行Web服务器
go run cmd/web/main.go -port=8080

# 运行API服务器
go run cmd/api/main.go

# 运行特定策略的测试
go test -v ./strategy/strategy1_test.go
go test -v ./strategy/strategy2_test.go

# 运行单个测试函数
go test -v ./strategy -run TestStrategy1
```

### 构建和部署

```bash
# 构建所有程序(生成Linux二进制文件)
./scripts/build.sh

# 构建后会在 bin/ 目录生成:
# - trade (主程序)
# - web_server (Web服务器)
# - get_hourly_history (历史数据获取工具)

# 部署到服务器(需要root权限)
sudo ./scripts/deploy.sh

# 查看服务状态
sudo systemctl status trade.service
sudo systemctl status trade-web.service

# 查看服务日志
sudo journalctl -u trade.service -f
sudo journalctl -u trade-web.service -f
```

### 配置文件

`config.json` 定义监控的交易对和时间周期:
```json
{
  "symbols": [
    {
      "symbol": "BTCUSDT",
      "intervals": ["1d", "8h", "4h", "2h", "1h"]
    }
  ]
}
```

## 策略说明

### 策略一: 历史同期涨跌分析 (strategy_1.go)

**分析维度**:
1. **跨年对比**: 所有年份相同日期(如 2018-10-30, 2019-10-30, 2020-10-30...)
2. **跨月对比**: 所有月份相同日期(如 01-30, 02-30, 03-30...)

**输出内容**:
- 上涨概率、涨跌次数统计
- 每年具体的开盘价、收盘价、价差
- 最佳/最差月份排名
- 当前月份在所有月份中的排名

**数据存储**: `strategy_1_results` 表(汇总)+ `strategy_1_detail_records` 表(明细)

### 策略二: 小时级别涨跌分析 (strategy_2.go)

**分析维度**: 24小时(0-23点)各时段的历史涨跌规律

**输出内容**:
- 每个小时的上涨概率
- 24小时可视化柱状图
- 最佳/最差交易时段
- 高胜率时段推荐(上涨率≥60%且样本≥10)

**数据存储**: `strategy_2_results` 表(汇总)+ `strategy_2_detail_records` 表(明细)

## K线数据更新逻辑

`kline/rest.go` 中的 `UpdateKline()` 实现智能增量更新:
1. 查询数据库最新记录的 `close_time`
2. 删除最新记录(可能不完整)
3. 从最新记录的 `open_time` 重新获取到昨日23:59:59
4. 每批次最多1500条,自动分页
5. 批次间延迟100ms避免API限流
6. 使用 `clause.OnConflict{DoNothing: true}` 避免重复数据

**时间维度解析** (`utils/parse.go`):
- `date`: YYYY-MM-DD (用于跨年对比)
- `day`: MM-DD (用于跨月对比)
- `hour`: 0-23 (用于小时分析)
- `week`: 1-7 (预留)

## 数据库配置

配置在 `db/postgreSql.go`:
- 默认使用 PostgreSQL
- GORM自动迁移表结构(`AutoMigrate`)
- 连接池: 最大空闲连接10,最大打开连接100
- 所有时间使用UTC时区

**主要表结构**:
- `klines` - K线原始数据
- `strategy_1_results` - 策略一汇总结果
- `strategy_1_detail_records` - 策略一详细记录
- `strategy_2_results` - 策略二汇总结果
- `strategy_2_detail_records` - 策略二详细记录

## 开发注意事项

1. **币安API初始化**: `db.InitBinance("", "")` 使用空密钥(仅需公开API)
2. **增量更新机制**: `UpdateKline()` 会删除最新记录重新获取,避免不完整数据
3. **定时调度**: cron表达式 `"0 0 0 * * *"` 使用秒级精度(第一个0是秒)
4. **构建目标**: `scripts/build.sh` 生成Linux amd64二进制文件,禁用CGO
5. **策略结果更新**: 使用 `FirstOrCreate` + `Assign` 实现upsert,避免重复记录
6. **时区一致性**: 所有时间处理使用 `time.UTC`,避免时区混淆
