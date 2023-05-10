package middlewares

import (
	"github.com/atomi-ai/atomi/app"
	"github.com/atomi-ai/atomi/middlewares"
	"github.com/atomi-ai/atomi/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// 初始化测试应用
	app, err := app.InitializeTestingApplication("auth")
	if err != nil {
		t.Fatalf("Failed to initialize testing application: %v", err)
	}

	// 创建一个用户
	user := &models.User{Name: "John Doe", Email: "john.doe@example.com"}
	if user, err = app.UserRepository.Save(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 初始化AuthMiddleware
	authMiddleware := middlewares.NewAuthMiddleware(app.UserRepository, app.AuthWrapper)

	// 准备一个测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置一个模拟ID令牌
	idToken := "test-id-token"
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+idToken)

	// 调用AuthMiddleware的handler
	authMiddleware.Handler()(c)

	// 检查用户是否正确设置
	if cUser, exists := c.Get("user"); !exists || cUser.(*models.User).Email != user.Email {
		t.Errorf("User not set correctly in context")
	}
}
