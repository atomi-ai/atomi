package middlewares

import (
	"context"
	"strings"

	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AuthMiddleware interface {
	Handler() gin.HandlerFunc
}

type authMiddlewareImpl struct {
	UserRepository repositories.UserRepository
	AuthWrapper    utils.AuthAppWrapper
}

func NewAuthMiddleware(userRepository repositories.UserRepository, authWrapper utils.AuthAppWrapper) AuthMiddleware {
	return &authMiddlewareImpl{
		UserRepository: userRepository,
		AuthWrapper:    authWrapper,
	}
}

func (a authMiddlewareImpl) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Auth: 0")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		log.Debugf("Auth: token: %v", idToken)
		ctx := context.Background()
		decodedToken, err := a.AuthWrapper.AuthAndDecode(ctx, idToken)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Errors in authenticating/decoding the context"})
			return
		}

		log.Debugf("Auth: decodedToken: %v", decodedToken)
		email := decodedToken.Claims["email"].(string)
		user, err := a.UserRepository.FindByEmail(email)
		if err == nil {
			c.Set("user", user)
		}
		log.Debugf("Auth: user: %v", user)

		c.Set("decodedToken", decodedToken)
		c.Next()
	}
}
