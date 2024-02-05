package routes

import (
	"diploma/internal/api/controllers"
	"diploma/internal/drivers"
)

type BalanceRoute struct {
	Controller controllers.BalanceController
	Handler    drivers.GinRouter
}

func NewBalanceRoute(controller controllers.BalanceController, handler drivers.GinRouter) BalanceRoute {
	return BalanceRoute{
		Controller: controller,
		Handler:    handler,
	}
}

func (r BalanceRoute) Setup() {
	r.Handler.Gin.GET("/api/user/balance", r.Controller.Get)
	r.Handler.Gin.POST("/api/user/balance/withdraw", r.Controller.Withdraw)
	r.Handler.Gin.GET("/api/user/withdrawals", r.Controller.Withdrawals)
}
