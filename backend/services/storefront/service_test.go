package storefront

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "storefrontapp/domain/storefront"
)

// region User Type

func TestServiceListAllUserTypes(t *testing.T) {
	testUserTypes := []domain.UserType{
		{ID: 1, Type: "testCustomer"},
		{ID: 2, Type: "testSeller"},
	}

	calls := 0

	repo := &fakeRepository{
		listUserTypesFunc: func(ctx context.Context) ([]domain.UserType, error) {
			calls++
			return testUserTypes, nil
		},
	}

	service := NewService(repo, fixedNow)

	userTypes, err := service.ListAllUserTypes(context.Background())
	if err != nil {
		t.Fatalf("list all user types %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListUserTypes to be called once, got %v calls instead", calls)
	}

	if len(userTypes) != len(testUserTypes) {
		t.Fatalf("expected %v user types, got %v instead", len(testUserTypes), len(userTypes))
	}

	for i, userType := range userTypes {
		if userType != testUserTypes[i] {
			t.Fatalf("expected user type %+v at index %v, got %+v instead", testUserTypes[i], i, userType)
		}
	}
}

func TestServiceListAllUserTypesError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		listUserTypesFunc: func(ctx context.Context) ([]domain.UserType, error) {
			return nil, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.ListAllUserTypes(context.Background())
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

// endregion

// region User Account

func TestServiceSignUpUserAccount(t *testing.T) {
	var receivedUserAccount domain.UserAccount // stores the user account that was received by the repo
	calls := 0

	repo := &fakeRepository{
		createUserAccountFunc: func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
			calls++
			receivedUserAccount = userAccount
			userAccount.ID = 1
			return userAccount, nil
		},
	}

	service := NewService(repo, fixedNow)

	created, err := service.SignUpUserAccount(context.Background(), "testUsername", "testPasswordHash", 2) // stores the user account created and sent back by the repo
	if err != nil {
		t.Fatalf("sign up user account %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository CreateUserAccount to be called once, got %v calls instead", calls)
	}

	if receivedUserAccount.Username != "testUsername" {
		t.Fatalf("expected repository to receive username %v, got %v instead", "testUsername", receivedUserAccount.Username)
	}

	if receivedUserAccount.PasswordHash != "testPasswordHash" {
		t.Fatalf("expected repository to receive password hash %v, got %v instead", "testPasswordHash", receivedUserAccount.PasswordHash)
	}

	if receivedUserAccount.UserTypeID != 2 {
		t.Fatalf("expected repository to receive user type id %v, got %v instead", 2, receivedUserAccount.UserTypeID)
	}

	if created.ID != 1 {
		t.Fatalf("expected created user account id %v, got %v instead", 1, created.ID)
	}

}

