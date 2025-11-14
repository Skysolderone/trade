# Trade 快速部署指南

本指南帮助您在10分钟内完成服务器部署。

## 一、服务器准备

### 1. 系统要求
- **操作系统**: Linux (Ubuntu 18.04+, CentOS 7+)
- **配置**: 2核CPU, 2GB内存, 20GB磁盘

### 2. 安装依赖

#### Ubuntu/Debian
```bash
# 更新系统
sudo apt-get update && sudo apt-get upgrade -y

# 安装PostgreSQL
sudo apt-get install -y postgresql postgresql-contrib

# 安装Nginx（可选）
sudo apt-get install -y nginx

# 安装其他工具
sudo apt-get install -y git curl wget
```

#### CentOS/RHEL
```bash
# 更新系统
sudo yum update -y

# 安装PostgreSQL
sudo yum install -y postgresql-server postgresql-contrib

# 初始化数据库
sudo postgresql-setup initdb

# 安装Nginx（可选）
sudo yum install -y nginx

# 安装其他工具
sudo yum install -y git curl wget
```

## 二、配置数据库

```bash
# 启动PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库和用户
sudo -u postgres psql << EOF
CREATE DATABASE trade;
CREATE USER tradeuser WITH PASSWORD 'YourStrongPassword123!';
GRANT ALL PRIVILEGES ON DATABASE trade TO tradeuser;
\q
EOF
```

## 三、上传项目文件

### 方法1: 使用Git（推荐）
```bash
# 克隆项目
cd ~
git clone your-repo-url trade
cd trade
```

### 方法2: 手动上传
```bash
# 在本地打包
tar -czf trade.tar.gz /path/to/trade

# 上传到服务器
scp trade.tar.gz user@server-ip:~/

# 在服务器上解压
ssh user@server-ip
tar -xzf trade.tar.gz
cd trade
```

## 四、配置项目

### 1. 修改数据库配置

编辑 `db/postgreSql.go`:
```go
dsn := "host=localhost user=tradeuser password=YourStrongPassword123! dbname=trade port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

或者创建环境变量文件（推荐）:
```bash
# 创建配置文件
cat > .env << EOF
DB_HOST=localhost
DB_USER=tradeuser
DB_PASSWORD=YourStrongPassword123!
DB_NAME=trade
DB_PORT=5432
EOF
```

### 2. 配置交易对

编辑 `config.json`:
```bash
cat > config.json << 'EOF'
{
  "symbols": [
    {
      "symbol": "BTCUSDT",
      "intervals": ["1d", "8h", "4h", "2h", "1h"]
    },
    {
      "symbol": "ETHUSDT",
      "intervals": ["1d", "8h", "4h", "2h", "1h"]
    }
  ]
}
EOF
```

## 五、构建和部署

### 一键部署
```bash
# 赋予执行权限
chmod +x scripts/build.sh scripts/deploy.sh

# 构建项目
./scripts/build.sh

# 部署到服务器（需要root权限）
sudo ./scripts/deploy.sh
```

完成！服务已自动启动。

## 六、验证部署

### 1. 检查服务状态
```bash
# 查看定时任务服务
sudo systemctl status trade.service

# 查看Web服务
sudo systemctl status trade-web.service
```

输出应该显示 **active (running)**

### 2. 测试Web访问
```bash
# 本地测试
curl http://localhost:8080

# 测试API
curl http://localhost:8080/api/strategy1
```

### 3. 查看日志
```bash
# 查看实时日志
sudo journalctl -u trade.service -f
sudo journalctl -u trade-web.service -f

# 或查看日志文件
tail -f /var/log/trade/trade.log
tail -f /var/log/trade/web.log
```

## 七、配置Nginx（可选）

如果您想通过域名访问，配置Nginx反向代理：

```bash
# 复制配置文件
sudo cp deploy/nginx/trade.conf /etc/nginx/sites-available/trade.conf

