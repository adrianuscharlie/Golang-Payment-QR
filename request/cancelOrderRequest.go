package request

type CancelOrderRequest struct {
	OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo" binding:"required"`
	MerchantId                 string `json:"merchantId" binding:"required"`
}
