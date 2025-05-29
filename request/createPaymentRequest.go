package request

type CreatePaymentRequest struct {
	TrxId       string  `json:"trxId" binding:"required"`
	Amount      float64 `json:"amount,string" binding:"required"`
	QrContent   string  `json:"qrContent" binding:"required"`
	ProductCode string  `json:"productCode" binding:"required"`
	BranchId    string  `json:"branchId" binding:"required"`
	CounterId   string  `json:"counterId" binding:"required"`
	CaCode      string  `json:"caCode" binding:"required"`
}
