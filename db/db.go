package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sqlx.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("❌ DB_URL не указан в .env")
	}

	var err error
	DB, err = sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("❌ Пинг БД не прошёл: %v", err)
	}

	log.Println("✅ БД подключена")
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("⚠️ Ошибка при закрытии БД: %v", err)
		} else {
			log.Println("🛑 БД соединение закрыто")
		}
	}
}
