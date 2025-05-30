package services

import (
	"encoding/json"
	"fmt"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/response"
	"qr-payment/utils"
)

type QueryPaymentService interface {
	CheckStatusPayment(trxId string, referenceNo string, merchantId string, productCode string) (*response.QueryPaymentResponse, error)
}

type queryPaymentService struct {
	transactionRepo repository.TransactionRepository
	tracelogRepo    repository.TracelogRepository
	productRepo     repository.ProductConfigRepository
}

func NewQueryPaymentService(transactionRepo repository.TransactionRepository, tracelogRepo repository.TracelogRepository, productRepo repository.ProductConfigRepository) QueryPaymentService {
	return &queryPaymentService{transactionRepo, tracelogRepo, productRepo}
}
func (s *queryPaymentService) logAndWrapError(trxId, message, stage string, err error) error {
	s.tracelogRepo.Log(trxId, message+": "+err.Error(), stage)
	return fmt.Errorf("%s: %w", message, err)
}
func (q queryPaymentService) CheckStatusPayment(trxId string, referenceNo string, merchantId string, productCode string) (*response.QueryPaymentResponse, error) {
	// Get Product Configuration
	productConfig, err := q.productRepo.GetConfig(productCode)
	if err != nil {
		return nil, q.logAndWrapError(trxId, "Error Get Product Config", "Product Config", err)
	}
	qrStatus := request.QueryPaymentRequest{
		OriginalPartnerReferenceNo: trxId,
		OriginalReferenceNo:        referenceNo,
		MerchantId:                 merchantId,
		ServiceCode:                "60",
	}
	urlConfig, err := q.productRepo.GetUrlConfig(productCode, "STATUS")
	if err != nil {
		return nil, q.logAndWrapError(trxId, "Error Get URL for Product Code: "+productCode, "URL Config", err)
	}
	meta := utils.RequestMeta{
		ClientSecret: productConfig.ClientSecret,
		ExtraParam1:  productConfig.ExtraParam1,
		ExtraParam2:  productConfig.ExtraParam2,
		ExtraParam3:  productConfig.ExtraParam3,
		Token:        "your-token-here",
	}

	body, headerStr, err := utils.SendRequest("POST", urlConfig.Url, qrStatus, meta)
	if err != nil {
		return nil, q.logAndWrapError(trxId, "Error Sending HTTP POST", "SEND HTTP POST", err)
	}

	q.tracelogRepo.Log(qrStatus.OriginalPartnerReferenceNo, "Sending HTTP POST : \n"+headerStr+"\n"+string(body), "SEND HTTP POST")

	var qrResponse response.QueryPaymentResponse
	err = json.Unmarshal(body, &qrResponse)
	if err != nil {
		return nil, q.logAndWrapError(trxId, "failed to unmarshal response", "Deserialize Json", err)
	}

	return &qrResponse, nil
}
