package repositories

import (
	"github.com/atomi-ai/atomi/models" // Replace with your actual repo import path
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *userRepositoryImpl) Save(user *models.User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "role", "phone", "name", "default_shipping_address_id", "default_billing_address_id", "stripe_customer_id", "payment_method_id"}),
	}).Save(user).Error
}
