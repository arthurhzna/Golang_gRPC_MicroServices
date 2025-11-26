package handler

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/cart"
)

type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

type ICartService interface {
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
	UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error)
}

func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

func (ch *cartHandler) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &cart.AddProductToCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.AddProductToCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &cart.ListCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {

	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &cart.DeleteCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ch.cartService.DeleteCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error) {

	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &cart.UpdateCartQuantityResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.UpdateCartQuantity(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil

}
