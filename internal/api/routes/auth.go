package routes

import (
	"diploma/internal/api/controllers"
	"diploma/internal/drivers"
)

type AuthRoute struct {
	Controller controllers.AuthController
	Handler    drivers.GinRouter
}

func NewAuthRoute(controller controllers.AuthController, handler drivers.GinRouter) AuthRoute {
	return AuthRoute{
		Controller: controller,
		Handler:    handler,
	}
}

func (r AuthRoute) Setup() {
	a := r.Handler.Gin.Group("/api/user")
	a.POST("/login", r.Controller.Login)
	a.POST("/register", r.Controller.Register)
}