func TestServiceSignUpUserAccountInvalid(t *testing.T) {
	repo := &fakeRepository{
		createUserAccountFunc: func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
			t.Fatalf("expected repository CreateUserAccount not to be called")
			return domain.UserAccount{}, nil
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.SignUpUserAccount(context.Background(), "", "testPasswordHash", 2)
	if err == nil {
		t.Fatalf("expected an error for an empty username, got nil instead")
	}
}

func TestServiceSignUpUserAccountError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		createUserAccountFunc: func(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
			return domain.UserAccount{}, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.SignUpUserAccount(context.Background(), "testUsername", "testPasswordHash", 2)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceListAllUserAccounts(t *testing.T) {
	testUserAccounts := []domain.UserAccount{
		newTestUserAccount(),
		{ID: 2, Username: "testUsername2", PasswordHash: "testPasswordHash2", UserTypeID: 2, UserType: "testUserType2"},
	}

	calls := 0

	repo := &fakeRepository{
		listUserAccountsFunc: func(ctx context.Context) ([]domain.UserAccount, error) {
			calls++
			return testUserAccounts, nil
		},
	}

	service := NewService(repo, fixedNow)

	userAccounts, err := service.ListAllUserAccounts(context.Background())
	if err != nil {
		t.Fatalf("list all user accounts %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListUserAccounts to be called once, got %v calls instead", calls)
	}

	if len(userAccounts) != len(testUserAccounts) {
		t.Fatalf("expected %v user accounts, got %v instead", len(testUserAccounts), len(userAccounts))
	}

	for i, userAccount := range userAccounts {
		if userAccount != testUserAccounts[i] {
			t.Fatalf("expected user account %+v at index %v, got %+v instead", testUserAccounts[i], i, userAccount)
		}
	}
}

func TestServiceListAllUserAccountsError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		listUserAccountsFunc: func(ctx context.Context) ([]domain.UserAccount, error) {
			return nil, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.ListAllUserAccounts(context.Background())
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceLogIntoUserAccount(t *testing.T) {
	testUserAccount := newTestUserAccount()

	var receivedUsername string // stores the username received by the repo
	calls := 0

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			calls++
			receivedUsername = username
			return testUserAccount, nil
		},
	}

	service := NewService(repo, fixedNow)

	userAccount, ok, err := service.LogIntoUserAccount(context.Background(), testUserAccount.Username, testUserAccount.PasswordHash)
	if err != nil {
		t.Fatalf("log into user account %v", err)
	}

	if !ok {
		t.Fatalf("expected login to succeed, got false instead")
	}

	if calls != 1 {
		t.Fatalf("expected repository GetUserAccount to be called once, got %v calls instead", calls)
	}

	if receivedUsername != testUserAccount.Username {
		t.Fatalf("expected repository to receive username %v, got %v instead", testUserAccount.Username, receivedUsername)
	}

	if userAccount != testUserAccount {
		t.Fatalf("expected user account %+v, got %+v instead", testUserAccount, userAccount)
	}
}

func TestServiceLogIntoUserAccountWrongPassword(t *testing.T) {
	testUserAccount := newTestUserAccount()

	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return testUserAccount, nil
		},
	}

	service := NewService(repo, fixedNow)

	userAccount, ok, err := service.LogIntoUserAccount(context.Background(), testUserAccount.Username, "testWrongPasswordHash")
	if !errors.Is(err, ErrIncorrectUsernameOrPassword) {
		t.Fatalf("expected error %v, got %v instead", ErrIncorrectUsernameOrPassword, err)
	}

	if ok {
		t.Fatalf("expected login to fail, got true instead")
	}

	if userAccount != (domain.UserAccount{}) {
		t.Fatalf("expected empty user account, got %+v instead", userAccount)
	}
}

func TestServiceLogIntoUserAccountNotFound(t *testing.T) {
	repo := &fakeRepository{
		getUserAccountFunc: func(ctx context.Context, username string) (domain.UserAccount, error) {
			return domain.UserAccount{}, ErrUserAccountNotFound
		},
	}

	service := NewService(repo, fixedNow)

	userAccount, ok, err := service.LogIntoUserAccount(context.Background(), "testUsername", "testPasswordHash")
	if !errors.Is(err, ErrUserAccountNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrUserAccountNotFound, err)
	}

	if ok {
		t.Fatalf("expected login to fail, got true instead")
	}

	if userAccount != (domain.UserAccount{}) {
		t.Fatalf("expected empty user account, got %+v instead", userAccount)
	}
}

func TestServiceDeleteUserAccount(t *testing.T) {
	var receivedID int64
	calls := 0

	repo := &fakeRepository{
		deleteUserAccountFunc: func(ctx context.Context, id int64) error {
			calls++
			receivedID = id
			return nil
		},
	}

	service := NewService(repo, fixedNow)

	err := service.DeleteUserAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("delete user account %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository DeleteUserAccount to be called once, got %v calls instead", calls)
	}

	if receivedID != 1 {
		t.Fatalf("expected repository to receive user account id %v, got %v instead", 1, receivedID)
	}
}

func TestServiceDeleteUserAccountNotFound(t *testing.T) {
	repo := &fakeRepository{
		deleteUserAccountFunc: func(ctx context.Context, id int64) error {
			return ErrUserAccountNotFound
		},
	}

	service := NewService(repo, fixedNow)

	err := service.DeleteUserAccount(context.Background(), 1)
	if !errors.Is(err, ErrUserAccountNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrUserAccountNotFound, err)
	}
}

// endregion

// region Category

func TestListCategories(t *testing.T) {
	testCategories := []domain.Category{
		{ID: 1, Name: "testCategory1"},
		{ID: 2, Name: "testCategory2"},
	}

	calls := 0

	repo := &fakeRepository{
		listCategoriesFunc: func(ctx context.Context) ([]domain.Category, error) {
			calls++
			return testCategories, nil
		},
	}

	service := NewService(repo, fixedNow)

	categories, err := service.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("list all categories %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListCategories to be called once, got %v calls instead", calls)
	}

	if len(categories) != len(testCategories) {
		t.Fatalf("expected %v categories, got %v instead", len(testCategories), len(categories))
	}

	for i, category := range categories {
		if category != testCategories[i] {
			t.Fatalf("expected category %+v at index %v, got %+v instead", testCategories[i], i, category)
		}
	}
}

