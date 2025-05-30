package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"qr-payment/dummy"
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
}

func NewPaymentService(productRepo repository.ProductConfigRepository, transactionRepo repository.TransactionRepository, tracelogRepo repository.TracelogRepository, queryPaymentService QueryPaymentService) PaymentService {
	return &paymentService{productRepo, transactionRepo, tracelogRepo, queryPaymentService}
}

func (s *paymentService) Payment(req request.CreatePaymentRequest) (*response.CreatePaymentResponse, error) {
	s.tracelogRepo.Log(req.TrxId, "Transaction Received with TrxId: "+req.TrxId, "Service IN")
	timeStamp := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	// Get Product Configuration
	productConfig, err := s.productRepo.GetConfig(req.ProductCode)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Error Get Product Config : "+err.Error(), "Product Config")
		return nil, err
	}
	urlConfig, err := s.productRepo.GetUrlConfig(req.ProductCode, "PAYMENT")
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Error Get URL for Product Code: "+req.ProductCode, "URL Config")
		return nil, err
	}
	// Create request requirements
	paymentRequest := request.PaymentRequest{
		PartnerReferenceNo: req.TrxId,
		QrContent:          req.QrContent,
		Amount: model.Amount{
			Value:    strconv.Itoa(int(req.Amount)),
			Currency: "IDR",
		},
		MerchantId:      "TP18",
		Title:           "QRIS CPM Payment via Dana",
		ExpiryTime:      timeStamp,
		ExternalStoreId: "M2351",
		ScannerInfo: request.ScannerInfo{
			DeviceId: "46252",
			DeviceIp: "172.24.281.24",
		},
		AdditionnalInfo: request.AdditionnalInfo{
			QrContentType: "BAR_CODE",
			Mcc:           "5732",
			ProductCode:   req.ProductCode,
			Envinfo: request.Envinfo{
				SourcePplatform:   "IPG",
				TerminalType:      "SYSTEM",
				OrderterminalType: "WEB",
			},
		},
	}

	jsonData, err := json.Marshal(paymentRequest)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Error Creating Request and Serialize into JSON ", "Create Request")
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", urlConfig.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	signature, err := utils.SignatureHeader(productConfig.ClientSecret, timeStamp)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Failed Creating Signature Header for Request", "Signature")
		return nil, fmt.Errorf("failed to create Signature request: %w", err)
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
		s.tracelogRepo.Log(req.TrxId, "Error Insert transaction into table transaction", "Signature")
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-TIMESTAMP", timeStamp)
	httpReq.Header.Set("X-PARTNER-ID", productConfig.ExtraParam1)
	httpReq.Header.Set("X-EXTERNAL-ID", productConfig.ExtraParam2)
	httpReq.Header.Set("CHANNEL-ID", productConfig.ExtraParam3)
	httpReq.Header.Set("Authorization", "Bearer your-token-here")
	httpReq.Header.Set("X-SIGNATURE", signature)

	var headerStr string
	for key, values := range httpReq.Header {
		for _, value := range values {
			headerStr += key + ":" + value + "\n"
		}
	}
	s.tracelogRepo.Log(req.TrxId, "Sending HTTP POST : \n"+headerStr+"\n"+string(jsonData), "SEND HTTP POST")
	// Send Request to Partner
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId, "Error Sending HTTP POST :"+err.Error(), "SEND HTTP POST")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		s.tracelogRepo.Log(req.TrxId, "Error :Parse response Body :"+err.Error(), "Parse Response")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	partnerResponse := dummy.DummyPaymentResponse(paymentRequest)
	resJson, err := json.Marshal(partnerResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	s.tracelogRepo.Log(req.TrxId, "Response String :"+string(resJson), "Response String")
	transaction.Response = string(resJson)
	transaction.ResponseCode = partnerResponse.ResponseCode
	transaction.ResponseMessage = partnerResponse.ResponseMessage
	transaction.TrxConfirm = partnerResponse.ReferenceNo
	var paymentResponse response.CreatePaymentResponse

	// Handle pending state that requires a query to determine final status
	if partnerResponse.ResponseCode == "2026000" {
		checkStatus, err := s.queryPaymentService.CheckStatusPayment(
			partnerResponse.PartnerReferenceNo,
			partnerResponse.ReferenceNo,
			paymentRequest.MerchantId,
			req.ProductCode,
		)
		if err != nil {
			return nil, err
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
			// If check status fails but we must still return something
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
		// Directly use CPM response if it's already final (00, 05, etc.)
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

	// Try to update the transaction
	err = s.transactionRepo.Update(&transaction)
	if err != nil {
		s.tracelogRepo.Log(req.TrxId,
			"Error Update Transaction with RC:"+transaction.ResponseCode+" Response Message:"+transaction.ResponseMessage+" Error: "+err.Error(),
			"Update Transaction",
		)
		return nil, err
	}

	// Final exit log
	s.tracelogRepo.Log(req.TrxId,
		"Transaction Exit "+paymentResponse.ResponseCode+";"+paymentResponse.ResponseMessage,
		"Service EXIT",
	)

	return &paymentResponse, nil
}
