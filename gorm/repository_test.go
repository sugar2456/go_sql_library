package gorm

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("root:password@tcp(%s:%s)/testdb?parseTime=true&charset=utf8mb4",
		dbHost, dbPort)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("データベース接続エラー: %v", err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("データベースPingエラー: %v", err)
	}

	return db
}

func cleanupTestData(t *testing.T, db *gorm.DB) {
	db.Exec("DELETE FROM users WHERE email LIKE 'test%@example.com'")
}

func TestUserRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	repo := NewUserRepository(db)

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll エラー: %v", err)
	}

	if len(users) == 0 {
		t.Error("ユーザーが取得できませんでした")
	}

	t.Logf("取得したユーザー数: %d", len(users))
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	defer cleanupTestData(t, db)

	repo := NewUserRepository(db)

	user, err := repo.Create("GORMテストユーザー", "test_gorm@example.com")
	if err != nil {
		t.Fatalf("Create エラー: %v", err)
	}

	if user.ID == 0 {
		t.Error("ユーザーIDが設定されていません")
	}

	if user.Name != "GORMテストユーザー" {
		t.Errorf("期待する名前: GORMテストユーザー, 実際: %s", user.Name)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	defer cleanupTestData(t, db)

	repo := NewUserRepository(db)

	created, err := repo.Create("GORMテストユーザー2", "test_gorm2@example.com")
	if err != nil {
		t.Fatalf("Create エラー: %v", err)
	}

	user, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID エラー: %v", err)
	}

	if user.ID != created.ID {
		t.Errorf("期待するID: %d, 実際: %d", created.ID, user.ID)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	defer cleanupTestData(t, db)

	repo := NewUserRepository(db)

	created, err := repo.Create("GORM更新前", "test_gorm3@example.com")
	if err != nil {
		t.Fatalf("Create エラー: %v", err)
	}

	err = repo.Update(created.ID, "GORM更新後", "test_gorm3_updated@example.com")
	if err != nil {
		t.Fatalf("Update エラー: %v", err)
	}

	user, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID エラー: %v", err)
	}

	if user.Name != "GORM更新後" {
		t.Errorf("期待する名前: GORM更新後, 実際: %s", user.Name)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	repo := NewUserRepository(db)

	created, err := repo.Create("GORM削除用", "test_gorm4@example.com")
	if err != nil {
		t.Fatalf("Create エラー: %v", err)
	}

	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete エラー: %v", err)
	}

	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("削除したユーザーが取得できてしまいました")
	}
}
