package model

import "time"

// User ユーザーモデル
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	// GetAll 全ユーザーを取得
	GetAll() ([]User, error)

	// GetByID IDでユーザーを取得
	GetByID(id int) (*User, error)

	// Create 新規ユーザーを作成
	Create(name, email string) (*User, error)

	// Update ユーザー情報を更新
	Update(id int, name, email string) error

	// Delete ユーザーを削除
	Delete(id int) error

	// Close データベース接続を閉じる
	Close() error
}
