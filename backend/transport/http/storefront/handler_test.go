package storefront

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	domain "storefrontapp/domain/storefront"
	storefrontservice "storefrontapp/services/storefront"
)

func TestListAllUserTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUserTypes := []domain.UserType{
		{ID: 1, Type: "testCustomer"},
		{ID: 2, Type: "testSeller"},
	}

	repo := &fakeRepository{
		listUserTypesFunc: func(ctx context.Context) ([]domain.UserType, error) {
			return testUserTypes, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/user_type/list", handler.ListAllUserTypes)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/user_type/list",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got []domain.UserType

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Type != "testCustomer" {
		t.Fatalf("expected type testCustomer, got %v", got[0].Type)
	}

	if got[1].ID != 2 {
		t.Fatalf("expected ID 2, got %d", got[1].ID)
	}

	if got[1].Type != "testSeller" {
		t.Fatalf("expected type testSeller, got %v", got[1].Type)
	}
}

func TestSignUpUserAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createUserAccountFunc: func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
			userAccount.ID = 1
			return userAccount, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	body := `{"username":"testUsername","passwordHash":"testPasswordHash","userTypeId":1}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/user_account/create",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var got domain.UserAccount

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != 1 {
		t.Fatalf("expected ID 1, got %d", got.ID)
	}

	if got.Username != "testUsername" {
		t.Fatalf("expected username testUsername, got %v", got.Username)
	}

	if got.PasswordHash != "testPasswordHash" {
		t.Fatalf("expected passwordHash testPasswordHash, got %v", got.PasswordHash)
	}

	if got.UserTypeID != 1 {
		t.Fatalf("expected user type id 1, got %v", got.UserTypeID)
	}
}

func TestListAllUserAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUserAccounts := []domain.UserAccount{
		{ID: 1, Username: "testUsername1", PasswordHash: "testPasswordHash1", UserTypeID: 1, UserType: "testUserType1"},
		{ID: 2, Username: "testUsername2", PasswordHash: "testPasswordHash2", UserTypeID: 2, UserType: "testUserType2"},
	}

	repo := &fakeRepository{
		listUserAccountsFunc: func(ctx context.Context) ([]domain.UserAccount, error) {
			return testUserAccounts, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/user_account/list", handler.ListAllUserAccounts)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/user_account/list",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got []domain.UserAccount

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Username != "testUsername1" {
		t.Fatalf("expected username testUsername1, got %v", got[0].Username)
	}

	if got[0].PasswordHash != "testPasswordHash1" {
		t.Fatalf("expected password hash testPasswordHash1, got %v", got[0].PasswordHash)
	}

	if got[0].UserTypeID != 1 {
		t.Fatalf("expected user type id 1, got %d", got[0].UserTypeID)
	}

	if got[0].UserType != "testUserType1" {
		t.Fatalf("expected user type testUserType1, got %v", got[0].UserType)
	}
}

func TestLogIntoUserAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUserAccount := domain.UserAccount{ID: 1, Username: "testUsername", PasswordHash: "testPasswordHash", UserTypeID: 1, UserType: "testUserType"}

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return testUserAccount, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)

	body := `{"username":"testUsername","passwordHash":"testPasswordHash"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/user_account/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got LogIntoAccountResponse

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.SuccessStatus != true {
		t.Fatalf("expected success status true, got %v", got.SuccessStatus)
	}

	if got.UserAccount.ID != 1 {
		t.Fatalf("expected id 1, got %v", got.UserAccount.ID)
	}

	if got.UserAccount.Username != "testUsername" {
		t.Fatalf("expected username testUsername, got %v", got.UserAccount.Username)
	}

	if got.UserAccount.UserTypeID != 1 {
		t.Fatalf("expected user type id 1, got %v", got.UserAccount.UserTypeID)
	}

	if got.UserAccount.UserType != "testUserType" {
		t.Fatalf("expected user type testUserType, got %v", got.UserAccount.UserType)
	}
}

func TestDeleteUserAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUserAccounts := []domain.UserAccount{
		{ID: 1, Username: "testUsername1", PasswordHash: "testPasswordHash1", UserTypeID: 1, UserType: "testUserType1"},
		{ID: 2, Username: "testUsername2", PasswordHash: "testPasswordHash2", UserTypeID: 2, UserType: "testUserType2"},
	}

	repo := &fakeRepository{
		deleteUserAccountFunc: func(ctx context.Context, id int64) error {
			for i := range testUserAccounts {
				if testUserAccounts[i].ID == id {
					testUserAccounts = append(testUserAccounts[:i], testUserAccounts[i+1:]...)
					return nil
				}
			}
			return nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/api/user_account/delete/:id", handler.DeleteUserAccount)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/user_account/delete/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if len(testUserAccounts) != 1 {
		t.Fatalf("expected number of accounts to be 1, got %v", len(testUserAccounts))
	}
}

// region Helper

var testNow = time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)