// endregion

// region Product

func TestServiceListProductCatalogue(t *testing.T) {
	testProducts := []domain.Product{
		newTestProduct(),
		{ID: 43, Name: "testOtherProduct", Price: 30, SellerID: 7, SellerUsername: "testSeller", CategoryID: 3, CategoryName: "testCategory"},
	}

	calls := 0

	repo := &fakeRepository{
		listUnsoldProductsFunc: func(ctx context.Context) ([]domain.Product, error) {
			calls++
			return testProducts, nil
		},
	}

	service := NewService(repo, fixedNow)

	products, err := service.ListProductCatalogue(context.Background())
	if err != nil {
		t.Fatalf("list product catalogue %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListUnsoldProducts to be called once, got %v calls instead", calls)
	}

	if len(products) != len(testProducts) {
		t.Fatalf("expected %v products, got %v instead", len(testProducts), len(products))
	}

	for i, product := range products {
		if product != testProducts[i] {
			t.Fatalf("expected product %+v at index %v, got %+v instead", testProducts[i], i, product)
		}
	}
}

func TestServiceListProductCatalogueError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		listUnsoldProductsFunc: func(ctx context.Context) ([]domain.Product, error) {
			return nil, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.ListProductCatalogue(context.Background())
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceGetSingleProductDetails(t *testing.T) {
	testProduct := newTestProduct()

	var receivedID int64
	calls := 0

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			calls++
			receivedID = id
			return testProduct, nil
		},
	}

	service := NewService(repo, fixedNow)

	product, err := service.GetSingleProductDetails(context.Background(), testProduct.ID)
	if err != nil {
		t.Fatalf("get single product details %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository GetProduct to be called once, got %v calls instead", calls)
	}

	if receivedID != testProduct.ID {
		t.Fatalf("expected repository to receive product id %v, got %v instead", testProduct.ID, receivedID)
	}

	if product != testProduct {
		t.Fatalf("expected product %+v, got %+v instead", testProduct, product)
	}
}

func TestServiceGetSingleProductDetailsNotFound(t *testing.T) {
	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{}, ErrProductNotFound
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.GetSingleProductDetails(context.Background(), 42)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrProductNotFound, err)
	}
}

func TestServiceListAllSellerProducts(t *testing.T) {
	testProducts := []domain.Product{newTestProduct()}

	var receivedSellerID int64
	calls := 0

	repo := &fakeRepository{
		listOwnProductsFunc: func(ctx context.Context, sellerId int64) ([]domain.Product, error) {
			calls++
			receivedSellerID = sellerId
			return testProducts, nil
		},
	}

	service := NewService(repo, fixedNow)

	products, err := service.ListAllSellerProducts(context.Background(), testProducts[0].SellerID)
	if err != nil {
		t.Fatalf("list all seller products %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListOwnProducts to be called once, got %v calls instead", calls)
	}

	if receivedSellerID != testProducts[0].SellerID {
		t.Fatalf("expected repository to receive seller id %v, got %v instead", testProducts[0].SellerID, receivedSellerID)
	}

	if len(products) != len(testProducts) {
		t.Fatalf("expected %v products, got %v instead", len(testProducts), len(products))
	}

	for i, product := range products {
		if product != testProducts[i] {
			t.Fatalf("expected product %+v at index %v, got %+v instead", testProducts[i], i, product)
		}
	}
}

