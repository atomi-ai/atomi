package services

import (
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

// ProductStoreService interface
type ProductStoreService interface {
	ConnectStoreAndProducts(store *models.Store, products []*models.Product) error
}

// productStoreServiceImpl represents the implementation of ProductStoreService
type productStoreServiceImpl struct {
	productStoreRepository repositories.ProductStoreRepository
}

// NewProductStoreService creates a new ProductStoreService instance
func NewProductStoreService(productStoreRepo repositories.ProductStoreRepository) ProductStoreService {
	return &productStoreServiceImpl{
		productStoreRepository: productStoreRepo,
	}
}

// ConnectStoreAndProducts connects a store with a list of products
func (s *productStoreServiceImpl) ConnectStoreAndProducts(store *models.Store, products []*models.Product) error {
	productStores := make([]*models.ProductStore, len(products))

	for i, product := range products {
		productStore := &models.ProductStore{
			Store:     store,
			Product:   product,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsEnable:  true,
		}
		productStores[i] = productStore
	}

	return s.productStoreRepository.SaveAll(productStores)
}
