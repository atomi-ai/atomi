package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindByUserID(userID int64) ([]models.Order, error)
	GetByID(orderID int64) (*models.Order, error)
	Save(order *models.Order) error
}

type OrderItemRepository interface {
	Save(orderItem *models.OrderItem) error
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

type orderItemRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepositoryImpl{db: db}
}

func (repo *orderRepositoryImpl) FindByUserID(userID int64) ([]models.Order, error) {
	var orders []models.Order
	err := repo.db.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (repo *orderRepositoryImpl) GetByID(orderID int64) (*models.Order, error) {
	var order models.Order
	err := repo.db.Preload("OrderItems.Product").First(&order, orderID).Error
	return &order, err
}

func (repo *orderRepositoryImpl) Save(order *models.Order) error {
	return repo.db.Save(order).Error
}

func (repo *orderItemRepositoryImpl) Save(orderItem *models.OrderItem) error {
	return repo.db.Save(orderItem).Error
}
