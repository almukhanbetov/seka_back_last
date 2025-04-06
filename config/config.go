package config

import (
	"log"
	

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env не найден, используются переменные окружения по умолчанию")
	}
}
