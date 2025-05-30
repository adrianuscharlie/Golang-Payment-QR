package routes

import (
	"qr-payment/handler"
	"qr-payment/repository"
	"qr-payment/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	productRepository := repository.NewProductConfigRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	tracelogRepository := repository.NewTracelogRepository(db)
	queryService := services.NewQueryPaymentService(transactionRepository, tracelogRepository, productRepository)
	cancelOrderService := services.NewCancelOrderService(transactionRepository, tracelogRepository, productRepository)
	paymentService := services.NewPaymentService(productRepository, transactionRepository, tracelogRepository, queryService, cancelOrderService)
	paymentHandler := handler.NewPaymenthandler(paymentService, tracelogRepository)

	api := router.Group("/api")
	{
		api.POST("/qris/payment", paymentHandler.HandleCPMPayment)
	}
}
