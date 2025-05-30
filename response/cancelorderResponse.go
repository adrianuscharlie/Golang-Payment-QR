package response

type CancelOrderResponse struct {
	ResponseCode               string `json:"responseCode" binding:"required"`
	ResponseMessage            string `json:"responseMessage" binding:"required"`
	OriginalReferenceNo        string `json:"originalReferenceNo"`
	OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo" binding:"required"`
	CancelTime                 string `json:"cancelTime"`
	TransactionDate            string `json:"transactionDate"`
}
