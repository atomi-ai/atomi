package controllers

import (
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddressController interface {
	GetAllAddressesForUser(c *gin.Context)
	AddAddressForUser(c *gin.Context)
	DeleteAddressForUser(c *gin.Context)
	SetDefaultShippingAddress(c *gin.Context)
	SetDefaultBillingAddress(c *gin.Context)
	GetDefaultShippingAddress(c *gin.Context)
	GetDefaultBillingAddress(c *gin.Context)
	DeleteAllAddressesForUser(c *gin.Context)
}

type AddressControllerImpl struct {
	AddressRepo    repositories.AddressRepository
	AddressService services.AddressService
	UserService    services.UserService
}

func NewAddressControl(addressServoce services.AddressService, userService services.UserService, addressRepo repositories.AddressRepository) AddressController {
	return &AddressControllerImpl{
		AddressRepo:    addressRepo,
		AddressService: addressServoce,
		UserService:    userService,
	}
}

func (ac *AddressControllerImpl) GetAllAddressesForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	addresses, err := ac.AddressService.GetAddressesByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (ac *AddressControllerImpl) AddAddressForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedAddress, err := ac.AddressService.AddAddressForUser(user, &address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, savedAddress)
}

func (ac *AddressControllerImpl) DeleteAddressForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	err := ac.AddressService.DeleteAddressForUser(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AddressControllerImpl) SetDefaultShippingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	updatedUser, err := ac.UserService.SetDefaultShippingAddress(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ac *AddressControllerImpl) SetDefaultBillingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	updatedUser, err := ac.UserService.SetDefaultBillingAddress(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ac *AddressControllerImpl) GetDefaultShippingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	address, err := ac.AddressRepo.FindByID(user.DefaultShippingAddressID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, address)
}

func (ac *AddressControllerImpl) GetDefaultBillingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	address, err := ac.AddressRepo.FindByID(user.DefaultBillingAddressID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, address)
}

func (ac *AddressControllerImpl) DeleteAllAddressesForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	err := ac.AddressService.DeleteAllAddressesForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
