package db

import (
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations() {
	if DB == nil {
		log.Fatal("❌ DB не инициализирована")
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	log.Println("🔄 Применение миграций...")
	if err := goose.Up(DB.DB, "db/migrations"); err != nil {
		log.Fatalf("❌ Ошибка миграций: %v", err)
	}
	log.Println("✅ Миграции успешно применены")
}
