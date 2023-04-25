package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	UserStoreRepo    repositories.UserStoreRepository
	StoreRepo        repositories.StoreRepository
	ProductStoreRepo repositories.ProductStoreRepository
}

func (ctrl *StoreController) GetDefaultStore(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStore, err := ctrl.UserStoreRepo.FindDefaultUserStore(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Default store not found"})
		return
	}

	c.JSON(http.StatusOK, userStore.Store)
}

func (ctrl *StoreController) SetDefaultStore(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	store, err := ctrl.StoreRepo.FindByID(storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	userStore, err := ctrl.UserStoreRepo.FindDefaultUserStore(userID)
	if err == nil {
		userStore.IsEnable = false
		_ = ctrl.UserStoreRepo.Save(userStore)
	}

	newUserStore := &models.UserStore{
		UserID:   userID,
		StoreID:  storeID,
		IsEnable: true,
	}

	err = ctrl.UserStoreRepo.Save(newUserStore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting default store"})
		return
	}

	c.JSON(http.StatusOK, store)
}

func (ctrl *StoreController) GetAllStores(c *gin.Context) {
	stores, err := ctrl.StoreRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving all stores"})
		return
	}

	c.JSON(http.StatusOK, stores)
}

func (sc *StoreController) DeleteDefaultStore(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)
	err := sc.UserStoreRepo.DisableDefaultStore(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting default store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default store deleted successfully"})
}

func (sc *StoreController) GetProductsByStoreID(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	productStores, err := sc.ProductStoreRepo.FindAllByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}

	products := make([]*models.Product, len(productStores))
	for i, productStore := range productStores {
		products[i] = productStore.Product
	}

	c.JSON(http.StatusOK, products)
}

func getUserIDFromContext(c *gin.Context) (int64, error) {
	user, exists := c.MustGet("user").(*models.User)
	if !exists {
		return 0, errors.New("User not found in context")
	}

	return user.ID, nil
}
