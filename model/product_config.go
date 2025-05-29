package model

type ProductConfig struct {
	ProductCode  string `gorm:"primaryKey;size:255"`
	ProductType  string `gorm:"size:255;not null"`
	Partner      string `gorm:"size:255;not null"`
	Connection   string `gorm:"size:255;not null"`
	ClientId     string `gorm:"size:255;not null"`
	ClientSecret string `gorm:"size:255; not null"`
	UseToken     bool   `gorm:"not null"`
	ExtraParam1  string `gorm:"size:255;not null"`
	ExtraParam2  string `gorm:"size:255;not null"`
	ExtraParam3  string `gorm:"size:255;not null"`
}

type UrlConfig struct {
	ProductCode string `gorm:"not null"`
	Url         string `gorm:"not null"`
	TrxType     string `gorm:"not null"`
}
