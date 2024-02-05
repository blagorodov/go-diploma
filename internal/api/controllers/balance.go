package controllers

import (
	"diploma/internal/api/services"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/logger"
	"diploma/internal/requests"
	"diploma/internal/responses"
	"diploma/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BalanceController struct {
	service services.BalanceService
}

func NewBalanceController(s services.BalanceService) BalanceController {
	return BalanceController{
		service: s,
	}
}

func (c BalanceController) Get(ctx *gin.Context) {
	logger.Log("BalanceController::Get")

	// Check user auth
	if _, err := auth.GetID(ctx); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, nil)
		return
	}

	balance, withdrawn, err := c.service.Get(ctx)

	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, nil)
		return
	}

	response := responses.Balance{
		Current:   balance,
		Withdrawn: withdrawn,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c BalanceController) Withdraw(ctx *gin.Context) {
	logger.Log("BalanceController::Withdraw")

	// Check user auth
	if _, err := auth.GetID(ctx); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, nil)
		return
	}

	// Parse request
	var withdrawRequest requests.Withdraw
	if err := ctx.ShouldBindJSON(&withdrawRequest); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := c.service.Withdraw(ctx, withdrawRequest.Order, withdrawRequest.Sum)

	// обработка 402
	if errors.Is(err, errs.ErrBalanceNotEnoughFunds) {
		utils.ErrorJSON(ctx, http.StatusPaymentRequired, nil)
		return
	}
	// обработка 422
	if errors.Is(err, errs.ErrOrderNotFound) {
		utils.ErrorJSON(ctx, http.StatusUnprocessableEntity, nil)
		return
	}
	// обработка 500
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, nil)
	}
}

func (c BalanceController) Withdrawals(ctx *gin.Context) {
	logger.Log("BalanceController::Withdrawals")

	// Check user auth
	if _, err := auth.GetID(ctx); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, nil)
		return
	}

	withdrawals, err := c.service.Withdrawals(ctx)
	if len(withdrawals) == 0 {
		utils.ErrorJSON(ctx, http.StatusNoContent, nil)
	}
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, nil)
		return
	}

	var response []responses.Withdrawal
	for _, withdrawal := range withdrawals {
		item := responses.Withdrawal{
			Order:       withdrawal.OrderNumber,
			Sum:         float64(withdrawal.Amount) / 100,
			ProcessedAt: withdrawal.ProcessedAt,
		}
		response = append(response, item)
	}
	ctx.JSON(http.StatusOK, response)
}
