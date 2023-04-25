package models

import (
	"time"
)

type Address struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Line1      string    `json:"line1"`
	Line2      string    `json:"line2"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `gorm:"column:postal_code" json:"postal_code"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"-"`
}