func TestServiceListAllSellerProductsError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		listOwnProductsFunc: func(ctx context.Context, sellerId int64) ([]domain.Product, error) {
			return nil, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.ListAllSellerProducts(context.Background(), 7)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceCreateProduct(t *testing.T) {
	var receivedProduct domain.Product
	calls := 0

	repo := &fakeRepository{
		createProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			calls++
			receivedProduct = product
			product.ID = 42
			return product, nil
		},
	}

	service := NewService(repo, fixedNow)

	created, err := service.CreateProduct(context.Background(), "testProduct", 20, "This is a test product.", 7, 3)
	if err != nil {
		t.Fatalf("create product %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository CreateProduct to be called once, got %v calls instead", calls)
	}

	if receivedProduct.Name != "testProduct" {
		t.Fatalf("expected repository to receive product name %v, got %v instead", "testProduct", receivedProduct.Name)
	}

	if receivedProduct.Price != 20 {
		t.Fatalf("expected repository to receive product price %v, got %v instead", 20, receivedProduct.Price)
	}

	if receivedProduct.Description != "This is a test product." {
		t.Fatalf("expected repository to receive product description %v, got %v instead", "This is a test product.", receivedProduct.Description)
	}

	if !receivedProduct.ListedAt.Equal(testNow) {
		t.Fatalf("expected repository to receive product listed at %v, got %v instead", testNow, receivedProduct.ListedAt)
	}

	if !receivedProduct.UpdatedAt.Equal(testNow) {
		t.Fatalf("expected repository to receive product updated at %v, got %v instead", testNow, receivedProduct.UpdatedAt)
	}

	if receivedProduct.SellerID != 7 {
		t.Fatalf("expected repository to receive product seller id %v, got %v instead", 7, receivedProduct.SellerID)
	}

	if receivedProduct.CategoryID != 3 {
		t.Fatalf("expected repository to receive product category id %v, got %v instead", 3, receivedProduct.CategoryID)
	}

	if created.ID != 42 {
		t.Fatalf("expected created product id %v, got %v instead", 42, created.ID)
	}
}

func TestServiceCreateProductInvalid(t *testing.T) {
	repo := &fakeRepository{
		createProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			t.Fatalf("expected repository CreateProduct not to be called")
			return domain.Product{}, nil
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.CreateProduct(context.Background(), "", 20, "This is a test product.", 7, 3)
	if err == nil {
		t.Fatalf("expected an error for an empty product name, got nil instead")
	}
}

func TestServiceCreateProductError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		createProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			return domain.Product{}, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.CreateProduct(context.Background(), "testProduct", 20, "This is a test product.", 7, 3)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceDeleteProduct(t *testing.T) {
	var receivedID int64
	calls := 0

	repo := &fakeRepository{
		deleteProductFunc: func(ctx context.Context, id int64) error {
			calls++
			receivedID = id
			return nil
		},
	}

	service := NewService(repo, fixedNow)

	err := service.DeleteProduct(context.Background(), 42)
	if err != nil {
		t.Fatalf("delete product %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository DeleteProduct to be called once, got %v calls instead", calls)
	}

	if receivedID != 42 {
		t.Fatalf("expected repository to receive product id %v, got %v instead", 42, receivedID)
	}
}

func TestServiceDeleteProductNotFound(t *testing.T) {
	repo := &fakeRepository{
		deleteProductFunc: func(ctx context.Context, id int64) error {
			return ErrProductNotFound
		},
	}

	service := NewService(repo, fixedNow)

	err := service.DeleteProduct(context.Background(), 42)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrProductNotFound, err)
	}
}

func TestServiceUpdateProduct(t *testing.T) {
	testProduct := newTestProduct()
	testProduct.ListedAt = testNow.Add(-24 * time.Hour)
	testProduct.UpdatedAt = testNow.Add(-24 * time.Hour)

	var receivedID int64
	var receivedProduct domain.Product
	getCalls := 0
	updateCalls := 0

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			getCalls++
			receivedID = id
			return testProduct, nil
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			updateCalls++
			receivedProduct = product
			return product, nil
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.UpdateProduct(context.Background(), testProduct.ID, "testUpdatedProduct", 25, "This is an updated test product.", 4)
	if err != nil {
		t.Fatalf("update product %v", err)
	}

	if getCalls != 1 {
		t.Fatalf("expected repository GetProduct to be called once, got %v calls instead", getCalls)
	}

	if receivedID != testProduct.ID {
		t.Fatalf("expected repository to receive product id %v, got %v instead", testProduct.ID, receivedID)
	}

	if updateCalls != 1 {
		t.Fatalf("expected repository UpdateProduct to be called once, got %v calls instead", updateCalls)
	}

	if receivedProduct.ID != testProduct.ID {
		t.Fatalf("expected updated product id %v, got %v instead", testProduct.ID, receivedProduct.ID)
	}

	if receivedProduct.Name != "testUpdatedProduct" {
		t.Fatalf("expected updated product name %v, got %v instead", "testUpdatedProduct", receivedProduct.Name)
	}

	if receivedProduct.Price != 25 {
		t.Fatalf("expected updated product price %v, got %v instead", 25, receivedProduct.Price)
	}

	if receivedProduct.Description != "This is an updated test product." {
		t.Fatalf("expected updated product description %v, got %v instead", "This is an updated test product.", receivedProduct.Description)
	}

	if !receivedProduct.ListedAt.Equal(testProduct.ListedAt) {
		t.Fatalf("expected product listed at %v to be unchanged, got %v instead", testProduct.ListedAt, receivedProduct.ListedAt)
	}

	if !receivedProduct.UpdatedAt.Equal(testNow) {
		t.Fatalf("expected updated product updated at %v, got %v instead", testNow, receivedProduct.UpdatedAt)
	}

	if receivedProduct.SellerID != testProduct.SellerID {
		t.Fatalf("expected product seller id %v to be unchanged, got %v instead", testProduct.SellerID, receivedProduct.SellerID)
	}

	if receivedProduct.CategoryID != 4 {
		t.Fatalf("expected updated product category id %v, got %v instead", 4, receivedProduct.CategoryID)
	}
}

