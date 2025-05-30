package handler

import (
	"net/http"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService services.PaymentService
	tracelogRepo   repository.TracelogRepository
}

func NewPaymenthandler(s services.PaymentService, t repository.TracelogRepository) *PaymentHandler {
	return &PaymentHandler{paymentService: s, tracelogRepo: t}
}

func (h *PaymentHandler) HandleCPMPayment(c *gin.Context) {
	var req request.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request in PurchaseQR : "+req.ProductCode+";"+req.QrContent, "HANDLER IN")
	resp, err := h.paymentService.Payment(req)
	if err != nil {
		h.tracelogRepo.Log(req.TrxId, "ERROR Processing Request PurchaseQR : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request OUT PurchaseQR : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
	c.JSON(http.StatusOK, resp)
}
