package repositories

import (
	"github.com/atomi-ai/atomi/models" // Replace with your actual repo import path
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	GetByID(userID uint64) (*models.User, error)
	Save(user *models.User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (repo *userRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repo *userRepositoryImpl) GetByID(userID uint64) (*models.User, error) {
	var user models.User
	err := repo.db.First(&user, userID).Error
	return &user, err
}

func (repo *userRepositoryImpl) Save(user *models.User) error {
	return repo.db.Save(user).Error
}
