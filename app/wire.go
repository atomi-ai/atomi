//go:build wireinject
// +build wireinject

package app

import (
	"github.com/atomi-ai/atomi/controllers"
	"github.com/atomi-ai/atomi/middlewares"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/atomi-ai/atomi/utils"
	"github.com/google/wire"
	"gorm.io/gorm"
)

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

func InitializeApplication(db *gorm.DB, authWrapper utils.AuthAppWrapper, stripeWrapper utils.StripeWrapper) (*Application, error) {
	wire.Build(
		middlewares.NewAuthMiddleware,

		repositories.NewUserRepository,
		repositories.NewStoreRepository,
		repositories.NewUserStoreRepository,
		repositories.NewProductRepository,
		repositories.NewProductStoreRepository,
		repositories.NewAddressRepository,
		repositories.NewUserAddressRepository,
		repositories.NewOrderRepository,
		repositories.NewOrderItemRepository,
		services.NewAddressService,
		services.NewUserService,
		services.NewStripeService,
		services.NewOrderService,
		controllers.NewAddressControl,
		controllers.NewLoginController,
		controllers.NewOrderController,
		controllers.NewStoreController,
		controllers.NewStripeController,
		controllers.NewUserController,
		wire.Struct(new(Application), "*"),
	)

	return nil, nil
}
