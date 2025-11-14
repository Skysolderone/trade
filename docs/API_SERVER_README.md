# Trade Strategy API Server

基于 Hertz 框架构建的加密货币策略分析 API 服务。

## 快速开始

### 1. 启动 API 服务

```bash
# 进入项目目录
cd /home/wws/trade

# 运行 API 服务器
go run cmd/api/main.go
```

服务将在 `http://localhost:8080` 启动。

### 2. 测试接口

#### 健康检查

```bash
curl http://localhost:8080/health
```

响应：
```json
{
  "status": "ok"
}
```

#### 获取 API 信息

```bash
curl http://localhost:8080/
```

响应：
```json
{
  "name": "Trade Strategy API",
  "version": "1.0.0",
  "endpoints": [
    "GET  /health",
    "GET  /api/v1/strategy/analyze",
    "POST /api/v1/strategy/analyze"
  ]
}
```

## API 接口详情

### 策略分析接口

**接口地址**: `/api/v1/strategy/analyze`

**请求方法**: `GET` 或 `POST`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| strategy_type | string | 是 | 策略类型 | strategy_1 或 strategy_2 |
| symbol | string | 是 | 交易对 | BTCUSDT |
| interval | string | 是 | K线周期 | 1d, 1h, 4h 等 |
| date | string | 否 | 日期(策略一) | 2024-10-30 |
| hour | int | 否 | 小时(策略二) | 14 (0-23) |

**interval 支持的值**:
- `1m`, `5m`, `15m`, `30m` (分钟级)
- `1h`, `2h`, `4h`, `8h` (小时级)
- `1d` (日线)
- `1w` (周线)

## 使用示例

### 策略一：历史同期涨跌分析

分析特定日期在历年的涨跌表现。

#### GET 请求示例

```bash
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d&date=2024-10-30'
```

#### POST 请求示例

```bash
curl -X POST http://localhost:8080/api/v1/strategy/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "strategy_type": "strategy_1",
    "symbol": "BTCUSDT",
    "interval": "1d",
    "date": "2024-10-30"
  }'
```

#### 响应示例

```json
{
  "code": 0,
  "message": "成功",
  "timestamp": 1698764800,
  "data": {
    "strategy_info": {
      "strategy_type": "strategy_1",
      "strategy_name": "历史同期涨跌分析",
      "description": "基于历史K线数据，统计特定日期在历年的涨跌表现，通过跨年、跨月对比分析价格走势规律",
      "analysis_method": "从PostgreSQL数据库查询历史K线数据，按日期(day)字段分组统计涨跌次数和概率"
    },
    "analysis_target": {
      "symbol": "BTCUSDT",
      "interval": "1d",
      "analysis_date": "2024-10-30",
      "target_period": "10-30"
    },
    "current_period_result": {
      "period_label": "10月30日",
      "sample_count": 7,
      "up_count": 5,
      "down_count": 2,
      "flat_count": 0,
      "up_rate": 71.43,
      "down_rate": 28.57,
      "reliability": "medium",
      "reliability_note": "样本数量适中，统计结果具有一定参考价值"
    },
    "cross_year_analysis": {
      "title": "跨年对比",
      "description": "分析历年同一日期(10-30)的涨跌情况",
      "years_analyzed": 7,
      "overall_up_rate": 71.43,
      "trend": "bullish"
    },
    "cross_month_analysis": {
      "title": "跨月对比",
      "description": "对比所有月份的30号，找出历史表现最好和最差的月份"
    },
    "trading_recommendation": {
      "signal": "bullish",
      "confidence_level": "medium",
      "confidence_score": 71.43,
      "main_reason": "历史数据显示上涨概率为71.43%"
    },
    "risk_warning": {
      "level": "medium",
      "warnings": [
        "历史数据不代表未来表现，仅供参考",
        "请结合实时行情、技术指标、资金管理等综合判断"
      ]
    }
  }
}
```

### 策略二：小时级别涨跌分析

分析特定小时的历史涨跌表现，找出最佳交易时段。

#### GET 请求示例

```bash
# 分析14点的历史表现
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h&hour=14'

# 不指定hour参数，默认使用当前小时
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h'
```

#### POST 请求示例

```bash
curl -X POST http://localhost:8080/api/v1/strategy/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "strategy_type": "strategy_2",
    "symbol": "BTCUSDT",
    "interval": "1h",
    "hour": 14
  }'
```

#### 响应示例

