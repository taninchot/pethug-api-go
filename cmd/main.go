package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"pethug-api-go/controllers"
	"pethug-api-go/db"
	"pethug-api-go/repositories"
	"pethug-api-go/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	db.ConnectDB()
	defer db.CloseDB()

	router := gin.Default()

	userRepository := repositories.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	api := router.Group("/api/v1")
	userController.RegisterRoutes(api)

	err = router.Run(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Unable to run server: %v\n", err)
		return
	}
}
