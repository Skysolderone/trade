# Trade 项目部署文档

## 目录
- [系统要求](#系统要求)
- [快速部署](#快速部署)
- [详细步骤](#详细步骤)
- [Nginx配置](#nginx配置)
- [服务管理](#服务管理)
- [监控和日志](#监控和日志)
- [故障排查](#故障排查)

## 系统要求

### 硬件要求
- CPU: 2核心及以上
- 内存: 2GB及以上
- 磁盘: 20GB及以上

### 软件要求
- 操作系统: Linux (CentOS 7+, Ubuntu 18.04+)
- Go: 1.18+ (编译时需要)
- PostgreSQL: 12+
- Nginx: 1.18+ (可选，用于反向代理)

## 快速部署

### 1. 准备工作

```bash
# 克隆项目（或上传到服务器）
cd /path/to/trade

# 确保配置文件存在
# 编辑 config.json 配置交易对和时间周期
vi config.json

# 编辑数据库配置
# 修改 db/postgreSql.go 中的数据库连接信息
vi db/postgreSql.go
```

### 2. 构建程序

```bash
# 赋予执行权限
chmod +x scripts/build.sh

# 执行构建
./scripts/build.sh
```

### 3. 部署到服务器

```bash
# 赋予执行权限
chmod +x scripts/deploy.sh

# 使用 root 权限执行部署
sudo ./scripts/deploy.sh
```

部署完成后，服务会自动启动并在每天零点执行策略更新。

## 详细步骤

### 步骤1: 准备数据库

```bash
# 安装 PostgreSQL（如果未安装）
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib

# CentOS/RHEL
sudo yum install postgresql-server postgresql-contrib

# 初始化数据库
sudo postgresql-setup initdb

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库和用户
sudo -u postgres psql
```

在 PostgreSQL 命令行中：
```sql
-- 创建数据库
CREATE DATABASE trade;

-- 创建用户
CREATE USER tradeuser WITH PASSWORD 'your_password';

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE trade TO tradeuser;

-- 退出
\q
```

### 步骤2: 配置数据库连接

编辑 `db/postgreSql.go` 文件，修改数据库连接信息：

```go
dsn := "host=localhost user=tradeuser password=your_password dbname=trade port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

### 步骤3: 配置交易对

编辑 `config.json` 文件：

```json
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
```

### 步骤4: 构建项目

```bash
./scripts/build.sh
```

这会在 `bin/` 目录下生成以下文件：
- `trade` - 主程序（策略分析+定时任务）
- `web_server` - Web服务器
- `get_hourly_history` - 历史数据获取工具

### 步骤5: 部署到服务器

```bash
sudo ./scripts/deploy.sh
```

部署脚本会自动：
1. 创建 `trade` 用户
2. 创建必要的目录 (`/opt/trade`, `/var/log/trade`)
3. 复制程序文件和配置文件
4. 安装 systemd 服务
5. 启动服务

### 步骤6: 验证部署

```bash
# 检查服务状态
sudo systemctl status trade.service
sudo systemctl status trade-web.service

# 查看日志
sudo journalctl -u trade.service -n 50
sudo journalctl -u trade-web.service -n 50

# 测试Web服务
curl http://localhost:8080
curl http://localhost:8080/api/strategy1
```

## Nginx配置

### 安装Nginx

```bash
# Ubuntu/Debian
sudo apt-get install nginx

# CentOS/RHEL
sudo yum install nginx
```

### 配置Nginx

```bash
# 复制配置文件
sudo cp deploy/nginx/trade.conf /etc/nginx/sites-available/trade.conf

# 修改域名
sudo vi /etc/nginx/sites-available/trade.conf
# 将 your-domain.com 改为您的域名

# 创建软链接（Ubuntu/Debian）
sudo ln -s /etc/nginx/sites-available/trade.conf /etc/nginx/sites-enabled/

# CentOS/RHEL 直接复制到 conf.d
sudo cp deploy/nginx/trade.conf /etc/nginx/conf.d/

# 测试配置
sudo nginx -t

# 重启Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx
```

### 配置防火墙

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 配置HTTPS（可选）

使用 Let's Encrypt 免费证书：

```bash
# 安装 Certbot
# Ubuntu/Debian
sudo apt-get install certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

## 服务管理

### systemd 命令

```bash
# 查看服务状态
sudo systemctl status trade.service
sudo systemctl status trade-web.service

# 启动服务
sudo systemctl start trade.service
sudo systemctl start trade-web.service

# 停止服务
sudo systemctl stop trade.service
sudo systemctl stop trade-web.service

# 重启服务
sudo systemctl restart trade.service
sudo systemctl restart trade-web.service

# 设置开机自启
sudo systemctl enable trade.service
sudo systemctl enable trade-web.service

# 禁用开机自启
sudo systemctl disable trade.service
sudo systemctl disable trade-web.service
```

### 手动运行

```bash
# 切换到 trade 用户
sudo su - trade

# 进入工作目录
cd /opt/trade

# 单次运行模式
./trade -mode=once

# 定时任务模式
./trade -mode=daemon

# 启动Web服务器
./web_server -port=8080
```

## 监控和日志

### 查看日志

```bash
# systemd 日志
sudo journalctl -u trade.service -f
sudo journalctl -u trade-web.service -f

# 查看最近100行
sudo journalctl -u trade.service -n 100

# 查看今天的日志
sudo journalctl -u trade.service --since today

# 应用日志文件
tail -f /var/log/trade/trade.log
tail -f /var/log/trade/trade-error.log
tail -f /var/log/trade/web.log
tail -f /var/log/trade/web-error.log
```

### 日志轮转

创建 `/etc/logrotate.d/trade` 文件：

```
/var/log/trade/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 trade trade
    sharedscripts
    postrotate
        systemctl reload trade.service > /dev/null 2>&1 || true
        systemctl reload trade-web.service > /dev/null 2>&1 || true
    endscript
}
```

### 监控指标

```bash
# CPU和内存使用情况
ps aux | grep -E "trade|web_server"

# 监控进程
top -u trade

# 网络连接
sudo netstat -tlnp | grep -E "trade|8080"

# 磁盘使用
df -h /opt/trade
du -sh /opt/trade/*
```

## 故障排查

### 服务启动失败

```bash
# 查看详细错误
sudo journalctl -u trade.service -xe

# 检查配置文件
cat /opt/trade/config.json

# 检查数据库连接
psql -h localhost -U tradeuser -d trade -c "SELECT 1;"

# 检查文件权限
ls -la /opt/trade
ls -la /var/log/trade
```

### Web服务无法访问

```bash
# 检查端口是否监听
sudo netstat -tlnp | grep 8080
sudo ss -tlnp | grep 8080

# 检查防火墙
sudo ufw status
sudo firewall-cmd --list-all

# 测试本地连接
curl http://localhost:8080

# 检查Nginx配置
sudo nginx -t
sudo systemctl status nginx
```

### 数据库连接问题

```bash
# 检查PostgreSQL服务
sudo systemctl status postgresql

# 查看PostgreSQL日志
sudo tail -f /var/log/postgresql/postgresql-*.log

# 测试数据库连接
psql -h localhost -U tradeuser -d trade

# 检查PostgreSQL配置
sudo cat /etc/postgresql/*/main/pg_hba.conf
```

### 定时任务未执行

```bash
# 查看定时任务日志
sudo journalctl -u trade.service --since "00:00:00" --until "00:10:00"

# 手动触发一次
cd /opt/trade
sudo -u trade ./trade -mode=daemon -now

# 查看cron任务状态（应用内部使用robfig/cron）
# 检查程序日志确认调度器是否正常启动
```

### 性能问题

```bash
# 查看慢查询
# 在PostgreSQL中启用慢查询日志
# 编辑 postgresql.conf
sudo vi /etc/postgresql/*/main/postgresql.conf

# 添加或修改
log_min_duration_statement = 1000  # 记录超过1秒的查询

# 重启PostgreSQL
sudo systemctl restart postgresql

# 监控系统资源
htop
iotop -o
```

## 更新部署

### 更新应用

```bash
# 1. 在本地构建新版本
./scripts/build.sh

# 2. 上传到服务器
scp -r bin/ user@server:/path/to/trade/

# 3. 在服务器上执行部署
ssh user@server
cd /path/to/trade
sudo ./scripts/deploy.sh
```

### 数据库迁移

如果需要修改数据库结构，GORM会自动迁移（AutoMigrate）。手动迁移：

```bash
# 备份数据库
pg_dump -h localhost -U tradeuser trade > trade_backup_$(date +%Y%m%d).sql

# 如果需要恢复
psql -h localhost -U tradeuser trade < trade_backup_20250101.sql
```

## 常见配置

### 修改Web端口

```bash
# 方法1: 修改systemd服务文件
sudo vi /etc/systemd/system/trade-web.service
# 修改 ExecStart 行的 -port 参数

# 方法2: 使用环境变量
sudo systemctl edit trade-web.service
# 添加
[Service]
Environment="WEB_PORT=9000"

# 重新加载并重启
sudo systemctl daemon-reload
sudo systemctl restart trade-web.service
```

### 修改定时任务时间

编辑 `scheduler/scheduler.go`，修改cron表达式：

```go
// 每天零点: "0 0 0 * * *"
// 每天上午8点: "0 0 8 * * *"
// 每6小时: "0 0 */6 * * *"
```

重新构建和部署：
```bash
./scripts/build.sh
sudo ./scripts/deploy.sh
```

## 安全建议

1. **数据库安全**
   - 使用强密码
   - 限制数据库访问IP
   - 定期备份数据

2. **系统安全**
   - 启用防火墙
   - 及时更新系统补丁
   - 使用SSH密钥认证

3. **应用安全**
   - 不要在公网暴露数据库端口
   - 使用HTTPS加密传输
   - 定期查看日志

4. **监控告警**
   - 设置磁盘空间告警
   - 监控服务运行状态
   - 关注异常日志

## 技术支持

如有问题，请查看：
- 应用日志: `/var/log/trade/`
- systemd日志: `sudo journalctl -u trade.service`
- 数据库日志: `/var/log/postgresql/`
