package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type TaxRateRepository interface {
	FindByZipCode(zipCode string) (*models.TaxRate, error)
	FindByZipCodeAndState(zipCode, state string) (*models.TaxRate, error)
}

type taxRateRepositoryImpl struct {
	db *gorm.DB
}

func NewTaxRateRepository(db *gorm.DB) TaxRateRepository {
	return &taxRateRepositoryImpl{db: db}
}

func (repo *taxRateRepositoryImpl) FindByZipCode(zipCode string) (*models.TaxRate, error) {
	var taxRate models.TaxRate
	err := repo.db.Where("zip_code = ?", zipCode).First(&taxRate).Error
	return &taxRate, err
}

func (repo *taxRateRepositoryImpl) FindByZipCodeAndState(zipCode, state string) (*models.TaxRate, error) {
	var taxRate models.TaxRate
	err := repo.db.Where("zip_code = ? AND tax_state = ?", zipCode, state).First(&taxRate).Error
	return &taxRate, err
}
