#!/bin/bash

# ベンチマークテスト実行スクリプト

set -e

echo "========================================"
echo "Go SQL Library パフォーマンステスト"
echo "========================================"

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

# ベンチマーク回数
BENCHTIME=${BENCHTIME:-5s}
echo ""
echo "ベンチマーク設定:"
echo "  実行時間: $BENCHTIME"
echo "  データベース: $DB_HOST:$DB_PORT"
echo ""

# 結果保存ディレクトリ
RESULTS_DIR="benchmark_results"
mkdir -p "$RESULTS_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "結果保存先: $RESULTS_DIR/benchmark_$TIMESTAMP.txt"
echo ""

# 全体のベンチマーク結果ファイル
RESULT_FILE="$RESULTS_DIR/benchmark_$TIMESTAMP.txt"

{
    echo "========================================"
    echo "Go SQL Library パフォーマンステスト"
    echo "実行日時: $(date)"
    echo "========================================"
    echo ""

    # 各パッケージのベンチマーク実行
    for pkg in standard sqlx gorm ent; do
        echo "========================================"
        echo "$pkg ベンチマーク"
        echo "========================================"
        go test -tags=benchmark -bench=. -benchtime="$BENCHTIME" -benchmem "./$pkg/..." 2>&1 || echo "エラー: $pkg のベンチマーク実行に失敗"
        echo ""
    done

    echo "========================================"
    echo "比較ベンチマーク"
    echo "========================================"
    echo ""
    echo "全ライブラリの比較:"
    go test -tags=benchmark -bench=BenchmarkGetAll -benchtime="$BENCHTIME" -benchmem ./standard/... ./sqlx/... ./gorm/... ./ent/... 2>&1 | grep -E "Benchmark|PASS|FAIL|ok"

} | tee "$RESULT_FILE"

echo ""
echo "========================================"
echo "ベンチマーク完了！"
echo "結果: $RESULT_FILE"
echo "========================================"
echo ""
echo "結果の見方:"
echo "  ns/op    - 1操作あたりのナノ秒"
echo "  B/op     - 1操作あたりのバイト数"
echo "  allocs/op - 1操作あたりのメモリアロケーション数"
