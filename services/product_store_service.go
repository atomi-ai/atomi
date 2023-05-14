package services

import (
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

// ProductStoreService interface
type ProductStoreService interface {
	ConnectStoreAndProducts(store *models.Store, products []*models.Product) error
	CreateProductInStore(user *models.User, storeID int64, product *models.Product) (*models.Product, error)
}

// productStoreServiceImpl represents the implementation of ProductStoreService
type productStoreServiceImpl struct {
	productStoreRepository repositories.ProductStoreRepository
	productRepository      repositories.ProductRepository
}

// NewProductStoreService creates a new ProductStoreService instance
func NewProductStoreService(productRepository repositories.ProductRepository, productStoreRepo repositories.ProductStoreRepository) ProductStoreService {
	return &productStoreServiceImpl{
		productRepository:      productRepository,
		productStoreRepository: productStoreRepo,
	}
}

// ConnectStoreAndProducts connects a store with a list of products
func (s *productStoreServiceImpl) ConnectStoreAndProducts(store *models.Store, products []*models.Product) error {
	productStores := make([]*models.ProductStore, len(products))

	for i, product := range products {
		productStore := &models.ProductStore{
			BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Store:     store,
			Product:   product,
			IsEnable:  true,
		}
		productStores[i] = productStore
	}

	return s.productStoreRepository.SaveAll(productStores)
}

func (s *productStoreServiceImpl) CreateProductInStore(user *models.User, storeID int64, product *models.Product) (*models.Product, error) {
	product.CreatorID = user.ID
	if err := s.productRepository.Save(product); err != nil {
		return nil, err
	}

	if err := s.productStoreRepository.AddProductToStore(storeID, product.ID); err != nil {
		return nil, err
	}

	return product, nil
}
