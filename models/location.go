package handlers

import (
	"net/http"
	"seka_back_last/db"

	"github.com/gin-gonic/gin"
)

func LogLocation(c *gin.Context) {
	var req struct {
		UserID    int     `json:"user_id"`
		TripID    int     `json:"trip_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.BindJSON(&req); err != nil ||
		req.UserID == 0 || req.TripID == 0 ||
		req.Latitude == 0 || req.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO location_logs (user_id, trip_id, latitude, longitude)
		VALUES ($1, $2, $3, $4)`,
		req.UserID, req.TripID, req.Latitude, req.Longitude)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении координат"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Координата сохранена"})
}

func GetTripRoute(c *gin.Context) {
	tripID := c.Param("id")
	type Point struct {
		Latitude   float64 `db:"latitude" json:"latitude"`
		Longitude  float64 `db:"longitude" json:"longitude"`
		RecordedAt string  `db:"recorded_at" json:"recorded_at"`
	}

	var points []Point
	err := db.DB.Select(&points, `
		SELECT latitude, longitude, recorded_at
		FROM location_logs
		WHERE trip_id = $1
		ORDER BY recorded_at ASC`, tripID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения маршрута"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trip_id": tripID,
		"route":   points,
	})
}
