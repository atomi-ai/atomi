package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type UserStoreRepository interface {
	FindDefaultUserStore(userID int64) (*models.UserStore, error)
	Save(userStore *models.UserStore) error
	DisableDefaultStore(userID int64) error
}

type userStoreRepositoryImpl struct {
	db *gorm.DB
}

func NewUserStoreRepository(db *gorm.DB) UserStoreRepository {
	return &userStoreRepositoryImpl{db: db}
}

func (r *userStoreRepositoryImpl) FindDefaultUserStore(userID int64) (*models.UserStore, error) {
	var userStore models.UserStore
	err := r.db.Preload("Store").Where("user_id = ? AND is_enable = ?", userID, true).First(&userStore).Error
	if err != nil {
		return nil, err
	}
	return &userStore, nil
}

func (r *userStoreRepositoryImpl) Save(userStore *models.UserStore) error {
	return r.db.Save(userStore).Error
}

func (r *userStoreRepositoryImpl) DisableDefaultStore(userID int64) error {
	return r.db.Model(&models.UserStore{}).Where("user_id = ? AND is_enable = ?", userID, true).Update("is_enable", false).Error
}
