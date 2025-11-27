package service

import (
	"context"
	"errors"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/dto"
	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
)

type IWebhookService interface {
	ReceiveInvoice(ctx context.Context, req *dto.XenditInvoiceRequest) error
}

type webhookService struct {
	orderRepository repository.IOrderRepository
}

func NewWebhookService(orderRepository repository.IOrderRepository) IWebhookService {
	return &webhookService{
		orderRepository: orderRepository,
	}
}

func (ws *webhookService) ReceiveInvoice(ctx context.Context, req *dto.XenditInvoiceRequest) error {

	orderEntity, err := ws.orderRepository.GetOrderById(ctx, req.ExternalID)
	if err != nil {
		return err
	}
	if orderEntity == nil {
		return errors.New("order not found")
	}

	now := time.Now()
	updatedBy := "System"
	orderEntity.OrderStatusCode = entity.OrderStatusCodePaid
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &updatedBy
	orderEntity.XenditPaidAt = &now
	orderEntity.XenditPaymentChannel = &req.PaymentChannel
	orderEntity.XenditPaymentMethod = &req.PaymentMethod

	err = ws.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return err
	}
	return nil
}
