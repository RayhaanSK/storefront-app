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
	ListUserTypes(ctx context.Context) ([]domain.UserType, error)                                      // PASSED
	CreateUserAccount(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) // PASSED
	ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error)                                // PASSED
	GetUserAccount(ctx context.Context, id int64) (domain.UserAccount, error)                          // PASSED
	DeleteUserAccount(ctx context.Context, id int64) error                                             // PASSED
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

/*

// region User Type

func (s *Service) ListUserTypes(ctx context.Context) ([]domain.UserType, error) {
	return s.repo.ListUserTypes(ctx)
}

// endregion

// region User Account

func (s *Service) SignUpUserAccount(ctx context.Context, username string, passwordHash string, userTypeId int64) (domain.UserAccount, error) {
	created, err := domain.NewUserAccount(username, passwordHash, userTypeId)
	if err != nil {
		return created, err
	}

	return s.repo.CreateUserAccount(ctx, created)
}

func (s *Service) ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error) {
	return s.repo.ListUserAccounts(ctx)
}

func (s *Service) LogIntoUserAccount(ctx context.Context, username string, passwordHash string) (domain.UserAccount, error) {

}

func (s *Service) DeleteUserAccount()

// endregion

// region Product

func (s *Service) ListProductCatalogue()
func (s *Service) GetProductDetails()
func (s *Service) ListOwnProducts()
func (s *Service) CreateProduct()
func (s *Service) DeleteProduct()
func (s *Service) UpdateProduct()

// endregion

// region Cart Item

func (s *Service) ListCartItems()
func (s *Service) AddProductToCart()
func (s *Service) RemoveProductFromCart()
func (s *Service) CheckoutCart()

// endregion

*/
