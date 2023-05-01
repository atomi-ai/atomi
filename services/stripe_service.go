package services

import (
	"errors"

	"github.com/atomi-ai/atomi/models"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

type StripeService interface {
	CreateStripeCustomer(email string) (string, error)
	AttachPaymentMethodToCustomer(stripeCustomerId, paymentMethodId string) (*stripe.PaymentMethod, error)
	DeletePaymentMethod(paymentMethodId string) (*stripe.PaymentMethod, error)
	ListPaymentMethods(stripeCustomerId string) (*paymentmethod.Iter, error)
	CreatePaymentIntent(user *models.User, piRequest *models.PaymentIntentRequest, shippingAddr *models.Address) (*stripe.PaymentIntent, error)
	GetLatestCustomerIdByEmail(email string) (string, error)
	ListPaymentIntents(stripeCustomerID string) (*paymentintent.Iter, error)
	RetrievePaymentIntent(intent string) (*stripe.PaymentIntent, error)
}

type StripeServiceImpl struct {
	sc *client.API
}

func NewStripeService() StripeService {
	return &StripeServiceImpl{
		sc: &client.API{},
	}
}

func (s *StripeServiceImpl) CreateStripeCustomer(email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (s *StripeServiceImpl) AttachPaymentMethodToCustomer(stripeCustomerId, paymentMethodId string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(stripeCustomerId),
	}
	updatedPaymentMethod, err := paymentmethod.Attach(paymentMethodId, params)
	if err != nil {
		return nil, err
	}
	return updatedPaymentMethod, nil
}

func (s *StripeServiceImpl) DeletePaymentMethod(paymentMethodId string) (*stripe.PaymentMethod, error) {
	return paymentmethod.Detach(paymentMethodId, nil)
}

func (s *StripeServiceImpl) ListPaymentMethods(stripeCustomerId string) (*paymentmethod.Iter, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeCustomerId),
		Type:     stripe.String("card"),
	}
	return paymentmethod.List(params), nil
}

func (s *StripeServiceImpl) CreatePaymentIntent(user *models.User, piRequest *models.PaymentIntentRequest, shippingAddr *models.Address) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(piRequest.Amount),
		Currency:           stripe.String(piRequest.Currency),
		Customer:           stripe.String(user.StripeCustomerID),
		PaymentMethod:      stripe.String(piRequest.PaymentMethodID),
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
		Confirm:            stripe.Bool(true),
	}

	if shippingAddr != nil {
		params.Shipping = &stripe.ShippingDetailsParams{
			Name: stripe.String(user.Email),
			Address: &stripe.AddressParams{
				Line1:      stripe.String(shippingAddr.Line1),
				Line2:      stripe.String(shippingAddr.Line2),
				City:       stripe.String(shippingAddr.City),
				State:      stripe.String(shippingAddr.State),
				Country:    stripe.String(shippingAddr.Country),
				PostalCode: stripe.String(shippingAddr.PostalCode),
			},
		}
	}

	return paymentintent.New(params)
}

func (s *StripeServiceImpl) GetLatestCustomerIdByEmail(email string) (string, error) {
	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}
	params.Filters.AddFilter("limit", "", "1")

	i := s.sc.Customers.List(params)

	if i.Err() != nil {
		return "", i.Err()
	}

	if !i.Next() {
		return "", errors.New("no customer found")
	}

	return i.Customer().ID, nil
}

func (s *StripeServiceImpl) ListPaymentIntents(stripeCustomerID string) (*paymentintent.Iter, error) {
	params := &stripe.PaymentIntentListParams{
		Customer: stripe.String(stripeCustomerID),
	}
	params.AddExpand("data.latest_charge")

	return paymentintent.List(params), nil
}

func (s *StripeServiceImpl) RetrievePaymentIntent(intent string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{}
	params.AddExpand("latest_charge")

	return paymentintent.Get(intent, params)
}
