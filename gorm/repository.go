package gorm

import (
	"go_sql_library/model"
	"time"

	"gorm.io/gorm"
)

// User GORMモデル（model.Userとは別に定義）
type User struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}

// UserRepository GORMを使ったユーザーリポジトリ
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository リポジトリの初期化
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// toModelUser GORM UserをmodelのUserに変換
func toModelUser(u *User) *model.User {
	return &model.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// GetAll 全ユーザーを取得
func (r *UserRepository) GetAll() ([]model.User, error) {
	var gormUsers []User
	if err := r.db.Order("id").Find(&gormUsers).Error; err != nil {
		return nil, err
	}

	users := make([]model.User, len(gormUsers))
	for i, u := range gormUsers {
		users[i] = *toModelUser(&u)
	}
	return users, nil
}

// GetByID IDでユーザーを取得
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return toModelUser(&u), nil
}

// Create 新規ユーザーを作成
func (r *UserRepository) Create(name, email string) (*model.User, error) {
	u := User{
		Name:  name,
		Email: email,
	}
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return toModelUser(&u), nil
}

// Update ユーザー情報を更新
func (r *UserRepository) Update(id int, name, email string) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(User{
		Name:  name,
		Email: email,
	}).Error
}

// Delete ユーザーを削除
func (r *UserRepository) Delete(id int) error {
	return r.db.Delete(&User{}, id).Error
}

// Close データベース接続を閉じる
func (r *UserRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
