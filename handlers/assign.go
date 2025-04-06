package handlers

import (
	"seka_back_last/db"

	"github.com/gin-gonic/gin"
)

func AssignCarToDriver(c *gin.Context) {
	var req struct {
		DriverID int `json:"driver_id"`
		CarID    int `json:"car_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.DriverID == 0 || req.CarID == 0 {
		c.JSON(400, gin.H{"error": "Нужен driver_id и car_id"})
		return
	}

	_, err := db.DB.Exec(`INSERT INTO driver_car (driver_id, car_id) VALUES ($1, $2)`, req.DriverID, req.CarID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при назначении машины"})
		return
	}

	c.JSON(200, gin.H{"message": "Машина назначена водителю"})
}
