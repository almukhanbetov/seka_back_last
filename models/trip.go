package handlers

import (
	"net/http"
	"seka_back_last/db"

	"github.com/gin-gonic/gin"
)

func StartTrip(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id обязателен"})
		return
	}

	var tripID int
	err := db.DB.QueryRowx(
		`INSERT INTO trips (user_id) VALUES ($1) RETURNING id`, req.UserID,
	).Scan(&tripID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании поездки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip_id": tripID})
}

func EndTrip(c *gin.Context) {
	var req struct {
		TripID int `json:"trip_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.TripID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip_id обязателен"})
		return
	}

	res, err := db.DB.Exec(
		`UPDATE trips SET end_time = NOW(), status = 'completed' WHERE id = $1`, req.TripID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка завершения поездки"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Поездка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Поездка завершена"})
}

func RateTrip(c *gin.Context) {
	var req struct {
		TripID int `json:"trip_id"`
		Rating int `json:"rating"`
	}
	if err := c.BindJSON(&req); err != nil || req.TripID == 0 || req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Оценка от 1 до 5 обязательна"})
		return
	}

	var status string
	err := db.DB.Get(&status, `SELECT status FROM trips WHERE id = $1`, req.TripID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Поездка не найдена"})
		return
	}

	if status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Оценку можно ставить только завершённым поездкам"})
		return
	}

	_, err = db.DB.Exec(`UPDATE trips SET rating = $1 WHERE id = $2`, req.Rating, req.TripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения оценки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Оценка сохранена"})
}
