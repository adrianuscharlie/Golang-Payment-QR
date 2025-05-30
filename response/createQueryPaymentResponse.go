package response

type CreateQueryPaymentResponse struct {
	ResponseCode            string `json:"responseCode" binding:"required"`
	ResponseMessage         string `json:"responseMessage" binding:"required"`
	TrxId                   string `json:"trxId" binding:"required"`
	LatestTransactionStatus string `json:"latestTransactionStatus" binding:"required"`
	PaidTime                string `json:"paidTime"`
	Amount                  string `json:"amount" binding:"required"`
}
