package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/atomi-ai/atomi/utils"

	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/services"
	"github.com/spf13/viper"
)

const (
	adminEmail = "admin@atomi.ai"
	userEmail  = "user@atomi.ai"
)

type TestEnvSetup struct {
	ConfigRepository       repositories.ConfigRepository
	ManagerStoreRepository repositories.ManagerStoreRepository
	OrderRepository        repositories.OrderRepository
	OrderItemRepository    repositories.OrderItemRepository
	ProductRepository      repositories.ProductRepository
	UserRepository         repositories.UserRepository

	ProductStoreService services.ProductStoreService
}

func LoadConfig() {
	configFile := os.Getenv("CONFIG_FILE")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}
}

// TODO(lamuguo): rename the directory as init-db.
func main() {
	// TODO(lamuguo): Please wire to inject.
	// App system initialization
	LoadConfig()
	db := models.InitDB()
	models.AutoMigrate(db)
	utils.InitStripe(viper.GetString("stripeKey"))
	firebaseApp := utils.FirebaseAppProvider()

	// Create an auth client.
	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		fmt.Println("error initializing auth client:", err)
		os.Exit(1)
	}

	testEnvSetup := &TestEnvSetup{
		ConfigRepository:       repositories.NewConfigRepository(db),
		OrderRepository:        repositories.NewOrderRepository(db),
		OrderItemRepository:    repositories.NewOrderItemRepository(db),
		ProductRepository:      repositories.NewProductRepository(db),
		ProductStoreService:    services.NewProductStoreService(repositories.NewProductRepository(db), repositories.NewProductStoreRepository(db)),
		ManagerStoreRepository: repositories.NewManagerStoreRepository(db),
		UserRepository:         repositories.NewUserRepository(db),
	}
	testEnvSetup.run(authClient)
}

func (t *TestEnvSetup) run(authClient *auth.Client) {
	// Initialize the database here

	// 1. Check and create users in Firebase
	admin := t.checkOrCreateUserInFirebase(authClient, adminEmail, "Admin", models.RoleAdmin)
	user := t.checkOrCreateUserInFirebase(authClient, userEmail, "User", models.RoleUser)
	manager := t.checkOrCreateUserInFirebase(authClient, "mgr@atomi.ai", "Manager", models.RoleMgr)

	// 2. Add products to the database
	products := t.addProducts(admin)

	// 3. Add stores to the database
	store1 := &models.Store{
		Name:    "first_store",
		Address: "1 Washington St",
		City:    "San Jose",
		ZipCode: "95192",
		State:   "CA",
		Phone:   "5103490111",
	}
	store2 := &models.Store{
		Name:    "second_store",
		Address: "1450 W Pleasant Run Rd",
		City:    "Lancaster",
		ZipCode: "75146",
		State:   "TX",
		Phone:   "5103490222",
	}
	t.ManagerStoreRepository.Save(store1)
	t.ManagerStoreRepository.Save(store2)
	log.Infof("store1 = %v, store2 = %v", store1, store2)
	err1 := t.ManagerStoreRepository.AssignStoreToUser(store1.ID, manager.ID)
	err2 := t.ManagerStoreRepository.AssignStoreToUser(store2.ID, manager.ID)
	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("Errors in assign store to manager, %v, %v", err1, err2))
	}

	// 4. Connect products and stores
	t.ProductStoreService.ConnectStoreAndProducts(store1, products)

	// 5. Add orders for testing
	t.addOrders(store1, products, user)

	// 6. Set testenv_status to "initialized"
	// 如果您有一个类似于Java代码中的ConfigRepository，请在此处将 testenv_status 设置为 "initialized"
	// 如果没有，请根据您的具体实现进行修改。
	t.ConfigRepository.Save(&models.Config{Key: "testenv_status", Value: "initialized"})
	fmt.Println("Finished initializing the test environment.")
}

func (t *TestEnvSetup) checkOrCreateUserInFirebase(authClient *auth.Client, email string, displayName string, role models.Role) *models.User {
	userRecord, err := authClient.GetUserByEmail(context.Background(), email)
	if err != nil {
		// 如果找不到用户，创建一个新的
		params := (&auth.UserToCreate{}).
			Email(email).
			EmailVerified(true).
			Password("password123").
			DisplayName(displayName).
			Disabled(false)
		userRecord, err = authClient.CreateUser(context.Background(), params)
		if err != nil {
			panic(fmt.Sprintf("Error creating user in Firebase: %v", err))
		}
	}

	// 在数据库中查找用户
	user, err := t.UserRepository.FindByEmail(email)
	if err != nil {
		// 如果找不到用户，创建一个新的
		user = &models.User{
			Email: userRecord.Email,
			Role:  role,
		}
		_, err = t.UserRepository.Save(user)
		if err != nil {
			panic(fmt.Sprintf("Error saving user to database: %v", err))
		}
	}
	return user
}

func (t *TestEnvSetup) addProducts(admin *models.User) []*models.Product {
	products := []*models.Product{
		{
			Creator:  admin,
			Name:     "Hamburger",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/3.png",
			Price:    25,
			Discount: 10,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Pasta",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/5.png",
			Price:    150,
			Discount: 7.8,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Akara",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/2.png",
			Price:    10.99,
			Discount: 0,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Strawberry",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/1.png",
			Price:    50,
			Discount: 14,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Coca-Cola",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/6.png",
			Price:    45.12,
			Discount: 2,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Lemonade",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/7.png",
			Price:    28,
			Discount: 5.2,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Vodka",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/8.png",
			Price:    78.99,
			Discount: 0,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Tequila",
			ImageURL: "https://atomidrone.blob.core.windows.net/images/9.png",
			Price:    1234567,
			Discount: 3.4,
			Category: models.ProductCategoryDrink,
		},
	}

	for _, product := range products {
		err := t.ProductRepository.Save(product)
		if err != nil {
			fmt.Printf("Error saving product to database: %v\n", err)
		}
	}

	return products
}

func (t *TestEnvSetup) saveOrder(order *models.Order) error {
	if err := t.OrderRepository.Save(order); err != nil {
		return err
	}

	for i := range order.OrderItems {
		orderItem := &order.OrderItems[i]
		orderItem.OrderID = order.ID
		if err := t.OrderItemRepository.Save(orderItem); err != nil {
			return err
		}
	}

	return nil
}
func (t *TestEnvSetup) addOrders(store *models.Store, products []*models.Product, user *models.User) {
	orders := []*models.Order{
		{
			UserID:        user.ID,
			StoreID:       store.ID,
			DisplayStatus: models.OrderStatusPaid,
			OrderItems: []models.OrderItem{
				{
					Product:   products[0],
					ProductID: products[0].ID,
					Quantity:  2,
				},
				{
					Product:   products[1],
					ProductID: products[1].ID,
					Quantity:  1,
				},
			},
		},
		{
			UserID:        user.ID,
			StoreID:       store.ID,
			DisplayStatus: models.OrderStatusInProduction,
			OrderItems: []models.OrderItem{
				{
					Product:   products[2],
					ProductID: products[2].ID,
					Quantity:  3,
				},
			},
		},
	}

	for _, order := range orders {
		err := t.saveOrder(order)
		if err != nil {
			fmt.Printf("Error saving order to database: %v\n", err)
		}
	}
}
