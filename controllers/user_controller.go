package controllers

import (
	"github.com/atomi-ai/atomi/repositories"
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SetCurrentPaymentMethod(c *gin.Context)
	GetUser(c *gin.Context)
}

type UserControllerImpl struct {
	UserRepo repositories.UserRepository
}

func NewUserController(userRepo repositories.UserRepository) UserController {
	return &UserControllerImpl{
		UserRepo: userRepo,
	}
}

func (sc *UserControllerImpl) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	c.JSON(http.StatusOK, user)
}

func (uc *UserControllerImpl) SetCurrentPaymentMethod(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	paymentMethodID := c.Param("paymentMethodId")

	// 更新用户的支付方法 ID
	user.PaymentMethodID = paymentMethodID

	// 保存更新后的用户
	err := uc.UserRepo.Save(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
