package repositories

import (
	"github.com/atomi-ai/atomi/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManagerStoreRepository interface {
	Save(store *models.Store) error
	DeleteStore(storeID int64) error
	AssignStoreToUser(storeID, userID int64) error
	GetStoresByManagerID(mgrID int64) ([]models.Store, error)
}

type managerStoreRepositoryImpl struct {
	db *gorm.DB
}

func NewManagerStoreRepository(db *gorm.DB) ManagerStoreRepository {
	return &managerStoreRepositoryImpl{
		db: db,
	}
}

func (r *managerStoreRepositoryImpl) Save(store *models.Store) error {
	log.Infof("xfguo: before saving store: %v", store)
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "address", "city", "state", "zip_code", "phone", "updated_at"}),
	}).Save(store).Error
	log.Infof("xfguo: saved store: %v", store)
	return err
}

func (r *managerStoreRepositoryImpl) DeleteStore(storeID int64) error {
	return r.db.Where("id = ?", storeID).Delete(&models.Store{}).Error
}

func (r *managerStoreRepositoryImpl) AssignStoreToUser(storeID, userID int64) error {
	return r.db.Create(&models.ManagerStores{UserID: userID, StoreID: storeID}).Error
}

func (r *managerStoreRepositoryImpl) GetStoresByManagerID(mgrID int64) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Table("stores").
		Joins("INNER JOIN manager_stores ON manager_stores.store_id = stores.id").
		Where("manager_stores.user_id = ?", mgrID).
		Find(&stores).Error
	return stores, err
}
