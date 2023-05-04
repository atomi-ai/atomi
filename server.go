package main

import (
	"context"
	"errors"
	"github.com/atomi-ai/atomi/utils"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/controllers"
	"github.com/atomi-ai/atomi/middlewares"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	firebase "firebase.google.com/go/v4"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
)

var (
	db          *gorm.DB
	firebaseApp *firebase.App

	UserRepository         repositories.UserRepository
	StoreRepository        repositories.StoreRepository
	UserStoreRepository    repositories.UserStoreRepository
	ProductStoreRepository repositories.ProductStoreRepository
	AddressRepository      repositories.AddressRepository
	UserAddressRepository  repositories.UserAddressRepository
	OrderRepository        repositories.OrderRepository
	OrderItemRepository    repositories.OrderItemRepository

	UserService    services.UserService
	StripeService  services.StripeService
	AddressService services.AddressService
	OrderService   services.OrderService

	err error
)

func main() {
	utils.LoadConfig()

	db = models.InitDB()
	models.AutoMigrate(db)
	utils.InitStripe(viper.GetString("stripeKey"))
	utils.InitFirebase()

	UserRepository = repositories.NewUserRepository(db)
	StoreRepository = repositories.NewStoreRepository(db)
	UserStoreRepository = repositories.NewUserStoreRepository(db)
	ProductStoreRepository = repositories.NewProductStoreRepository(db)
	AddressRepository = repositories.NewAddressRepository(db)
	UserAddressRepository = repositories.NewUserAddressRepository(db)
	OrderRepository = repositories.NewOrderRepository(db)
	OrderItemRepository = repositories.NewOrderItemRepository(db)

	AddressService = services.NewAddressService(UserRepository, AddressRepository, UserAddressRepository)
	UserService = services.NewUserService(UserRepository)
	StripeService = services.NewStripeService()
	OrderService = services.NewOrderService(OrderRepository, OrderItemRepository, StripeService)

	storeController := controllers.NewStoreController(ProductStoreRepository, StoreRepository, UserStoreRepository)
	addressController := controllers.NewAddressControl(AddressService, UserService, AddressRepository)
	stripeController := controllers.NewStripeController(UserService, StripeService, OrderService, AddressRepository)
	userController := controllers.NewUserController(UserService)
	orderController := controllers.NewOrderController(OrderService)

	r := gin.Default()
	r.Use(AuthMiddleware())
	r.Use(middlewares.RequestResponseLogger()) // 添加自定义的请求/响应日志中间件
	r.GET("/api/login", Login)

	// Add StoreController endpoints here
	r.GET("/api/default-store", storeController.GetDefaultStore)
	r.PUT("/api/default-store/:store_id", storeController.SetDefaultStore)
	r.GET("/api/stores", storeController.GetAllStores)
	r.DELETE("/api/default-store", storeController.DeleteDefaultStore)
	r.GET("/api/products/:store_id", storeController.GetProductsByStoreID)

	// Add AddressController endpoints here
	r.GET("/api/addresses", addressController.GetAllAddressesForUser)
	r.POST("/api/addresses", addressController.AddAddressForUser)
	r.DELETE("/api/addresses/:addressId", addressController.DeleteAddressForUser)
	r.POST("/api/addresses/shipping/:addressId", addressController.SetDefaultShippingAddress)
	r.POST("/api/addresses/billing/:addressId", addressController.SetDefaultBillingAddress)
	r.GET("/api/addresses/shipping", addressController.GetDefaultShippingAddress)
	r.GET("/api/addresses/billing", addressController.GetDefaultBillingAddress)
	r.DELETE("/api/addresses", addressController.DeleteAllAddressesForUser)

	// Add payment endpoints here.
	r.PUT("/api/payment-methods/:paymentMethodId", stripeController.AttachPaymentMethodToCustomer)
	r.GET("/api/payment-methods", stripeController.ListPaymentMethods)
	r.DELETE("/api/payment-methods/:paymentMethodId", stripeController.DeletePaymentMethod)
	r.POST("/api/pay", stripeController.Pay)
	r.DELETE("/api/payment-methods", stripeController.DeleteAllPaymentMethods)
	r.GET("/api/payment-intents", stripeController.ListPaymentIntents)
	r.GET("/api/payment-intent/:paymentIntentId", stripeController.PaymentIntent)

	// Add user endpoints here
	r.GET("/api/user", userController.GetUser)
	r.PUT("/api/user/current-payment-method/:paymentMethodId", userController.SetCurrentPaymentMethod)

	// Add order endpoints here
	r.GET("/api/orders", orderController.GetUserOrders)
	r.POST("/api/order", orderController.AddOrderForUser)

	// APIs below are not tested by flutter tests yet.
	r.Run(":8081")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		ctx := context.Background()
		client, err := firebaseApp.Auth(ctx)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Error getting Auth client"})
			return
		}

		decodedToken, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		email := decodedToken.Claims["email"].(string)
		user, err := UserRepository.FindByEmail(email)
		if err == nil {
			c.Set("user", user)
		}
		c.Set("decodedToken", decodedToken)

		c.Next()
	}
}

func Login(c *gin.Context) {
	token, _ := c.MustGet("decodedToken").(*auth.Token)

	email, _ := token.Claims["email"].(string)
	user, err := UserRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &models.User{
				Email: email,
				Role:  "USER",
			}
		} else {
			c.JSON(500, gin.H{"error": "Error fetching user"})
			return
		}
	}

	dirty := false
	if user.StripeCustomerID == "" {
		dirty = true
		stripeCustomer, err := createStripeCustomer(email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error creating Stripe customer"})
			return
		}
		user.StripeCustomerID = stripeCustomer.ID
	}

	if dirty {
		UserRepository.Save(user)
	}

	c.JSON(200, user)
}

func createStripeCustomer(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	return customer.New(params)
}
