package model

import "time"

type Transaction struct {
	TrxId           string `gorm:"primaryKey"`
	ProductCode     string `gorm:"size=255;not null"`
	TimeStamp       time.Time
	Amount          float64 `gorm:"type:decimal(10,2);not null"`
	QrContent       string  `gorm:"type:text;not null"`
	BranchId        string  `gorm:"size=255;not null"`
	CounterId       string  `gorm:"size=255;not null"`
	CaCode          string  `gorm:"size=255; not null"`
	Request         string  `gorm:"type:text"`
	Response        string  `gorm:"type:text"`
	ResponseCode    string  `gorm:"size=255"`
	ResponseMessage string  `gorm:"size=255"`
	TrxConfirm      string  `gorm:"size=255"`
	ReversalDate    time.Time
	ReversalCode    string `gorm:"size=255"`
	ReversalMessage string `gorm:"size=255"`
}
