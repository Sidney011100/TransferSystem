package main

import (
	"log"
	"os"
	db "transferSystem/database"
	"transferSystem/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	h := gin.Default()
	RegisterRoutes(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db.InitDatabase()
	defer db.CloseDatabase()

	log.Printf("Starting server on port %s", port)
	if err := h.Run(); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}

func RegisterRoutes(h *gin.Engine) {
	api := h.Group("/api/v1")

	accountGroup := api.Group("/accounts")
	accountGroup.POST("", handler.UserCreateAccount)
	accountGroup.GET("/:account_id", handler.UserGetAccount)

	transactionGroup := api.Group("/transactions")
	transactionGroup.POST("", handler.UserTransaction)
}
