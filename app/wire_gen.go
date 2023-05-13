// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/atomi-ai/atomi/controllers"
	"github.com/atomi-ai/atomi/middlewares"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/atomi-ai/atomi/utils"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitializeApplication(db *gorm.DB, authWrapper utils.AuthAppWrapper, stripeWrapper utils.StripeWrapper) (*Application, error) {
	userRepository := repositories.NewUserRepository(db)
	authMiddleware := middlewares.NewAuthMiddleware(userRepository, authWrapper)
	addressRepository := repositories.NewAddressRepository(db)
	orderRepository := repositories.NewOrderRepository(db)
	orderItemRepository := repositories.NewOrderItemRepository(db)
	productRepository := repositories.NewProductRepository(db)
	productStoreRepository := repositories.NewProductStoreRepository(db)
	storeRepository := repositories.NewStoreRepository(db)
	userAddressRepository := repositories.NewUserAddressRepository(db)
	userStoreRepository := repositories.NewUserStoreRepository(db)
	addressService := services.NewAddressService(userRepository, addressRepository, userAddressRepository)
	stripeService := services.NewStripeService()
	orderService := services.NewOrderService(orderRepository, orderItemRepository, stripeService)
	userService := services.NewUserService(userRepository)
	uberService := services.NewUberService()
	addressController := controllers.NewAddressControl(addressService, userService, addressRepository)
	loginController := controllers.NewLoginController(userRepository, stripeWrapper)
	orderController := controllers.NewOrderController(orderService, uberService)
	storeController := controllers.NewStoreController(productStoreRepository, storeRepository, userStoreRepository)
	stripeController := controllers.NewStripeController(userService, stripeService, orderService, uberService, addressRepository)
	userController := controllers.NewUserController(userService)
	application := &Application{
		AuthWrapper:            authWrapper,
		StripeWrapper:          stripeWrapper,
		AuthMiddleware:         authMiddleware,
		AddressRepository:      addressRepository,
		OrderRepository:        orderRepository,
		OrderItemRepository:    orderItemRepository,
		ProductRepository:      productRepository,
		ProductStoreRepository: productStoreRepository,
		StoreRepository:        storeRepository,
		UserAddressRepository:  userAddressRepository,
		UserRepository:         userRepository,
		UserStoreRepository:    userStoreRepository,
		AddressService:         addressService,
		OrderService:           orderService,
		StripeService:          stripeService,
		UserService:            userService,
		AddressController:      addressController,
		LoginController:        loginController,
		OrderController:        orderController,
		StoreController:        storeController,
		StripeController:       stripeController,
		UserController:         userController,
	}
	return application, nil
}

// wire.go:

type Application struct {
	AuthWrapper    utils.AuthAppWrapper
	StripeWrapper  utils.StripeWrapper
	AuthMiddleware middlewares.AuthMiddleware

	AddressRepository      repositories.AddressRepository
	OrderRepository        repositories.OrderRepository
	OrderItemRepository    repositories.OrderItemRepository
	ProductRepository      repositories.ProductRepository
	ProductStoreRepository repositories.ProductStoreRepository
	StoreRepository        repositories.StoreRepository
	UserAddressRepository  repositories.UserAddressRepository
	UserRepository         repositories.UserRepository
	UserStoreRepository    repositories.UserStoreRepository

	AddressService services.AddressService
	OrderService   services.OrderService
	StripeService  services.StripeService
	UserService    services.UserService

	AddressController controllers.AddressController
	LoginController   controllers.LoginController
	OrderController   controllers.OrderController
	StoreController   controllers.StoreController
	StripeController  controllers.StripeController
	UserController    controllers.UserController
}