```json
{
  "code": 0,
  "message": "成功",
  "timestamp": 1698764800,
  "data": {
    "strategy_info": {
      "strategy_type": "strategy_2",
      "strategy_name": "小时级别涨跌分析",
      "description": "分析特定小时的历史涨跌表现，通过24小时对比找出最佳交易时段"
    },
    "current_hour_result": {
      "hour_label": "14:00",
      "sample_count": 120,
      "up_count": 78,
      "down_count": 40,
      "flat_count": 2,
      "up_rate": 65.0,
      "reliability": "high"
    },
    "hourly_comparison": {
      "title": "24小时对比",
      "hours_analyzed": 24,
      "average_up_rate": 52.5
    },
    "high_win_hours": {
      "title": "高胜率时段推荐",
      "description": "上涨率≥60%且样本数≥10的时段",
      "hours": [
        {
          "hour_label": "14:00",
          "up_rate": 68.5,
          "sample_count": 120
        }
      ]
    },
    "trading_recommendation": {
      "signal": "bullish",
      "confidence_level": "high",
      "optimal_trading_hours": ["14:00", "09:00", "21:00"],
      "avoid_trading_hours": ["03:00", "06:00"]
    }
  }
}
```

## 状态码说明

| 状态码 | 说明 | 处理建议 |
|--------|------|----------|
| 0 | 成功 | 正常展示数据 |
| 1001 | 参数错误 | 检查请求参数 |
| 1002 | 数据不存在 | 更换查询条件 |
| 1003 | 样本量不足 | 显示警告，谨慎参考 |
| 5000 | 服务器内部错误 | 稍后重试 |

## 错误响应示例

### 参数错误

```json
{
  "code": 1001,
  "message": "参数错误：缺少symbol参数",
  "timestamp": 1698764800,
  "data": null
}
```

### 数据不存在

```json
{
  "code": 1002,
  "message": "未找到BTCUSDT在10-30的历史数据",
  "timestamp": 1698764800,
  "data": null
}
```

### 样本量不足

```json
{
  "code": 1003,
  "message": "样本量不足，统计结果可能不可靠",
  "timestamp": 1698764800,
  "data": {
    "strategy_info": {...},
    "current_period_result": {
      "sample_count": 2,
      "reliability": "low",
      "reliability_note": "样本数量过少(仅2条)，统计结果不可靠"
    }
  }
}
```

## 项目结构

```
trade/
├── api/
│   ├── handler/          # 请求处理器
│   │   └── strategy_handler.go
│   ├── response/         # 响应结构体
│   │   └── response.go
│   └── router/           # 路由配置
│       └── router.go
├── cmd/
│   └── api/              # API服务启动文件
│       └── main.go
├── db/                   # 数据库连接
├── model/                # 数据模型
├── strategy/             # 策略实现
└── kline/                # K线数据处理
```

## 开发说明

### 添加新的策略

1. 在 `strategy/` 目录下实现新的策略逻辑
2. 在 `api/handler/strategy_handler.go` 中添加新的处理函数
3. 在 `api/response/response.go` 中定义响应结构体
4. 更新路由配置（如需要）

### 修改端口

编辑 `cmd/api/main.go`，修改：

```go
h := server.Default(
    server.WithHostPorts(":8080"),  // 修改这里的端口号
)
```

### 启用日志

可以添加 Hertz 的日志中间件：

```go
import "github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"

h.Use(recovery.Recovery())
```

### 跨域配置

如果需要支持跨域请求，可以添加 CORS 中间件：

```go
import "github.com/hertz-contrib/cors"

h.Use(cors.Default())
```

## 部署建议

### 编译二进制文件

```bash
# Linux/Mac
go build -o trade-api cmd/api/main.go

# Windows
go build -o trade-api.exe cmd/api/main.go
```

### 使用 Docker

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o trade-api cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/trade-api .
EXPOSE 8080
CMD ["./trade-api"]
```

构建和运行：

```bash
docker build -t trade-api .
docker run -p 8080:8080 trade-api
```

### 使用 systemd (Linux)

创建服务文件 `/etc/systemd/system/trade-api.service`:

```ini
[Unit]
Description=Trade Strategy API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/home/wws/trade
ExecStart=/home/wws/trade/trade-api
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl start trade-api
sudo systemctl enable trade-api
```

## 监控和维护

### 性能监控

建议添加 Hertz 的性能监控中间件，监控请求延迟、QPS 等指标。

### 日志管理

建议配置日志轮转，避免日志文件过大。

### 数据库连接池

当前配置：
- 最大空闲连接: 10
- 最大打开连接: 100

可在 `db/postgreSql.go` 中调整。

## 常见问题

### Q: 如何查看当前支持的交易对？

A: 查看 `config.json` 文件中的配置。

### Q: 数据从哪里来？

A: 数据从币安合约 API 获取，存储在 PostgreSQL 数据库中。

### Q: 如何更新K线数据？

A: 运行主程序 `go run main.go` 会自动更新数据。

### Q: API 性能如何？

A: 单次查询响应时间通常在 10-50ms 之间，取决于数据量和数据库性能。

## 技术支持

如有问题，请查看：
- `API_FORMAT_README.md` - 详细的响应格式说明
- `api_error_response_examples.json` - 错误响应示例
