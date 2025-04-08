package handlers

import (
	"net/http"
	"seka_back_last/db"

	"github.com/gin-gonic/gin"
)
func GetTripHistory(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, user_id, start_time, end_time, duration, distance
		FROM trips
		WHERE status = 'completed'
		ORDER BY end_time DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка получения истории"})
		return
	}

	var trips []map[string]interface{}
	for rows.Next() {
		var id, userID, duration int
		var distance float64
		var startTime, endTime string
		_ = rows.Scan(&id, &userID, &startTime, &endTime, &duration, &distance)

		trips = append(trips, gin.H{
			"id":         id,
			"user_id":    userID,
			"start_time": startTime,
			"end_time":   endTime,
			"duration":   duration,
			"distance":   distance,
		})
	}

	c.JSON(200, gin.H{"trips": trips})
}
func GetTripSummary(c *gin.Context) {
	var totalTrips int
	var totalDuration int
	var totalDistance float64

	err := db.DB.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(duration), 0), COALESCE(SUM(distance), 0)
		FROM trips WHERE status = 'completed'
	`).Scan(&totalTrips, &totalDuration, &totalDistance)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при подсчёте статистики"})
		return
	}

	c.JSON(200, gin.H{
		"totalTrips":    totalTrips,
		"totalDuration": totalDuration,
		"totalDistance": totalDistance,
	})
}

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
		TripID   int     `json:"trip_id"`
		Duration int     `json:"duration"` // в секундах
		Distance float64 `json:"distance"` // в метрах
	}
	if err := c.BindJSON(&req); err != nil || req.TripID == 0 {
		c.JSON(400, gin.H{"error": "trip_id обязателен"})
		return
	}

	res, err := db.DB.Exec(
		`UPDATE trips SET end_time = NOW(), status = 'completed', duration = $2, distance = $3 WHERE id = $1`,
		req.TripID, req.Duration, req.Distance,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка завершения поездки"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(404, gin.H{"error": "Поездка не найдена"})
		return
	}

	c.JSON(200, gin.H{"message": "Поездка завершена"})
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
func GetActiveDrivers(c *gin.Context) {
	type DriverWithCar struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Status   int     `json:"status"`
		Image    string  `json:"image"`
		CarBrand string  `json:"car_brand"`
		CarModel string  `json:"car_model"`
	}

	rows, err := db.DB.Query(`
		SELECT u.id, u.name, u.latitude, u.longitude, u.status, u.image, c.brand, c.model
		FROM users u
		LEFT JOIN driver_car dc ON u.id = dc.driver_id
		LEFT JOIN cars c ON dc.car_id = c.id
		WHERE u.latitude IS NOT NULL AND u.longitude IS NOT NULL
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка запроса"})
		return
	}

	var drivers []DriverWithCar
	for rows.Next() {
		var d DriverWithCar
		_ = rows.Scan(&d.ID, &d.Name, &d.Latitude, &d.Longitude, &d.Status, &d.Image, &d.CarBrand, &d.CarModel)
		drivers = append(drivers, d)
	}

	c.JSON(200, gin.H{"drivers": drivers})
}

