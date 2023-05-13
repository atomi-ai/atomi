package models

import (
	"time"
)

type Order struct {
	ID              int64       `gorm:"primaryKey" json:"id"`
	UserID          int64       `gorm:"column:user_id" json:"user_id"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at" json:"updated_at"`
	PaymentIntentID *string     `gorm:"column:payment_intent_id;unique" json:"payment_intent_id"`
	DeliveryID      *string     `gorm:"column:delivery_id;unique" json:"delivery_id"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	DisplayStatus   string      `gorm:"-" json:"display_status"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID        int64    `gorm:"primaryKey" json:"id"`
	OrderID   int64    `gorm:"column:order_id" json:"order_id"`
	Product   *Product `gorm:"foreignKey:ProductID" json:"product"`
	ProductID int64    `gorm:"column:product_id" json:"product_id"`
	Quantity  int      `gorm:"column:quantity" json:"quantity"`
}
