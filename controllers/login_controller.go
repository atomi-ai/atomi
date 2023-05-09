package controllers

import (
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"github.com/atomi-ai/atomi/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	"gorm.io/gorm"
)

type LoginController interface {
	Login(c *gin.Context)
}

type LoginControllerImpl struct {
	UserRepository repositories.UserRepository
	StripeWrapper  utils.StripeWrapper
}

func NewLoginController(userRepo repositories.UserRepository, wrapper utils.StripeWrapper) LoginController {
	return &LoginControllerImpl{
		UserRepository: userRepo,
		StripeWrapper:  wrapper,
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
		stripeCustomer, err := l.createStripeCustomer(email)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error creating Stripe customer"})
			return
		}
		user.StripeCustomerID = stripeCustomer.ID
	}

	if dirty {
		if user, err = l.UserRepository.Save(user); err != nil {
			log.Errorf("Errors in saving user %v, err: \n%v", user, err)
		}
	}

	c.JSON(200, user)
}

func (l *LoginControllerImpl) createStripeCustomer(email string) (*stripe.Customer, error) {
	return l.StripeWrapper.CreateCustomer(email)
}
