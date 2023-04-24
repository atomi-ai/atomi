package main

import (
	"context"
	"fmt"
	"github.com/atomi-ai/atomi/services"
	"google.golang.org/api/option"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

const (
	adminEmail = "admin@atomi.ai"
	userEmail  = "user@atomi.ai"
)

type TestEnvSetup struct {
	ConfigRepository    repositories.ConfigRepository
	ProductRepository   repositories.ProductRepository
	ProductStoreService services.ProductStoreService
	StoreRepository     repositories.StoreRepository
	UserRepository      repositories.UserRepository
}

func main() {
	db := models.InitDB()

	testEnvSetup := &TestEnvSetup{
		ConfigRepository:    repositories.NewConfigRepository(db),
		ProductRepository:   repositories.NewProductRepository(db),
		ProductStoreService: services.NewProductStoreService(repositories.NewProductStoreRepository(db)),
		StoreRepository:     repositories.NewStoreRepository(db),
		UserRepository:      repositories.NewUserRepository(db),
	}

	// Initialize Firebase app, set your Firebase local emulator URL for testing.
	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	opt := option.WithCredentialsFile("testing/testing-firebase-secret.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error initializing firebase app:", err)
		os.Exit(1)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		fmt.Println("error initializing auth client:", err)
		os.Exit(1)
	}

	testEnvSetup.run(authClient)
}

func (t *TestEnvSetup) run(authClient *auth.Client) {
	// Initialize the database here

	// 1. Check and create users in Firebase
	admin, err := t.checkOrCreateUserInFirebase(authClient, adminEmail, "Admin", models.RoleAdmin)
	if err != nil {
		fmt.Printf("Error creating admin user: %v\n", err)
		return
	}

	_, err = t.checkOrCreateUserInFirebase(authClient, userEmail, "User", models.RoleUser)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}

	// 2. Add products to the database
	products := t.addProducts(admin)

	// 3. Add stores to the database
	store1 := &models.Store{
		Name:    "first_store",
		Address: "1200 Fremont Blvd",
		City:    "Fremont",
		ZipCode: "94555",
		State:   "CA",
	}
	store2 := &models.Store{
		Name:    "second_store",
		Address: "1800 Washington Blvd",
		City:    "Fremont",
		ZipCode: "94536",
		State:   "CA",
	}
	t.StoreRepository.Save(store1)
	t.StoreRepository.Save(store2)

	// 4. Connect products and stores
	t.ProductStoreService.ConnectStoreAndProducts(store1, products)

	// 5. Set testenv_status to "initialized"
	// 如果您有一个类似于Java代码中的ConfigRepository，请在此处将 testenv_status 设置为 "initialized"
	// 如果没有，请根据您的具体实现进行修改。
	t.ConfigRepository.Save(&models.Config{Key: "testenv_status", Value: "initialized"})

	fmt.Println("Finished initializing the test environment.")
}

func (t *TestEnvSetup) checkOrCreateUserInFirebase(authClient *auth.Client, email string, displayName string, role models.Role) (*models.User, error) {
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
			return nil, fmt.Errorf("Error creating user in Firebase: %w", err)
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
		err = t.UserRepository.Save(user)
		if err != nil {
			return nil, fmt.Errorf("Error saving user to database: %w", err)
		}
	}

	return user, nil
}

func (t *TestEnvSetup) addProducts(admin *models.User) []*models.Product {
	products := []*models.Product{
		{
			Creator:  admin,
			Name:     "Hamburger",
			ImageURL: "images/3.png",
			Price:    25,
			Discount: 10,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Pasta",
			ImageURL: "images/5.png",
			Price:    150,
			Discount: 7.8,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Akara",
			ImageURL: "images/2.png",
			Price:    10.99,
			Discount: 0,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Strawberry",
			ImageURL: "images/1.png",
			Price:    50,
			Discount: 14,
			Category: models.ProductCategoryFood,
		},
		{
			Creator:  admin,
			Name:     "Coca-Cola",
			ImageURL: "images/6.png",
			Price:    45.12,
			Discount: 2,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Lemonade",
			ImageURL: "images/7.png",
			Price:    28,
			Discount: 5.2,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Vodka",
			ImageURL: "images/8.png",
			Price:    78.99,
			Discount: 0,
			Category: models.ProductCategoryDrink,
		},
		{
			Creator:  admin,
			Name:     "Tequila",
			ImageURL: "images/9.png",
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
