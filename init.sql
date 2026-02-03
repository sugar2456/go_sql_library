-- データベース初期化スクリプト
USE testdb;

-- サンプルテーブルの作成
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- サンプルデータの挿入
INSERT INTO users (name, email) VALUES
    ('山田太郎', 'yamada@example.com'),
    ('佐藤花子', 'sato@example.com'),
    ('鈴木一郎', 'suzuki@example.com')
ON DUPLICATE KEY UPDATE name=name;

-- インデックスの作成
CREATE INDEX idx_email ON users(email);
