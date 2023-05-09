package utils

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
)

type StripeWrapper interface {
	CreateCustomer(email string) (*stripe.Customer, error)
}

type StripeWrapperImpl struct{}

func NewStripeWrapper() StripeWrapper {
	return &StripeWrapperImpl{}
}

func (s *StripeWrapperImpl) CreateCustomer(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	return customer.New(params)
}
