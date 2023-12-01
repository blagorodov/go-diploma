package controllers

import (
	"diploma/internal/api/services"
	"diploma/internal/logger"
	"diploma/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrdersController struct {
	service services.OrdersService
}

func NewOrdersController(s services.OrdersService) OrdersController {
	return OrdersController{
		service: s,
	}
}

func (c OrdersController) Add(ctx *gin.Context) {
	logger.Log("OrdersController::Add")

	// Check for 400 - content-type
	if ctx.Request.Header.Get("Content-Type") != "text/plain" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, nil)
		return
	}

	// Get order number from post body
	data, err := ctx.GetRawData()
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, nil)
		return
	}
	orderNumber := string(data)
	fmt.Println(orderNumber)
	//
}

func (c OrdersController) List(ctx *gin.Context) {
	logger.Log("OrdersController::List")

}
