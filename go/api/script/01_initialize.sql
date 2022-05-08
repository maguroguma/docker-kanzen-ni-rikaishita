-- DB作成
CREATE DATABASE testdb;

-- 作成したDBへ切り替え
\c testdb

-- スキーマ作成
CREATE SCHEMA testschema;

-- ロールの作成
CREATE ROLE hoge WITH LOGIN PASSWORD 'passw0rd';

-- 権限追加
GRANT ALL PRIVILEGES ON SCHEMA testschema TO hoge;
