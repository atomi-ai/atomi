package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/tests"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserOrders(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("order")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个产品
	product := &models.Product{Name: "Test Product", Price: 9.99, Description: "Test product description"}
	if err = app.ProductRepository.Save(product); err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// 创建两个订单
	order1 := &models.Order{UserID: user.ID}
	order2 := &models.Order{UserID: user.ID}
	if err = app.OrderRepository.Save(order1); err != nil {
		t.Fatalf("Failed to create order1: %v", err)
	}
	if err = app.OrderRepository.Save(order2); err != nil {
		t.Fatalf("Failed to create order2: %v", err)
	}

	// 创建订单项
	orderItem1 := &models.OrderItem{OrderID: order1.ID, ProductID: product.ID, Quantity: 1}
	orderItem2 := &models.OrderItem{OrderID: order2.ID, ProductID: product.ID, Quantity: 2}
	if err = app.OrderItemRepository.Save(orderItem1); err != nil {
		t.Fatalf("Failed to create orderItem1: %v", err)
	}
	if err = app.OrderItemRepository.Save(orderItem2); err != nil {
		t.Fatalf("Failed to create orderItem2: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 GetUserOrders
	app.OrderController.GetUserOrders(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respOrders []*models.Order
	err = json.Unmarshal(w.Body.Bytes(), &respOrders)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(respOrders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(respOrders))
	}

	expectedOrderIDs := []int64{order1.ID, order2.ID}
	for _, order := range respOrders {
		if !containsOrderID(expectedOrderIDs, order.ID) {
			t.Errorf("Unexpected order in response: %v", order)
		}

		if len(order.OrderItems) != 1 {
			t.Errorf("Expected 1 order item, got %d", len(order.OrderItems))
		}

		orderItem := order.OrderItems[0]
		if orderItem.ProductID != product.ID {
			t.Errorf("Expected product ID %d, got %d", product.ID, orderItem.ProductID)
		}

		if orderItem.Quantity != 1 && orderItem.Quantity != 2 {
			t.Errorf("Expected quantity 1 or 2, got %d", orderItem.Quantity)
		}
	}
}

func containsOrderID(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestAddOrderForUser(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("order")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个产品
	product := &models.Product{Name: "Test Product", Price: 9.99, Description: "Test product description"}
	if err = app.ProductRepository.Save(product); err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 准备一个要添加的订单
	order := &models.Order{
		UserID: user.ID,
		OrderItems: []models.OrderItem{
			{
				ProductID: product.ID,
				Quantity:  2,
			},
		},
	}

	// 将订单转换为JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// 设置请求体
	c.Request, _ = http.NewRequest("POST", "/orders", bytes.NewReader(orderJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用 AddOrderForUser
	app.OrderController.AddOrderForUser(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var savedOrder models.Order
	err = json.Unmarshal(w.Body.Bytes(), &savedOrder)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if savedOrder.UserID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, savedOrder.UserID)
	}

	if len(savedOrder.OrderItems) != 1 {
		t.Fatalf("Expected 1 order item, got %d", len(savedOrder.OrderItems))
	}

	orderItem := savedOrder.OrderItems[0]
	if orderItem.ProductID != product.ID {
		t.Errorf("Expected product ID %d, got %d", product.ID, orderItem.ProductID)
	}

	if orderItem.Quantity != 2 {
		t.Errorf("Expected quantity 2, got %d", orderItem.Quantity)
	}
}
