package sqlx

import (
	"go_sql_library/model"

	"github.com/jmoiron/sqlx"
)

// UserRepository sqlxを使ったユーザーリポジトリ
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository リポジトリの初期化
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll 全ユーザーを取得
func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY id"
	err := r.db.Select(&users, query)
	return users, err
}

// GetByID IDでユーザーを取得
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	var u model.User
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	err := r.db.Get(&u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create 新規ユーザーを作成
func (r *UserRepository) Create(name, email string) (*model.User, error) {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	result, err := r.db.Exec(query, name, email)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

// Update ユーザー情報を更新
func (r *UserRepository) Update(id int, name, email string) error {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	_, err := r.db.Exec(query, name, email, id)
	return err
}

// Delete ユーザーを削除
func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// Close データベース接続を閉じる
func (r *UserRepository) Close() error {
	return r.db.Close()
}
