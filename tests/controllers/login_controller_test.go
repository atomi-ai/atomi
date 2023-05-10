package controllers

import (
	"github.com/atomi-ai/atomi/tests"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func TestLoginController_Login(t *testing.T) {
	app, err := tests.Setup("login")
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}

	// 创建一个模拟的 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", nil)

	// 为模拟的 gin.Context 设置 decodedToken
	mockToken := &auth.Token{
		Claims: map[string]interface{}{
			"email": "test@example.com",
		},
	}
	c.Set("decodedToken", mockToken)

	// 调用 LoginController 的 Login 方法
	app.LoginController.Login(c)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 检查响应体
	respBody := w.Body.String()

	expectedFields := []string{
		`"email":"test@example.com"`,
		`"role":"USER"`,
		`"phone":""`,
		`"name":""`,
		`"default_shipping_address_id":0`,
		`"default_billing_address_id":0`,
		`"stripe_customer_id":"cus_mock_id"`,
		`"payment_method_id":null`,
	}

	for _, expectedField := range expectedFields {
		if !strings.Contains(respBody, expectedField) {
			t.Errorf("Expected response body to contain %q, but got %q", expectedField, respBody)
		}
	}
}
