package models

import (
	"time"
)

// ProductCategory represents the product category
type ProductCategory string

const (
	ProductCategoryFood  ProductCategory = "FOOD"
	ProductCategoryDrink ProductCategory = "DRINK"
	ProductCategoryOther ProductCategory = "OTHER"
)

// Product represents the product entity
type Product struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"unique" json:"name"`
	Creator     *User           `gorm:"foreignKey:CreatorID" json:"-"`
	CreatorID   int64           `json:"creator_id"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Discount    float64         `json:"discount"`
	Category    ProductCategory `json:"category"`
	ImageURL    string          `gorm:"column:image_url" json:"image_url"`
}

// Store represents the store entity
type Store struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `gorm:"column:zip_code" json:"zip_code"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ProductStore represents the product store entity
type ProductStore struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Store     *Store    `gorm:"foreignKey:StoreID" json:"store"`
	StoreID   int64     `gorm:"uniqueIndex:idx_store_product" json:"store_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product"`
	ProductID int64     `gorm:"uniqueIndex:idx_store_product" json:"product_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsEnable  bool      `gorm:"column:is_enable" json:"is_enable"`
}
