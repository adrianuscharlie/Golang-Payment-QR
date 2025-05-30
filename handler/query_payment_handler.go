package handler

import (
	"net/http"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/services"

	"github.com/gin-gonic/gin"
)

type QueryPaymentHandler struct {
	queryPaymentService services.QueryPaymentService
	tracelogRepo        repository.TracelogRepository
}

func NewQueryPaymentHandler(s services.QueryPaymentService, t repository.TracelogRepository) *QueryPaymentHandler {
	return &QueryPaymentHandler{s, t}
}

func (h *QueryPaymentHandler) HandleCheckStatus(c *gin.Context) {
	var req request.CreateQueryPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request in Query Payment Check Status : "+req.TrxId, "HANDLER IN")
	resp, err := h.queryPaymentService.CheckStatusPayment(req)
	if err != nil {
		h.tracelogRepo.Log(req.TrxId, "ERROR Processing Request PurchaseQR : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request OUT Query Payment Check Status : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
	c.JSON(http.StatusOK, resp)
}
