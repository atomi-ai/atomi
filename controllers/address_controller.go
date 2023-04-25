package controllers

import (
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddressController struct {
	AddressService services.AddressService
	UserService    services.UserService
	AddressRepo    repositories.AddressRepository
}

func (ac *AddressController) GetAllAddressesForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	addresses, err := ac.AddressService.GetAddressesByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (ac *AddressController) AddAddressForUser(c *gin.Context) {
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

func (ac *AddressController) DeleteAddressForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	err := ac.AddressService.DeleteAddressForUser(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AddressController) SetDefaultShippingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	updatedUser, err := ac.UserService.SetDefaultShippingAddress(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ac *AddressController) SetDefaultBillingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	addressId, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)

	updatedUser, err := ac.UserService.SetDefaultBillingAddress(user, addressId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ac *AddressController) GetDefaultShippingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	address, err := ac.AddressRepo.FindByID(user.DefaultShippingAddressID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, address)
}

func (ac *AddressController) GetDefaultBillingAddress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	address, err := ac.AddressRepo.FindByID(user.DefaultBillingAddressID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, address)
}

func (ac *AddressController) DeleteAllAddressesForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	err := ac.AddressService.DeleteAllAddressesForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
