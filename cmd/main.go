package main

import (
	"log"
	"qr-payment/config"
	"qr-payment/model"
	"qr-payment/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&model.ProductConfig{},
		&model.UrlConfig{},
		&model.Transaction{})
	r := gin.Default()
	routes.RegisterRoutes(r, config.DB)
	r.Run(":8080")
}
