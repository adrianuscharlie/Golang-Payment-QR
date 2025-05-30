package response

import "qr-payment/model"

type QueryPaymentResponse struct {
	ResponseCode               string       `json:"responseCode" binding:"required"`
	ResponseMessage            string       `json:"responseMessage" binding:"required"`
	OriginalPartnerReferenceNo string       `json:"originalPartnerReferenceNo" binding:"required"`
	OriginalReferenceNo        string       `json:"originalReferenceNo" binding:"required"`
	ServiceCode                string       `json:"serviceCode" binding:"required"`
	LatestTransactionStatus    string       `json:"latestTransactionStatus" binding:"required"`
	TransAmount                model.Amount `json:"transAmount"`
	Amount                     model.Amount `json:"amount"`
	PaidTime                   string       `json:"paidTime"`
	Title                      string       `json:"title"`
}
