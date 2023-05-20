package controllers

import (
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type OrderController interface {
	GetUserOrders(c *gin.Context)
	AddOrderForUser(c *gin.Context)
	UberQuote(c *gin.Context)
	GetDelivery(c *gin.Context)
	CreateDelivery(c *gin.Context)
	GetTaxRate(c *gin.Context)
}

type OrderControllerImpl struct {
	OrderService   services.OrderService
	UberService    services.UberService
	TaxRateService services.TaxRateService
}

func NewOrderController(orderService services.OrderService, uberService services.UberService, taxRateService services.TaxRateService) OrderController {
	return &OrderControllerImpl{
		OrderService:   orderService,
		UberService:    uberService,
		TaxRateService: taxRateService,
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

func (oc *OrderControllerImpl) UberQuote(c *gin.Context) {
	var requestBody models.QuoteRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := oc.UberService.Quote(&requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (oc *OrderControllerImpl) GetDelivery(c *gin.Context) {
	deliveryID := c.Param("deliveryID")
	response, err := oc.UberService.GetDelivery(deliveryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (oc *OrderControllerImpl) CreateDelivery(c *gin.Context) {
	var requestBody models.DeliveryData
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据你的配置文件设置测试模式
	testMode := viper.GetBool("testMode")

	// 如果处于测试模式，则在requestBody中插入字段
	if testMode {
		requestBody.TestSpecifications = &models.TestSpecifications{
			RoboCourierSpecification: models.RoboCourierSpecification{
				Mode: "auto",
			},
		}
	}

	response, err := oc.UberService.CreateDelivery(&requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (oc *OrderControllerImpl) GetTaxRate(c *gin.Context) {
	var address models.Address
	if err := c.BindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taxRate, err := oc.TaxRateService.GetTaxRateByZipCodeAndState(&address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, taxRate)
}
