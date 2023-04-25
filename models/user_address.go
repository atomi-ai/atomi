package models

import (
	"time"
)

type UserAddress struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	User      *User     `gorm:"foreignKey:UserID" json:"-"`
	UserID    int64     `json:"user_id"`
	Address   *Address  `gorm:"foreignKey:AddressID" json:"-"`
	AddressID int64     `json:"address_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}
