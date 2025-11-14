# Web分析结果展示系统

## 概述

这是一个基于Web的交易策略分析结果展示系统，可以通过浏览器查看策略分析结果，避免每次都重新计算。

## 功能特性

### 1. 数据持久化
- 策略分析结果自动保存到PostgreSQL数据库
- 支持结果更新（upsert），避免重复数据
- 详细记录按结果ID关联存储

### 2. Web界面
- **策略一**：历史同期涨跌分析（跨年、跨月对比）
- **策略二**：小时级别涨跌分析（日内时段对比）
- 支持按交易对、时间周期筛选
- 美观的卡片式布局
- 可视化进度条展示上涨率
- 可展开查看详细历史记录

### 3. API接口
- `GET /api/strategy1` - 获取策略一结果
- `GET /api/strategy2` - 获取策略二结果
- 支持查询参数过滤

## 使用说明

### 1. 运行策略分析并保存结果

首先运行主程序，执行策略分析并保存结果到数据库：

```bash
# 构建主程序
go build -o trade main.go

# 运行策略分析
./trade
```

这会：
- 更新K线数据
- 运行策略一和策略二
- 自动将分析结果保存到数据库

### 2. 启动Web服务器

```bash
# 构建Web服务器
go build -o web_server ./cmd/web/main.go

# 启动服务器（默认端口8080）
./web_server

# 或指定端口
./web_server -port 9000
```

### 3. 访问Web界面

在浏览器中打开：
```
http://localhost:8080
```

## API使用示例

### 获取策略一结果

```bash
# 获取所有结果
curl http://localhost:8080/api/strategy1

# 按交易对筛选
curl http://localhost:8080/api/strategy1?symbol=BTCUSDT

# 按时间周期筛选
curl http://localhost:8080/api/strategy1?interval=1d

# 组合筛选
curl http://localhost:8080/api/strategy1?symbol=BTCUSDT&interval=1d
```

### 获取策略二结果

```bash
# 获取所有结果
curl http://localhost:8080/api/strategy2

# 按交易对筛选
curl http://localhost:8080/api/strategy2?symbol=BTCUSDT

# 按时间周期筛选
curl http://localhost:8080/api/strategy2?interval=1h

# 按小时筛选
curl http://localhost:8080/api/strategy2?hour=15

# 组合筛选
curl http://localhost:8080/api/strategy2?symbol=BTCUSDT&interval=1h&hour=15
```

## 数据库表结构

### 策略一结果表 (strategy1_results)
- symbol: 交易对
- interval: 时间周期
- analyze_day: 分析日期(MM-DD)
- total_count: 总样本数
- up_count: 上涨次数
- down_count: 下跌次数
- flat_count: 平盘次数
- up_rate: 上涨概率
- best_month: 最佳月份
- worst_month: 最差月份

### 策略一详细记录表 (strategy1_detail_records)
- result_id: 关联结果ID
- year: 年份
- open_price: 开盘价
- close_price: 收盘价
- price_diff: 价差
- is_up: 是否上涨

### 策略二结果表 (strategy2_results)
- symbol: 交易对
- interval: 时间周期
- hour: 小时(0-23)
- total_count: 总样本数
- up_count: 上涨次数
- down_count: 下跌次数
- flat_count: 平盘次数
- up_rate: 上涨概率

### 策略二详细记录表 (strategy2_detail_records)
- result_id: 关联结果ID
- date: 日期
- open_price: 开盘价
- close_price: 收盘价
- price_diff: 价差
- is_up: 是否上涨

## 工作流程

1. **数据采集** → 运行 `main.go` 更新K线数据
2. **策略分析** → 运行策略一和策略二，结果自动保存到数据库
3. **Web展示** → 运行 `web_server` 启动Web服务
4. **浏览查看** → 通过浏览器查看分析结果

## 优势

- ✅ 避免重复计算，提高效率
- ✅ 历史结果可追溯，支持对比
- ✅ Web界面友好，方便查看
- ✅ 支持多维度筛选
- ✅ 详细记录可展开查看
- ✅ 响应式设计，支持移动端

## 注意事项

1. 首次使用前需运行 `main.go` 生成分析数据
2. 数据库连接信息在 `db/postgreSql.go` 中配置
3. Web服务器默认端口为8080，可通过 `-port` 参数修改
4. 支持跨域请求(CORS)，可集成到其他前端项目
