package controllers

import (
	"net/http"

	"github.com/atomi-ai/atomi/services"

	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SetCurrentPaymentMethod(c *gin.Context)
	GetUser(c *gin.Context)
}

type UserControllerImpl struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (uc *UserControllerImpl) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	c.JSON(http.StatusOK, user)
}

func (uc *UserControllerImpl) SetCurrentPaymentMethod(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	paymentMethodID := c.Param("paymentMethodId")

	updatedUser, err := uc.UserService.SetCurrentPaymentMethod(user, &paymentMethodID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
