package routes

import (
	"diploma/internal/api/controllers"
	"diploma/internal/drivers"
)

type OrdersRoute struct {
	Controller controllers.OrdersController
	Handler    drivers.GinRouter
}

func NewOrdersRoute(controller controllers.OrdersController, handler drivers.GinRouter) OrdersRoute {
	return OrdersRoute{
		Controller: controller,
		Handler:    handler,
	}
}

func (r OrdersRoute) Setup() {
	r.Handler.Gin.POST("/api/user/orders", r.Controller.Add)
	r.Handler.Gin.GET("/api/user/orders", r.Controller.List)
}
