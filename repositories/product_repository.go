package repositories

import (
	"github.com/atomi-ai/atomi/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	Save(product *models.Product) error
	FindByID(id int64) (*models.Product, error)
	FindAll() ([]models.Product, error)
	Update(product *models.Product) error
	Delete(product *models.Product) error
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

// Save saves the given product in the database
func (r *productRepositoryImpl) Save(product *models.Product) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "creator_id", "description", "price", "discount", "category", "image_url"}),
	}).Save(product).Error
}

// FindByID finds a product by its ID
func (r *productRepositoryImpl) FindByID(id int64) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindAll retrieves all products from the database
func (r *productRepositoryImpl) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Update updates the given product in the database
func (r *productRepositoryImpl) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes the given product from the database
func (r *productRepositoryImpl) Delete(product *models.Product) error {
	return r.db.Delete(product).Error
}
