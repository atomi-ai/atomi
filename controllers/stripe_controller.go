package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
)

type StripeController interface {
	AttachPaymentMethodToCustomer(c *gin.Context)
	DeleteAllPaymentMethods(c *gin.Context)
	DeletePaymentMethod(c *gin.Context)
	ListPaymentMethods(c *gin.Context)
	Pay(c *gin.Context)
	ListPaymentIntents(c *gin.Context)
	PaymentIntent(c *gin.Context)
}

type StripeControllerImpl struct {
	UserService   services.UserService
	StripeService services.StripeService
	OrderService  services.OrderService
	UberService   services.UberService
	AddressRepo   repositories.AddressRepository
}

func NewStripeController(userService services.UserService, stripeService services.StripeService, orderService services.OrderService, uberService services.UberService, addressRepo repositories.AddressRepository) StripeController {
	return &StripeControllerImpl{
		UserService:   userService,
		StripeService: stripeService,
		OrderService:  orderService,
		UberService:   uberService,
		AddressRepo:   addressRepo,
	}
}

func (sc *StripeControllerImpl) AttachPaymentMethodToCustomer(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	paymentMethodID := c.Param("paymentMethodId")

	pm, err := sc.StripeService.AttachPaymentMethodToCustomer(user.StripeCustomerID, paymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pm)
}

func (sc *StripeControllerImpl) ListPaymentMethods(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	iter, err := sc.StripeService.ListPaymentMethods(user.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 初始化一个空的PaymentMethod切片
	var paymentMethods []*stripe.PaymentMethod

	// 遍历迭代器并将每个PaymentMethod添加到切片中
	for iter.Next() {
		paymentMethod := iter.PaymentMethod()
		paymentMethods = append(paymentMethods, paymentMethod)
	}

	// 检查迭代器中是否存在错误
	if err := iter.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}

func (sc *StripeControllerImpl) DeletePaymentMethod(c *gin.Context) {
	paymentMethodID := c.Param("paymentMethodId")

	pm, err := sc.StripeService.DeletePaymentMethod(paymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 检查要删除的 paymentMethodID 是否是用户表中的 PaymentMethodID
	user := c.MustGet("user").(*models.User)
	if user.PaymentMethodID != nil && *user.PaymentMethodID == paymentMethodID {
		// 如果匹配，将 PaymentMethodID 设置为 nil 并保存更改
		_, err = sc.UserService.SetCurrentPaymentMethod(user, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, pm)
}

func (sc *StripeControllerImpl) Pay(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var piRequest models.PaymentIntentRequest

	err := json.NewDecoder(c.Request.Body).Decode(&piRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if piRequest.OrderID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order id must be greater than 0"})
		return
	}

	order, err := sc.OrderService.FindOrderByID(piRequest.OrderID)
	if err != nil || order == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found"})
		return
	}

	shippingAddrID := piRequest.ShippingAddressID
	if shippingAddrID <= 0 {
		shippingAddrID = user.DefaultShippingAddressID
	}

	shippingAddr, err := sc.AddressRepo.FindByID(shippingAddrID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pi, err := sc.StripeService.CreatePaymentIntent(user, &piRequest, shippingAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = sc.OrderService.UpdatePaymentIntentID(piRequest.OrderID, pi.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if piRequest.DeliveryData == nil {
		c.JSON(http.StatusOK, pi)
		return
	}

	deliveryRequest := *piRequest.DeliveryData
	// 根据你的配置文件设置测试模式
	testMode := viper.GetBool("testMode")

	// 如果处于测试模式，则在requestBody中插入字段
	if testMode {
		deliveryRequest.TestSpecifications = &models.TestSpecifications{
			RoboCourierSpecification: models.RoboCourierSpecification{
				Mode: "auto",
			},
		}
	}

	// TODO: 改成由后台手动创建Delivery订单？
	deliveryResponse, err := sc.UberService.CreateDelivery(&deliveryRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = sc.OrderService.UpdateDeliveryID(piRequest.OrderID, deliveryResponse.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pi)
}

func (sc *StripeControllerImpl) DeleteAllPaymentMethods(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	iter, err := sc.StripeService.ListPaymentMethods(user.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 遍历迭代器并删除每个PaymentMethod
	for iter.Next() {
		paymentMethod := iter.PaymentMethod()
		_, err := sc.StripeService.DeletePaymentMethod(paymentMethod.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 检查迭代器中是否存在错误
	if err := iter.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = sc.UserService.SetCurrentPaymentMethod(user, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All payment methods deleted"})
}

func (sc *StripeControllerImpl) ListPaymentIntents(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	iter, err := sc.StripeService.ListPaymentIntents(user.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var paymentIntents []*stripe.PaymentIntent
	for iter.Next() {
		paymentIntent := iter.PaymentIntent()
		paymentIntents = append(paymentIntents, paymentIntent)
	}

	if err := iter.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentIntents)
}

func (sc *StripeControllerImpl) PaymentIntent(c *gin.Context) {
	paymentIntentID := c.Param("paymentIntentId")
	paymentIntent, err := sc.StripeService.RetrievePaymentIntent(paymentIntentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentIntent)
}
