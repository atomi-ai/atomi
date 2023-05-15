package controllers

import (
	"net/http"
	"strconv"

	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	log "github.com/sirupsen/logrus"

	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
)

type ManagerStoreController interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type ManagerStoreControllerImpl struct {
	managerStoreRepository repositories.ManagerStoreRepository
	orderRepository        repositories.OrderRepository
	productRepository      repositories.ProductRepository
	productStoreRepository repositories.ProductStoreRepository
	storeRepository        repositories.StoreRepository
	productStoreService    services.ProductStoreService
}

func NewManagerStoreController(
	managerStoreRepository repositories.ManagerStoreRepository,
	orderRepository repositories.OrderRepository,
	productRepository repositories.ProductRepository,
	productStoreRepository repositories.ProductStoreRepository,
	storeRepository repositories.StoreRepository,
	productStoreService services.ProductStoreService) ManagerStoreController {
	return &ManagerStoreControllerImpl{
		managerStoreRepository: managerStoreRepository,
		orderRepository:        orderRepository,
		productRepository:      productRepository,
		productStoreRepository: productStoreRepository,
		storeRepository:        storeRepository,
		productStoreService:    productStoreService,
	}
}

func (msc *ManagerStoreControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	// TODO(lamuguo): Please be consistent to use camelCase or snake_case. ("order_id" vs. "storeId")
	router.GET("/stores", msc.getStoresForMgr)
	router.POST("/store", msc.createStore)
	router.DELETE("/store/:store_id", msc.deleteStore)
	router.PUT("/store/:store_id", msc.assignStoreToUser)
	router.GET("/products", msc.getProductsForMgr)
	router.PUT("/store/add/:storeId/product/:productId", msc.AddProductToStore)
	router.DELETE("/store/remove/:storeId/product/:productId", msc.RemoveProductFromStore)
	router.POST("/store/:storeId/product", msc.CreateProductInStore)
	router.PUT("/orders/:order_id/status", msc.UpdateOrderStatus)
	router.GET("/store/:storeId/orders", msc.GetOrdersByStoreID)
}

func isAuthorized(user *models.User) bool {
	if user == nil {
		return false
	}
	return user.Role == models.RoleAdmin || user.Role == models.RoleMgr
}
func (msc *ManagerStoreControllerImpl) getStoresForMgr(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	stores, err := msc.managerStoreRepository.GetStoresByManagerID(manager.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stores)
}

func (msc *ManagerStoreControllerImpl) getProductsForMgr(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	products, err := msc.productRepository.FindAllProductsForMgr(manager.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error in queries all products for one user"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (msc *ManagerStoreControllerImpl) createStore(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := msc.managerStoreRepository.Save(&store); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := msc.managerStoreRepository.AssignStoreToUser(store.ID, manager.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, store)
}

func (msc *ManagerStoreControllerImpl) deleteStore(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, err := strconv.ParseInt(ctx.Param("store_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}
	if err := msc.managerStoreRepository.DeleteStore(storeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (msc *ManagerStoreControllerImpl) assignStoreToUser(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, err := strconv.ParseInt(ctx.Param("store_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}
	if err := msc.managerStoreRepository.AssignStoreToUser(storeID, manager.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (msc *ManagerStoreControllerImpl) AddProductToStore(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, _ := strconv.ParseInt(ctx.Param("storeId"), 10, 64)
	productID, _ := strconv.ParseInt(ctx.Param("productId"), 10, 64)

	if !msc.storeRepository.CheckUserHasAccessToStore(manager, storeID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to manage this store"})
		return
	}

	err := msc.productStoreRepository.AddProductToStore(storeID, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to store"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (msc *ManagerStoreControllerImpl) RemoveProductFromStore(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, _ := strconv.ParseInt(ctx.Param("storeId"), 10, 64)
	productID, _ := strconv.ParseInt(ctx.Param("productId"), 10, 64)

	if !msc.storeRepository.CheckUserHasAccessToStore(manager, storeID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to manage this store"})
		return
	}

	err := msc.productStoreRepository.RemoveProductFromStore(storeID, productID)
	if err != nil {
		log.Errorf("Failed to remove product from store: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from store"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (msc *ManagerStoreControllerImpl) CreateProductInStore(c *gin.Context) {
	user, _ := c.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, _ := strconv.ParseInt(c.Param("storeId"), 10, 64)
	if !msc.storeRepository.CheckUserHasAccessToStore(manager, storeID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to manage this store"})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := msc.productStoreService.CreateProductInStore(user.(*models.User), storeID, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

func (msc *ManagerStoreControllerImpl) GetOrdersByStoreID(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	storeID, _ := strconv.ParseInt(ctx.Param("storeId"), 10, 64)
	if !msc.storeRepository.CheckUserHasAccessToStore(manager, storeID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to manage this store"})
		return
	}

	orders, err := msc.orderRepository.GetOrdersByStoreID(storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (msc *ManagerStoreControllerImpl) UpdateOrderStatus(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	orderID, err := strconv.ParseInt(ctx.Param("order_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var input struct {
		Status models.OrderStatus `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = msc.orderRepository.UpdateOrderStatus(orderID, input.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
