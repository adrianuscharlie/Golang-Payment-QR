package handler

import (
	"net/http"
	"qr-payment/request"
	"qr-payment/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymenthandler(s services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: s}
}

func (h *PaymentHandler) HandleCPMPayment(c *gin.Context) {
	var req request.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.paymentService.Payment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
