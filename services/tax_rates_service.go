package services

import (
	"errors"
	"strings"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"gorm.io/gorm"
)

type TaxRateService interface {
	GetTaxRateByZipCodeAndState(address *models.Address) (*models.TaxRate, error)
}

type taxRateServiceImpl struct {
	TaxRateRepo repositories.TaxRateRepository
}

func NewTaxRateService(taxRateRepo repositories.TaxRateRepository) TaxRateService {
	return &taxRateServiceImpl{
		TaxRateRepo: taxRateRepo,
	}
}

func (s *taxRateServiceImpl) GetTaxRateByZipCodeAndState(address *models.Address) (*models.TaxRate, error) {
	state := strings.ToUpper(address.State)
	if state == "" || len(state) != 2 {
		return s.getTaxRateByZipCode(address.PostalCode)
	}

	taxRate, err := s.TaxRateRepo.FindByZipCodeAndState(address.PostalCode, state)
	if err == nil {
		return taxRate, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return s.getTaxRateByZipCode(address.PostalCode)
}

func (s *taxRateServiceImpl) getTaxRateByZipCode(zipCode string) (*models.TaxRate, error) {
	taxRate, err := s.TaxRateRepo.FindByZipCode(zipCode)
	if err != nil {
		return nil, err
	}

	return taxRate, nil
}