# 修改域名
sudo sed -i 's/your-domain.com/yourdomain.com/g' /etc/nginx/sites-available/trade.conf

# 启用配置（Ubuntu/Debian）
sudo ln -s /etc/nginx/sites-available/trade.conf /etc/nginx/sites-enabled/

# CentOS/RHEL直接复制
sudo cp deploy/nginx/trade.conf /etc/nginx/conf.d/

# 测试配置
sudo nginx -t

# 重启Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx

# 开放防火墙端口
sudo ufw allow 80/tcp   # Ubuntu/Debian
sudo ufw allow 443/tcp
sudo firewall-cmd --permanent --add-service=http    # CentOS/RHEL
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

现在可以通过 `http://yourdomain.com` 访问Web界面。

## 八、配置HTTPS（推荐）

使用Let's Encrypt免费SSL证书：

```bash
# 安装Certbot
# Ubuntu/Debian
sudo apt-get install -y certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install -y certbot python3-certbot-nginx

# 获取证书（自动配置Nginx）
sudo certbot --nginx -d yourdomain.com

# 设置自动续期
sudo certbot renew --dry-run
```

完成！现在可以通过 `https://yourdomain.com` 安全访问。

## 九、日常管理

### 查看服务状态
```bash
sudo systemctl status trade.service
sudo systemctl status trade-web.service
```

### 重启服务
```bash
sudo systemctl restart trade.service
sudo systemctl restart trade-web.service
```

### 查看日志
```bash
# 查看最近50条日志
sudo journalctl -u trade.service -n 50

# 实时查看日志
sudo journalctl -u trade.service -f

# 查看应用日志
tail -f /var/log/trade/trade.log
```

### 手动触发策略更新
```bash
# 切换到trade用户
sudo su - trade

# 单次运行
cd /opt/trade
./trade -mode=once
```

## 十、功能说明

### 自动定时任务
- 程序会**每天00:00:00自动执行**策略更新
- 包括：更新K线数据 → 运行策略一 → 运行策略二
- 所有结果自动保存到数据库

### Web界面
访问 `http://服务器IP:8080` 或 `http://yourdomain.com`

功能：
- **策略一**：查看历史同期涨跌分析（跨年、跨月对比）
- **策略二**：查看小时级别涨跌分析（日内时段对比）
- 支持按交易对、时间周期筛选
- 可展开查看详细历史记录

### API接口
```bash
# 获取策略一结果
curl http://localhost:8080/api/strategy1
curl http://localhost:8080/api/strategy1?symbol=BTCUSDT
curl http://localhost:8080/api/strategy1?symbol=BTCUSDT&interval=1d

# 获取策略二结果
curl http://localhost:8080/api/strategy2
curl http://localhost:8080/api/strategy2?symbol=BTCUSDT&hour=15
```

## 故障排查

### 服务启动失败
```bash
# 查看详细错误
sudo journalctl -u trade.service -xe

# 检查配置文件
cat /opt/trade/config.json

# 测试数据库连接
psql -h localhost -U tradeuser -d trade -c "SELECT 1;"
```

### Web无法访问
```bash
# 检查端口
sudo netstat -tlnp | grep 8080

# 检查防火墙
sudo ufw status
sudo firewall-cmd --list-all

# 测试本地连接
curl http://localhost:8080
```

### 数据库连接失败
```bash
# 检查PostgreSQL服务
sudo systemctl status postgresql

# 查看日志
sudo tail -f /var/log/postgresql/postgresql-*.log

# 测试连接
psql -h localhost -U tradeuser -d trade
```

## 进阶配置

详细配置和高级功能请参考 [DEPLOY.md](./DEPLOY.md)

## 帮助和支持

- 查看完整部署文档: [DEPLOY.md](./DEPLOY.md)
- 查看Web功能说明: [WEB_README.md](./WEB_README.md)
- 查看项目说明: [CLAUDE.md](./CLAUDE.md)
