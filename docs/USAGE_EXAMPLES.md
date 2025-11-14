# API 使用示例

## 前提条件

1. 确保数据库已初始化并包含K线数据
2. API 服务正在运行（`go run cmd/api/main.go`）

## cURL 示例

### 1. 基础示例

#### 策略一：分析今天的历史表现

```bash
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d'
```

#### 策略一：分析指定日期

```bash
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d&date=2024-10-30'
```

#### 策略二：分析当前小时

```bash
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h'
```

#### 策略二：分析指定小时

```bash
curl 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h&hour=14'
```

### 2. POST 请求示例

#### 策略一

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

#### 策略二

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

### 3. 使用 jq 格式化输出

```bash
curl -s 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d' | jq '.'
```

### 4. 只获取特定字段

#### 只获取交易建议

```bash
curl -s 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d' \
  | jq '.data.trading_recommendation'
```

#### 只获取上涨率

```bash
curl -s 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=BTCUSDT&interval=1d' \
  | jq '.data.current_period_result.up_rate'
```

#### 只获取高胜率时段（策略二）

```bash
curl -s 'http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_2&symbol=BTCUSDT&interval=1h' \
  | jq '.data.high_win_hours.hours'
```

## JavaScript/TypeScript 示例

### 使用 Fetch API

```javascript
// 策略一
async function analyzeStrategy1(symbol, interval, date) {
  const params = new URLSearchParams({
    strategy_type: 'strategy_1',
    symbol: symbol,
    interval: interval,
    date: date
  });

  const response = await fetch(`http://localhost:8080/api/v1/strategy/analyze?${params}`);
  const data = await response.json();

  if (data.code === 0) {
    console.log('分析成功:', data.data);
    console.log('上涨概率:', data.data.current_period_result.up_rate + '%');
    console.log('交易信号:', data.data.trading_recommendation.signal);
  } else {
    console.error('分析失败:', data.message);
  }

  return data;
}

// 调用示例
analyzeStrategy1('BTCUSDT', '1d', '2024-10-30');
```

```javascript
// 策略二
async function analyzeStrategy2(symbol, interval, hour) {
  const params = new URLSearchParams({
    strategy_type: 'strategy_2',
    symbol: symbol,
    interval: interval,
    hour: hour
  });

  const response = await fetch(`http://localhost:8080/api/v1/strategy/analyze?${params}`);
  const data = await response.json();

  if (data.code === 0) {
    console.log('分析成功:', data.data);
    console.log('当前小时上涨率:', data.data.current_hour_result.up_rate + '%');
    console.log('最佳交易时段:', data.data.trading_recommendation.optimal_trading_hours);
  } else {
    console.error('分析失败:', data.message);
  }

  return data;
}

// 调用示例
analyzeStrategy2('BTCUSDT', '1h', 14);
```

### 使用 Axios

```javascript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

// 策略一
async function analyzeStrategy1(symbol, interval, date) {
  try {
    const response = await axios.get(`${API_BASE_URL}/strategy/analyze`, {
      params: {
        strategy_type: 'strategy_1',
        symbol,
        interval,
        date
      }
    });

    if (response.data.code === 0) {
      return response.data.data;
    } else {
      throw new Error(response.data.message);
    }
  } catch (error) {
    console.error('请求失败:', error);
    throw error;
  }
}

// 策略二
async function analyzeStrategy2(symbol, interval, hour) {
  try {
    const response = await axios.post(`${API_BASE_URL}/strategy/analyze`, {
      strategy_type: 'strategy_2',
      symbol,
      interval,
      hour
    });

    if (response.data.code === 0) {
      return response.data.data;
    } else {
      throw new Error(response.data.message);
    }
  } catch (error) {
    console.error('请求失败:', error);
    throw error;
  }
}
```

## Python 示例

### 使用 requests

```python
import requests
import json

API_BASE_URL = 'http://localhost:8080/api/v1'

def analyze_strategy1(symbol, interval, date=None):
    """策略一：历史同期涨跌分析"""
    params = {
        'strategy_type': 'strategy_1',
        'symbol': symbol,
        'interval': interval
    }
    if date:
        params['date'] = date

    response = requests.get(f'{API_BASE_URL}/strategy/analyze', params=params)
    data = response.json()

    if data['code'] == 0:
        print(f"分析成功")
        print(f"上涨概率: {data['data']['current_period_result']['up_rate']:.2f}%")
        print(f"交易信号: {data['data']['trading_recommendation']['signal']}")
        return data['data']
    else:
        print(f"分析失败: {data['message']}")
        return None

def analyze_strategy2(symbol, interval, hour=None):
    """策略二：小时级别涨跌分析"""
    payload = {
        'strategy_type': 'strategy_2',
        'symbol': symbol,
        'interval': interval
    }
    if hour is not None:
        payload['hour'] = hour

    response = requests.post(
        f'{API_BASE_URL}/strategy/analyze',
        json=payload
    )
    data = response.json()

    if data['code'] == 0:
        print(f"分析成功")
        print(f"当前小时上涨率: {data['data']['current_hour_result']['up_rate']:.2f}%")
        print(f"最佳交易时段: {data['data']['trading_recommendation']['optimal_trading_hours']}")
        return data['data']
    else:
        print(f"分析失败: {data['message']}")
        return None

# 使用示例
if __name__ == '__main__':
    # 策略一
    result1 = analyze_strategy1('BTCUSDT', '1d', '2024-10-30')

    # 策略二
    result2 = analyze_strategy2('BTCUSDT', '1h', 14)
