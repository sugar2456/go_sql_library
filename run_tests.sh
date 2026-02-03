#!/bin/bash

# テスト実行スクリプト

set -e

echo "================================"
echo "Go SQL Library テスト実行"
echo "================================"

# 環境変数設定
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-3306}

# MySQLが起動しているか確認
echo ""
echo "MySQLの接続確認..."
if ! docker exec go_mysql mysqladmin ping -h localhost -u root -ppassword > /dev/null 2>&1; then
    echo "エラー: MySQLコンテナが起動していません"
    echo "docker compose up -d を実行してください"
    exit 1
fi
echo "✓ MySQL接続OK"

# 各パッケージのテスト実行
echo ""
echo "1. 標準SQLライブラリのテスト"
echo "----------------------------"
go test -v ./standard/...

echo ""
echo "2. sqlxのテスト"
echo "----------------------------"
go test -v ./sqlx/...

echo ""
echo "3. GORMのテスト"
echo "----------------------------"
go test -v ./gorm/...

echo ""
echo "4. entのテスト"
echo "----------------------------"
go test -v ./ent/...

echo ""
echo "================================"
echo "すべてのテストが完了しました！"
echo "================================"
