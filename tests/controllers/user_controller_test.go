package controllers

import (
	"github.com/atomi-ai/atomi/tests"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
)

func TestUserController_GetUser(t *testing.T) {
	app, err := tests.Setup("user")
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	//c.Request = httptest.NewRequest("GET", "/user", nil)

	mockUser := &models.User{
		BaseModel: models.BaseModel{ID: 1},
	}
	c.Set("user", mockUser)

	app.UserController.GetUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	respBody := w.Body.String()

	expectedFields := []string{
		`"id":1`,
	}

	for _, expectedField := range expectedFields {
		if !strings.Contains(respBody, expectedField) {
			t.Errorf("Expected response body to contain %q, but got %q", expectedField, respBody)
		}
	}
}

func TestUserController_SetCurrentPaymentMethod(t *testing.T) {
	app, err := tests.Setup("user")
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	//c.Request = httptest.NewRequest("PUT", "/user/current-payment-method/123", nil)
	// 添加路径参数
	c.Params = []gin.Param{
		{
			Key:   "paymentMethodId",
			Value: "123",
		},
	}

	mockUser := &models.User{
		BaseModel: models.BaseModel{ID: 1},
	}
	c.Set("user", mockUser)

	app.UserController.SetCurrentPaymentMethod(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	respBody := w.Body.String()

	expectedFields := []string{
		`"id":1`,
		`"payment_method_id":"123"`,
	}

	for _, expectedField := range expectedFields {
		if !strings.Contains(respBody, expectedField) {
			t.Errorf("Expected response body to contain %q, but got %q", expectedField, respBody)
		}
	}
}
