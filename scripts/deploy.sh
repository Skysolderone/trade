#!/bin/bash

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          Trade 项目部署脚本                                    ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# 配置变量（可以通过环境变量覆盖）
DEPLOY_USER="${DEPLOY_USER:-trade}"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/trade}"
LOG_DIR="${LOG_DIR:-/var/log/trade}"
WEB_PORT="${WEB_PORT:-8080}"

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "❌ 请使用 root 权限运行此脚本"
    exit 1
fi

echo "📋 部署配置:"
echo "  - 用户: $DEPLOY_USER"
echo "  - 目录: $DEPLOY_DIR"
echo "  - 日志目录: $LOG_DIR"
echo "  - Web端口: $WEB_PORT"
echo ""

# 创建用户（如果不存在）
if ! id "$DEPLOY_USER" &>/dev/null; then
    echo "👤 创建用户: $DEPLOY_USER"
    useradd -r -s /bin/bash -d "$DEPLOY_DIR" "$DEPLOY_USER"
fi

# 创建目录
echo "📁 创建目录结构..."
mkdir -p "$DEPLOY_DIR"
mkdir -p "$LOG_DIR"
mkdir -p "$DEPLOY_DIR/config"
mkdir -p "$DEPLOY_DIR/web"

# 停止现有服务
echo ""
echo "🛑 停止现有服务..."
systemctl stop trade.service 2>/dev/null || true
systemctl stop trade-web.service 2>/dev/null || true

# 复制二进制文件
echo ""
echo "📦 复制程序文件..."
if [ -f "bin/trade" ]; then
    cp bin/trade "$DEPLOY_DIR/"
    chmod +x "$DEPLOY_DIR/trade"
    echo "  ✅ trade"
else
    echo "  ❌ bin/trade 不存在，请先运行 ./scripts/build.sh"
    exit 1
fi

if [ -f "bin/web_server" ]; then
    cp bin/web_server "$DEPLOY_DIR/"
    chmod +x "$DEPLOY_DIR/web_server"
    echo "  ✅ web_server"
fi

if [ -f "bin/get_hourly_history" ]; then
    cp bin/get_hourly_history "$DEPLOY_DIR/"
    chmod +x "$DEPLOY_DIR/get_hourly_history"
    echo "  ✅ get_hourly_history"
fi

# 复制配置文件
echo ""
echo "⚙️  复制配置文件..."
if [ -f "config.json" ]; then
    cp config.json "$DEPLOY_DIR/"
    echo "  ✅ config.json"
else
    echo "  ⚠️  config.json 不存在，请手动创建"
fi

# 复制Web资源
echo ""
echo "🌐 复制Web资源..."
if [ -d "web" ]; then
    cp -r web/* "$DEPLOY_DIR/web/"
    echo "  ✅ web/"
fi

# 安装systemd服务
echo ""
echo "🔧 安装systemd服务..."
if [ -f "deploy/systemd/trade.service" ]; then
    cp deploy/systemd/trade.service /etc/systemd/system/
    echo "  ✅ trade.service"
fi

if [ -f "deploy/systemd/trade-web.service" ]; then
    cp deploy/systemd/trade-web.service /etc/systemd/system/
    echo "  ✅ trade-web.service"
fi

# 设置权限
echo ""
echo "🔐 设置文件权限..."
chown -R "$DEPLOY_USER:$DEPLOY_USER" "$DEPLOY_DIR"
chown -R "$DEPLOY_USER:$DEPLOY_USER" "$LOG_DIR"

# 重新加载systemd
echo ""
echo "🔄 重新加载systemd..."
systemctl daemon-reload

# 启动服务
echo ""
echo "🚀 启动服务..."
systemctl enable trade.service
systemctl start trade.service
echo "  ✅ trade.service 已启动"

systemctl enable trade-web.service
systemctl start trade-web.service
echo "  ✅ trade-web.service 已启动"

# 等待服务启动
sleep 2

# 检查服务状态
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 服务状态:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
systemctl status trade.service --no-pager -l || true
echo ""
systemctl status trade-web.service --no-pager -l || true
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          部署完成！                                            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "📝 常用命令:"
echo "  # 查看服务状态"
echo "  sudo systemctl status trade.service"
echo "  sudo systemctl status trade-web.service"
echo ""
echo "  # 查看日志"
echo "  sudo journalctl -u trade.service -f"
echo "  sudo journalctl -u trade-web.service -f"
echo "  tail -f $LOG_DIR/trade.log"
echo "  tail -f $LOG_DIR/web.log"
echo ""
echo "  # 重启服务"
echo "  sudo systemctl restart trade.service"
echo "  sudo systemctl restart trade-web.service"
echo ""
echo "  # 停止服务"
echo "  sudo systemctl stop trade.service"
echo "  sudo systemctl stop trade-web.service"
echo ""
echo "🌐 Web访问地址: http://服务器IP:$WEB_PORT"
echo ""
