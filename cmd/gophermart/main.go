package main

import (
	"diploma/internal/api/controllers"
	"diploma/internal/api/repositories"
	"diploma/internal/api/routes"
	"diploma/internal/api/services"
	"diploma/internal/config"
	"diploma/internal/drivers"
	"diploma/internal/logger"
	"diploma/internal/models"
)

func main() {
	logger.Init()
	router := drivers.NewGinRouter()
	db := drivers.NewDatabase()

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)
	authRoute := routes.NewAuthRoute(authController, router)
	authRoute.Setup()

	ordersRepository := repositories.NewOrdersRepository(db)
	ordersService := services.NewOrdersService(ordersRepository)
	ordersController := controllers.NewOrdersController(ordersService)
	ordersRoute := routes.NewOrdersRoute(ordersController, router)
	ordersRoute.Setup()

	balanceRepository := repositories.NewBalanceRepository(db)
	balanceService := services.NewBalanceService(balanceRepository, ordersRepository)
	balanceController := controllers.NewBalanceController(balanceService)
	balanceRoute := routes.NewBalanceRoute(balanceController, router)
	balanceRoute.Setup()

	if err := db.DB.AutoMigrate(&models.User{}, &models.Order{}, &models.Withdrawal{}); err != nil {
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
