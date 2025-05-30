package repository

import (
	"qr-payment/model"
	"time"

	"gorm.io/gorm"
)

type TracelogRepository interface {
	Log(trxId string, message string, process string) error
}

type tracelogRepository struct {
	db *gorm.DB
}

func NewTracelogRepository(db *gorm.DB) TracelogRepository {
	return &tracelogRepository{db}
}

func (r tracelogRepository) Log(trxId string, message string, process string) error {
	trace := model.Tracelog{
		Trxid:     trxId,
		Message:   message,
		Process:   process,
		TimeTrace: time.Now(),
	}
	err := r.db.Create(&trace).Error
	return err
}
