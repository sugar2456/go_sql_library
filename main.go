package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_sql_library/model"
	entRepo "go_sql_library/ent"
	gormRepo "go_sql_library/gorm"
	standardRepo "go_sql_library/standard"
	sqlxRepo "go_sql_library/sqlx"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var repo model.UserRepository
var libraryType string

func main() {
	// データベース接続設定
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	libraryType = os.Getenv("LIBRARY_TYPE")
	if libraryType == "" {
		libraryType = "standard" // デフォルトは標準ライブラリ
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 使用するライブラリによって初期化方法を変更
	var err error
	switch libraryType {
	case "standard":
		repo, err = initStandard(dsn)
	case "sqlx":
		repo, err = initSqlx(dsn)
	case "gorm":
		repo, err = initGorm(dsn)
	case "ent":
		repo, err = initEnt(dsn)
	default:
		log.Fatalf("未対応のライブラリタイプ: %s", libraryType)
	}

	if err != nil {
		log.Fatal("データベース初期化エラー:", err)
	}
	defer repo.Close()

	log.Printf("データベース接続成功！（ライブラリ: %s）\n", libraryType)

	// ルーティング設定
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)

	log.Println("サーバーを起動します: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initStandard(dsn string) (model.UserRepository, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("データベース接続待機中... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return standardRepo.NewUserRepository(db), nil
}

func initSqlx(dsn string) (model.UserRepository, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sqlx.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("データベース接続待機中... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return sqlxRepo.NewUserRepository(db), nil
}

func initGorm(dsn string) (model.UserRepository, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			err = sqlDB.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("データベース接続待機中... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return gormRepo.NewUserRepository(db), nil
}

func initEnt(dsn string) (model.UserRepository, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("データベース接続待機中... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return entRepo.NewUserRepository(db), nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go + MySQL アプリケーションへようこそ！\n\n")
	fmt.Fprintf(w, "使用中のライブラリ: %s\n\n", libraryType)
	fmt.Fprintf(w, "利用可能なエンドポイント:\n")
	fmt.Fprintf(w, "  GET  /           - このメッセージ\n")
	fmt.Fprintf(w, "  GET  /ping       - ヘルスチェック\n")
	fmt.Fprintf(w, "  GET  /users      - 全ユーザー取得\n")
	fmt.Fprintf(w, "  GET  /users/{id} - 特定ユーザー取得\n")
	fmt.Fprintf(w, "  POST /users      - ユーザー作成（name, email必須）\n")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case "GET":
		users, err := repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

	case "POST":
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repo.Create(input.Name, input.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// IDを抽出（/users/123 から 123 を取得）
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		user, err := repo.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)

	case "PUT":
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.Update(id, input.Name, input.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, _ := repo.GetByID(id)
		json.NewEncoder(w).Encode(user)

	case "DELETE":
		if err := repo.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
