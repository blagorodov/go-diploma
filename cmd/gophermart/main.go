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
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	srv := runServer(router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log("Server shutdown failure:")
		logger.Log(err.Error())
		return
	}
	select {
	case <-ctx.Done():
		logger.Log("Timeout exceeded, shutting down server...")
	}
	logger.Log("Server exiting")
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
	go func() {
		if err := ordersService.RunPollingStatuses(mainCtx); err != nil {
			logger.Log("Failed polling statuses")
		}
	}()
}

func runServer(router drivers.GinRouter) *http.Server {
	srv := &http.Server{
		Addr:    config.Options.ServerAddress,
		Handler: router.Gin.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log("Error starting server")
			logger.Log(err.Error())
		}
	}()

	return srv
}
