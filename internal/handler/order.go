package handler

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/order"
)

type orderHandler struct {
	// NOTE: this should be embedded by value instead of pointer to avoid a nil
	// pointer dereference when methods are called.
	order.UnimplementedOrderServiceServer

	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &order.CreateOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (oh *orderHandler) ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &order.ListOrderAdminResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.ListOrderAdmin(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (oh *orderHandler) ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &order.ListOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (oh *orderHandler) DetailOrder(ctx context.Context, request *order.DetailOrderRequest) (*order.DetailOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.DetailOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.DetailOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (oh *orderHandler) UpdateOrderStatus(ctx context.Context, request *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.UpdateOrderStatusResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.UpdateOrderStatus(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
