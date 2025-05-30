package handler

import (
	"net/http"
	"qr-payment/repository"
	"qr-payment/request"
	"qr-payment/services"

	"github.com/gin-gonic/gin"
)

type CancelOrderHandler struct {
	cancelOrderService services.CancelOrderService
	tracelogRepo       repository.TracelogRepository
}

func NewCancelOrderHandler(c services.CancelOrderService, t repository.TracelogRepository) *CancelOrderHandler {
	return &CancelOrderHandler{c, t}
}

func (h *CancelOrderHandler) CancelOrder(c *gin.Context) {
	var req request.CreateCancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request in Cancel Order : "+req.TrxId, "HANDLER IN")
	resp, err := h.cancelOrderService.CancelOrder(req)
	if err != nil {
		h.tracelogRepo.Log(req.TrxId, "ERROR Processing Cancel Order Request : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.tracelogRepo.Log(req.TrxId, "Request OUT Cancel Order : "+resp.ResponseCode+";"+resp.ResponseMessage, "HANDLER OUT")
	c.JSON(http.StatusOK, resp)
}
