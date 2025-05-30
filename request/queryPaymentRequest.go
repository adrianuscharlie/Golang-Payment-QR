package request

type QueryPaymentRequest struct {
	OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo" binding:"required"`
	OriginalReferenceNo        string `json:"originalReferenceNo" binding:"required"`
	ServiceCode                string `json:"serviceCode" binding:"required"`
	MerchantId                 string `json:"merchantId" binding:"required"`
}
