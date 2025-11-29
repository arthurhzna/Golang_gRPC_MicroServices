package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/product"
	"github.com/google/uuid"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error)
	DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error)
	EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error)
	DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
	ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error)
	HighlightProduct(ctx context.Context, req *product.HighlightProductRequest) (*product.HighlightProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join("storage", "product", req.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("image file not found"),
			}, nil
		}
		return nil, err
	}

	NewProduct := entity.Product{
		Id:            uuid.New().String(),
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageFileName: req.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}

	err = ps.productRepository.CreateNewProduct(ctx, &NewProduct)
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product created successfully"),
		Id:   NewProduct.Id,
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {

	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	return &product.DetailProductResponse{
		Base:        utils.SuccessResponse("Product detail retrieved successfully"),
		Id:          productEntity.Id,
		Name:        productEntity.Name,
		Description: productEntity.Description,
		Price:       productEntity.Price,
		ImageUrl:    fmt.Sprintf("%s/storage/product/%s", os.Getenv("STORAGE_SERVICE_URL"), productEntity.ImageFileName),
	}, nil
}

func (ps *productService) EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	if productEntity.ImageFileName != req.ImageFileName {
		imagePath := filepath.Join("storage", "product", req.ImageFileName)
		_, err = os.Stat(imagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &product.EditProductResponse{
					Base: utils.BadRequestResponse("Image file not found"),
				}, nil
			}
			return nil, err
		}
		oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}

	}
	/*
		product := &entity.Product{}
		Yang terjadi:
		Go membuat struct baru di memory dengan semua zero values
		Go mengembalikan pointer (alamat memory) ke struct tersebut

		Memory Address: 0x1000
		┌─────────────────────────────────┐
		│  entity.Product (struct)        │
		│  Id: ""                         │
		│  Name: ""                       │
		│  Description: ""                │
		│  Price: 0.0                     │
		│  ImageFileName: ""              │
		│  CreatedAt: 0001-01-01 UTC      │
		│  CreatedBy: ""                  │
		│  UpdatedAt: 0001-01-01 UTC      │
		│  UpdatedBy: nil                 │
		│  DeletedAt: 0001-01-01 UTC      │
		│  DeletedBy: nil                 │
		│  IsDeleted: false               │
		└─────────────────────────────────┘
				↑
				│
		product = 0x1000 (pointer ke struct)

		--------------------------------------------------------

		var product *entity.Product

		Yang terjadi:
		product adalah pointer bertipe *entity.Product
		Nilainya nil (tidak menunjuk ke struct apapun)
		TIDAK ada struct di memory

		product = nil (tidak menunjuk ke mana-mana)
			│
			└──> ❌ TIDAK ADA STRUCT!

		product = nil
		Tipe: *entity.Product (pointer)
		TIDAK ada struct di memory
		TIDAK bisa digunakan (akan panic jika akses field)

		-----------------------------------------------------------

		newProduct := entity.Product{}

		Hasil:
		newProduct bertipe entity.Product (struct)
		Ini adalah data asli (bukan pointer)

		Memory Address: 0x1000
		┌─────────────────────────────────┐
		│  newProduct (struct langsung)   │
		│  Id: "123"                      │
		│  Name: "Product A"              │
		│  Description: "..."             │
		│  Price: 100.0                   │
		│  ...                            │
		└─────────────────────────────────┘

		newProduct = struct langsung (tipe: entity.Product)

		-----------------------------------------------------------
		func UpdateProduct(product entity.Product) error {
			// Function menerima struct (bukan pointer)
		}

		// Panggil dengan:
		newProduct := entity.Product{...}
		UpdateProduct(newProduct)  // ✅ OK

		func UpdateProduct(product *entity.Product) error {
			// Function menerima pointer
		}

		// Panggil dengan:
		newProduct := &entity.Product{...}
		UpdateProduct(newProduct)  // ✅ OK
	*/

	newProduct := entity.Product{
		Id:            productEntity.Id,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageFileName: req.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = ps.productRepository.EditProduct(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Product detail retrieved successfully"),
		Id:   req.Id,
	}, nil
}

func (ps *productService) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(claims.Role)

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnaunthorizedResponse()
	}

	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DeleteProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	err = ps.productRepository.DeleteProduct(ctx, req.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	// Don't remove image files from storage, because if a user has already ordered a product, the image files could be lost.

	//remove image file from storage
	// imagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
	// err = os.Remove(imagePath)
	// if err != nil {
	// 	return nil, err
	// }

	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Delete product successfully"),
	}, nil

}

func (ps *productService) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	products, paginationResponse, err := ps.productRepository.GetProductsByPagination(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductResponseItem = make([]*product.ListProductResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.ListProductResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/storage/product/%s", os.Getenv("STORAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}
	return &product.ListProductResponse{
		Base:       utils.SuccessResponse("List product successfully"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil
}

func (ps *productService) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnaunthorizedResponse()
	}

	products, paginationResponse, err := ps.productRepository.GetProductsByPaginationAdmin(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductAdminResponseItem = make([]*product.ListProductAdminResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.ListProductAdminResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/storage/product/%s", os.Getenv("STORAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}
	return &product.ListProductAdminResponse{
		Base:       utils.SuccessResponse("List product successfully"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil
}

func (ps *productService) HighlightProduct(ctx context.Context, req *product.HighlightProductRequest) (*product.HighlightProductResponse, error) {

	products, err := ps.productRepository.GetProductsHighlight(ctx)
	if err != nil {
		return nil, err
	}

	var data []*product.HighlightProductResponseItem = make([]*product.HighlightProductResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.HighlightProductResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/storage/product/%s", os.Getenv("STORAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}
	return &product.HighlightProductResponse{
		Base: utils.SuccessResponse("get list product highlight successfully"),
		Data: data,
	}, nil
}
