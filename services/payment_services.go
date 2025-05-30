package services

import (
	"encoding/json"
	"fmt"
	"qr-payment/model"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/response"
	"qr-payment/utils"
	"strconv"
	"time"
)

type PaymentService interface {
	Payment(req request.CreatePaymentRequest) (*response.CreatePaymentResponse, error)
}

type paymentService struct {
	productRepo         repository.ProductConfigRepository
	transactionRepo     repository.TransactionRepository
	tracelogRepo        repository.TracelogRepository
	queryPaymentService QueryPaymentService
	cancelOrderService  CancelOrderService
}

func NewPaymentService(productRepo repository.ProductConfigRepository, transactionRepo repository.TransactionRepository, tracelogRepo repository.TracelogRepository, queryPaymentService QueryPaymentService, cancelOrderService CancelOrderService) PaymentService {
	return &paymentService{productRepo, transactionRepo, tracelogRepo, queryPaymentService, cancelOrderService}
}
func (s *paymentService) logAndWrapError(trxId, message, stage string, err error) error {
	s.tracelogRepo.Log(trxId, message+": "+err.Error(), stage)
	return fmt.Errorf("%s: %w", message, err)
}

func (s *paymentService) cancelOrderAndRespond(trxId string, merchantId string, productCode string) *response.CreatePaymentResponse {
	s.tracelogRepo.Log(trxId, "Triggering Cancel Order To Partner", "SEND HTTP POST")
	_, err := s.cancelOrderService.CancelOrder(request.CreateCancelOrderRequest{TrxId: trxId, ProductCode: productCode})
	if err != nil {
		s.tracelogRepo.Log(trxId, "Failed cancel order: "+err.Error(), "Cancel Order")
	}
	return &response.CreatePaymentResponse{
		TrxId:           trxId,
		TrxConfirm:      "",
		ResponseCode:    "21",
		ResponseMessage: "Order was canceled because payment failed!",
	}
}

const (
	defaultMerchantID     = "TP18"
	defaultStoreID        = "M2351"
	defaultQrContentType  = "BAR_CODE"
	defaultSourcePlatform = "IPG"
	defaultOrderTerminal  = "WEB"
)

