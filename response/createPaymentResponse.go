package response

type CreatePaymentResponse struct {
	TrxId           string `json:"trxId" binding:"required"`
	TrxConfirm      string `json:"trxConfirm" binding:"required"`
	ResponseCode    string `json:"responseCode" binding:"required"`
	ResponseMessage string `json:"responseMessage" binding:"required"`
	PaidAt          string `json:"paidAt" binding:"required"`
}
