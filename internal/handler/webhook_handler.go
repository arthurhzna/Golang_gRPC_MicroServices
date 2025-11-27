package handler

import (
	"log"
	"net/http"

	"github.com/arthurhzna/Golang_gRPC/internal/dto"
	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
	webhookService service.IWebhookService
}

func NewWebhookHandler(webhookService service.IWebhookService) *webhookHandler {
	return &webhookHandler{
		webhookService: webhookService,
	}
}

func (wh *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	var request dto.XenditInvoiceRequest

	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}

	err = wh.webhookService.ReceiveInvoice(c.Context(), &request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}
