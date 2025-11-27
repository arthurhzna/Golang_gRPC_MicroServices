package main

import (
	"context"
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/arthurhzna/Golang_gRPC/internal/handler"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func handleGetFileName(c *fiber.Ctx) error {
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", "product", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("Not Found")
		}
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	ext := path.Ext(filePath)
	mineType := mime.TypeByExtension(ext)

	c.Set("Content-Type", mineType)

	return c.SendStream(file)

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	app := fiber.New()
	app.Use(cors.New())

	db := database.ConnectDb(ctx, os.Getenv("DB_URL"))

	orderRepository := repository.NewOrderRepository(db)
	webhookService := service.NewWebhookService(orderRepository)
	webhookHandler := handler.NewWebhookHandler(webhookService)

	app.Get("/storage/product/:filename", handleGetFileName)
	app.Post("/product/upload", handler.UploadHandler)
	app.Post("webhook/xendit/invoice", webhookHandler.ReceiveInvoice)

	app.Listen(":3000")
}
