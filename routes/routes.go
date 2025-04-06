package routes

import (
	"github.com/gin-gonic/gin"
	"seka_back_last/handlers"
	"seka_back_last/ws"
	"github.com/gin-contrib/cors"
)

func SetupRouter(hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.GET("/ws", func(c *gin.Context) {
		handlers.HandleWebSocket(hub, c)
	})
	r.Use(cors.Default())

	api := r.Group("/api")
	{
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		api.POST("/trip/start", handlers.StartTrip)
		api.POST("/trip/end", handlers.EndTrip)
		api.POST("/trip/rate", handlers.RateTrip)
		api.GET("/trip/:id/route", handlers.GetTripRoute)

		api.POST("/location", handlers.LogLocation)

		api.POST("/driver", handlers.CreateDriver)
		api.GET("/drivers", handlers.GetDrivers)

		api.POST("/car", handlers.CreateCar)
		api.POST("/assign-car", handlers.AssignCarToDriver)
	}

	return r
}
