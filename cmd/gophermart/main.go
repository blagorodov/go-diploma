package main

import (
	authcontrollers "diploma/internal/api/controllers"
	authrepositories "diploma/internal/api/repositories"
	authroutes "diploma/internal/api/routes"
	authservices "diploma/internal/api/services"
	"diploma/internal/config"
	"diploma/internal/drivers"
	"diploma/internal/logger"
	"diploma/internal/models"
)

func main() {
	logger.Init()
	router := drivers.NewGinRouter()
	db := drivers.NewDatabase()
	authRepository := authrepositories.NewRepository(db)
	authService := authservices.NewService(authRepository)
	authController := authcontrollers.NewController(authService)
	authRoute := authroutes.NewRoute(authController, router)
	authRoute.Setup()

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		logger.Log("Error when automigrate")
		logger.Log(err.Error())
		return
	}
	if err := router.Gin.Run(config.Options.ServerAddress); err != nil {
		logger.Log("Error starting server")
		logger.Log(err.Error())
		return
	}
}
