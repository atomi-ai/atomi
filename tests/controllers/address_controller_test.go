package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/tests"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAddressGetAllAddressesForUser(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建两个地址
	addresses := []*models.Address{
		{
			Line1:      "123 Main St",
			City:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			PostalCode: "12345",
		},
		{
			Line1:      "456 Market St",
			City:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			PostalCode: "67890",
		},
	}

	// 使用 AddressService 为用户添加地址
	for _, address := range addresses {
		_, err := app.AddressService.AddAddressForUser(user, address)
		if err != nil {
			t.Fatalf("Failed to add address for user: %v", err)
		}
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 GetAllAddressesForUser
	app.AddressController.GetAllAddressesForUser(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respAddresses []*models.Address
	err = json.Unmarshal(w.Body.Bytes(), &respAddresses)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(respAddresses) != len(addresses) {
		t.Errorf("Expected %d addresses, got %d", len(addresses), len(respAddresses))
	}

	for i, respAddress := range respAddresses {
		if respAddress.Line1 != addresses[i].Line1 || respAddress.City != addresses[i].City || respAddress.State != addresses[i].State || respAddress.Country != addresses[i].Country || respAddress.PostalCode != addresses[i].PostalCode {
			t.Errorf("Expected address %v, got %v", addresses[i], respAddress)
		}
	}
}

func TestAddressAddAddressForUser(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 准备请求数据
	address := models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}
	reqBody, err := json.Marshal(address)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	c.Request, err = http.NewRequest("POST", "/addresses", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用 AddAddressForUser
	app.AddressController.AddAddressForUser(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respAddress models.Address
	err = json.Unmarshal(w.Body.Bytes(), &respAddress)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if respAddress.Line1 != address.Line1 || respAddress.City != address.City || respAddress.State != address.State || respAddress.Country != address.Country || respAddress.PostalCode != address.PostalCode {
		t.Errorf("Expected address %v, got %v", address, respAddress)
	}

	// 确保地址已保存在数据库中
	dbAddresses, err := app.AddressService.GetAddressesByUserID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get addresses from database: %v", err)
	}

	if len(dbAddresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(dbAddresses))
	}

	if dbAddresses[0].Line1 != address.Line1 || dbAddresses[0].City != address.City || dbAddresses[0].State != address.State || dbAddresses[0].Country != address.Country || dbAddresses[0].PostalCode != address.PostalCode {
		t.Errorf("Expected address %v, got %v", address, dbAddresses[0])
	}
}

func TestAddressDeleteAddressForUser(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个地址
	address := &models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}

	// 使用 AddressService 为用户添加地址
	address, err = app.AddressService.AddAddressForUser(user, address)
	if err != nil {
		t.Fatalf("Failed to add address for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 设置请求参数
	c.Params = []gin.Param{
		{Key: "addressId", Value: strconv.FormatInt(address.ID, 10)},
	}

	// 调用 DeleteAddressForUser
	app.AddressController.DeleteAddressForUser(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var respStatus map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &respStatus)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if respStatus["status"] != "success" {
		t.Errorf("Expected status 'success', got %s", respStatus["status"])
	}

	// 确保地址已从数据库中删除
	dbAddresses, err := app.AddressService.GetAddressesByUserID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get addresses from database: %v", err)
	}

	if len(dbAddresses) != 0 {
		t.Errorf("Expected 0 addresses, got %d", len(dbAddresses))
	}
}

