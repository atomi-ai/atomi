package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductStoreRepository interface
type ProductStoreRepository interface {
	Save(productStore *models.ProductStore) error
	FindAllByStoreID(storeID int64) ([]*models.ProductStore, error)
	FindByStoreAndProduct(store *models.Store, product *models.Product) (*models.ProductStore, error)
	SaveAll(productStores []*models.ProductStore) error
	AddProductToStore(storeID, productID int64) error
	RemoveProductFromStore(storeID, productID int64) error
}

// productStoreRepositoryImpl represents the implementation of ProductStoreRepository
type productStoreRepositoryImpl struct {
	db *gorm.DB
}

// NewProductStoreRepository creates a new ProductStoreRepository instance
func NewProductStoreRepository(db *gorm.DB) ProductStoreRepository {
	return &productStoreRepositoryImpl{
		db: db,
	}
}

// Save saves the given product store in the database
func (r *productStoreRepositoryImpl) Save(productStore *models.ProductStore) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_enable", "updated_at"}),
	}).Create(productStore).Error
}

// FindAllByStoreID retrieves all product stores by a given store ID
func (r *productStoreRepositoryImpl) FindAllByStoreID(storeID int64) ([]*models.ProductStore, error) {
	var productStores []*models.ProductStore
	err := r.db.Preload("Product").Where("store_id = ?", storeID).Find(&productStores).Error
	if err != nil {
		return nil, err
	}
	return productStores, nil
}

// FindByStoreAndProduct finds a product store by a given store and product
func (r *productStoreRepositoryImpl) FindByStoreAndProduct(store *models.Store, product *models.Product) (*models.ProductStore, error) {
	var productStore models.ProductStore
	err := r.db.Where("store_id = ? AND product_id = ?", store.ID, product.ID).First(&productStore).Error
	if err != nil {
		return nil, err
	}
	return &productStore, nil
}

// SaveAll saves a list of ProductStore instances
func (r *productStoreRepositoryImpl) SaveAll(productStores []*models.ProductStore) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, productStore := range productStores {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "store_id"}, {Name: "product_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_enable", "updated_at"}),
		}).Create(productStore).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *productStoreRepositoryImpl) AddProductToStore(storeID, productID int64) error {
	productStore := &models.ProductStore{StoreID: storeID, ProductID: productID}
	return r.db.Create(productStore).Error
}

func (r *productStoreRepositoryImpl) RemoveProductFromStore(storeID, productID int64) error {
	productStore := &models.ProductStore{StoreID: storeID, ProductID: productID}
	return r.db.Where("store_id = ? AND product_id = ?", storeID, productID).Delete(productStore).Error
}
