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
	CheckStatusPayment(request.CreateQueryPaymentRequest) (*response.CreateQueryPaymentResponse, error)
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

func (q queryPaymentService) CheckStatusPayment(req request.CreateQueryPaymentRequest) (*response.CreateQueryPaymentResponse, error) {
	var responseStatus response.CreateQueryPaymentResponse
	// Get Product Configuration
	productConfig, err := q.productRepo.GetConfig(req.ProductCode)
	if err != nil {
		return nil, q.logAndWrapError(req.TrxId, "Error Get Product Config", "Product Config", err)
	}
	qrStatus := request.QueryPaymentRequest{
		OriginalPartnerReferenceNo: req.TrxId,
		OriginalReferenceNo:        "",
		MerchantId:                 defaultMerchantID,
		ServiceCode:                "60",
	}
	urlConfig, err := q.productRepo.GetUrlConfig(req.ProductCode, "STATUS")
	if err != nil {
		return nil, q.logAndWrapError(req.TrxId, "Error Get URL for Product Code: "+req.ProductCode, "URL Config", err)
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
		return nil, q.logAndWrapError(req.TrxId, "Error Sending HTTP POST", "SEND HTTP POST", err)
	}

	q.tracelogRepo.Log(qrStatus.OriginalPartnerReferenceNo, "Sending HTTP POST : \n"+headerStr+"\n"+string(body), "SEND HTTP POST")

	var qrResponse response.QueryPaymentResponse
	err = json.Unmarshal(body, &qrResponse)
	if err != nil {
		return nil, q.logAndWrapError(req.TrxId, "failed to unmarshal response", "Deserialize Json", err)
	}

	responseStatus = response.CreateQueryPaymentResponse{
		TrxId:                   req.TrxId,
		ResponseCode:            qrResponse.ResponseCode,
		ResponseMessage:         qrResponse.ResponseMessage,
		LatestTransactionStatus: qrResponse.LatestTransactionStatus,
		PaidTime:                qrResponse.PaidTime,
		Amount:                  qrResponse.Amount.Value,
	}

	return &responseStatus, nil
}
