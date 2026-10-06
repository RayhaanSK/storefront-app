package storefront

import (
	"context"
	"errors"
	"time"

	domain "storefrontapp/domain/storefront"
)

var ErrUserAccountNotFound = errors.New("User account not found.")
var ErrProductNotFound = errors.New("Product not found.")
var ErrCartItemNotFound = errors.New("Cart item not found")

type Repository interface {
	CreateUserAccount(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error)
	ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error)
	GetUserAccount(ctx context.Context, id int64) (domain.UserAccount, error)
	DeleteUserAccount(ctx context.Context, id int64) error
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	ListUnsoldProducts(ctx context.Context) ([]domain.Product, error)
	ListOwnProducts(ctx context.Context, sellerId int64) ([]domain.Product, error)
	ListProductsInCart(ctx context.Context, customerId int64) ([]domain.Product, error)
	GetProduct(ctx context.Context, id int64) (domain.Product, error)
	UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	CreateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	DeleteCartItem(ctx context.Context, customerId int64, productId int64) error
	DeleteCartItems(ctx context.Context, customerId int64) error
}

type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService constructs the storefront service
func NewService(repo Repository, now func() time.Time) *Service {
	return &Service{
		repo: repo,
		now:  now,
	}
}
