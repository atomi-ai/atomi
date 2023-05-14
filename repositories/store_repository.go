package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type StoreRepository interface {
	FindAll() ([]*models.Store, error)
	FindByID(id int64) (*models.Store, error)
	CheckUserHasAccessToStore(mgr *models.User, storeID int64) bool
}

type storeRepositoryImpl struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (s *storeRepositoryImpl) FindAll() ([]*models.Store, error) {
	var stores []*models.Store
	err := s.db.Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *storeRepositoryImpl) FindByID(id int64) (*models.Store, error) {
	var store models.Store
	err := s.db.First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// CheckUserHasAccessToStore 查找用户与商店之间的关联关系
func (s *storeRepositoryImpl) CheckUserHasAccessToStore(mgr *models.User, storeID int64) bool {
	var mgrStore models.ManagerStores
	err := s.db.Where(&models.ManagerStores{UserID: mgr.ID, StoreID: storeID}).First(&mgrStore).Error
	return err == nil
}
