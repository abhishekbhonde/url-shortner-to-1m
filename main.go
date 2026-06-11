package main

import (
	"log"

	"c/database"
	"c/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	router := gin.Default()

	routes.SetupRoutes(router)

	log.Println("Server running on port 8080")

	router.Run(":8080")
}