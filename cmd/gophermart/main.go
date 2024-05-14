package main

import (
	"context"
	"diploma/internal/api/clients"
	"diploma/internal/api/controllers"
	"diploma/internal/api/repositories"
	"diploma/internal/api/routes"
	"diploma/internal/api/services"
	"diploma/internal/config"
	"diploma/internal/drivers"
	"diploma/internal/logger"
	"diploma/internal/models"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()
	router := drivers.NewGinRouter()
	db := drivers.NewDatabase()

	ordersServices := initApp(db, router)

	if err := db.DB.AutoMigrate(&models.User{}, &models.Order{}, &models.Withdrawal{}); err != nil {
		logger.Log("Error when automigrate")
		logger.Log(err.Error())
		return
	}

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runPolling(ordersServices, mainCtx)

	if err := router.Gin.Run(config.Options.ServerAddress); err != nil {
		logger.Log("Error starting server")
		logger.Log(err.Error())
		return
	}
}

func initApp(db drivers.Database, router drivers.GinRouter) services.OrdersService {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)
	authRoute := routes.NewAuthRoute(authController, router)
	authRoute.Setup()

	accrualClient := clients.NewAccrualClient(config.Options.AccrualAddress)

	ordersRepository := repositories.NewOrdersRepository(db)
	ordersService := services.NewOrdersService(ordersRepository, accrualClient)
	ordersController := controllers.NewOrdersController(ordersService)
	ordersRoute := routes.NewOrdersRoute(ordersController, router)
	ordersRoute.Setup()

	balanceRepository := repositories.NewBalanceRepository(db)
	balanceService := services.NewBalanceService(balanceRepository, ordersRepository)
	balanceController := controllers.NewBalanceController(balanceService)
	balanceRoute := routes.NewBalanceRoute(balanceController, router)
	balanceRoute.Setup()

	return ordersService
}

func runPolling(ordersService services.OrdersService, mainCtx context.Context) {
	go func() error {
		if err := ordersService.RunPollingStatuses(mainCtx); err != nil {
			logger.Log("Failed polling statuses")
			return err
		}
		return nil
	}()
}
