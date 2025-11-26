package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/cart"
	"github.com/google/uuid"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
	UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}

func (cs *cartService) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	productEntity, err := cs.productRepository.GetProductById(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	cartEntity, err := cs.cartRepository.GetCartByProductAndUserId(ctx, req.ProductId, claims.Subject)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &claims.Subject

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			return nil, err
		}

		return &cart.AddProductToCartResponse{
			Base: utils.SuccessResponse("Product added to cart successfully"),
			Id:   cartEntity.Id,
		}, nil
	}

	newCartEntity := entity.Cart{
		Id:        uuid.NewString(),
		UserId:    claims.Subject,
		ProductId: req.ProductId,
		Quantity:  cartEntity.Quantity + 1,
		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}

	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse("Product added to cart successfully"),
		Id:   productEntity.Id,
	}, nil
}

func (cs *cartService) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	carts, err := cs.cartRepository.GetListCart(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	if carts == nil {
		return &cart.ListCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	var items []*cart.ListCartResponseItem = make([]*cart.ListCartResponseItem, 0)
	for _, cartEntity := range carts {
		item := cart.ListCartResponseItem{
			CartId:          cartEntity.Id,
			ProductId:       cartEntity.ProductId,
			ProductName:     cartEntity.Product.Name,
			ProductImageUrl: fmt.Sprintf("%s/storage/product/%s", os.Getenv("STORAGE_SERVICE_URL"), cartEntity.Product.ImageFileName),
			ProductPrice:    cartEntity.Product.Price,
			Quantity:        int64(cartEntity.Quantity),
		}
		items = append(items, &item)
	}

	return &cart.ListCartResponse{
		Base:  utils.SuccessResponse("Cart list retrieved successfully"),
		Items: items,
	}, nil

}

func (cs *cartService) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cartEntity, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.DeleteCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity.UserId != claims.Subject {
		return nil, utils.UnaunthorizedResponse()
	}

	err = cs.cartRepository.DeleteCart(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	return &cart.DeleteCartResponse{
		Base: utils.SuccessResponse("Cart deleted successfully"),
	}, nil
}

func (cs *cartService) UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cartEntity, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.UpdateCartQuantityResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity.UserId != claims.Subject {
		return nil, utils.UnaunthorizedResponse()
	}

	if req.NewQuantity <= 0 {
		cs.cartRepository.DeleteCart(ctx, req.CartId)
		if err != nil {
			return nil, err
		}
		return &cart.UpdateCartQuantityResponse{
			Base: utils.SuccessResponse("Cart deleted successfully"),
		}, nil
	}

	now := time.Now()
	cartEntity.Quantity = int(req.NewQuantity)
	cartEntity.UpdatedAt = &now
	cartEntity.UpdatedBy = &claims.FullName

	err = cs.cartRepository.UpdateCart(ctx, cartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.UpdateCartQuantityResponse{
		Base: utils.SuccessResponse("Cart quantity updated successfully"),
	}, nil

}
