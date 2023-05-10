package services

import (
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

type UserService interface {
	SetDefaultShippingAddress(user *models.User, addressID int64) (*models.User, error)
	SetDefaultBillingAddress(user *models.User, addressID int64) (*models.User, error)
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

func (us *userService) SetDefaultShippingAddress(user *models.User, addressID int64) (*models.User, error) {
	user.DefaultShippingAddressID = addressID
	return us.UserRepo.Save(user)
}

func (us *userService) SetDefaultBillingAddress(user *models.User, addressID int64) (*models.User, error) {
	user.DefaultBillingAddressID = addressID
	return us.UserRepo.Save(user)
}

func (us *userService) SetCurrentPaymentMethod(user *models.User, paymentMethodID *string) (*models.User, error) {
	user.PaymentMethodID = paymentMethodID
	return us.UserRepo.Save(user)
}
