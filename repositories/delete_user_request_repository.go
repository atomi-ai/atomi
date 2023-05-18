package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type DeleteUserRequestRepository interface {
	AddRequest(userID int64) error
}

type deleteUserRequestRepositoryImpl struct {
	db *gorm.DB
}

func NewDeleteUserRequestRepository(db *gorm.DB) DeleteUserRequestRepository {
	return &deleteUserRequestRepositoryImpl{db: db}
}

func (r *deleteUserRequestRepositoryImpl) AddRequest(userID int64) error {
	return r.db.Create(&models.DeleteUserRequest{UserID: userID}).Error
}
