package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
)

type ICartRepository interface {
	GetCartByProductAndUserId(ctx context.Context, productId string, userId string) (*entity.Cart, error)
	CreateNewCart(ctx context.Context, cart *entity.Cart) error
	UpdateCart(ctx context.Context, cart *entity.Cart) error
	GetListCart(ctx context.Context, userId string) ([]*entity.Cart, error)
	GetCartById(ctx context.Context, cartId string) (*entity.Cart, error)
	DeleteCart(ctx context.Context, cartId string) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) ICartRepository {
	return &cartRepository{db: db}
}

func (cr *cartRepository) GetCartByProductAndUserId(ctx context.Context, productId string, userId string) (*entity.Cart, error) {

	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE product_id = $1 AND user_id = $2",
		productId,
		userId,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var cart entity.Cart
	err := row.Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.UserId,
		&cart.Quantity,
		&cart.CreatedAt,
		&cart.CreatedBy,
		&cart.UpdatedAt,
		&cart.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (cr *cartRepository) CreateNewCart(ctx context.Context, cart *entity.Cart) error {
	_, err := cr.db.ExecContext(
		ctx,
		`INSERT INTO "user_cart" (id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		cart.Id,
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.CreatedAt,
		cart.CreatedBy,
		cart.UpdatedAt,
		cart.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) UpdateCart(ctx context.Context, cart *entity.Cart) error {
	_, err := cr.db.ExecContext(
		ctx,
		`UPDATE "user_cart" SET product_id = $1, user_id = $2, quantity = $3, updated_at = $4, updated_by = $5 WHERE id = $6`,
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.UpdatedAt,
		cart.UpdatedBy,
		cart.Id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) GetListCart(ctx context.Context, userId string) ([]*entity.Cart, error) {

	rows, err := cr.db.QueryContext(
		ctx,
		"SELECT uc.id, uc.product_id, uc.user_id, uc.quantity, uc.created_at, uc.created_by, uc.updated_at, uc.updated_by, p.id, p.name, p.image_file_name, p.price FROM user_cart uc JOIN product p ON uc.product_id = p.id WHERE uc.user_id = $1 AND p.is_deleted = false",
		userId,
	)
	if err != nil {
		return nil, err
	}

	var carts []*entity.Cart = make([]*entity.Cart, 0)
	for rows.Next() {
		var cart entity.Cart // ini struct non pointer
		/*
			Note:
			- 'cart' dibuat ulang setiap iterasi.
			- Karena &cart digunakan, variabel cart otomatis dialokasikan di heap (escape).
			- Setiap append(&cart) menyimpan alamat heap yang berbeda → data sebelumnya aman.
			- Updating cart di iterasi berikutnya tidak mengubah data yang sudah diappend.
			- GC akan menghapus objek jika tidak ada lagi pointer yang mereferensikannya.
		*/

		/*
			catatan:
			- Menggunakan `var cart entity.Cart` aman karena setiap iterasi struct baru dibuat,
			lalu dialokasikan ke heap saat &cart digunakan.
			- Nilai yang sudah di-append tidak akan tertimpa oleh iterasi berikutnya.
			- Alternatifnya: `cart := &entity.Cart{}` (langsung pointer), hasilnya sama saja --->> ini sama aja membuat struct baru , dan mereferencekan ke pointer, mendung langsung buat var cart entity.Cart
			tetapi lebih eksplisit bahwa objek dibuat di heap.

			Rekomendasi: tetap gunakan `var cart entity.Cart` karena lebih ringkas dan idiomatic.
		*/

		cart.Product = &entity.Product{} // ini pointer karena product adalah reference dari product table
		err = rows.Scan(
			&cart.Id,
			&cart.ProductId,
			&cart.UserId,
			&cart.Quantity,
			&cart.CreatedAt,
			&cart.CreatedBy,
			&cart.UpdatedAt,
			&cart.UpdatedBy,
			&cart.Product.Id,
			&cart.Product.Name,
			&cart.Product.ImageFileName,
			&cart.Product.Price,
		)
		if err != nil {
			return nil, err
		}
		carts = append(carts, &cart)
	}
	return carts, nil

}

func (cr *cartRepository) GetCartById(ctx context.Context, cartId string) (*entity.Cart, error) {
	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE id = $1",
		cartId,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var cart entity.Cart
	err := row.Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.UserId,
		&cart.Quantity,
		&cart.CreatedAt,
		&cart.CreatedBy,
		&cart.UpdatedAt,
		&cart.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil

}

func (cr *cartRepository) DeleteCart(ctx context.Context, cartId string) error {
	_, err := cr.db.ExecContext(
		ctx,
		"DELETE FROM user_cart WHERE id = $1",
		cartId,
	)
	if err != nil {
		return err
	}
	return nil
}
