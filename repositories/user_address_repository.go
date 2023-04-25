package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type UserAddressRepository interface {
	FindAddressesByUserID(userID int64) ([]*models.Address, error)
	FindByUserIDAndAddressID(userID, addressID int64) (*models.UserAddress, error)
	Save(userAddress *models.UserAddress) (*models.UserAddress, error)
	Delete(userAddress *models.UserAddress) error
	DeleteAllByUserID(userID int64) error
}

type userAddressRepository struct {
	db *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) UserAddressRepository {
	return &userAddressRepository{db}
}

func (uar *userAddressRepository) FindAddressesByUserID(userID int64) ([]*models.Address, error) {
	var addresses []*models.Address
	err := uar.db.Table("user_addresses").Select("addresses.*").
		Joins("JOIN addresses ON user_addresses.address_id = addresses.id").
		Where("user_addresses.user_id = ?", userID).
		Scan(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (uar *userAddressRepository) FindByUserIDAndAddressID(userID, addressID int64) (*models.UserAddress, error) {
	var userAddress models.UserAddress
	err := uar.db.Where("user_id = ? AND address_id = ?", userID, addressID).First(&userAddress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userAddress, nil
}

func (uar *userAddressRepository) Save(userAddress *models.UserAddress) (*models.UserAddress, error) {
	err := uar.db.Save(userAddress).Error
	if err != nil {
		return nil, err
	}
	return userAddress, nil
}

func (uar *userAddressRepository) Delete(userAddress *models.UserAddress) error {
	return uar.db.Delete(userAddress).Error
}

func (uar *userAddressRepository) DeleteAllByUserID(userID int64) error {
	err := uar.db.Where("user_id = ?", userID).Delete(models.UserAddress{}).Error
	return err
}