func fixedNow() time.Time {
	return testNow
}

type fakeRepository struct {
	listUserTypesFunc      func(ctx context.Context) ([]domain.UserType, error)
	createUserAccountFunc  func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error)
	listUserAccountsFunc   func(ctx context.Context) ([]domain.UserAccount, error)
	getUserAccountFunc     func(ctx context.Context, username string) (domain.UserAccount, error)
	deleteUserAccountFunc  func(ctx context.Context, id int64) error
	createProductFunc      func(ctx context.Context, product domain.Product) (domain.Product, error)
	listUnsoldProductsFunc func(ctx context.Context) ([]domain.Product, error)
	listOwnProductsFunc    func(ctx context.Context, sellerId int64) ([]domain.Product, error)
	listProductsInCartFunc func(ctx context.Context, customerId int64) ([]domain.Product, error)
	getProductFunc         func(ctx context.Context, id int64) (domain.Product, error)
	updateProductFunc      func(ctx context.Context, product domain.Product) (domain.Product, error)
	deleteProductFunc      func(ctx context.Context, id int64) error
	createCartItemFunc     func(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	deleteCartItemFunc     func(ctx context.Context, customerId int64, productId int64) error
	deleteCartItemsFunc    func(ctx context.Context, customerId int64) error
}

func (f *fakeRepository) ListUserTypes(ctx context.Context) ([]domain.UserType, error) {
	return f.listUserTypesFunc(ctx)
}

func (f *fakeRepository) CreateUserAccount(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
	return f.createUserAccountFunc(ctx, userAccount)
}

func (f *fakeRepository) ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error) {
	return f.listUserAccountsFunc(ctx)
}

func (f *fakeRepository) GetUserAccount(ctx context.Context, username string) (domain.UserAccount, error) {
	return f.getUserAccountFunc(ctx, username)
}

func (f *fakeRepository) DeleteUserAccount(ctx context.Context, id int64) error {
	return f.deleteUserAccountFunc(ctx, id)
}

func (f *fakeRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return f.createProductFunc(ctx, product)
}

func (f *fakeRepository) ListUnsoldProducts(ctx context.Context) ([]domain.Product, error) {
	return f.listUnsoldProductsFunc(ctx)
}

func (f *fakeRepository) ListOwnProducts(ctx context.Context, sellerId int64) ([]domain.Product, error) {
	return f.listOwnProductsFunc(ctx, sellerId)
}

func (f *fakeRepository) ListProductsInCart(ctx context.Context, customerId int64) ([]domain.Product, error) {
	return f.listProductsInCartFunc(ctx, customerId)
}

func (f *fakeRepository) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	return f.getProductFunc(ctx, id)
}

func (f *fakeRepository) UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return f.updateProductFunc(ctx, product)
}

func (f *fakeRepository) DeleteProduct(ctx context.Context, id int64) error {
	return f.deleteProductFunc(ctx, id)
}

func (f *fakeRepository) CreateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
	return f.createCartItemFunc(ctx, cartItem)
}

func (f *fakeRepository) DeleteCartItem(ctx context.Context, customerId int64, productId int64) error {
	return f.deleteCartItemFunc(ctx, customerId, productId)
}

func (f *fakeRepository) DeleteCartItems(ctx context.Context, customerId int64) error {
	return f.deleteCartItemsFunc(ctx, customerId)
}

// endregion
