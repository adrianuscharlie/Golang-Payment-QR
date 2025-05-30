package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/response"
	"qr-payment/utils"
	"time"
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

func (q queryPaymentService) CheckStatusPayment(trxId string, referenceNo string, merchantId string, productCode string) (*response.QueryPaymentResponse, error) {
	timeStamp := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	// Get Product Configuration
	productConfig, err := q.productRepo.GetConfig(productCode)
	if err != nil {
		return nil, err
	}
	qrStatus := request.QueryPaymentRequest{
		OriginalPartnerReferenceNo: trxId,
		OriginalReferenceNo:        referenceNo,
		MerchantId:                 merchantId,
		ServiceCode:                "60",
	}
	urlConfig, err := q.productRepo.GetUrlConfig(productCode, "STATUS")
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(qrStatus)
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
	q.tracelogRepo.Log(qrStatus.OriginalPartnerReferenceNo, "Sending HTTP POST : \n"+headerStr+"\n"+string(jsonData), "SEND HTTP POST")
	// Send Request to Partner
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		q.tracelogRepo.Log(qrStatus.OriginalPartnerReferenceNo, "Error Sending HTTP POST :"+err.Error(), "SEND HTTP POST")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		q.tracelogRepo.Log(qrStatus.OriginalPartnerReferenceNo, "Error :Parse response Body :"+err.Error(), "Parse Response")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var qrResponse response.QueryPaymentResponse
	err = json.Unmarshal(body, &qrResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &qrResponse, nil
}
