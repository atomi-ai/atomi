package app

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/atomi-ai/atomi/models"
	"github.com/stripe/stripe-go/v74"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockFirebaseApp struct{}

func (m *MockFirebaseApp) Auth(ctx context.Context) (*auth.Client, error) {
	// 在这里返回一个Mock的auth.Client
	// 或者如果你想模拟一个错误，你可以返回一个错误，例如：
	return nil, errors.New("Not implemented yet")
}

type MockStripeWrapper struct{}

func (m *MockStripeWrapper) CreateCustomer(email string) (*stripe.Customer, error) {
	mockCustomer := &stripe.Customer{
		ID:    "cus_mock_id",
		Email: email,
	}
	return mockCustomer, nil
}

func InitializeTestingApplication() (*Application, error) {
	// 创建一个内存中的SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to in-memory SQLite database: %v", err)
	}

	// 自动迁移模型
	db.AutoMigrate(&models.User{})

	// 使用Mock替换Firebase和Stripe等外部服务
	mockFirebaseApp := new(MockFirebaseApp)
	mockStripeWrapper := new(MockStripeWrapper)

	// 创建一个用于测试的 *Application 实例
	app, err := InitializeApplication(db, mockFirebaseApp, mockStripeWrapper)
	if err != nil {
		return nil, err
	}

	return app, nil
}
