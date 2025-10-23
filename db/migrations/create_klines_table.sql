-- 创建K线数据表
CREATE TABLE IF NOT EXISTS klines (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,                    -- 交易对，如 BTC/USDT
    timeframe VARCHAR(10) NOT NULL,                 -- K线周期，如 1m, 5m, 1h, 1d
    open_time BIGINT NOT NULL,                      -- 开盘时间戳（毫秒）
    open_time_dt TIMESTAMP NOT NULL,                -- 开盘时间（可读格式）
    open_price DECIMAL(20, 8) NOT NULL,             -- 开盘价
    high_price DECIMAL(20, 8) NOT NULL,             -- 最高价
    low_price DECIMAL(20, 8) NOT NULL,              -- 最低价
    close_price DECIMAL(20, 8) NOT NULL,            -- 收盘价
    volume DECIMAL(20, 8) NOT NULL,                 -- 成交量
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间

    -- 唯一约束：同一交易对、同一周期、同一时间只能有一条记录
    CONSTRAINT unique_kline UNIQUE (symbol, timeframe, open_time)
);

-- 创建索引以提高查询性能
CREATE INDEX idx_klines_symbol ON klines(symbol);
CREATE INDEX idx_klines_timeframe ON klines(timeframe);
CREATE INDEX idx_klines_open_time ON klines(open_time);
CREATE INDEX idx_klines_symbol_timeframe_time ON klines(symbol, timeframe, open_time DESC);

-- 创建更新时间的触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_klines_updated_at
    BEFORE UPDATE ON klines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加表注释
COMMENT ON TABLE klines IS '币安永续合约K线数据表';
COMMENT ON COLUMN klines.symbol IS '交易对符号';
COMMENT ON COLUMN klines.timeframe IS 'K线时间周期';
COMMENT ON COLUMN klines.open_time IS '开盘时间戳（毫秒）';
COMMENT ON COLUMN klines.open_time_dt IS '开盘时间（易读格式）';
COMMENT ON COLUMN klines.open_price IS '开盘价';
COMMENT ON COLUMN klines.high_price IS '最高价';
COMMENT ON COLUMN klines.low_price IS '最低价';
COMMENT ON COLUMN klines.close_price IS '收盘价';
COMMENT ON COLUMN klines.volume IS '成交量';
