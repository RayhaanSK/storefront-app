package storefront

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	domain "storefrontapp/domain/storefront"
	storefrontservice "storefrontapp/services/storefront"
)

// region User Type

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

	recorder := serve(router, http.MethodGet, "/api/user_type/list", "")

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

func TestListAllUserTypesInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		listUserTypesFunc: func(ctx context.Context) ([]domain.UserType, error) {
			return nil, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/user_type/list", handler.ListAllUserTypes)

	recorder := serve(router, http.MethodGet, "/api/user_type/list", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

// endregion

// region User Account

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

	recorder := serve(router, http.MethodPost, "/api/user_account/create", body)

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

func TestSignUpUserAccountInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	// userTypeId sent as a string, which cannot bind to int64
	body := `{"username":"testUsername","passwordHash":"testPasswordHash","userTypeId":"1"}`

	recorder := serve(router, http.MethodPost, "/api/user_account/create", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestSignUpUserAccountMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	recorder := serve(router, http.MethodPost, "/api/user_account/create", `{"username":`)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestSignUpUserAccountBlankUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	body := `{"username":"   ","passwordHash":"testPasswordHash","userTypeId":1}`

	recorder := serve(router, http.MethodPost, "/api/user_account/create", body)

	assertError(t, recorder, http.StatusBadRequest, "title is required")
}

func TestSignUpUserAccountBlankPasswordHash(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	body := `{"username":"testUsername","passwordHash":"   ","userTypeId":1}`

	recorder := serve(router, http.MethodPost, "/api/user_account/create", body)

	assertError(t, recorder, http.StatusBadRequest, "title is required")
}

func TestSignUpUserAccountInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createUserAccountFunc: func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
			return domain.UserAccount{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/create", handler.SignUpUserAccount)

	body := `{"username":"testUsername","passwordHash":"testPasswordHash","userTypeId":1}`

	recorder := serve(router, http.MethodPost, "/api/user_account/create", body)

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
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

	recorder := serve(router, http.MethodGet, "/api/user_account/list", "")

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

func TestListAllUserAccountsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		listUserAccountsFunc: func(ctx context.Context) ([]domain.UserAccount, error) {
			return nil, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/user_account/list", handler.ListAllUserAccounts)

	recorder := serve(router, http.MethodGet, "/api/user_account/list", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
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

	recorder := serve(router, http.MethodPost, "/api/user_account/login", body)

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

func TestLogIntoUserAccountInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)

	// username sent as a number, which cannot bind to string
	body := `{"username":123,"passwordHash":"testPasswordHash"}`

	recorder := serve(router, http.MethodPost, "/api/user_account/login", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestLogIntoUserAccountNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return domain.UserAccount{}, storefrontservice.ErrUserAccountNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)

	body := `{"username":"testUsername","passwordHash":"testPasswordHash"}`

	recorder := serve(router, http.MethodPost, "/api/user_account/login", body)

	assertError(t, recorder, http.StatusNotFound, "user account not found")
}

func TestLogIntoUserAccountIncorrectPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return domain.UserAccount{ID: 1, Username: "testUsername", PasswordHash: "correctHash", UserTypeID: 1}, nil
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)

	body := `{"username":"testUsername","passwordHash":"wrongHash"}`

	recorder := serve(router, http.MethodPost, "/api/user_account/login", body)

	assertError(t, recorder, http.StatusUnauthorized, "incorrect username or password")
}

func TestLogIntoUserAccountInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return domain.UserAccount{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/user_account/login", handler.LogIntoUserAccount)

	body := `{"username":"testUsername","passwordHash":"testPasswordHash"}`

	recorder := serve(router, http.MethodPost, "/api/user_account/login", body)

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
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

	recorder := serve(router, http.MethodDelete, "/api/user_account/delete/1", "")

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if len(testUserAccounts) != 1 {
		t.Fatalf("expected number of accounts to be 1, got %v", len(testUserAccounts))
	}
}

func TestDeleteUserAccountInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.DELETE("/api/user_account/delete/:id", handler.DeleteUserAccount)

	recorder := serve(router, http.MethodDelete, "/api/user_account/delete/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid user account id")
}

func TestDeleteUserAccountNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteUserAccountFunc: func(ctx context.Context, id int64) error {
			return storefrontservice.ErrUserAccountNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/user_account/delete/:id", handler.DeleteUserAccount)

	recorder := serve(router, http.MethodDelete, "/api/user_account/delete/1", "")

	assertError(t, recorder, http.StatusNotFound, "user account not found")
}

func TestDeleteUserAccountInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteUserAccountFunc: func(ctx context.Context, id int64) error {
			return errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/user_account/delete/:id", handler.DeleteUserAccount)

	recorder := serve(router, http.MethodDelete, "/api/user_account/delete/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

// endregion

// region Category

func TestListCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCategories := []domain.Category{
		{ID: 1, Name: "testCategory1"},
		{ID: 2, Name: "testCategory2"},
	}

	repo := &fakeRepository{
		listCategoriesFunc: func(ctx context.Context) ([]domain.Category, error) {
			return testCategories, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/category/list", handler.ListCategories)

	recorder := serve(router, http.MethodGet, "/api/category/list", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got []domain.Category

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Name != "testCategory1" {
		t.Fatalf("expected name testCategory1, got %v", got[0].Name)
	}

	if got[1].ID != 2 {
		t.Fatalf("expected ID 2, got %d", got[1].ID)
	}

	if got[1].Name != "testCategory2" {
		t.Fatalf("expected name testCategory2, got %v", got[1].Name)
	}
}

// endregion

// region Product

func TestListProductCatalogue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testProducts := []domain.Product{
		{ID: 1, Name: "testProduct1", Price: 10.5, SellerID: 1, SellerUsername: "testSeller1", CategoryID: 1, CategoryName: "testCategory1"},
		{ID: 2, Name: "testProduct2", Price: 20.5, SellerID: 2, SellerUsername: "testSeller2", CategoryID: 2, CategoryName: "testCategory2"},
	}

	repo := &fakeRepository{
		listUnsoldProductsFunc: func(ctx context.Context) ([]domain.Product, error) {
			return testProducts, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/product/list", handler.ListProductCatalogue)

	recorder := serve(router, http.MethodGet, "/api/product/list", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got []domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 products, got %d", len(got))
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Name != "testProduct1" {
		t.Fatalf("expected name testProduct1, got %v", got[0].Name)
	}

	if got[0].Price != 10.5 {
		t.Fatalf("expected price 10.5, got %v", got[0].Price)
	}

	if got[0].SellerUsername != "testSeller1" {
		t.Fatalf("expected seller username testSeller1, got %v", got[0].SellerUsername)
	}

	if got[0].CategoryName != "testCategory1" {
		t.Fatalf("expected category name testCategory1, got %v", got[0].CategoryName)
	}

	if got[1].ID != 2 {
		t.Fatalf("expected ID 2, got %d", got[1].ID)
	}

	if got[1].Name != "testProduct2" {
		t.Fatalf("expected name testProduct2, got %v", got[1].Name)
	}
}

func TestListProductCatalogueInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		listUnsoldProductsFunc: func(ctx context.Context) ([]domain.Product, error) {
			return nil, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/product/list", handler.ListProductCatalogue)

	recorder := serve(router, http.MethodGet, "/api/product/list", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestGetSingleProductDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testProduct := domain.Product{
		ID:             1,
		Name:           "testProduct",
		Price:          10.5,
		Description:    "testDescription",
		ListedAt:       testNow,
		SellerID:       2,
		SellerUsername: "testSeller",
		CategoryID:     3,
		CategoryName:   "testCategory",
	}

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return testProduct, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/product/:id", handler.GetSingleProductDetails)

	recorder := serve(router, http.MethodGet, "/api/product/1", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != 1 {
		t.Fatalf("expected ID 1, got %d", got.ID)
	}

	if got.Name != "testProduct" {
		t.Fatalf("expected name testProduct, got %v", got.Name)
	}

	if got.Price != 10.5 {
		t.Fatalf("expected price 10.5, got %v", got.Price)
	}

	if got.Description != "testDescription" {
		t.Fatalf("expected description testDescription, got %v", got.Description)
	}

	if !got.ListedAt.Equal(testNow) {
		t.Fatalf("expected listed at %v, got %v", testNow, got.ListedAt)
	}

	if got.SellerID != 2 {
		t.Fatalf("expected seller id 2, got %d", got.SellerID)
	}

	if got.SellerUsername != "testSeller" {
		t.Fatalf("expected seller username testSeller, got %v", got.SellerUsername)
	}

	if got.CategoryID != 3 {
		t.Fatalf("expected category id 3, got %d", got.CategoryID)
	}

	if got.CategoryName != "testCategory" {
		t.Fatalf("expected category name testCategory, got %v", got.CategoryName)
	}
}

func TestGetSingleProductDetailsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.GET("/api/product/:id", handler.GetSingleProductDetails)

	recorder := serve(router, http.MethodGet, "/api/product/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid product id")
}

func TestGetSingleProductDetailsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{}, storefrontservice.ErrProductNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/product/:id", handler.GetSingleProductDetails)

	recorder := serve(router, http.MethodGet, "/api/product/1", "")

	assertError(t, recorder, http.StatusNotFound, "product not found")
}

func TestGetSingleProductDetailsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/product/:id", handler.GetSingleProductDetails)

	recorder := serve(router, http.MethodGet, "/api/product/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestListAllSellerProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testProducts := []domain.Product{
		{ID: 1, Name: "testProduct1", Price: 10.5, Description: "testDescription1", ListedAt: testNow, UpdatedAt: testNow, Sold: false, CategoryID: 1, CategoryName: "testCategory1"},
		{ID: 2, Name: "testProduct2", Price: 20.5, Description: "testDescription2", ListedAt: testNow, UpdatedAt: testNow, Sold: true, CategoryID: 2, CategoryName: "testCategory2"},
	}
	repo := &fakeRepository{
		listOwnProductsFunc: func(ctx context.Context, sellerId int64) ([]domain.Product, error) {
			return testProducts, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/product/seller/:id", handler.ListAllSellerProducts)

	recorder := serve(router, http.MethodGet, "/api/product/seller/1", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got []domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 products, got %d", len(got))
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Name != "testProduct1" {
		t.Fatalf("expected name testProduct1, got %v", got[0].Name)
	}

	if got[0].Description != "testDescription1" {
		t.Fatalf("expected description testDescription1, got %v", got[0].Description)
	}

	if got[0].Sold != false {
		t.Fatalf("expected sold false, got %v", got[0].Sold)
	}

	if got[1].ID != 2 {
		t.Fatalf("expected ID 2, got %d", got[1].ID)
	}

	if got[1].Sold != true {
		t.Fatalf("expected sold true, got %v", got[1].Sold)
	}
}

func TestListAllSellerProductsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.GET("/api/product/seller/:id", handler.ListAllSellerProducts)

	recorder := serve(router, http.MethodGet, "/api/product/seller/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid seller id")
}

func TestListAllSellerProductsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		listOwnProductsFunc: func(ctx context.Context, sellerId int64) ([]domain.Product, error) {
			return nil, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/product/seller/:id", handler.ListAllSellerProducts)

	recorder := serve(router, http.MethodGet, "/api/product/seller/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			product.ID = 1
			return product, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/api/product/create", handler.CreateProduct)

	body := `{"name":"testProduct","price":10.5,"description":"testDescription","sellerId":2,"categoryId":3}`

	recorder := serve(router, http.MethodPost, "/api/product/create", body)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var got domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != 1 {
		t.Fatalf("expected ID 1, got %d", got.ID)
	}

	if got.Name != "testProduct" {
		t.Fatalf("expected name testProduct, got %v", got.Name)
	}

	if got.Price != 10.5 {
		t.Fatalf("expected price 10.5, got %v", got.Price)
	}

	if got.Description != "testDescription" {
		t.Fatalf("expected description testDescription, got %v", got.Description)
	}

	if !got.ListedAt.Equal(testNow) {
		t.Fatalf("expected listed at %v, got %v", testNow, got.ListedAt)
	}

	if !got.UpdatedAt.Equal(testNow) {
		t.Fatalf("expected updated at %v, got %v", testNow, got.UpdatedAt)
	}

	if got.SellerID != 2 {
		t.Fatalf("expected seller id 2, got %d", got.SellerID)
	}

	if got.CategoryID != 3 {
		t.Fatalf("expected category id 3, got %d", got.CategoryID)
	}
}

func TestCreateProductInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/product/create", handler.CreateProduct)

	// price sent as a string, which cannot bind to float64
	body := `{"name":"testProduct","price":"abc","description":"testDescription","sellerId":2,"categoryId":3}`

	recorder := serve(router, http.MethodPost, "/api/product/create", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestCreateProductBlankName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/product/create", handler.CreateProduct)

	body := `{"name":"   ","price":10.5,"description":"testDescription","sellerId":2,"categoryId":3}`

	recorder := serve(router, http.MethodPost, "/api/product/create", body)

	assertError(t, recorder, http.StatusBadRequest, "title is required")
}

func TestCreateProductInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			return domain.Product{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/product/create", handler.CreateProduct)

	body := `{"name":"testProduct","price":10.5,"description":"testDescription","sellerId":2,"categoryId":3}`

	recorder := serve(router, http.MethodPost, "/api/product/create", body)

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestUpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	existingListedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{
				ID:          id,
				Name:        "oldName",
				Price:       10,
				Description: "oldDescription",
				ListedAt:    existingListedAt,
				UpdatedAt:   existingListedAt,
				SellerID:    2,
				CategoryID:  1,
			}, nil
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			return product, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	body := `{"name":"newName","price":20.5,"description":"newDescription","categoryId":3}`

	recorder := serve(router, http.MethodPut, "/api/product/update/1", body)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != 1 {
		t.Fatalf("expected ID 1, got %d", got.ID)
	}

	if got.Name != "newName" {
		t.Fatalf("expected name newName, got %v", got.Name)
	}

	if got.Price != 20.5 {
		t.Fatalf("expected price 20.5, got %v", got.Price)
	}

	if got.Description != "newDescription" {
		t.Fatalf("expected description newDescription, got %v", got.Description)
	}

	if got.CategoryID != 3 {
		t.Fatalf("expected category id 3, got %d", got.CategoryID)
	}

	if !got.ListedAt.Equal(existingListedAt) {
		t.Fatalf("expected listed at to be unchanged (%v), got %v", existingListedAt, got.ListedAt)
	}

	if !got.UpdatedAt.Equal(testNow) {
		t.Fatalf("expected updated at %v, got %v", testNow, got.UpdatedAt)
	}

	if got.SellerID != 2 {
		t.Fatalf("expected seller id 2 to be unchanged, got %d", got.SellerID)
	}
}

func TestUpdateProductInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	body := `{"name":"newName","price":20.5,"description":"newDescription","categoryId":3}`

	recorder := serve(router, http.MethodPut, "/api/product/update/abc", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid product id")
}

func TestUpdateProductInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	// categoryId sent as a string, which cannot bind to int64
	body := `{"name":"newName","price":20.5,"description":"newDescription","categoryId":"3"}`

	recorder := serve(router, http.MethodPut, "/api/product/update/1", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestUpdateProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{}, storefrontservice.ErrProductNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	body := `{"name":"newName","price":20.5,"description":"newDescription","categoryId":3}`

	recorder := serve(router, http.MethodPut, "/api/product/update/1", body)

	assertError(t, recorder, http.StatusNotFound, "product not found")
}

func TestUpdateProductBlankName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{ID: id, Name: "oldName"}, nil
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	body := `{"name":"   ","price":20.5,"description":"newDescription","categoryId":3}`

	recorder := serve(router, http.MethodPut, "/api/product/update/1", body)

	assertError(t, recorder, http.StatusBadRequest, "title is required")
}

func TestUpdateProductInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{ID: id, Name: "oldName"}, nil
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			return domain.Product{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.PUT("/api/product/update/:id", handler.UpdateProduct)

	body := `{"name":"newName","price":20.5,"description":"newDescription","categoryId":3}`

	recorder := serve(router, http.MethodPut, "/api/product/update/1", body)

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestDeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testProducts := []domain.Product{
		{ID: 1, Name: "testProduct1"},
		{ID: 2, Name: "testProduct2"},
	}

	repo := &fakeRepository{
		deleteProductFunc: func(ctx context.Context, id int64) error {
			for i := range testProducts {
				if testProducts[i].ID == id {
					testProducts = append(testProducts[:i], testProducts[i+1:]...)
					return nil
				}
			}
			return storefrontservice.ErrProductNotFound
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/api/product/delete/:id", handler.DeleteProduct)

	recorder := serve(router, http.MethodDelete, "/api/product/delete/1", "")

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if len(testProducts) != 1 {
		t.Fatalf("expected number of products to be 1, got %v", len(testProducts))
	}

	if testProducts[0].ID != 2 {
		t.Fatalf("expected remaining product ID 2, got %d", testProducts[0].ID)
	}
}

func TestDeleteProductInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.DELETE("/api/product/delete/:id", handler.DeleteProduct)

	recorder := serve(router, http.MethodDelete, "/api/product/delete/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid product id")
}

func TestDeleteProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteProductFunc: func(ctx context.Context, id int64) error {
			return storefrontservice.ErrProductNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/product/delete/:id", handler.DeleteProduct)

	recorder := serve(router, http.MethodDelete, "/api/product/delete/1", "")

	assertError(t, recorder, http.StatusNotFound, "product not found")
}

func TestDeleteProductInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteProductFunc: func(ctx context.Context, id int64) error {
			return errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/product/delete/:id", handler.DeleteProduct)

	recorder := serve(router, http.MethodDelete, "/api/product/delete/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

// endregion

// region Cart Item

func TestListAllCustomerCartItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testProducts := []domain.Product{
		{ID: 1, Name: "testProduct1", Price: 10.5, SellerID: 1, SellerUsername: "testSeller1", CategoryID: 1, CategoryName: "testCategory1"},
		{ID: 2, Name: "testProduct2", Price: 20.5, SellerID: 2, SellerUsername: "testSeller2", CategoryID: 2, CategoryName: "testCategory2"},
	}

	var gotCustomerId int64

	repo := &fakeRepository{
		listProductsInCartFunc: func(ctx context.Context, customerId int64) ([]domain.Product, error) {
			gotCustomerId = customerId
			return testProducts, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/api/cart_item/list/:customerId", handler.ListAllCustomerCartItems)

	recorder := serve(router, http.MethodGet, "/api/cart_item/list/7", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if gotCustomerId != 7 {
		t.Fatalf("expected customer id 7 to be passed to repository, got %d", gotCustomerId)
	}

	var got []domain.Product

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 products, got %d", len(got))
	}

	if got[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", got[0].ID)
	}

	if got[0].Name != "testProduct1" {
		t.Fatalf("expected name testProduct1, got %v", got[0].Name)
	}

	if got[0].Price != 10.5 {
		t.Fatalf("expected price 10.5, got %v", got[0].Price)
	}

	if got[0].SellerUsername != "testSeller1" {
		t.Fatalf("expected seller username testSeller1, got %v", got[0].SellerUsername)
	}

	if got[1].ID != 2 {
		t.Fatalf("expected ID 2, got %d", got[1].ID)
	}

	if got[1].Name != "testProduct2" {
		t.Fatalf("expected name testProduct2, got %v", got[1].Name)
	}
}

func TestListAllCustomerCartItemsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.GET("/api/cart_item/list/:customerId", handler.ListAllCustomerCartItems)

	recorder := serve(router, http.MethodGet, "/api/cart_item/list/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid customer id")
}

func TestListAllCustomerCartItemsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		listProductsInCartFunc: func(ctx context.Context, customerId int64) ([]domain.Product, error) {
			return nil, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.GET("/api/cart_item/list/:customerId", handler.ListAllCustomerCartItems)

	recorder := serve(router, http.MethodGet, "/api/cart_item/list/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestAddProductToCustomerCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createCartItemFunc: func(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
			return cartItem, nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/api/cart_item/create", handler.AddProductToCustomerCart)

	body := `{"customerId":1,"productId":2}`

	recorder := serve(router, http.MethodPost, "/api/cart_item/create", body)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var got domain.CartItem

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.CustomerID != 1 {
		t.Fatalf("expected customer id 1, got %d", got.CustomerID)
	}

	if got.ProductID != 2 {
		t.Fatalf("expected product id 2, got %d", got.ProductID)
	}
}

func TestAddProductToCustomerCartInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.POST("/api/cart_item/create", handler.AddProductToCustomerCart)

	// ids sent as strings, which cannot bind to int64
	body := `{"customerId":"1","productId":"2"}`

	recorder := serve(router, http.MethodPost, "/api/cart_item/create", body)

	assertError(t, recorder, http.StatusBadRequest, "invalid request")
}

func TestAddProductToCustomerCartInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		createCartItemFunc: func(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
			return domain.CartItem{}, errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.POST("/api/cart_item/create", handler.AddProductToCustomerCart)

	body := `{"customerId":1,"productId":2}`

	recorder := serve(router, http.MethodPost, "/api/cart_item/create", body)

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestRemoveProductFromCustomerCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCartItems := []domain.CartItem{
		{CustomerID: 1, ProductID: 1},
		{CustomerID: 1, ProductID: 2},
	}

	repo := &fakeRepository{
		deleteCartItemFunc: func(ctx context.Context, customerId int64, productId int64) error {
			for i := range testCartItems {
				if testCartItems[i].CustomerID == customerId && testCartItems[i].ProductID == productId {
					testCartItems = append(testCartItems[:i], testCartItems[i+1:]...)
					return nil
				}
			}
			return storefrontservice.ErrCartItemNotFound
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/delete/1/2", "")

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if len(testCartItems) != 1 {
		t.Fatalf("expected number of cart items to be 1, got %v", len(testCartItems))
	}

	if testCartItems[0].ProductID != 1 {
		t.Fatalf("expected remaining product id 1, got %d", testCartItems[0].ProductID)
	}
}

func TestRemoveProductFromCustomerCartInvalidCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/delete/abc/2", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid customer id")
}

func TestRemoveProductFromCustomerCartInvalidProductID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/delete/1/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid product id")
}

func TestRemoveProductFromCustomerCartNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteCartItemFunc: func(ctx context.Context, customerId int64, productId int64) error {
			return storefrontservice.ErrCartItemNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/delete/1/2", "")

	assertError(t, recorder, http.StatusNotFound, "cart item not found")
}

func TestRemoveProductFromCustomerCartInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteCartItemFunc: func(ctx context.Context, customerId int64, productId int64) error {
			return errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/delete/:customerId/:productId", handler.RemoveProductFromCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/delete/1/2", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

func TestCheckoutCustomerCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCartItems := []domain.CartItem{
		{CustomerID: 1, ProductID: 1},
		{CustomerID: 1, ProductID: 2},
		{CustomerID: 2, ProductID: 3},
	}

	repo := &fakeRepository{
		deleteCartItemsFunc: func(ctx context.Context, customerId int64) error {
			remaining := []domain.CartItem{}
			for _, item := range testCartItems {
				if item.CustomerID != customerId {
					remaining = append(remaining, item)
				}
			}
			if len(remaining) == len(testCartItems) {
				return storefrontservice.ErrCartItemNotFound
			}
			testCartItems = remaining
			return nil
		},
	}

	service := storefrontservice.NewService(repo, fixedNow)

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/api/cart_item/checkout/:customerId", handler.CheckoutCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/checkout/1", "")

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if len(testCartItems) != 1 {
		t.Fatalf("expected number of cart items to be 1, got %v", len(testCartItems))
	}

	if testCartItems[0].CustomerID != 2 {
		t.Fatalf("expected remaining cart item to belong to customer 2, got %d", testCartItems[0].CustomerID)
	}
}

func TestCheckoutCustomerCartInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(storefrontservice.NewService(&fakeRepository{}, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/checkout/:customerId", handler.CheckoutCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/checkout/abc", "")

	assertError(t, recorder, http.StatusBadRequest, "invalid customer id")
}

func TestCheckoutCustomerCartNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteCartItemsFunc: func(ctx context.Context, customerId int64) error {
			return storefrontservice.ErrCartItemNotFound
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/checkout/:customerId", handler.CheckoutCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/checkout/1", "")

	assertError(t, recorder, http.StatusNotFound, "cart item not found")
}

func TestCheckoutCustomerCartInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeRepository{
		deleteCartItemsFunc: func(ctx context.Context, customerId int64) error {
			return errors.New("database failure")
		},
	}

	handler := NewHandler(storefrontservice.NewService(repo, fixedNow))

	router := gin.New()
	router.DELETE("/api/cart_item/checkout/:customerId", handler.CheckoutCustomerCart)

	recorder := serve(router, http.MethodDelete, "/api/cart_item/checkout/1", "")

	assertError(t, recorder, http.StatusInternalServerError, "internal server error")
}

// endregion

// region Helper

var testNow = time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)

func fixedNow() time.Time {
	return testNow
}

func serve(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	return recorder
}

func assertError(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, wantMessage string) {
	t.Helper()

	if recorder.Code != wantStatus {
		t.Fatalf("expected status %d, got %d", wantStatus, recorder.Code)
	}

	var got map[string]string

	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["error"] != wantMessage {
		t.Fatalf("expected error %q, got %q", wantMessage, got["error"])
	}
}

type fakeRepository struct {
	listUserTypesFunc      func(ctx context.Context) ([]domain.UserType, error)
	createUserAccountFunc  func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error)
	listUserAccountsFunc   func(ctx context.Context) ([]domain.UserAccount, error)
	getUserAccountFunc     func(ctx context.Context, username string) (domain.UserAccount, error)
	deleteUserAccountFunc  func(ctx context.Context, id int64) error
	listCategoriesFunc     func(ctx context.Context) ([]domain.Category, error)
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

func (f *fakeRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return f.listCategoriesFunc(ctx)
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
