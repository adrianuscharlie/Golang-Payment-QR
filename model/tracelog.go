package model

import "time"

type Tracelog struct {
	ID        uint   `gorm:"primarykey"`
	Trxid     string `gorm:"index"`
	Message   string `gorm:"text"`
	Process   string `gorm:"size=255;"`
	TimeTrace time.Time
}
