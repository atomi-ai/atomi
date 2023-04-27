package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

type StripeController interface {
	AttachPaymentMethodToCustomer(c *gin.Context)
	DeleteAllPaymentMethods(c *gin.Context)
	DeletePaymentMethod(c *gin.Context)
	ListPaymentMethods(c *gin.Context)
	Pay(c *gin.Context)
}

type StripeControllerImpl struct {
	UserService   services.UserService
	StripeService services.StripeService
	AddressRepo   repositories.AddressRepository
}

func NewStripeController(userService services.UserService, stripeService services.StripeService, addressRepo repositories.AddressRepository) StripeController {
	return &StripeControllerImpl{
		UserService:   userService,
		StripeService: stripeService,
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

	shippingAddrID := piRequest.ShippingAddressID
	if shippingAddrID == 0 {
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

	c.JSON(http.StatusOK, gin.H{"message": "All payment methods deleted"})
}
