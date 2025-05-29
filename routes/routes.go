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
	paymentService := services.NewPaymentService(productRepository, transactionRepository)
	paymentHandler := handler.NewPaymenthandler(paymentService)

	api := router.Group("/api")
	{
		api.POST("/qris/payment", paymentHandler.HandleCPMPayment)
	}
}
