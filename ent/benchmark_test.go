//go:build benchmark
// +build benchmark

package ent

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func setupBenchDB(b *testing.B) *sql.DB {
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

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		b.Fatalf("データベース接続エラー: %v", err)
	}

	if err := db.Ping(); err != nil {
		b.Fatalf("データベースPingエラー: %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return db
}

func cleanupBenchData(b *testing.B, db *sql.DB) {
	_, err := db.Exec("DELETE FROM users WHERE email LIKE 'bench%@example.com'")
	if err != nil {
		b.Errorf("ベンチマークデータクリーンアップエラー: %v", err)
	}
}

func BenchmarkGetAll(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()

	repo := NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetAll()
		if err != nil {
			b.Fatalf("GetAll エラー: %v", err)
		}
	}
}

func BenchmarkGetByID(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()

	repo := NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetByID(1)
		if err != nil {
			b.Fatalf("GetByID エラー: %v", err)
		}
	}
}

func BenchmarkCreate(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	defer cleanupBenchData(b, db)

	repo := NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("bench%d@example.com", i)
		_, err := repo.Create("entベンチ", email)
		if err != nil {
			b.Fatalf("Create エラー: %v", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	defer cleanupBenchData(b, db)

	repo := NewUserRepository(db)

	user, err := repo.Create("更新ベンチ", "bench_update@example.com")
	if err != nil {
		b.Fatalf("テストデータ作成エラー: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.Update(user.ID, fmt.Sprintf("更新%d", i), "bench_update@example.com")
		if err != nil {
			b.Fatalf("Update エラー: %v", err)
		}
	}
}

func BenchmarkConcurrentReads(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()

	repo := NewUserRepository(db)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.GetAll()
			if err != nil {
				b.Errorf("GetAll エラー: %v", err)
			}
		}
	})
}

func BenchmarkConcurrentWrites(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	defer cleanupBenchData(b, db)

	repo := NewUserRepository(db)

	var counter int
	var mu sync.Mutex

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			email := fmt.Sprintf("bench_concurrent%d@example.com", counter)
			mu.Unlock()

			_, err := repo.Create("並行ベンチ", email)
			if err != nil {
				b.Errorf("Create エラー: %v", err)
			}
		}
	})
}

func BenchmarkBulkInsert(b *testing.B) {
	db := setupBenchDB(b)
	defer db.Close()
	defer cleanupBenchData(b, db)

	repo := NewUserRepository(db)

	counts := []int{10, 100, 1000}
	for _, count := range counts {
		b.Run(fmt.Sprintf("Insert%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				cleanupBenchData(b, db)
				b.StartTimer()

				for j := 0; j < count; j++ {
					email := fmt.Sprintf("bench_bulk%d_%d@example.com", i, j)
					_, err := repo.Create("一括ベンチ", email)
					if err != nil {
						b.Fatalf("Create エラー: %v", err)
					}
				}
			}
		})
	}
}
