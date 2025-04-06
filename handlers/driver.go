package handlers

import (
	"fmt"
	"io"
	"os"
	"seka_back_last/db"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateDriver(c *gin.Context) {
	// Получаем обычные поля
	name := c.PostForm("name")
	email := c.PostForm("email")

	if name == "" {
		c.JSON(400, gin.H{"error": "Имя обязательно"})
		return
	}

	// Обработка изображения
	var imagePath string

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		imagePath = "uploads/no-image.png"
	} else {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)
		out, err := os.Create(filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка сохранения файла"})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		imagePath = filename
	}

	// Вставка в БД
	_, err = db.DB.Exec(`
		INSERT INTO users (name, email, image) 
		VALUES ($1, $2, $3)
	`, name, email, imagePath)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка создания водителя"})
		return
	}

	c.JSON(200, gin.H{"message": "Водитель добавлен"})
}
func GetDrivers(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, name FROM users ORDER BY id`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка получения водителей"})
		return
	}

	var drivers []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		_ = rows.Scan(&id, &name)
		drivers = append(drivers, gin.H{"id": id, "name": name})
	}
	c.JSON(200, gin.H{"drivers": drivers})
}
