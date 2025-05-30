package services

import (
	"encoding/json"
	"fmt"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/response"
	"qr-payment/utils"
)

type CancelOrderService interface {
	CancelOrder(req request.CreateCancelOrderRequest) (*response.CreateCancelOrderResponse, error)
}

type cancelOrderService struct {
	transactionRepo repository.TransactionRepository
	tracelogRepo    repository.TracelogRepository
	productRepo     repository.ProductConfigRepository
}

func NewCancelOrderService(transactionRepo repository.TransactionRepository,
	tracelogRepo repository.TracelogRepository,
	productRepo repository.ProductConfigRepository) CancelOrderService {
	return cancelOrderService{transactionRepo, tracelogRepo, productRepo}
}
func (s *cancelOrderService) logAndWrapError(trxId, message, stage string, err error) error {
	s.tracelogRepo.Log(trxId, message+": "+err.Error(), stage)
	return fmt.Errorf("%s: %w", message, err)
}

func (c cancelOrderService) CancelOrder(req request.CreateCancelOrderRequest) (*response.CreateCancelOrderResponse, error) {
	productConfig, err := c.productRepo.GetConfig(req.ProductCode)
	if err != nil {
		return nil, c.logAndWrapError(req.TrxId, "Error Get Product Config", "Product Config", err)
	}

	cancel := request.CancelOrderRequest{
		OriginalPartnerReferenceNo: req.TrxId,
		MerchantId:                 defaultMerchantID,
	}

	urlConfig, err := c.productRepo.GetUrlConfig(req.ProductCode, "CANCEL")
	if err != nil {
		return nil, c.logAndWrapError(req.TrxId, "Error Get URL for Product Code: "+req.ProductCode, "URL Config", err)
	}
	meta := utils.RequestMeta{
		ClientSecret: productConfig.ClientSecret,
		ExtraParam1:  productConfig.ExtraParam1,
		ExtraParam2:  productConfig.ExtraParam2,
		ExtraParam3:  productConfig.ExtraParam3,
		Token:        "your-token-here",
	}

	body, headerStr, err := utils.SendRequest("POST", urlConfig.Url, cancel, meta)
	if err != nil {
		return nil, c.logAndWrapError(req.TrxId, "Error Sending HTTP POST", "SEND HTTP POST", err)
	}
	c.tracelogRepo.Log(cancel.OriginalPartnerReferenceNo, "Sending HTTP POST : \n"+headerStr+"\n"+string(body), "SEND HTTP POST")

	var cancelResponse response.CancelOrderResponse

	err = json.Unmarshal(body, &cancelResponse)
	if err != nil {
		return nil, c.logAndWrapError(req.TrxId, "failed to unmarshal response", "Deserialize Json", err)
	}

	cResponse := response.CreateCancelOrderResponse{
		ResponseCode:    cancelResponse.ResponseCode,
		ResponseMessage: cancelResponse.ResponseMessage,
		TrxId:           cancel.OriginalPartnerReferenceNo,
		CancelTime:      cancelResponse.CancelTime,
		TransactionDate: cancelResponse.TransactionDate,
	}

	return &cResponse, nil
}
