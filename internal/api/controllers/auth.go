package controllers

import (
	"diploma/internal/api/services"
	"diploma/internal/errs"
	"diploma/internal/logger"
	"diploma/internal/requests"
	"diploma/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) AuthController {
	return AuthController{
		service: s,
	}
}

func (c AuthController) Login(ctx *gin.Context) {
	logger.Log("AuthController::Login")
	var userRequest requests.User

	// Parse request
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Try to login
	if err := c.service.Login(ctx, userRequest); err != nil {
		logger.Log(err.Error())
		if errors.Is(err, errs.ErrLoginOrPasswordNotFound) {
			utils.ErrorJSON(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		utils.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c AuthController) Register(ctx *gin.Context) {
	logger.Log("AuthController::Register")
	var userRequest requests.User

	// Parse request
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Register user
	if err := c.service.Register(ctx, userRequest); err != nil {
		logger.Log(err.Error())
		if errors.Is(err, errs.ErrLoginUniqueViolation) {
			utils.ErrorJSON(ctx, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
