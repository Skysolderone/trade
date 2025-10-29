/*
 Navicat Premium Dump SQL

 Source Server         : line_postgreSql
 Source Server Type    : PostgreSQL
 Source Server Version : 170006 (170006)
 Source Host           : pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com:5432
 Source Catalog        : trade
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170006 (170006)
 File Encoding         : 65001

 Date: 25/10/2025 16:04:56
*/


-- ----------------------------
-- Table structure for klines
-- ----------------------------
DROP TABLE IF EXISTS "public"."klines";
CREATE TABLE "public"."klines" (
  "id" int8 NOT NULL DEFAULT nextval('klines_id_seq'::regclass),
  "symbol" text COLLATE "pg_catalog"."default",
  "open" numeric,
  "close" numeric,
  "high" numeric,
  "low" numeric,
  "open_time" timestamptz(6),
  "close_time" timestamptz(6),
  "date" text COLLATE "pg_catalog"."default",
  "day" varchar(255) COLLATE "pg_catalog"."default",
  "hour" varchar(255) COLLATE "pg_catalog"."default",
  "week" varchar(255) COLLATE "pg_catalog"."default",
  "min" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."klines" OWNER TO "wws";

-- ----------------------------
-- Primary Key structure for table klines
-- ----------------------------
ALTER TABLE "public"."klines" ADD CONSTRAINT "klines_pkey" PRIMARY KEY ("id");
