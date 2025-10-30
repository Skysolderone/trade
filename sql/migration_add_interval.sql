-- ========================================
-- 数据库迁移脚本：添加 interval 字段和唯一索引
-- 创建日期: 2025-10-29
-- ========================================

-- 步骤 1: 添加 interval 字段（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'klines' AND column_name = 'interval'
    ) THEN
        ALTER TABLE public.klines ADD COLUMN interval text;
        RAISE NOTICE 'interval 字段添加成功';
    ELSE
        RAISE NOTICE 'interval 字段已存在，跳过';
    END IF;
END $$;

-- 步骤 2: 为现有数据填充默认 interval 值（假设都是日线数据）
-- ⚠️ 注意：如果你的历史数据包含不同时间周期，需要根据实际情况修改
UPDATE public.klines
SET interval = '1d'
WHERE interval IS NULL;

-- 步骤 3: 删除旧的可能存在的索引
DROP INDEX IF EXISTS idx_unique_kline;

-- 步骤 4: 创建唯一索引（symbol + interval + open_time）
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_kline
ON public.klines (symbol, interval, open_time);

-- 步骤 5: 修改 day 字段格式（统一为两位数格式）
-- 例如：从 "1-5" 改为 "01-05"
UPDATE public.klines
SET day =
    CASE
        WHEN day SIMILAR TO '[0-9]-[0-9]' THEN '0' || SUBSTRING(day, 1, 1) || '-0' || SUBSTRING(day, 3, 1)
        WHEN day SIMILAR TO '[0-9]{2}-[0-9]' THEN SUBSTRING(day, 1, 2) || '-0' || SUBSTRING(day, 4, 1)
        WHEN day SIMILAR TO '[0-9]-[0-9]{2}' THEN '0' || SUBSTRING(day, 1, 1) || '-' || SUBSTRING(day, 3, 2)
        ELSE day
    END
WHERE day NOT SIMILAR TO '[0-9]{2}-[0-9]{2}';

-- 步骤 6: 查看重复数据（运行后检查是否有重复）
SELECT symbol, interval, open_time, COUNT(*) as count
FROM public.klines
GROUP BY symbol, interval, open_time
HAVING COUNT(*) > 1
ORDER BY count DESC
LIMIT 10;

-- 步骤 7: 删除重复数据（保留 ID 最小的记录）
-- ⚠️ 警告：此操作会删除数据，建议先备份
/*
DELETE FROM public.klines a
USING public.klines b
WHERE a.id > b.id
  AND a.symbol = b.symbol
  AND a.interval = b.interval
  AND a.open_time = b.open_time;
*/

-- 步骤 8: 验证迁移结果
SELECT
    COUNT(*) as total_records,
    COUNT(DISTINCT symbol) as unique_symbols,
    COUNT(DISTINCT interval) as unique_intervals,
    MIN(open_time) as earliest_time,
    MAX(open_time) as latest_time
FROM public.klines;

-- 迁移完成！
