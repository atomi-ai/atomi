package services

import (
	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

type AddressService interface {
	GetAddressesByUserId(userID int64) ([]*models.Address, error)
	AddAddressForUser(user *models.User, address *models.Address) (*models.Address, error)
	DeleteAddressForUser(user *models.User, addressID int64) error
	DeleteAllAddressesForUser(user *models.User) error
}

type addressServiceImpl struct {
	UserRepo        repositories.UserRepository
	AddressRepo     repositories.AddressRepository
	UserAddressRepo repositories.UserAddressRepository
}

// TODO(lamuguo): 感觉没法忍下去了，这种一点一点传参的方式，实在是太需要dependency injection了。
func NewAddressService(userRepo repositories.UserRepository, addressRepo repositories.AddressRepository, userAddressRepo repositories.UserAddressRepository) AddressService {
	return &addressServiceImpl{
		UserRepo:        userRepo,
		AddressRepo:     addressRepo,
		UserAddressRepo: userAddressRepo,
	}
}

func (as *addressServiceImpl) GetAddressesByUserId(userID int64) ([]*models.Address, error) {
	return as.UserAddressRepo.FindAddressesByUserID(userID)
}

func (as *addressServiceImpl) AddAddressForUser(user *models.User, address *models.Address) (*models.Address, error) {
	savedAddr, err := as.AddressRepo.Save(address)
	if err != nil {
		return nil, err
	}

	userAddress := &models.UserAddress{
		Address:   savedAddr,
		User:      user,
		UserID:    user.ID,
		AddressID: savedAddr.ID,
	}

	_, err = as.UserAddressRepo.Save(userAddress)
	if err != nil {
		return nil, err
	}

	return savedAddr, nil
}

func (as *addressServiceImpl) DeleteAddressForUser(user *models.User, addressID int64) error {
	userAddress, err := as.UserAddressRepo.FindByUserIDAndAddressID(user.ID, addressID)
	if err != nil {
		return err
	}

	if userAddress != nil {
		err = as.UserAddressRepo.Delete(userAddress)
		if err != nil {
			return err
		}
	}

	dirty := false
	if user.DefaultShippingAddressID == addressID {
		dirty = true
		user.DefaultShippingAddressID = 0
	}
	if user.DefaultBillingAddressID == addressID {
		dirty = true
		user.DefaultBillingAddressID = 0
	}
	if dirty {
		err = as.UserRepo.Save(user)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAllAddressesForUser @TestingOnly
// Please don't use the function below, it is testing only. Please consider to rewrite it if you need the feature.
//
// 这个函数不会删跟用户相关的address，所以有可能造成很多不用了的address就留下来了。
func (as *addressServiceImpl) DeleteAllAddressesForUser(user *models.User) error {
	// Remove all user x address relations.
	err := as.UserAddressRepo.DeleteAllByUserID(user.ID)

	// Update the user's default addresses.
	dirty := false
	if user.DefaultShippingAddressID != 0 {
		dirty = true
		user.DefaultShippingAddressID = 0
	}
	if user.DefaultBillingAddressID != 0 {
		dirty = true
		user.DefaultBillingAddressID = 0
	}
	if dirty {
		err = as.UserRepo.Save(user)
		if err != nil {
			return err
		}
	}
	return err
}
