package controllers

import (
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserController interface {
	SetCurrentPaymentMethod(c *gin.Context)
	GetUser(c *gin.Context)
	SubmitDeleteUserRequest(c *gin.Context)
}

type UserControllerImpl struct {
	UserService services.UserService
	DeleteRepo  repositories.DeleteUserRequestRepository
}

func NewUserController(userService services.UserService, deleteRepo repositories.DeleteUserRequestRepository) UserController {
	return &UserControllerImpl{
		UserService: userService,
		DeleteRepo:  deleteRepo,
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
		log.Errorf("Errors in updating payment method of user: %v(%v), err: \n%v", user, paymentMethodID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (uc *UserControllerImpl) SubmitDeleteUserRequest(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	err := uc.DeleteRepo.AddRequest(user.ID)

	if err != nil {
		log.Errorf("Errors in submitting delete user request: %v, err: \n%v", user, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request submitted successfully"})
}
