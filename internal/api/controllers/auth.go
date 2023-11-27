package auth

import (
	authservices "diploma/internal/api/services"
	"diploma/internal/auth"
	"diploma/internal/logger"
	"diploma/internal/models"
	"diploma/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	service authservices.Service
}

func NewController(s authservices.Service) Controller {
	return Controller{
		service: s,
	}
}

func (c Controller) Login(ctx *gin.Context) {
	logger.Log("Controller::Login")
	var userRequest models.User

	// Parse request
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Find user
	userDb, err := c.service.Get(userRequest.Login)
	if errors.Is(err, models.ErrLoginOrPasswordNotFound) {
		logger.Log(err.Error())
		utils.ErrorJSON(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	if !userDb.CheckPasswordHash(userRequest.Password) {
		logger.Log(models.ErrLoginOrPasswordNotFound)
		utils.ErrorJSON(ctx, http.StatusUnauthorized, models.ErrLoginOrPasswordNotFound)
		return
	}

	// Set token & cookie
	auth.SetID(ctx, strconv.Itoa(userDb.ID))

	ctx.JSON(http.StatusOK, nil)
}

func (c Controller) Register(ctx *gin.Context) {
	logger.Log("Controller::Register")
	var user models.User

	// Parse request
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Register user
	if err := user.HashPassword(); err != nil {
		logger.Log(err.Error())
		utils.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.service.Register(&user); err != nil {
		logger.Log(err.Error())
		if errors.Is(err, models.ErrLoginUniqueViolation) {
			utils.ErrorJSON(ctx, http.StatusConflict, err.Error())
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Set token & cookie
	auth.SetID(ctx, strconv.Itoa(user.ID))

	ctx.JSON(http.StatusOK, nil)
}
