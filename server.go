package main

import (
	application "github.com/atomi-ai/atomi/app"
	"github.com/atomi-ai/atomi/middlewares"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	utils.LoadConfig()
	db := models.InitDB()
	models.AutoMigrate(db)
	utils.InitStripe(viper.GetString("stripeKey"))

	app, err := application.InitializeApplication(db, utils.NewFirebaseAppWrapper(utils.FirebaseAppProvider()), utils.NewStripeWrapper())
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.Use(app.AuthMiddleware.Handler())
	r.Use(middlewares.RequestResponseLogger()) // 添加自定义的请求/响应日志中间件
	r.GET("/api/login", app.LoginController.Login)

	// Add StoreController endpoints here
	r.GET("/api/default-store", app.StoreController.GetDefaultStore)
	r.PUT("/api/default-store/:store_id", app.StoreController.SetDefaultStore)
	r.GET("/api/stores", app.StoreController.GetAllStores)
	r.DELETE("/api/default-store", app.StoreController.DeleteDefaultStore)
	r.GET("/api/products/:store_id", app.StoreController.GetProductsByStoreID)

	// Add AddressController endpoints here
	r.GET("/api/addresses", app.AddressController.GetAllAddressesForUser)
	r.POST("/api/addresses", app.AddressController.AddAddressForUser)
	r.DELETE("/api/addresses/:addressId", app.AddressController.DeleteAddressForUser)
	r.POST("/api/addresses/shipping/:addressId", app.AddressController.SetDefaultShippingAddress)
	r.POST("/api/addresses/billing/:addressId", app.AddressController.SetDefaultBillingAddress)
	r.GET("/api/addresses/shipping", app.AddressController.GetDefaultShippingAddress)
	r.GET("/api/addresses/billing", app.AddressController.GetDefaultBillingAddress)
	r.DELETE("/api/addresses", app.AddressController.DeleteAllAddressesForUser)

	// Add payment endpoints here.
	r.PUT("/api/payment-methods/:paymentMethodId", app.StripeController.AttachPaymentMethodToCustomer)
	r.GET("/api/payment-methods", app.StripeController.ListPaymentMethods)
	r.DELETE("/api/payment-methods/:paymentMethodId", app.StripeController.DeletePaymentMethod)
	r.POST("/api/pay", app.StripeController.Pay)
	r.DELETE("/api/payment-methods", app.StripeController.DeleteAllPaymentMethods)
	r.GET("/api/payment-intents", app.StripeController.ListPaymentIntents)
	r.GET("/api/payment-intent/:paymentIntentId", app.StripeController.PaymentIntent)

	// Add user endpoints here
	r.GET("/api/user", app.UserController.GetUser)
	r.PUT("/api/user/current-payment-method/:paymentMethodId", app.UserController.SetCurrentPaymentMethod)

	// Add order endpoints here
	r.GET("/api/orders", app.OrderController.GetUserOrders)
	r.POST("/api/order", app.OrderController.AddOrderForUser)

	// APIs below are not tested by flutter tests yet.
	if err = r.Run(":8081"); err != nil {
		log.Fatal("Errors in running application on port 8081", err)
	}
}
