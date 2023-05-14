package controllers

import (
	"encoding/json"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/tests"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestStoreGetDefaultStore(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("store")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个商店
	store := &models.Store{Name: "Test Store", Address: "123 Main St", City: "New York", State: "NY", ZipCode: "10001", Phone: "555-1234"}
	if err = app.ManagerStoreRepository.Save(store); err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 为用户设置默认商店
	userStore := &models.UserStore{UserID: user.ID, StoreID: store.ID, IsEnable: true}
	if err = app.UserStoreRepository.Save(userStore); err != nil {
		t.Fatalf("Failed to set default store for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 GetDefaultStore
	app.StoreController.GetDefaultStore(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respStore models.Store
	err = json.Unmarshal(w.Body.Bytes(), &respStore)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if respStore.ID != store.ID || respStore.Name != store.Name || respStore.Address != store.Address || respStore.City != store.City || respStore.State != store.State || respStore.ZipCode != store.ZipCode || respStore.Phone != store.Phone {
		t.Errorf("Expected store %v, got %v", store, respStore)
	}
}

func TestStoreSetDefaultStore(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("store")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "Jane Doe", Email: "jane.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建两个商店
	store1 := &models.Store{Name: "Test Store 1", Address: "123 Main St", City: "New York", State: "NY", ZipCode: "10001", Phone: "555-1234"}
	if err = app.ManagerStoreRepository.Save(store1); err != nil {
		t.Fatalf("Failed to create store1: %v", err)
	}

	store2 := &models.Store{Name: "Test Store 2", Address: "456 Main St", City: "New York", State: "NY", ZipCode: "10002", Phone: "555-5678"}
	if err = app.ManagerStoreRepository.Save(store2); err != nil {
		t.Fatalf("Failed to create store2: %v", err)
	}

	// 为用户设置默认商店
	userStore := &models.UserStore{UserID: user.ID, StoreID: store1.ID, IsEnable: true}
	if err = app.UserStoreRepository.Save(userStore); err != nil {
		t.Fatalf("Failed to set default store for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 设置请求参数
	c.Params = []gin.Param{
		{
			Key:   "store_id",
			Value: strconv.FormatInt(store2.ID, 10),
		},
	}

	// 调用 SetDefaultStore
	app.StoreController.SetDefaultStore(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respStore models.Store
	err = json.Unmarshal(w.Body.Bytes(), &respStore)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if respStore.ID != store2.ID || respStore.Name != store2.Name || respStore.Address != store2.Address || respStore.City != store2.City || respStore.State != store2.State || respStore.ZipCode != store2.ZipCode || respStore.Phone != store2.Phone {
		t.Errorf("Expected store %v, got %v", store2, respStore)
	}

	// 验证用户的默认商店已更改
	newUserStore, err := app.UserStoreRepository.FindDefaultUserStore(user.ID)
	if err != nil {
		t.Fatalf("Failed to find default user store: %v", err)
	}

	if newUserStore.StoreID != store2.ID {
		t.Errorf("Expected default store ID to be %d, got %d", store2.ID, newUserStore.StoreID)
	}
}

func TestStoreGetAllStores(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("store2")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建两个商店
	store1 := &models.Store{Name: "Test Store 1", Address: "123 Main St", City: "New York", State: "NY", ZipCode: "10001", Phone: "555-1234"}
	if err = app.ManagerStoreRepository.Save(store1); err != nil {
		t.Fatalf("Failed to create store1: %v", err)
	}

	store2 := &models.Store{Name: "Test Store 2", Address: "456 Main St", City: "New York", State: "NY", ZipCode: "10002", Phone: "555-5678"}
	if err = app.ManagerStoreRepository.Save(store2); err != nil {
		t.Fatalf("Failed to create store2: %v", err)
	}

	// 准备一个测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 调用 GetAllStores
	app.StoreController.GetAllStores(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respStores []models.Store
	err = json.Unmarshal(w.Body.Bytes(), &respStores)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(respStores) != 2 {
		t.Errorf("Expected 2 stores, got %d", len(respStores))
	}

	expectedStores := []*models.Store{store1, store2}
	for i, respStore := range respStores {
		if respStore.ID != expectedStores[i].ID || respStore.Name != expectedStores[i].Name || respStore.Address != expectedStores[i].Address || respStore.City != expectedStores[i].City || respStore.State != expectedStores[i].State || respStore.ZipCode != expectedStores[i].ZipCode || respStore.Phone != expectedStores[i].Phone {
			t.Errorf("Expected store %v, got %v", expectedStores[i], respStore)
		}
	}
}

func TestStoreDeleteDefaultStore(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("store")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个商店
	store := &models.Store{Name: "Test Store", Address: "123 Main St", City: "New York", State: "NY", ZipCode: "10001", Phone: "555-1234"}
	if err = app.ManagerStoreRepository.Save(store); err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 为用户设置默认商店
	userStore := &models.UserStore{UserID: user.ID, StoreID: store.ID, IsEnable: true}
	if err = app.UserStoreRepository.Save(userStore); err != nil {
		t.Fatalf("Failed to set default store for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 DeleteDefaultStore
	app.StoreController.DeleteDefaultStore(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var response gin.H
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Default store deleted successfully" {
		t.Errorf("Expected message 'Default store deleted successfully', got '%v'", response["message"])
	}

	// 确保默认商店已被删除
	deletedUserStore, err := app.UserStoreRepository.FindDefaultUserStore(user.ID)
	if err == nil {
		t.Errorf("Expected default store to be deleted, got %+v", deletedUserStore)
	}
}

func TestGetProductsByStoreID(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("product")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个商店
	store := &models.Store{Name: "Test Store", Address: "123 Main St", City: "New York", State: "NY", ZipCode: "10001", Phone: "555-1234"}
	if err = app.ManagerStoreRepository.Save(store); err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 创建两个产品
	product1 := &models.Product{Name: "Product 1", Description: "Product 1 description", Price: 10.00, Category: models.ProductCategoryFood}
	product2 := &models.Product{Name: "Product 2", Description: "Product 2 description", Price: 5.00, Category: models.ProductCategoryDrink}
	if err = app.ProductRepository.Save(product1); err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	if err = app.ProductRepository.Save(product2); err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}

	// 将产品添加到商店
	productStore1 := &models.ProductStore{StoreID: store.ID, ProductID: product1.ID, IsEnable: true}
	productStore2 := &models.ProductStore{StoreID: store.ID, ProductID: product2.ID, IsEnable: true}
	if err = app.ProductStoreRepository.Save(productStore1); err != nil {
		t.Fatalf("Failed to add product1 to store: %v", err)
	}
	if err = app.ProductStoreRepository.Save(productStore2); err != nil {
		t.Fatalf("Failed to add product2 to store: %v", err)
	}

	// 准备一个测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置路由参数
	c.Params = []gin.Param{
		{Key: "store_id", Value: strconv.FormatInt(store.ID, 10)},
	}

	// 调用 GetProductsByStoreID
	app.StoreController.GetProductsByStoreID(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respProducts []*models.Product
	err = json.Unmarshal(w.Body.Bytes(), &respProducts)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(respProducts) != 2 {
		t.Fatalf("Expected 2 products, got %d", len(respProducts))
	}

	expectedProductIDs := []int64{product1.ID, product2.ID}
	for _, product := range respProducts {
		if !contains(expectedProductIDs, product.ID) {
			t.Errorf("Unexpected product in response: %v", product)
		}
	}
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
