package repository

import (
	"qr-payment/model"

	"gorm.io/gorm"
)

type ProductConfigRepository interface {
	GetConfig(productCode string) (*model.ProductConfig, error)
	GetUrlConfig(productCode string) (*model.UrlConfig, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductConfigRepository(db *gorm.DB) ProductConfigRepository {
	return &productRepository{db}
}

func (r *productRepository) GetConfig(productCode string) (*model.ProductConfig, error) {
	var productConfig model.ProductConfig
	err := r.db.First(&productConfig, productCode).Error
	return &productConfig, err
}

func (r *productRepository) GetUrlConfig(productCode string) (*model.UrlConfig, error) {
	var urlConfig model.UrlConfig
	err := r.db.First(&urlConfig, productCode).Error
	return &urlConfig, err
}
