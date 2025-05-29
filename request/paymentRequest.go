package request

import "qr-payment/model"

type PaymentRequest struct {
	PartnerReferenceNo string          `json:"partnerReferenceNo" binding:"required"`
	QrContent          string          `json:"qrContent" binding:"required"`
	Amount             model.Amount    `json:"amount" binding:"required"`
	MerchantId         string          `json:"merchantId" binding:"required"`
	Title              string          `json:"title" binding:"required"`
	ExpiryTime         string          `json:"expiryTime" binding:"required"`
	ExternalStoreId    string          `json:"externalStoreid" binding:"required"`
	ScannerInfo        ScannerInfo     `json:"scannerInfo" binding:"required"`
	AdditionnalInfo    AdditionnalInfo `json:"additionalInfo" binding:"required"`
}

type ScannerInfo struct {
	DeviceId string `json:"deviceId" binding:"required"`
	DeviceIp string `jjson:"deviceIp" binding:"required"`
}

type AdditionnalInfo struct {
	QrContentType string `json:"qrContentType" binding:"required"`
	Mcc           string `json:"mcc" binding:"required"`
	ProductCode   string `json:"productCode" binding:"required"`
	// NotifyUrl     string  `json:"notifyUrl"`
	Envinfo Envinfo `json:"envInfo" binding:"required"`
}

type Envinfo struct {
	SourcePplatform   string `json:"sourcePlatform" binding:"required"`
	TerminalType      string `json:"terminalType" binding:"required"`
	OrderterminalType string `json:"orderTerminalType" binding:"reqiired"`
}
