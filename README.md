# 概要
このリポジトリはgoのSQLライブラリ検証用のプロジェクトです

以下のライブラリの比較をします

- 標準SQLライブラリ
- sqlx
- GORM
- ent

## プロジェクト構造

```
go_sql_library/
├── main.go              # メインアプリケーション
├── model/
│   └── user.go         # 共通モデルとインターフェース定義
├── standard/
│   └── repository.go   # 標準database/sql実装
├── sqlx/
│   └── repository.go   # sqlx実装
├── gorm/
│   └── repository.go   # GORM実装
└── ent/
    └── repository.go   # ent実装（簡略版）
```

各パッケージは共通の`UserRepository`インターフェースを実装しており、環境変数で切り替えが可能です。

## 環境構築

### 必要なもの
- Docker
- Docker Compose

### 起動方法

```bash
# コンテナのビルドと起動
docker compose up -d

# ログの確認
docker compose logs -f app

# 動作確認
curl http://localhost:8081
curl http://localhost:8081/users
```

### 停止方法

```bash
# コンテナの停止
docker compose down

# ボリュームも削除する場合
docker compose down -v
```

## ライブラリの切り替え

[compose.yml](compose.yml)の`LIBRARY_TYPE`環境変数を変更することで、使用するライブラリを切り替えられます。

```yaml
environment:
  - LIBRARY_TYPE=standard  # standard, sqlx, gorm, ent から選択
```

変更後は、コンテナを再起動してください。

```bash
docker compose down
docker compose up -d
```

## 接続情報

- **アプリケーション**: http://localhost:8081
- **MySQL**:
  - ホスト: localhost
  - ポート: 3306
  - ユーザー: root
  - パスワード: password
  - データベース: testdb

## API エンドポイント

- `GET /` - ホーム（使用中のライブラリ表示）
- `GET /ping` - ヘルスチェック
- `GET /users` - 全ユーザー取得
- `GET /users/{id}` - 特定ユーザー取得
- `POST /users` - ユーザー作成
- `PUT /users/{id}` - ユーザー更新
- `DELETE /users/{id}` - ユーザー削除

### 使用例

```bash
# 全ユーザー取得
curl http://localhost:8081/users

# 新規ユーザー作成
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"name":"田中太郎","email":"tanaka@example.com"}'

# ユーザー取得
curl http://localhost:8081/users/1

# ユーザー更新
curl -X PUT http://localhost:8081/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"田中次郎","email":"tanaka2@example.com"}'

# ユーザー削除
curl -X DELETE http://localhost:8081/users/1
```

## テスト

各SQLライブラリの実装に対するテストコードが用意されています。

### テストの実行方法

```bash
# コンテナ内で全てのテストを実行
docker exec go_app sh -c "go test -v ./standard/... ./sqlx/... ./gorm/... ./ent/..."

# 個別のパッケージをテスト
docker exec go_app go test -v ./standard/...
docker exec go_app go test -v ./sqlx/...
docker exec go_app go test -v ./gorm/...
docker exec go_app go test -v ./ent/...
```

### テスト内容

各パッケージで以下の機能をテストしています：
- `GetAll()` - 全ユーザー取得
- `GetByID()` - ID指定でユーザー取得
- `Create()` - 新規ユーザー作成
- `Update()` - ユーザー情報更新
- `Delete()` - ユーザー削除

### カバレッジの確認

```bash
# カバレッジを確認
docker exec go_app go test -cover ./standard/... ./sqlx/... ./gorm/... ./ent/...
```

## パフォーマンステスト

各SQLライブラリのパフォーマンスを比較するベンチマークテストが用意されています。
通常のテストとは分離されており、`-tags=benchmark`オプションで実行します。

### ベンチマークの実行方法

```bash
# 全てのベンチマークを実行（コンテナ内）
docker exec go_app go test -tags=benchmark -bench=. -benchmem ./standard/... ./sqlx/... ./gorm/... ./ent/...

# 個別のパッケージをベンチマーク
docker exec go_app go test -tags=benchmark -bench=. -benchmem ./standard/...

# 特定のベンチマークのみ実行
docker exec go_app go test -tags=benchmark -bench=BenchmarkGetAll -benchmem ./standard/...

# 実行時間を指定（デフォルトは1秒）
docker exec go_app go test -tags=benchmark -bench=. -benchtime=5s -benchmem ./standard/...
```

### ベンチマーク内容

各パッケージで以下のベンチマークを実施：
- `BenchmarkGetAll` - 全ユーザー取得
- `BenchmarkGetByID` - ID指定取得
- `BenchmarkCreate` - ユーザー作成
- `BenchmarkUpdate` - ユーザー更新
- `BenchmarkConcurrentReads` - 並行読み取り
- `BenchmarkConcurrentWrites` - 並行書き込み
- `BenchmarkBulkInsert` - 大量挿入（10/100/1000件）

### 結果の見方

```
BenchmarkGetAll-8    5000    250000 ns/op    1024 B/op    20 allocs/op
```

- `5000` - 実行回数
- `250000 ns/op` - 1操作あたりのナノ秒（値が小さいほど高速）
- `1024 B/op` - 1操作あたりのメモリ使用量
- `20 allocs/op` - 1操作あたりのメモリアロケーション数

### 通常テストとの分離

ベンチマークテストは`//go:build benchmark`タグで分離されているため：
- 通常のテスト: `go test ./...` ではベンチマークは実行されない
- ベンチマーク: `go test -tags=benchmark -bench=. ./...` で実行