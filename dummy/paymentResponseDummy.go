package dummy

import (
	"qr-payment/request"
	"qr-payment/response"
	"time"
)

func DummyPaymentResponse(r request.PaymentRequest) response.PaymentResponse {
	paymentResponse := response.PaymentResponse{
		PartnerReferenceNo: r.PartnerReferenceNo,
		ResponseCode:       "00",
		ResponseMessage:    "Success",
		ReferenceNo:        r.PartnerReferenceNo + "-XX",
		TransactionDate:    time.Now().Format(time.RFC3339),
		AdditionnalInfo: response.AdditionnalInfo{
			Amount:   r.Amount,
			PaidTime: time.Now().Add(1 * time.Second).Format(time.RFC3339),
		},
	}
	return paymentResponse
}
