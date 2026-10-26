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

var ErrIncorrectUsernameOrPassword = errors.New("Incorrect username or password.")

type Repository interface {
	ListUserTypes(ctx context.Context) ([]domain.UserType, error)                                      // PASSED
	CreateUserAccount(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) // PASSED
	ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error)                                // PASSED
	GetUserAccount(ctx context.Context, username string) (domain.UserAccount, error)                   // PASSED
	DeleteUserAccount(ctx context.Context, id int64) error                                             // PASSED
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)                 // PASSED
	ListUnsoldProducts(ctx context.Context) ([]domain.Product, error)                                  // PASSED
	ListOwnProducts(ctx context.Context, sellerId int64) ([]domain.Product, error)                     // PASSED
	ListProductsInCart(ctx context.Context, customerId int64) ([]domain.Product, error)                // PASSED
	GetProduct(ctx context.Context, id int64) (domain.Product, error)                                  // PASSED
	UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error)                 // PASSED
	DeleteProduct(ctx context.Context, id int64) error                                                 // PASSED
	CreateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)             // PASSED
	DeleteCartItem(ctx context.Context, customerId int64, productId int64) error                       // PASSED
	DeleteCartItems(ctx context.Context, customerId int64) error                                       // PASSED
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

// region User Type

func (s *Service) ListAllUserTypes(ctx context.Context) ([]domain.UserType, error) {
	return s.repo.ListUserTypes(ctx)
}

// endregion

// region User Account

func (s *Service) SignUpUserAccount(ctx context.Context, username string, passwordHash string, userTypeId int64) (domain.UserAccount, error) {
	userAccount, err := domain.NewUserAccount(username, passwordHash, userTypeId)
	if err != nil {
		return domain.UserAccount{}, err
	}

	return s.repo.CreateUserAccount(ctx, userAccount)
}

func (s *Service) ListAllUserAccounts(ctx context.Context) ([]domain.UserAccount, error) {
	return s.repo.ListUserAccounts(ctx)
}

// Returns true if login is successful, returns false otherwise
func (s *Service) LogIntoUserAccount(ctx context.Context, username string, passwordHash string) (domain.UserAccount, bool, error) {
	userAccount, err := s.repo.GetUserAccount(ctx, username)
	if err != nil {
		return domain.UserAccount{}, false, err
	}

	if userAccount.PasswordHash != passwordHash {
		return domain.UserAccount{}, false, ErrIncorrectUsernameOrPassword
	}

	return userAccount, true, nil
}

func (s *Service) DeleteUserAccount(ctx context.Context, id int64) error {
	return s.repo.DeleteUserAccount(ctx, id)
}

// endregion

// region Product

func (s *Service) ListProductCatalogue(ctx context.Context) ([]domain.Product, error) {
	return s.repo.ListUnsoldProducts(ctx)
}

func (s *Service) GetSingleProductDetails(ctx context.Context, id int64) (domain.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *Service) ListAllSellerProducts(ctx context.Context, sellerId int64) ([]domain.Product, error) {
	return s.repo.ListOwnProducts(ctx, sellerId)
}

func (s *Service) CreateProduct(ctx context.Context, name string, price float64, description string, sellerId int64, categoryId int64) (domain.Product, error) {
	product, err := domain.NewProduct(name, price, description, s.now(), sellerId, categoryId)
	if err != nil {
		return domain.Product{}, err
	}

	return s.repo.CreateProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *Service) UpdateProduct(ctx context.Context, id int64, name string, price float64, description string, categoryId int64) (domain.Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	err = product.UpdateProduct(name, price, description, s.now(), categoryId)
	if err != nil {
		return domain.Product{}, err
	}

	return s.repo.UpdateProduct(ctx, product)
}

// endregion

// region Cart Item

func (s *Service) ListAllCustomerCartItems(ctx context.Context, customerId int64) ([]domain.Product, error) {
	return s.repo.ListProductsInCart(ctx, customerId)
}

func (s *Service) AddProductToCustomerCart(ctx context.Context, customerId int64, productId int64) (domain.CartItem, error) {
	cartItem, err := domain.NewCartItem(customerId, productId)
	if err != nil {
		return domain.CartItem{}, err
	}

	return s.repo.CreateCartItem(ctx, cartItem)
}

func (s *Service) RemoveProductFromCustomerCart(ctx context.Context, customerId int64, productId int64) error {
	return s.repo.DeleteCartItem(ctx, customerId, productId)
}

func (s *Service) CheckoutCustomerCart(ctx context.Context, customerId int64) error {
	return s.repo.DeleteCartItems(ctx, customerId)
}

// endregion
