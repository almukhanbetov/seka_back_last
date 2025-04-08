package handlers

import (
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"seka_back_last/db"
	"seka_back_last/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ✅ Создание водителя
func CreateDriver(c *gin.Context) {
	log.Println("📤 Пришёл POST /api/driver")

	name := c.PostForm("name")
	email := c.PostForm("email")

	if name == "" {
		c.JSON(400, gin.H{"error": "Имя обязательно"})
		return
	}

	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			c.JSON(400, gin.H{"error": "Некорректный email"})
			return
		}
	}

	var imagePath string
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		log.Printf("📎 Фото не передано, будет дефолт: %v", err)
		imagePath = "uploads/no-image.png"
	} else {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)

		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", os.ModePerm); err != nil {
				log.Printf("❌ Ошибка создания папки uploads: %v", err)
				c.JSON(500, gin.H{"error": "Ошибка подготовки папки"})
				return
			}
		}

		out, err := os.Create(filename)
		if err != nil {
			log.Printf("❌ Ошибка сохранения файла: %v", err)
			c.JSON(500, gin.H{"error": "Ошибка сохранения фото"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			log.Printf("❌ Ошибка копирования файла: %v", err)
			c.JSON(500, gin.H{"error": "Ошибка загрузки фото"})
			return
		}

		imagePath = filename
	}

	var driverID int
	err = db.DB.QueryRow(`
		INSERT INTO users (name, email, image)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, email, imagePath).Scan(&driverID)

	if err != nil {
		log.Printf("❌ Ошибка создания водителя: %v", err)
		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(400, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}
		c.JSON(500, gin.H{"error": "Ошибка при сохранении водителя"})
		return
	}

	c.JSON(200, gin.H{
		"message":   "✅ Водитель добавлен",
		"driver_id": driverID,
	})
}

// ✅ Получение всех водителей
func GetDrivers(c *gin.Context) {
	var drivers []models.Driver

	err := db.DB.Select(&drivers, `SELECT id, name, email, image, status FROM users ORDER BY id`)
	if err != nil {
		log.Printf("❌ Ошибка получения водителей: %v", err)
		c.JSON(500, gin.H{"error": "Ошибка получения водителей"})
		return
	}

	c.JSON(200, gin.H{
		"drivers": drivers,
	})
}


// ✅ Обновление водителя
func UpdateDriver(c *gin.Context) {
	idParam := c.Param("id")
	driverID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID водителя"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")

	if name == "" {
		c.JSON(400, gin.H{"error": "Имя обязательно"})
		return
	}

	// Валидация email, если передан
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			c.JSON(400, gin.H{"error": "Некорректный формат email"})
			return
		}
	}

	// Фото, если обновляется
	imagePath := ""
	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), header.Filename)

		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", os.ModePerm); err != nil {
				c.JSON(500, gin.H{"error": "Ошибка создания папки"})
				return
			}
		}

		out, err := os.Create(filename)
		if err != nil {
			log.Printf("❌ Ошибка записи файла: %v", err)
			c.JSON(500, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			c.JSON(500, gin.H{"error": "Ошибка записи изображения"})
			return
		}
		imagePath = filename
	}

	// Подготовка запроса
	var query string
	var args []interface{}

	if imagePath != "" {
		query = `
			UPDATE users SET name = $1, email = $2, image = $3 WHERE id = $4
		`
		args = []interface{}{name, email, imagePath, driverID}
	} else {
		query = `
			UPDATE users SET name = $1, email = $2 WHERE id = $3
		`
		args = []interface{}{name, email, driverID}
	}

	res, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("❌ Ошибка обновления водителя: %v", err)

		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(400, gin.H{"error": "Email уже используется"})
			return
		}

		c.JSON(500, gin.H{"error": "Ошибка при обновлении водителя"})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(404, gin.H{"error": "Водитель не найден"})
		return
	}

	c.JSON(200, gin.H{"message": "✅ Водитель обновлён"})
}

func AssignCarToDriver(c *gin.Context) {
	var req struct {
		DriverID int   `json:"driver_id"`
		CarIDs   []int `json:"car_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.DriverID == 0 || len(req.CarIDs) == 0 {
		c.JSON(400, gin.H{"error": "driver_id и car_ids обязательны"})
		return
	}

	// Сначала удалим все назначения
	_, _ = db.DB.Exec(`DELETE FROM driver_car WHERE driver_id = $1`, req.DriverID)

	// Вставим новые
	for _, carID := range req.CarIDs {
		_, err := db.DB.Exec(`INSERT INTO driver_car (driver_id, car_id) VALUES ($1, $2)`, req.DriverID, carID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка назначения машины"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Машины назначены водителю"})
}
