package controllers

import (
	"diploma/internal/api/services"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/logger"
	"diploma/internal/models"
	"diploma/internal/responses"
	"diploma/internal/utils"
	"errors"
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

	// Check user auth
	if _, err := auth.GetID(ctx); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, nil)
		return
	}

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

	err = c.service.Add(ctx, orderNumber)
	if errors.Is(err, errs.ErrOrderNumberFormat) {
		utils.ErrorJSON(ctx, http.StatusUnprocessableEntity, nil)
		return
	}
	if errors.Is(err, errs.ErrOrderOtherUserDuplicate) {
		utils.ErrorJSON(ctx, http.StatusConflict, nil)
		return
	}
	if errors.Is(err, errs.ErrOrderDuplicate) {
		utils.ErrorJSON(ctx, http.StatusOK, nil)
		return
	}
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, nil)
	}

	ctx.Status(http.StatusAccepted)
}

func (c OrdersController) List(ctx *gin.Context) {
	logger.Log("OrdersController::ListAll")
	if _, err := auth.GetID(ctx); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, nil)
		return
	}
	orders, err := c.service.List(ctx)
	if len(orders) == 0 {
		utils.ErrorJSON(ctx, http.StatusNoContent, nil)
		return
	}
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, nil)
		return
	}
	var response []responses.Order
	for _, order := range orders {
		item := responses.Order{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		}
		if order.Status == models.PROCESSED {
			item.Accrual = order.Accrual / 100
		}
		response = append(response, item)
	}
	ctx.JSON(http.StatusOK, response)
}