func (s *paymentService) Payment(req request.CreatePaymentRequest) (*response.CreatePaymentResponse, error) {
	var paymentResponse response.CreatePaymentResponse
	s.tracelogRepo.Log(req.TrxId, "Transaction Received with TrxId: "+req.TrxId, "Service IN")
	timeStamp := time.Now().Add(5 * time.Minute).Format(time.RFC3339)

	productConfig, err := s.productRepo.GetConfig(req.ProductCode)
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "Error Get Product Config", "Product Config", err)
	}
	urlConfig, err := s.productRepo.GetUrlConfig(req.ProductCode, "PAYMENT")
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "Error Get URL for Product Code: "+req.ProductCode, "URL Config", err)
	}

	paymentRequest := request.PaymentRequest{
		PartnerReferenceNo: req.TrxId,
		QrContent:          req.QrContent,
		Amount: model.Amount{
			Value:    strconv.Itoa(int(req.Amount)),
			Currency: "IDR",
		},
		MerchantId:      defaultMerchantID,
		Title:           "QRIS CPM Payment via Dana",
		ExpiryTime:      timeStamp,
		ExternalStoreId: defaultStoreID,
		ScannerInfo: request.ScannerInfo{
			DeviceId: "46252",
			DeviceIp: "172.24.281.24",
		},
		AdditionnalInfo: request.AdditionnalInfo{
			QrContentType: defaultQrContentType,
			Mcc:           "5732",
			ProductCode:   req.ProductCode,
			Envinfo: request.Envinfo{
				SourcePplatform:   defaultSourcePlatform,
				TerminalType:      "SYSTEM",
				OrderterminalType: defaultOrderTerminal,
			},
		},
	}

	jsonData, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "Error Creating Request and Serialize into JSON ", "Create Request", err)
	}

	transaction := model.Transaction{
		TrxId:     req.TrxId,
		QrContent: req.QrContent,
		TimeStamp: time.Now(),
		BranchId:  req.BranchId,
		CounterId: req.CounterId,
		Amount:    req.Amount,
		CaCode:    req.CaCode,
		Request:   string(jsonData),
	}
	err = s.transactionRepo.Insert(&transaction)
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "Error Insert transaction into table transaction", "Insert Transaction", err)
	}
	meta := utils.RequestMeta{
		ClientSecret: productConfig.ClientSecret,
		ExtraParam1:  productConfig.ExtraParam1,
		ExtraParam2:  productConfig.ExtraParam2,
		ExtraParam3:  productConfig.ExtraParam3,
		Token:        "your-token-here",
	}

	body, headerStr, err := utils.SendRequest("POST", urlConfig.Url, paymentRequest, meta)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Error Sending HTTP POST :"+err.Error(), "SEND HTTP POST")
		s.tracelogRepo.Log(req.TrxId, "Triggering Cancel Order To Partner :"+err.Error(), "SEND HTTP POST")
		paymentResponse = *s.cancelOrderAndRespond(paymentRequest.PartnerReferenceNo, paymentRequest.MerchantId, req.ProductCode)
		transaction.ResponseCode = paymentResponse.ResponseCode
		transaction.ResponseMessage = paymentResponse.ResponseMessage
		transaction.ReversalCode = paymentResponse.ResponseCode
		transaction.ReversalMessage = paymentResponse.ResponseMessage
		transaction.ReversalDate = time.Now()
		s.transactionRepo.Update(&transaction)
		return &paymentResponse, nil
	}
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "failed to unmarshal response", "Deserialize Json", err)
	}
	var partnerResponse response.PaymentResponse
	s.tracelogRepo.Log(partnerResponse.PartnerReferenceNo, "Sending HTTP POST : \n"+headerStr+"\n"+string(body), "SEND HTTP POST")
	transaction.Response = string(body)
	transaction.ResponseCode = partnerResponse.ResponseCode
	transaction.ResponseMessage = partnerResponse.ResponseMessage
	transaction.TrxConfirm = partnerResponse.ReferenceNo

	if partnerResponse.ResponseCode == "2026000" {
		s.tracelogRepo.Log(req.TrxId, "Transaction pending, check status payment to Partner", "CHECK STATUS")
		checkStatus, err := s.queryPaymentService.CheckStatusPayment(request.CreateQueryPaymentRequest{TrxId: req.TrxId, ProductCode: req.ProductCode})
		if err != nil {
			s.tracelogRepo.Log(req.TrxId, "Error Sending HTTP POST :"+err.Error(), "SEND HTTP POST")
			s.tracelogRepo.Log(req.TrxId, "Triggering Cancel Order To Partner :"+err.Error(), "SEND HTTP POST")
			paymentResponse = *s.cancelOrderAndRespond(paymentRequest.PartnerReferenceNo, paymentRequest.MerchantId, req.ProductCode)
			transaction.ResponseCode = paymentResponse.ResponseCode
			transaction.ResponseMessage = paymentResponse.ResponseMessage
			transaction.ReversalCode = paymentResponse.ResponseCode
			transaction.ReversalMessage = paymentResponse.ResponseMessage
			transaction.ReversalDate = time.Now()
			s.transactionRepo.Update(&transaction)
			return &paymentResponse, nil
		}

		if checkStatus.ResponseCode == "2005500" {
			switch checkStatus.LatestTransactionStatus {
			case "00":
				transaction.ResponseCode = "00"
				transaction.ResponseMessage = "SUCCESS"
			case "05":
				transaction.ResponseCode = "2005505"
				transaction.ResponseMessage = "FAILED"
			case "02":
				transaction.ResponseCode = "2005502"
				transaction.ResponseMessage = "PENDING"
			default:
				transaction.ResponseCode = "99"
				transaction.ResponseMessage = "UNKNOWN STATUS"
			}

			paymentResponse = response.CreatePaymentResponse{
				TrxId:           req.TrxId,
				TrxConfirm:      partnerResponse.ReferenceNo,
				ResponseCode:    transaction.ResponseCode,
				ResponseMessage: transaction.ResponseMessage,
				PaidAt:          checkStatus.PaidTime,
			}
		} else {
			transaction.ResponseCode = partnerResponse.ResponseCode
			transaction.ResponseMessage = partnerResponse.ResponseMessage

			paymentResponse = response.CreatePaymentResponse{
				TrxId:           req.TrxId,
				TrxConfirm:      partnerResponse.ReferenceNo,
				ResponseCode:    transaction.ResponseCode,
				ResponseMessage: transaction.ResponseMessage,
				PaidAt:          partnerResponse.AdditionnalInfo.PaidTime,
			}
		}
	} else {

		transaction.ResponseCode = partnerResponse.ResponseCode
		transaction.ResponseMessage = partnerResponse.ResponseMessage
		paymentResponse = response.CreatePaymentResponse{
			TrxId:           req.TrxId,
			TrxConfirm:      partnerResponse.ReferenceNo,
			ResponseCode:    transaction.ResponseCode,
			ResponseMessage: transaction.ResponseMessage,
			PaidAt:          partnerResponse.AdditionnalInfo.PaidTime,
		}
	}

	err = s.transactionRepo.Update(&transaction)
	if err != nil {
		return nil, s.logAndWrapError(req.TrxId, "Error Update Transaction with RC:"+transaction.ResponseCode+" Response Message:"+transaction.ResponseMessage, "Update Transaction", err)
	}

	s.tracelogRepo.Log(req.TrxId,
		"Transaction Exit "+paymentResponse.ResponseCode+";"+paymentResponse.ResponseMessage,
		"Service EXIT",
	)

	return &paymentResponse, nil
}
