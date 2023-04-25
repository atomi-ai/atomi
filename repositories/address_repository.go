package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByID(id int64) (*models.Address, error)
	Save(address *models.Address) (*models.Address, error)
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db}
}

func (ar *addressRepository) FindByID(id int64) (*models.Address, error) {
	var address models.Address
	err := ar.db.First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (ar *addressRepository) Save(address *models.Address) (*models.Address, error) {
	err := ar.db.Save(address).Error
	if err != nil {
		return nil, err
	}
	return address, nil
}
