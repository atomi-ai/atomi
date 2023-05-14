package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"

	"github.com/gin-gonic/gin"
)

type StoreController interface {
	GetDefaultStore(c *gin.Context)
	SetDefaultStore(c *gin.Context)
	GetAllStores(c *gin.Context)
	DeleteDefaultStore(c *gin.Context)
	GetProductsByStoreID(c *gin.Context)
	GetStoreInfo(c *gin.Context)
}

type StoreControllerImpl struct {
	ManagerStoreRepository repositories.ManagerStoreRepository
	ProductStoreRepo       repositories.ProductStoreRepository
	StoreRepo              repositories.StoreRepository
	UserStoreRepo          repositories.UserStoreRepository
}

func NewStoreController(
	managerStoreRep repositories.ManagerStoreRepository,
	psRepo repositories.ProductStoreRepository,
	storeRepository repositories.StoreRepository,
	usRepo repositories.UserStoreRepository) StoreController {
	return &StoreControllerImpl{
		ManagerStoreRepository: managerStoreRep,
		ProductStoreRepo:       psRepo,
		StoreRepo:              storeRepository,
		UserStoreRepo:          usRepo,
	}
}

func (sc *StoreControllerImpl) GetDefaultStore(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStore, err := sc.UserStoreRepo.FindDefaultUserStore(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Default store not found"})
		return
	}

	c.JSON(http.StatusOK, userStore.Store)
}

func (sc *StoreControllerImpl) SetDefaultStore(c *gin.Context) {
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

	store, err := sc.StoreRepo.FindByID(storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	userStore, err := sc.UserStoreRepo.FindDefaultUserStore(userID)
	if err == nil {
		userStore.IsEnable = false
		_ = sc.UserStoreRepo.Save(userStore)
	}

	newUserStore := &models.UserStore{
		UserID:   userID,
		StoreID:  storeID,
		IsEnable: true,
	}

	err = sc.UserStoreRepo.Save(newUserStore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting default store"})
		return
	}

	c.JSON(http.StatusOK, store)
}

func (sc *StoreControllerImpl) GetAllStores(c *gin.Context) {
	stores, err := sc.StoreRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving all stores"})
		return
	}

	c.JSON(http.StatusOK, stores)
}

func (sc *StoreControllerImpl) DeleteDefaultStore(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)
	err := sc.UserStoreRepo.DisableDefaultStore(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting default store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default store deleted successfully"})
}

func (sc *StoreControllerImpl) GetProductsByStoreID(c *gin.Context) {
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

func (sc *StoreControllerImpl) GetStoreInfo(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	store, err := sc.StoreRepo.FindByID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}
