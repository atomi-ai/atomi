package app

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/models"
	"github.com/stripe/stripe-go/v74"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockAuthApp struct{}

func (m *MockAuthApp) AuthAndDecode(_ context.Context, _ string) (*auth.Token, error) {
	mockDecodedToken := &auth.Token{
		Claims: map[string]interface{}{
			"email": "john.doe@example.com",
		},
	}
	return mockDecodedToken, nil
}

type MockStripeWrapper struct{}

func (m *MockStripeWrapper) CreateCustomer(email string) (*stripe.Customer, error) {
	mockCustomer := &stripe.Customer{
		ID:    "cus_mock_id",
		Email: email,
	}
	return mockCustomer, nil
}

func InitializeTestingApplication(dbName string) (*Application, error) {
	// 创建一个内存中的SQLite数据库
	dsn := fmt.Sprintf("file:%v?mode=memory&cache=shared", dbName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to in-memory SQLite database: %w", err)
	}

	// 自动迁移模型
	models.AutoMigrate(db)

	// 使用Mock替换Firebase和Stripe等外部服务
	mockAuthApp := new(MockAuthApp)
	mockStripeWrapper := new(MockStripeWrapper)

	// 创建一个用于测试的 *Application 实例
	app, err := InitializeApplication(db, mockAuthApp, mockStripeWrapper)
	if err != nil {
		return nil, err
	}

	return app, nil
}
