package response

type CreateCancelOrderResponse struct {
	ResponseCode    string `json:"responseCode" binding:"required"`
	ResponseMessage string `json:"responseMessage" binding:"required"`
	TrxId           string `json:"trxId"`
	CancelTime      string `json:"cancelTime"`
	TransactionDate string `json:"transactionDate"`
}
