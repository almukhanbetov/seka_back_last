package utils

import (
	"log"
	"os"
)

func EnsureUploadsFolder() {
	path := "uploads"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatalf("❌ Ошибка создания папки uploads: %v", err)
		}
		log.Println("📁 Папка uploads создана")
	}
}
