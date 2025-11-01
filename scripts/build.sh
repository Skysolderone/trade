#!/bin/bash

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          开始构建 Trade 项目                                   ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未安装 Go 环境"
    exit 1
fi

echo "Go 版本: $(go version)"
echo ""

# 创建输出目录
mkdir -p bin

# 构建主程序（策略分析和定时任务）
echo "📦 构建主程序: trade"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/trade -ldflags="-s -w" main.go
echo "✅ 主程序构建完成: bin/trade"
echo ""

# 构建Web服务器
echo "📦 构建Web服务器: web_server"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/web_server -ldflags="-s -w" ./cmd/web/main.go
echo "✅ Web服务器构建完成: bin/web_server"
echo ""

# 构建小时历史数据获取工具
echo "📦 构建工具: get_hourly_history"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/get_hourly_history -ldflags="-s -w" ./cmd/get_hourly_history/main.go
echo "✅ 工具构建完成: bin/get_hourly_history"
echo ""

# 显示文件大小
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "构建结果:"
ls -lh bin/
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          构建完成！                                            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "💡 提示:"
echo "  - 主程序: ./bin/trade -mode=daemon"
echo "  - Web服务: ./bin/web_server -port=8080"
echo "  - 历史数据工具: ./bin/get_hourly_history"
