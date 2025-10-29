# CLAUDE.md

这个文件为 Claude Code (claude.ai/code) 提供在此代码库中工作的指导。

## 项目概述

这是一个基于 Go 的加密货币自动策略交易系统 (Ws_trade),主要用于币安合约交易的 K 线数据采集、存储和策略分析。

## 核心功能

1. **K 线数据获取**: 通过币安合约 REST API 批量获取历史 K 线数据
2. **WebSocket 数据流**: 订阅币安合约 WebSocket 实时接收 K 线数据更新
3. **数据存储**: 使用 PostgreSQL 数据库存储所有 K 线数据
4. **策略分析**: 基于历史数据分析特定日期的涨跌概率

## 架构说明

### 目录结构

- `db/` - 数据库连接层,使用 GORM 连接 PostgreSQL
- `kline/` - K 线数据处理
  - `rest.go` - 通过 REST API 批量获取历史数据
  - `ws.go` - WebSocket 实时数据流处理
- `model/` - 数据模型定义
- `strategy/` - 交易策略实现
- `utils/` - 工具函数
- `test/` - 测试文件
- `sql/` - 数据库表结构定义

### 数据流程

1. 程序启动时通过 `db.InitPostgreSql()` 初始化数据库连接
2. 使用 `kline.GetKline()` 批量获取历史 K 线数据并存储到数据库
3. K 线数据包含时间维度字段 (date, day, hour, week, min) 用于策略分析
4. 策略函数从数据库查询特定时间维度的数据进行分析

### 关键依赖

- `github.com/adshao/go-binance/v2` - 币安 API 客户端
- `gorm.io/gorm` - ORM 框架
- `gorm.io/driver/postgres` - PostgreSQL 驱动

## 常用命令

### 运行程序
```bash
go run main.go
```

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./test
go test ./strategy

# 运行特定测试文件
go test ./test/db_test.go
go test ./test/kline_test.go
```

### 构建
```bash
go build -o trade.exe
```

### 依赖管理
```bash
# 下载依赖
go mod download

# 清理未使用的依赖
go mod tidy
```

## K 线数据获取逻辑

`kline/rest.go` 中的 `getAllKlines()` 函数实现了分页获取历史数据:
- 默认从 2018-01-01 开始获取数据
- 每次请求最多 1500 条数据
- 自动计算批次结束时间,避免超过 API 限制
- 每批次之间有 100ms 延迟避免触发 API 频率限制
- 数据存储到 PostgreSQL 时会解析时间维度字段 (year, day, hour, week, min)

## 策略说明

`strategy/strategy_1.go` 实现了基于历史同期数据的涨跌概率分析:
- 根据当前日期 (月-日) 查询历史所有年份该日期的 K 线数据
- 统计涨跌次数计算概率
- 输出每个历史年份该日期的具体涨跌幅度

## 数据库配置

数据库连接信息在 `db/postgreSql.go` 中配置:
- 使用阿里云 RDS PostgreSQL
- GORM 自动迁移表结构
- 连接池配置: 最大空闲连接 10,最大打开连接 100

## 开发注意事项

1. **API 密钥配置**: `ccxt-accounts.json` 文件包含币安 API 密钥配置,实际部署时需要替换真实密钥
2. **时区处理**: 所有时间数据使用 UTC 时区存储
3. **数据去重**: 当前实现没有数据去重逻辑,重复运行会插入重复数据
4. **错误处理**: K 线数据插入失败时会跳过该条继续处理下一条
