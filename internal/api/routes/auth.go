package auth

import (
	auth "diploma/internal/api/controllers"
	"diploma/internal/drivers"
)

type AuthRoute struct {
	Controller auth.Controller
	Handler    drivers.GinRouter
}

func NewRoute(controller auth.Controller, handler drivers.GinRouter) AuthRoute {
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
