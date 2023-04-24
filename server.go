package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

var db *gorm.DB
var err error
var firebaseApp *firebase.App
var UserRepository repositories.UserRepository

func main() {
	db = models.InitDB()
	UserRepository = repositories.NewUserRepository(db)

	// Initialize Firebase app, set your Firebase local emulator URL for testing.
	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	opt := option.WithCredentialsFile("testing/testing-firebase-secret.json")
	firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error initializing firebase app:", err)
		os.Exit(1)
	}

	r := gin.Default()
	r.Use(AuthMiddleware())
	r.GET("/api/login", Login)
	// Add other endpoints here

	r.Run(":8081")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		ctx := context.Background()
		client, err := firebaseApp.Auth(ctx)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Error getting Auth client"})
			return
		}

		decodedToken, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// Set the decoded token in the request context
		c.Set("decodedToken", decodedToken)
		c.Next()
	}
}

func Login(c *gin.Context) {
	decodedToken, _ := c.Get("decodedToken")
	token := decodedToken.(*auth.Token)

	email, _ := token.Claims["email"].(string)
	user, err := UserRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &models.User{
				Email: email,
				Role:  "USER",
			}
		} else {
			c.JSON(500, gin.H{"error": "Error fetching user"})
			return
		}
	}

	dirty := false
	if user.StripeCustomerID == "" {
		dirty = true
		stripeCustomer, err := createStripeCustomer(email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error creating Stripe customer"})
			return
		}
		user.StripeCustomerID = stripeCustomer.ID
	}

	if dirty {
		UserRepository.Save(user)
	}

	c.JSON(200, user)
}

func createStripeCustomer(email string) (*stripe.Customer, error) {
	stripe.Key = "sk_test_x7J2qxqTLBNo4WQoYkRNMEGx"

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	return customer.New(params)
}
