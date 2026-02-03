package ent

import (
	"database/sql"
	"errors"
	"go_sql_library/model"
)

// UserRepository entを使ったユーザーリポジトリ
// 注意: 実際のentではコード生成を使用しますが、
// ここではデモ用に簡略化した実装を提供しています
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository リポジトリの初期化
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll 全ユーザーを取得
func (r *UserRepository) GetAll() ([]model.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// GetByID IDでユーザーを取得
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	var u model.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
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
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	return r.db.Close()
}
