package services

import (
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

type UserService interface {
	SetDefaultShippingAddress(user *models.User, addressID int64) (error, error)
	SetDefaultBillingAddress(user *models.User, addressID int64) (error, error)
	SetCurrentPaymentMethod(user *models.User, paymentMethodID *string) (*models.User, error)
}

type userService struct {
	UserRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		UserRepo: userRepo,
	}
}

func (us *userService) SetDefaultShippingAddress(user *models.User, addressID int64) (error, error) {
	user.DefaultShippingAddressID = addressID
	_, err := us.UserRepo.Save(user)
	return err, nil
}

func (us *userService) SetDefaultBillingAddress(user *models.User, addressID int64) (error, error) {
	user.DefaultBillingAddressID = addressID
	_, err := us.UserRepo.Save(user)
	return err, nil
}

func (us *userService) SetCurrentPaymentMethod(user *models.User, paymentMethodID *string) (*models.User, error) {
	user.PaymentMethodID = paymentMethodID
	return us.UserRepo.Save(user)
}
