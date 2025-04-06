package handlers

import (
	"fmt"
	"io"
	"os"
	"seka_back_last/db"
	"time"

	"github.com/gin-gonic/gin"
)

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
		c.JSON(500, gin.H{"error": "Ошибка добавления машины"})
		return
	}

	c.JSON(200, gin.H{"message": "Машина добавлена"})
}