package handlers

import (
	"fmt"
	"io"
	"log"
	"os"
	"seka_back_last/db"
	"time"

	"github.com/gin-gonic/gin"
)

// 🚗 Создание машины
func CreateCar(c *gin.Context) {
	brand := c.PostForm("brand")
	model := c.PostForm("model")

	if brand == "" || model == "" {
		c.JSON(400, gin.H{"error": "Марка и модель обязательны"})
		return
	}

	var imagePath string
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		log.Println("📎 Нет фото — используем placeholder")
		imagePath = "uploads/no-image.png"
	} else {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)
		out, err := os.Create(filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		imagePath = filename
	}

	_, err = db.DB.Exec(`
		INSERT INTO cars (brand, model, image)
		VALUES ($1, $2, $3)
	`, brand, model, imagePath)

	if err != nil {
		log.Printf("❌ Ошибка при создании машины: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка при создании машины"})
		return
	}

	c.JSON(200, gin.H{"message": "Машина добавлена"})
}

// 📋 Получение всех машин
func GetCars(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, brand, model, image FROM cars ORDER BY id`)
	if err != nil {
		log.Printf("❌ Ошибка получения машин: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка получения машин"})
		return
	}

	var cars []map[string]interface{}
	for rows.Next() {
		var id int
		var brand, model, image string

		if err := rows.Scan(&id, &brand, &model, &image); err != nil {
			log.Printf("⚠️ Ошибка Scan: %v", err)
			continue
		}

		cars = append(cars, gin.H{
			"id":    id,
			"brand": brand,
			"model": model,
			"image": image,
		})
	}

	c.JSON(200, gin.H{"cars": cars})
}

// ✏️ Обновление машины
func UpdateCar(c *gin.Context) {
	id := c.Param("id")
	brand := c.PostForm("brand")
	model := c.PostForm("model")

	if brand == "" || model == "" {
		c.JSON(400, gin.H{"error": "Марка и модель обязательны"})
		return
	}

	var imagePath string
	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)
		out, err := os.Create(filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		imagePath = filename
	}

	if imagePath != "" {
		_, err = db.DB.Exec(`
			UPDATE cars SET brand = $1, model = $2, image = $3 WHERE id = $4
		`, brand, model, imagePath, id)
	} else {
		_, err = db.DB.Exec(`
			UPDATE cars SET brand = $1, model = $2 WHERE id = $3
		`, brand, model, id)
	}

	if err != nil {
		log.Printf("❌ Ошибка обновления машины: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка обновления машины"})
		return
	}

	c.JSON(200, gin.H{"message": "Машина обновлена"})
}

// 🗑 Удаление машины
func DeleteCar(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec(`DELETE FROM cars WHERE id = $1`, id)
	if err != nil {
		log.Printf("❌ Ошибка удаления машины: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка удаления машины"})
		return
	}

	c.JSON(200, gin.H{"message": "Машина удалена"})
}
