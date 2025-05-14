package main

import (	
	"seka_back_last/config"
	"seka_back_last/db"
	"seka_back_last/routes"
	"seka_back_last/utils"
	"seka_back_last/ws"
	
	// "seka_backend2/routes"
	// "seka_backend2/utils"
	// "seka_backend2/ws"
)

func main() {
	config.LoadEnv()
	utils.EnsureUploadsFolder()

	db.InitDB()
	defer db.CloseDB()

	db.RunMigrations()
	
	hub := ws.NewHub()
	go hub.Run()

	router := routes.SetupRouter(hub)
	utils.StartServerGracefully(router, "8383")
}
