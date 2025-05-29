package response

import "qr-payment/model"

type PaymentResponse struct {
	ResponseCode       string          `json:"responseCode"`
	ResponseMessage    string          `json:"responseMessage"`
	ReferenceNo        string          `json:"referenceNo"`
	PartnerReferenceNo string          `json:"partnerReferenceNo"`
	TransactionDate    string          `json:"transactionDate"`
	AdditionnalInfo    AdditionnalInfo `json:"additionalInfo"`
}

type AdditionnalInfo struct {
	Amount   model.Amount `json:"amount"`
	PaidTime string       `json:"paidTime"`
}
