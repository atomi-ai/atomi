package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StoreRepository interface
type StoreRepository interface {
	Save(store *models.Store) error
	FindAll() ([]*models.Store, error)
	FindByID(id int64) (*models.Store, error)
}

// storeRepositoryImpl represents the implementation of StoreRepository
type storeRepositoryImpl struct {
	db *gorm.DB
}

// NewStoreRepository creates a new StoreRepository instance
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) Save(store *models.Store) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "address", "city", "state", "zip_code", "phone", "updated_at"}),
	}).Save(store).Error
}

// FindAll retrieves all stores from the database
func (r *storeRepositoryImpl) FindAll() ([]*models.Store, error) {
	var stores []*models.Store
	err := r.db.Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

// FindByID finds a store by its ID
func (r *storeRepositoryImpl) FindByID(id int64) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}
