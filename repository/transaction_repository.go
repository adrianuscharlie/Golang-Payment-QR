package repository

import (
	"qr-payment/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Insert(*model.Transaction) error
	Update(*model.Transaction) error
	GetByIdQr(id string, qr string) (model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Insert(m *model.Transaction) error {
	return r.db.Create(m).Error
}

func (r *transactionRepository) Update(m *model.Transaction) error {
	return r.db.Model(&model.Transaction{}).
		Where("trx_id = ? AND qr_content = ?", m.TrxId, m.QrContent).
		Updates(m).Error
}

func (r *transactionRepository) GetByIdQr(id string, qr string) (model.Transaction, error) {
	var trx model.Transaction
	err := r.db.Where("trx_id = ? AND qr_content = ?", id, qr).First(&trx).Error
	return trx, err
}
