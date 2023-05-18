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
	BlobStorage    utils.BlobStorage
	StripeWrapper  utils.StripeWrapper
	AuthMiddleware middlewares.AuthMiddleware

	AddressController      controllers.AddressController
	ImageController        controllers.ImageController
	LoginController        controllers.LoginController
	ManagerStoreController controllers.ManagerStoreController
	OrderController        controllers.OrderController
	StoreController        controllers.StoreController
	StripeController       controllers.StripeController
	UserController         controllers.UserController

	AddressRepository           repositories.AddressRepository
	DeleteUserRequestRepository repositories.DeleteUserRequestRepository
	ManagerStoreRepository      repositories.ManagerStoreRepository
	OrderRepository             repositories.OrderRepository
	OrderItemRepository         repositories.OrderItemRepository
	ProductRepository           repositories.ProductRepository
	ProductStoreRepository      repositories.ProductStoreRepository
	StoreRepository             repositories.StoreRepository
	UserAddressRepository       repositories.UserAddressRepository
	UserRepository              repositories.UserRepository
	UserStoreRepository         repositories.UserStoreRepository

	AddressService      services.AddressService
	OrderService        services.OrderService
	ProductStoreService services.ProductStoreService
	StripeService       services.StripeService
	UserService         services.UserService
}

func InitializeApplication(db *gorm.DB, authWrapper utils.AuthAppWrapper, blobStorage utils.BlobStorage, stripeWrapper utils.StripeWrapper) (*Application, error) {
	wire.Build(
		middlewares.NewAuthMiddleware,

		controllers.NewAddressControl,
		controllers.NewImageController,
		controllers.NewLoginController,
		controllers.NewManagerStoreController,
		controllers.NewOrderController,
		controllers.NewStoreController,
		controllers.NewStripeController,
		controllers.NewUserController,
		repositories.NewAddressRepository,
		repositories.NewDeleteUserRequestRepository,
		repositories.NewManagerStoreRepository,
		repositories.NewOrderItemRepository,
		repositories.NewOrderRepository,
		repositories.NewProductRepository,
		repositories.NewProductStoreRepository,
		repositories.NewStoreRepository,
		repositories.NewUserAddressRepository,
		repositories.NewUserRepository,
		repositories.NewUserStoreRepository,
		services.NewAddressService,
		services.NewOrderService,
		services.NewProductStoreService,
		services.NewStripeService,
		services.NewUserService,
		services.NewUberService,

		wire.Struct(new(Application), "*"),
	)

	return nil, nil
}
