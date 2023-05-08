package controllers

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"gorm.io/gorm"
)

type LoginController interface {
	Login(c *gin.Context)
}

type LoginControllerImpl struct {
	UserRepository repositories.UserRepository
}

func NewLoginController(userRepo repositories.UserRepository) LoginController {
	return &LoginControllerImpl{
		UserRepository: userRepo,
	}
}

func (l *LoginControllerImpl) Login(c *gin.Context) {
	token, _ := c.MustGet("decodedToken").(*auth.Token)

	email, _ := token.Claims["email"].(string)
	user, err := l.UserRepository.FindByEmail(email)
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
		l.UserRepository.Save(user)
	}

	c.JSON(200, user)
}

func createStripeCustomer(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	return customer.New(params)
}
