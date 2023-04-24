package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

// ConfigRepository interface
type ConfigRepository interface {
	Save(config *models.Config) error
	FindAll() ([]*models.Config, error)
	FindByKey(key string) (*models.Config, error)
}

// configRepositoryImpl represents the implementation of ConfigRepository
type configRepositoryImpl struct {
	db *gorm.DB
}

// NewConfigRepository creates a new ConfigRepository instance
func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepositoryImpl{
		db: db,
	}
}

// Save saves a Config instance
func (r *configRepositoryImpl) Save(config *models.Config) error {
	return r.db.Save(config).Error
}

// FindAll finds all Config instances
func (r *configRepositoryImpl) FindAll() ([]*models.Config, error) {
	var configs []*models.Config
	err := r.db.Find(&configs).Error
	return configs, err
}

// FindByKey finds a Config instance by key
func (r *configRepositoryImpl) FindByKey(key string) (*models.Config, error) {
	var config models.Config
	err := r.db.First(&config, "config_key = ?", key).Error
	return &config, err
}