```

## Go 示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

const APIBaseURL = "http://localhost:8080/api/v1"

type APIResponse struct {
    Code      int             `json:"code"`
    Message   string          `json:"message"`
    Timestamp int64           `json:"timestamp"`
    Data      json.RawMessage `json:"data"`
}

// 策略一
func analyzeStrategy1(symbol, interval, date string) (*APIResponse, error) {
    params := url.Values{}
    params.Add("strategy_type", "strategy_1")
    params.Add("symbol", symbol)
    params.Add("interval", interval)
    if date != "" {
        params.Add("date", date)
    }

    resp, err := http.Get(fmt.Sprintf("%s/strategy/analyze?%s", APIBaseURL, params.Encode()))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result APIResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

// 策略二
func analyzeStrategy2(symbol, interval string, hour *int) (*APIResponse, error) {
    params := url.Values{}
    params.Add("strategy_type", "strategy_2")
    params.Add("symbol", symbol)
    params.Add("interval", interval)
    if hour != nil {
        params.Add("hour", fmt.Sprintf("%d", *hour))
    }

    resp, err := http.Get(fmt.Sprintf("%s/strategy/analyze?%s", APIBaseURL, params.Encode()))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result APIResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func main() {
    // 策略一
    result1, err := analyzeStrategy1("BTCUSDT", "1d", "2024-10-30")
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    fmt.Printf("策略一结果: %+v\n", result1)

    // 策略二
    hour := 14
    result2, err := analyzeStrategy2("BTCUSDT", "1h", &hour)
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    fmt.Printf("策略二结果: %+v\n", result2)
}
```

## 移动端示例

### Android (Kotlin + Retrofit)

```kotlin
import retrofit2.http.*

interface StrategyApi {
    @GET("strategy/analyze")
    suspend fun analyzeStrategy(
        @Query("strategy_type") strategyType: String,
        @Query("symbol") symbol: String,
        @Query("interval") interval: String,
        @Query("date") date: String? = null,
        @Query("hour") hour: Int? = null
    ): ApiResponse<StrategyData>
}

data class ApiResponse<T>(
    val code: Int,
    val message: String,
    val timestamp: Long,
    val data: T?
)

// 使用示例
class StrategyRepository(private val api: StrategyApi) {
    suspend fun analyzeStrategy1(symbol: String, interval: String, date: String): Result<StrategyData> {
        return try {
            val response = api.analyzeStrategy(
                strategyType = "strategy_1",
                symbol = symbol,
                interval = interval,
                date = date
            )
            if (response.code == 0) {
                Result.success(response.data!!)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
```

### iOS (Swift + Alamofire)

```swift
import Alamofire

class StrategyAPI {
    private let baseURL = "http://localhost:8080/api/v1"

    func analyzeStrategy1(symbol: String, interval: String, date: String, completion: @escaping (Result<StrategyResponse, Error>) -> Void) {
        let parameters: [String: Any] = [
            "strategy_type": "strategy_1",
            "symbol": symbol,
            "interval": interval,
            "date": date
        ]

        AF.request("\(baseURL)/strategy/analyze", parameters: parameters)
            .validate()
            .responseDecodable(of: APIResponse<StrategyResponse>.self) { response in
                switch response.result {
                case .success(let data):
                    if data.code == 0 {
                        completion(.success(data.data!))
                    } else {
                        completion(.failure(NSError(domain: "", code: data.code, userInfo: [NSLocalizedDescriptionKey: data.message])))
                    }
                case .failure(let error):
                    completion(.failure(error))
                }
            }
    }
}

struct APIResponse<T: Codable>: Codable {
    let code: Int
    let message: String
    let timestamp: Int64
    let data: T?
}
```

## 批量请求示例

### 分析多个交易对

```bash
#!/bin/bash

symbols=("BTCUSDT" "ETHUSDT" "BNBUSDT")

for symbol in "${symbols[@]}"; do
    echo "分析 $symbol..."
    curl -s "http://localhost:8080/api/v1/strategy/analyze?strategy_type=strategy_1&symbol=$symbol&interval=1d" \
      | jq "{symbol: \"$symbol\", up_rate: .data.current_period_result.up_rate, signal: .data.trading_recommendation.signal}"
    echo ""
done
```

### Python 并发请求

```python
import asyncio
import aiohttp

async def analyze_multiple_symbols(symbols):
    async with aiohttp.ClientSession() as session:
        tasks = []
        for symbol in symbols:
            task = analyze_symbol(session, symbol)
            tasks.append(task)

        results = await asyncio.gather(*tasks)
        return results

async def analyze_symbol(session, symbol):
    url = f"http://localhost:8080/api/v1/strategy/analyze"
    params = {
        'strategy_type': 'strategy_1',
        'symbol': symbol,
        'interval': '1d'
    }

    async with session.get(url, params=params) as response:
        data = await response.json()
        return {
            'symbol': symbol,
            'up_rate': data['data']['current_period_result']['up_rate'],
            'signal': data['data']['trading_recommendation']['signal']
        }

# 使用
symbols = ['BTCUSDT', 'ETHUSDT', 'BNBUSDT']
results = asyncio.run(analyze_multiple_symbols(symbols))
print(results)
```

## 运行测试脚本

```bash
# 确保API服务正在运行
go run cmd/api/main.go

# 在另一个终端运行测试脚本
./test_api.sh
```

## 注意事项

1. **样本量**: 当返回 `code: 1003` 时，表示样本量不足，建议谨慎参考
2. **数据可用性**: 确保数据库中有对应交易对和周期的数据
3. **日期格式**: date 参数格式必须为 `YYYY-MM-DD`
4. **小时范围**: hour 参数必须在 0-23 之间
5. **网络延迟**: 建议设置合理的超时时间（如 5-10 秒）
