package controllers

import (
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
)

type OrderController interface {
	GetUserOrders(c *gin.Context)
	AddOrderForUser(c *gin.Context)
}

type OrderControllerImpl struct {
	OrderService services.OrderService
}

func NewOrderController(orderService services.OrderService) OrderController {
	return &OrderControllerImpl{
		OrderService: orderService,
	}
}

func (oc *OrderControllerImpl) GetUserOrders(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	orders, err := oc.OrderService.GetUserOrders(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (oc *OrderControllerImpl) AddOrderForUser(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedOrder, err := oc.OrderService.AddOrderForUser(user, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, savedOrder)
}
