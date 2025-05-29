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
	productRepo     repository.ProductConfigRepository
	transactionRepo repository.TransactionRepository
}

func NewPaymentService(productRepo repository.ProductConfigRepository, transactionRepo repository.TransactionRepository) PaymentService {
	return &paymentService{productRepo, transactionRepo}
}

func (s *paymentService) Payment(req request.CreatePaymentRequest) (*response.CreatePaymentResponse, error) {
	timeStamp := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	// Get Product Configuration
	productConfig, err := s.productRepo.GetConfig(req.ProductCode)
	if err != nil {
		return nil, err
	}
	urlConfig, err := s.productRepo.GetUrlConfig(req.ProductCode)
	if err != nil {
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
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", urlConfig.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	signature, err := utils.SignatureHeader(productConfig.ClientSecret, timeStamp)
	if err != nil {
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
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-TIMESTAMP", timeStamp)
	httpReq.Header.Set("X-PARTNER-ID", productConfig.ExtraParam1)
	httpReq.Header.Set("X-EXTERNAL-ID", productConfig.ExtraParam2)
	httpReq.Header.Set("CHANNEL-ID", productConfig.ExtraParam3)
	httpReq.Header.Set("Authorization", "Bearer your-token-here")
	httpReq.Header.Set("X-SIGNATURE", signature)

	// Send Request to Partner
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	partnerResponse := dummy.DummyPaymentResponse(paymentRequest)
	resJson, err := json.Marshal(partnerResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	transaction.Response = string(resJson)
	transaction.ResponseCode = partnerResponse.ResponseCode
	transaction.ResponseMessage = partnerResponse.ResponseMessage
	transaction.TrxConfirm = partnerResponse.ReferenceNo
	s.transactionRepo.Update(&transaction)

	paymentResponse := response.CreatePaymentResponse{
		TrxId:           req.TrxId,
		TrxConfirm:      partnerResponse.ReferenceNo,
		ResponseCode:    partnerResponse.ResponseCode,
		ResponseMessage: partnerResponse.ResponseMessage,
		PaidAt:          partnerResponse.AdditionnalInfo.PaidTime,
	}

	return &paymentResponse, nil
}
