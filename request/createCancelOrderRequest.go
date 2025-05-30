package request

type CreateCancelOrderRequest struct {
	TrxId       string `json:"trxId" binding:"required"`
	ProductCode string `json:"productCode" binding:"required"`
}