func TestServiceUpdateProductNotFound(t *testing.T) {
	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return domain.Product{}, ErrProductNotFound
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			t.Fatalf("expected repository UpdateProduct not to be called")
			return domain.Product{}, nil
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.UpdateProduct(context.Background(), 42, "testUpdatedProduct", 25, "This is an updated test product.", 4)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrProductNotFound, err)
	}
}

func TestServiceUpdateProductInvalid(t *testing.T) {
	testProduct := newTestProduct()

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return testProduct, nil
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			t.Fatalf("expected repository UpdateProduct not to be called")
			return domain.Product{}, nil
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.UpdateProduct(context.Background(), testProduct.ID, "", 25, "This is an updated test product.", 4)
	if err == nil {
		t.Fatalf("expected an error for an empty product name, got nil instead")
	}
}

func TestServiceUpdateProductError(t *testing.T) {
	testErr := errors.New("test repository error")
	testProduct := newTestProduct()

	repo := &fakeRepository{
		getProductFunc: func(ctx context.Context, id int64) (domain.Product, error) {
			return testProduct, nil
		},
		updateProductFunc: func(ctx context.Context, product domain.Product) (domain.Product, error) {
			return domain.Product{}, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.UpdateProduct(context.Background(), testProduct.ID, "testUpdatedProduct", 25, "This is an updated test product.", 4)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

// endregion

// region Cart Item

func TestServiceListAllCustomerCartItems(t *testing.T) {
	testProducts := []domain.Product{newTestProduct()}

	var receivedCustomerID int64
	calls := 0

	repo := &fakeRepository{
		listProductsInCartFunc: func(ctx context.Context, customerId int64) ([]domain.Product, error) {
			calls++
			receivedCustomerID = customerId
			return testProducts, nil
		},
	}

	service := NewService(repo, fixedNow)

	products, err := service.ListAllCustomerCartItems(context.Background(), 5)
	if err != nil {
		t.Fatalf("list all customer cart items %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository ListProductsInCart to be called once, got %v calls instead", calls)
	}

	if receivedCustomerID != 5 {
		t.Fatalf("expected repository to receive customer id %v, got %v instead", 5, receivedCustomerID)
	}

	if len(products) != len(testProducts) {
		t.Fatalf("expected %v products, got %v instead", len(testProducts), len(products))
	}

	for i, product := range products {
		if product != testProducts[i] {
			t.Fatalf("expected product %+v at index %v, got %+v instead", testProducts[i], i, product)
		}
	}
}

func TestServiceListAllCustomerCartItemsError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		listProductsInCartFunc: func(ctx context.Context, customerId int64) ([]domain.Product, error) {
			return nil, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.ListAllCustomerCartItems(context.Background(), 5)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceAddProductToCustomerCart(t *testing.T) {
	var receivedCartItem domain.CartItem
	calls := 0

	repo := &fakeRepository{
		createCartItemFunc: func(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
			calls++
			receivedCartItem = cartItem
			return cartItem, nil
		},
	}

	service := NewService(repo, fixedNow)

	created, err := service.AddProductToCustomerCart(context.Background(), 5, 42)
	if err != nil {
		t.Fatalf("add product to customer cart %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository CreateCartItem to be called once, got %v calls instead", calls)
	}

	if receivedCartItem.CustomerID != 5 {
		t.Fatalf("expected repository to receive customer id %v, got %v instead", 5, receivedCartItem.CustomerID)
	}

	if receivedCartItem.ProductID != 42 {
		t.Fatalf("expected repository to receive product id %v, got %v instead", 42, receivedCartItem.ProductID)
	}

	if created != receivedCartItem {
		t.Fatalf("expected cart item %+v, got %+v instead", receivedCartItem, created)
	}
}

func TestServiceAddProductToCustomerCartError(t *testing.T) {
	testErr := errors.New("test repository error")

	repo := &fakeRepository{
		createCartItemFunc: func(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
			return domain.CartItem{}, testErr
		},
	}

	service := NewService(repo, fixedNow)

	_, err := service.AddProductToCustomerCart(context.Background(), 5, 42)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v instead", testErr, err)
	}
}

func TestServiceRemoveProductFromCustomerCart(t *testing.T) {
	var receivedCustomerID int64
	var receivedProductID int64
	calls := 0

	repo := &fakeRepository{
		deleteCartItemFunc: func(ctx context.Context, customerId int64, productId int64) error {
			calls++
			receivedCustomerID = customerId
			receivedProductID = productId
			return nil
		},
	}

	service := NewService(repo, fixedNow)

	err := service.RemoveProductFromCustomerCart(context.Background(), 5, 42)
	if err != nil {
		t.Fatalf("remove product from customer cart %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository DeleteCartItem to be called once, got %v calls instead", calls)
	}

	if receivedCustomerID != 5 {
		t.Fatalf("expected repository to receive customer id %v, got %v instead", 5, receivedCustomerID)
	}

	if receivedProductID != 42 {
		t.Fatalf("expected repository to receive product id %v, got %v instead", 42, receivedProductID)
	}
}

func TestServiceRemoveProductFromCustomerCartNotFound(t *testing.T) {
	repo := &fakeRepository{
		deleteCartItemFunc: func(ctx context.Context, customerId int64, productId int64) error {
			return ErrCartItemNotFound
		},
	}

	service := NewService(repo, fixedNow)

	err := service.RemoveProductFromCustomerCart(context.Background(), 5, 42)
	if !errors.Is(err, ErrCartItemNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrCartItemNotFound, err)
	}
}

func TestServiceCheckoutCustomerCart(t *testing.T) {
	var receivedCustomerID int64
	calls := 0

	repo := &fakeRepository{
		deleteCartItemsFunc: func(ctx context.Context, customerId int64) error {
			calls++
			receivedCustomerID = customerId
			return nil
		},
	}

	service := NewService(repo, fixedNow)

	err := service.CheckoutCustomerCart(context.Background(), 5)
	if err != nil {
		t.Fatalf("checkout customer cart %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected repository DeleteCartItems to be called once, got %v calls instead", calls)
	}

	if receivedCustomerID != 5 {
		t.Fatalf("expected repository to receive customer id %v, got %v instead", 5, receivedCustomerID)
	}
}

func TestServiceCheckoutCustomerCartNotFound(t *testing.T) {
	repo := &fakeRepository{
		deleteCartItemsFunc: func(ctx context.Context, customerId int64) error {
			return ErrCartItemNotFound
		},
	}

	service := NewService(repo, fixedNow)

	err := service.CheckoutCustomerCart(context.Background(), 5)
	if !errors.Is(err, ErrCartItemNotFound) {
		t.Fatalf("expected error %v, got %v instead", ErrCartItemNotFound, err)
	}
}

// endregion

// region Helper

var testNow = time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)

func fixedNow() time.Time {
	return testNow
}

func newTestUserAccount() domain.UserAccount {
	return domain.UserAccount{
		ID:           1,
		Username:     "testUsername",
		PasswordHash: "testPasswordHash",
		UserTypeID:   2,
		UserType:     "testCustomer",
	}
}

func newTestProduct() domain.Product {
	return domain.Product{
		ID:             42,
		Name:           "testProduct",
		Price:          20,
		Description:    "This is a test product.",
		ListedAt:       testNow,
		UpdatedAt:      testNow,
		SellerID:       7,
		SellerUsername: "testSeller",
		CategoryID:     3,
		CategoryName:   "testCategory",
	}
}

type fakeRepository struct {
	Repository // embedded field - gives the struct every method of the interface, therefore making it valid

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