func TestAddressSetDefaultShippingAddress(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个地址
	address := &models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}

	// 使用 AddressService 为用户添加地址
	address, err = app.AddressService.AddAddressForUser(user, address)
	if err != nil {
		t.Fatalf("Failed to add address for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 设置请求参数
	c.Params = []gin.Param{
		{Key: "addressId", Value: strconv.FormatInt(address.ID, 10)},
	}

	// 调用 SetDefaultShippingAddress
	app.AddressController.SetDefaultShippingAddress(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var updatedUser *models.User
	err = json.Unmarshal(w.Body.Bytes(), &updatedUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedUser.DefaultShippingAddressID != address.ID {
		t.Errorf("Expected default shipping address ID %d, got %d", address.ID, updatedUser.DefaultShippingAddressID)
	}
}

func TestAddressSetDefaultBillingAddress(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个地址
	address := &models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}

	// 使用 AddressService 为用户添加地址
	address, err = app.AddressService.AddAddressForUser(user, address)
	if err != nil {
		t.Fatalf("Failed to add address for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 设置请求参数
	c.Params = []gin.Param{
		{Key: "addressId", Value: strconv.FormatInt(address.ID, 10)},
	}

	// 调用 SetDefaultBillingAddress
	app.AddressController.SetDefaultBillingAddress(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var updatedUser *models.User
	err = json.Unmarshal(w.Body.Bytes(), &updatedUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedUser.DefaultBillingAddressID != address.ID {
		t.Errorf("Expected default billing address ID %d, got %d", address.ID, updatedUser.DefaultBillingAddressID)
	}
}

func TestAddressGetDefaultShippingAddress(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个地址
	address := &models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}

	// 使用 AddressService 为用户添加地址
	address, err = app.AddressService.AddAddressForUser(user, address)
	if err != nil {
		t.Fatalf("Failed to add address for user: %v", err)
	}

	// 使用 UserService 为用户设置默认收货地址
	user, err = app.UserService.SetDefaultShippingAddress(user, address.ID)
	if err != nil {
		t.Fatalf("Failed to set default shipping address for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 GetDefaultShippingAddress
	app.AddressController.GetDefaultShippingAddress(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var defaultShippingAddress *models.Address
	err = json.Unmarshal(w.Body.Bytes(), &defaultShippingAddress)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if defaultShippingAddress.ID != address.ID {
		t.Errorf("Expected default shipping address ID %d, got %d", address.ID, defaultShippingAddress.ID)
	}
}

func TestAddressGetDefaultBillingAddress(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建一个地址
	address := &models.Address{
		Line1:      "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		Country:    "USA",
		PostalCode: "12345",
	}

	// 使用 AddressService 为用户添加地址
	address, err = app.AddressService.AddAddressForUser(user, address)
	if err != nil {
		t.Fatalf("Failed to add address for user: %v", err)
	}

	// 使用 UserService 为用户设置默认账单地址
	user, err = app.UserService.SetDefaultBillingAddress(user, address.ID)
	if err != nil {
		t.Fatalf("Failed to set default billing address for user: %v", err)
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 GetDefaultBillingAddress
	app.AddressController.GetDefaultBillingAddress(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	var defaultBillingAddress *models.Address
	err = json.Unmarshal(w.Body.Bytes(), &defaultBillingAddress)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if defaultBillingAddress.ID != address.ID {
		t.Errorf("Expected default billing address ID %d, got %d", address.ID, defaultBillingAddress.ID)
	}
}

func TestAddressDeleteAllAddressesForUser(t *testing.T) {
	// 初始化测试应用
	app, err := tests.Setup("address")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建两个地址
	addresses := []*models.Address{
		{
			Line1:      "123 Main St",
			City:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			PostalCode: "12345",
		},
		{
			Line1:      "456 Market St",
			City:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			PostalCode: "67890",
		},
	}

	// 使用 AddressService 为用户添加地址
	for _, address := range addresses {
		_, err := app.AddressService.AddAddressForUser(user, address)
		if err != nil {
			t.Fatalf("Failed to add address for user: %v", err)
		}
	}

	// 准备一个测试上下文并设置用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", user)

	// 调用 DeleteAllAddressesForUser
	app.AddressController.DeleteAllAddressesForUser(c)

	// 检查响应
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", c.Writer.Status())
	}

	// 确保所有地址都已删除
	respAddresses, err := app.AddressService.GetAddressesByUserID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get addresses for user: %v", err)
	}

	if len(respAddresses) != 0 {
		t.Errorf("Expected 0 addresses, got %d", len(respAddresses))
	}
}
