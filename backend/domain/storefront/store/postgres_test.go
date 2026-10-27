package store

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	domain "storefrontapp/domain/storefront"
)

// region User Type

func TestPostgresRepositoryListUserTypes(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")

	userTypes, err := repo.ListUserTypes(context.Background())
	if err != nil {
		t.Fatalf("list user types: %v", err)
	}

	foundTestType := false
	for _, userType := range userTypes {
		if userType.ID == 0 {
			t.Fatalf("expected user type id to be set, got 0")
		}
		if userType.Type == "" {
			t.Fatalf("expected user type %d to have a type name", userType.ID)
		}
		if userType.ID == testUserTypeID {
			foundTestType = true
		}
	}

	if !foundTestType {
		t.Fatalf("expected user type %d in list", testUserTypeID)
	}

}

// endregion

// region User Account
func TestPostgresRepositoryCreateUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")

	userAccount, err := domain.NewUserAccount("testUsername", "testPasswordHash", testUserTypeID)
	if err != nil {
		t.Fatalf("create domain user account %v", err)
	}

	created, err := repo.CreateUserAccount(context.Background(), userAccount)
	if err != nil {
		t.Fatalf("create user account %v", err)
	}

	cleanupRows(t, db, "DELETE FROM user_account WHERE id = $1", created.ID)

	if created.ID == 0 {
		t.Fatalf("expected database to generate id")
	}

	if created.Username != "testUsername" {
		t.Fatalf("expected username testUsername, got %s", created.Username)
	}

	if created.PasswordHash != "testPasswordHash" {
		t.Fatalf("expected password hash testPasswordHash, got %s", created.PasswordHash)
	}

	if created.UserTypeID != testUserTypeID {
		t.Fatalf("expected user type 1, got %v", created.UserTypeID)
	}
}

func TestPostgresRepositoryListUserAccounts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")

	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)

	userAccounts, err := repo.ListUserAccounts(context.Background())
	if err != nil {
		t.Fatalf("list user accounts: %v", err)
	}

	foundTestAccount := false
	for _, userAccount := range userAccounts {
		if userAccount.ID == 0 {
			t.Fatalf("expected user account id to be set, got 0")
		}
		if userAccount.Username == "" {
			t.Fatalf("expected user account %d to have a username", userAccount.ID)
		}
		if userAccount.PasswordHash == "" {
			t.Fatalf("expected user account %d to have a password hash", userAccount.ID)
		}
		if userAccount.UserTypeID == 0 {
			t.Fatalf("expected user account %d to have a user type id", userAccount.ID)
		}
		if userAccount.ID == testUserAccountID {
			foundTestAccount = true

			if userAccount.Username != "testUsername" {
				t.Fatalf("expected username %v, got %v instead", "testUsername", userAccount.Username)
			}

			if userAccount.PasswordHash != "testPasswordHash" {
				t.Fatalf("expected password hash %v, got %v instead", "testPasswordHash", userAccount.PasswordHash)
			}

			if userAccount.UserTypeID != testUserTypeID {
				t.Fatalf("expected user type id %v, got %v instead", testUserTypeID, userAccount.UserTypeID)
			}

			if userAccount.UserType != "test" {
				t.Fatalf("expected user type %v, got %v instead", "test", userAccount.UserType)
			}

		}
	}

	if !foundTestAccount {
		t.Fatalf("expected user account %d in list", testUserAccountID)
	}
}

func TestPostgresRepositoryGetUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")

	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)

	userAccount, err := repo.GetUserAccount(context.Background(), "testUsername")
	if err != nil {
		t.Fatalf("get user account %v", err)
	}

	if userAccount.ID != testUserAccountID {
		t.Fatalf("expected user account id %v, got %v instead", testUserAccountID, userAccount.ID)
	}

	if userAccount.Username != "testUsername" {
		t.Fatalf("expected username %v, got %v instead", "testUsername", userAccount.Username)
	}

	if userAccount.PasswordHash != "testPasswordHash" {
		t.Fatalf("expected password hash %v, got %v instead", "testPasswordHash", userAccount.PasswordHash)
	}

	if userAccount.UserTypeID != testUserTypeID {
		t.Fatalf("expected user type id %v, got %v instead", testUserTypeID, userAccount.UserTypeID)
	}

	if userAccount.UserType != "test" {
		t.Fatalf("expected user type %v, got %v instead", "test", userAccount.UserType)
	}
}

func TestPostgresRepositoryDeleteUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")

	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)

	err := repo.DeleteUserAccount(context.Background(), testUserAccountID)
	if err != nil {
		t.Fatalf("get user account %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM user_account WHERE id = $1", testUserAccountID); n != 0 {
		t.Fatalf("expected user account %d to be deleted, found %d rows", testUserAccountID, n)
	}
}

// endregion

// region Product

func TestPostgresRepositoryCreateProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	now := time.Now()

	product, err := domain.NewProduct("testProduct", 20, "This is a test product.", now, testUserAccountID, testCategoryID)
	if err != nil {
		t.Fatalf("create domain user account %v", err)
	}

	created, err := repo.CreateProduct(context.Background(), product)
	if err != nil {
		t.Fatalf("create product %v", err)
	}

	cleanupRows(t, db, "DELETE FROM product WHERE id = $1", created.ID)

	if created.ID == 0 {
		t.Fatalf("expected database to generate id, got 0 instead")
	}

	if created.Name != product.Name {
		t.Fatalf("expected product name %v, got %v instead", product.Name, created.Name)
	}

	if created.Price != product.Price {
		t.Fatalf("expected product price %v, got %v instead", product.Price, created.Price)
	}

	if created.Description != product.Description {
		t.Fatalf("expected product description %v, got %v instead", product.Description, created.Description)
	}

	if created.ListedAt.Round(1) != product.ListedAt.Round(1) {
		t.Fatalf("expected product listed at %v, got %v instead", product.ListedAt.Round(1), created.ListedAt.Round(1))
	}

	if created.UpdatedAt.Round(1) != product.UpdatedAt.Round(1) {
		t.Fatalf("expected product updated at %v, got %v instead", product.UpdatedAt.Round(1), created.UpdatedAt.Round(1))
	}

	if created.Sold != product.Sold {
		t.Fatalf("expected product sold status %v, got %v instead", product.Sold, created.Sold)
	}

	if created.SellerID != product.SellerID {
		t.Fatalf("expected product seller id %v, got %v instead", product.SellerID, created.SellerID)
	}

	if created.CategoryID != product.CategoryID {
		t.Fatalf("expected product category id %v, got %v instead", product.CategoryID, created.CategoryID)
	}
}

func TestPostgresRepositoryListUnsoldProducts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testUnsoldProduct := createTestProduct(t, db, "testUnsoldProduct", 20, "This is an unsold test product.", false, testUserAccountID, testCategoryID)
	testSoldProduct := createTestProduct(t, db, "testSoldProduct", 30, "This is a sold test product.", true, testUserAccountID, testCategoryID)

	products, err := repo.ListUnsoldProducts(context.Background())
	if err != nil {
		t.Fatalf("list unsold products: %v", err)
	}

	foundUnsoldProduct := false
	for _, product := range products {
		if product.ID == 0 {
			t.Fatalf("expected product id to be set, got 0")
		}
		if product.ID == testSoldProduct.ID {
			t.Fatalf("expected sold product %d not to be in list", testSoldProduct.ID)
		}
		if product.ID == testUnsoldProduct.ID {
			foundUnsoldProduct = true

			if product.Name != testUnsoldProduct.Name {
				t.Fatalf("expected product name %v, got %v instead", testUnsoldProduct.Name, product.Name)
			}

			if product.Price != testUnsoldProduct.Price {
				t.Fatalf("expected product price %v, got %v instead", testUnsoldProduct.Price, product.Price)
			}

			if product.SellerID != testUnsoldProduct.SellerID {
				t.Fatalf("expected product seller id %v, got %v instead", testUnsoldProduct.SellerID, product.SellerID)
			}

			if product.SellerUsername != "testUsername" {
				t.Fatalf("expected product seller username %v, got %v instead", "testUsername", product.SellerUsername)
			}

			if product.CategoryID != testUnsoldProduct.CategoryID {
				t.Fatalf("expected product category id %v, got %v instead", testUnsoldProduct.CategoryID, product.CategoryID)
			}

			if product.CategoryName != "testCategory" {
				t.Fatalf("expected product category name %v, got %v instead", "testCategory", product.CategoryName)
			}
		}
	}

	if !foundUnsoldProduct {
		t.Fatalf("expected unsold product %d in list", testUnsoldProduct.ID)
	}
}

func TestPostgresRepositoryListOwnProducts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testSellerID := createTestUserAccount(t, db, "testSeller", "testPasswordHash", testUserTypeID)
	testOtherSellerID := createTestUserAccount(t, db, "testOtherSeller", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testOwnProduct := createTestProduct(t, db, "testOwnProduct", 20, "This is the seller's own test product.", false, testSellerID, testCategoryID)
	testOtherProduct := createTestProduct(t, db, "testOtherProduct", 30, "This is another seller's test product.", false, testOtherSellerID, testCategoryID)

	products, err := repo.ListOwnProducts(context.Background(), testSellerID)
	if err != nil {
		t.Fatalf("list own products: %v", err)
	}

	foundOwnProduct := false
	for _, product := range products {
		if product.ID == 0 {
			t.Fatalf("expected product id to be set, got 0")
		}
		if product.ID == testOtherProduct.ID {
			t.Fatalf("expected other seller's product %d not to be in list", testOtherProduct.ID)
		}
		if product.ID == testOwnProduct.ID {
			foundOwnProduct = true

			if product.Name != testOwnProduct.Name {
				t.Fatalf("expected product name %v, got %v instead", testOwnProduct.Name, product.Name)
			}

			if product.Price != testOwnProduct.Price {
				t.Fatalf("expected product price %v, got %v instead", testOwnProduct.Price, product.Price)
			}

			if product.Description != testOwnProduct.Description {
				t.Fatalf("expected product description %v, got %v instead", testOwnProduct.Description, product.Description)
			}

			if product.ListedAt.Round(1) != testOwnProduct.ListedAt.Round(1) {
				t.Fatalf("expected product listed at %v, got %v instead", testOwnProduct.ListedAt.Round(1), product.ListedAt.Round(1))
			}

			if product.UpdatedAt.Round(1) != testOwnProduct.UpdatedAt.Round(1) {
				t.Fatalf("expected product updated at %v, got %v instead", testOwnProduct.UpdatedAt.Round(1), product.UpdatedAt.Round(1))
			}

			if product.Sold != testOwnProduct.Sold {
				t.Fatalf("expected product sold status %v, got %v instead", testOwnProduct.Sold, product.Sold)
			}

			if product.CategoryID != testOwnProduct.CategoryID {
				t.Fatalf("expected product category id %v, got %v instead", testOwnProduct.CategoryID, product.CategoryID)
			}

			if product.CategoryName != "testCategory" {
				t.Fatalf("expected product category name %v, got %v instead", "testCategory", product.CategoryName)
			}
		}
	}

	if !foundOwnProduct {
		t.Fatalf("expected own product %d in list", testOwnProduct.ID)
	}
}

func TestPostgresRepositoryListProductsInCart(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testSellerID := createTestUserAccount(t, db, "testSeller", "testPasswordHash", testUserTypeID)
	testCustomerID := createTestUserAccount(t, db, "testCustomer", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testCartProduct := createTestProduct(t, db, "testCartProduct", 20, "This is a test product in the cart.", false, testSellerID, testCategoryID)
	testOtherProduct := createTestProduct(t, db, "testOtherProduct", 30, "This is a test product not in the cart.", false, testSellerID, testCategoryID)

	createTestCartItem(t, db, testCustomerID, testCartProduct.ID)

	products, err := repo.ListProductsInCart(context.Background(), testCustomerID)
	if err != nil {
		t.Fatalf("list products in cart: %v", err)
	}

	foundCartProduct := false
	for _, product := range products {
		if product.ID == 0 {
			t.Fatalf("expected product id to be set, got 0")
		}
		if product.ID == testOtherProduct.ID {
			t.Fatalf("expected product %d not to be in cart", testOtherProduct.ID)
		}
		if product.ID == testCartProduct.ID {
			foundCartProduct = true

			if product.Name != testCartProduct.Name {
				t.Fatalf("expected product name %v, got %v instead", testCartProduct.Name, product.Name)
			}

			if product.Price != testCartProduct.Price {
				t.Fatalf("expected product price %v, got %v instead", testCartProduct.Price, product.Price)
			}

			if product.SellerID != testCartProduct.SellerID {
				t.Fatalf("expected product seller id %v, got %v instead", testCartProduct.SellerID, product.SellerID)
			}

			if product.SellerUsername != "testSeller" {
				t.Fatalf("expected product seller username %v, got %v instead", "testSeller", product.SellerUsername)
			}

			if product.CategoryID != testCartProduct.CategoryID {
				t.Fatalf("expected product category id %v, got %v instead", testCartProduct.CategoryID, product.CategoryID)
			}

			if product.CategoryName != "testCategory" {
				t.Fatalf("expected product category name %v, got %v instead", "testCategory", product.CategoryName)
			}
		}
	}

	if !foundCartProduct {
		t.Fatalf("expected cart product %d in list", testCartProduct.ID)
	}
}

func TestPostgresRepositoryGetProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testProduct := createTestProduct(t, db, "testProduct", 20, "This is a test product.", false, testUserAccountID, testCategoryID)

	product, err := repo.GetProduct(context.Background(), testProduct.ID)
	if err != nil {
		t.Fatalf("get product %v", err)
	}

	if product.ID != testProduct.ID {
		t.Fatalf("expected product id %v, got %v instead", testProduct.ID, product.ID)
	}

	if product.Name != testProduct.Name {
		t.Fatalf("expected product name %v, got %v instead", testProduct.Name, product.Name)
	}

	if product.Price != testProduct.Price {
		t.Fatalf("expected product price %v, got %v instead", testProduct.Price, product.Price)
	}

	if product.Description != testProduct.Description {
		t.Fatalf("expected product description %v, got %v instead", testProduct.Description, product.Description)
	}

	if product.ListedAt.Round(1) != testProduct.ListedAt.Round(1) {
		t.Fatalf("expected product listed at %v, got %v instead", testProduct.ListedAt.Round(1), product.ListedAt.Round(1))
	}

	if product.SellerID != testProduct.SellerID {
		t.Fatalf("expected product seller id %v, got %v instead", testProduct.SellerID, product.SellerID)
	}

	if product.SellerUsername != "testUsername" {
		t.Fatalf("expected product seller username %v, got %v instead", "testUsername", product.SellerUsername)
	}

	if product.CategoryID != testProduct.CategoryID {
		t.Fatalf("expected product category id %v, got %v instead", testProduct.CategoryID, product.CategoryID)
	}

	if product.CategoryName != "testCategory" {
		t.Fatalf("expected product category name %v, got %v instead", "testCategory", product.CategoryName)
	}
}

func TestPostgresRepositoryUpdateProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")
	testUpdatedCategoryID := createTestCategory(t, db, "testUpdatedCategory")

	testProduct := createTestProduct(t, db, "testProduct", 20, "This is a test product.", false, testUserAccountID, testCategoryID)

	product := testProduct
	product.Name = "testUpdatedProduct"
	product.Price = 25
	product.Description = "This is an updated test product."
	product.UpdatedAt = time.Now()
	product.CategoryID = testUpdatedCategoryID

	updated, err := repo.UpdateProduct(context.Background(), product)
	if err != nil {
		t.Fatalf("update product %v", err)
	}

	if updated.ID != testProduct.ID {
		t.Fatalf("expected product id %v, got %v instead", testProduct.ID, updated.ID)
	}

	if updated.Name != product.Name {
		t.Fatalf("expected product name %v, got %v instead", product.Name, updated.Name)
	}

	if updated.Price != product.Price {
		t.Fatalf("expected product price %v, got %v instead", product.Price, updated.Price)
	}

	if updated.Description != product.Description {
		t.Fatalf("expected product description %v, got %v instead", product.Description, updated.Description)
	}

	if updated.ListedAt.Round(1) != testProduct.ListedAt.Round(1) {
		t.Fatalf("expected product listed at %v, got %v instead", testProduct.ListedAt.Round(1), updated.ListedAt.Round(1))
	}

	if updated.UpdatedAt.Round(1) != product.UpdatedAt.Round(1) {
		t.Fatalf("expected product updated at %v, got %v instead", product.UpdatedAt.Round(1), updated.UpdatedAt.Round(1))
	}

	if updated.Sold != testProduct.Sold {
		t.Fatalf("expected product sold status %v, got %v instead", testProduct.Sold, updated.Sold)
	}

	if updated.CategoryID != testUpdatedCategoryID {
		t.Fatalf("expected product category id %v, got %v instead", testUpdatedCategoryID, updated.CategoryID)
	}

	if updated.CategoryName != "testUpdatedCategory" {
		t.Fatalf("expected product category name %v, got %v instead", "testUpdatedCategory", updated.CategoryName)
	}
}

func TestPostgresRepositoryDeleteProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testUserAccountID := createTestUserAccount(t, db, "testUsername", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testProduct := createTestProduct(t, db, "testProduct", 20, "This is a test product.", false, testUserAccountID, testCategoryID)

	err := repo.DeleteProduct(context.Background(), testProduct.ID)
	if err != nil {
		t.Fatalf("delete product %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM product WHERE id = $1", testProduct.ID); n != 0 {
		t.Fatalf("expected product %d to be deleted, found %d rows", testProduct.ID, n)
	}
}

// endregion

// region Category

func TestPostgresRepositoryListCategories(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testCategoryId := createTestCategory(t, db, "testCategory")

	categories, err := repo.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}

	foundTestCategory := false
	for _, category := range categories {
		if category.ID == 0 {
			t.Fatalf("expected category id to be set, got 0")
		}
		if category.Name == "" {
			t.Fatalf("expected category %d to have a name", category.ID)
		}
		if category.ID == testCategoryId {
			foundTestCategory = true
		}
	}

	if !foundTestCategory {
		t.Fatalf("expected category %d in list", testCategoryId)
	}

}

// endregion

// region Cart Item

func TestPostgresRepositoryCreateCartItem(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testSellerID := createTestUserAccount(t, db, "testSeller", "testPasswordHash", testUserTypeID)
	testCustomerID := createTestUserAccount(t, db, "testCustomer", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testProduct := createTestProduct(t, db, "testProduct", 20, "This is a test product.", false, testSellerID, testCategoryID)

	cartItem := domain.CartItem{
		CustomerID: testCustomerID,
		ProductID:  testProduct.ID,
	}

	created, err := repo.CreateCartItem(context.Background(), cartItem)
	if err != nil {
		t.Fatalf("create cart item %v", err)
	}

	cleanupRows(t, db, "DELETE FROM cart_item WHERE customer_id = $1 AND product_id = $2", created.CustomerID, created.ProductID)

	if created.CustomerID != testCustomerID {
		t.Fatalf("expected cart item customer id %v, got %v instead", testCustomerID, created.CustomerID)
	}

	if created.ProductID != testProduct.ID {
		t.Fatalf("expected cart item product id %v, got %v instead", testProduct.ID, created.ProductID)
	}
}

func TestPostgresRepositoryDeleteCartItem(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testSellerID := createTestUserAccount(t, db, "testSeller", "testPasswordHash", testUserTypeID)
	testCustomerID := createTestUserAccount(t, db, "testCustomer", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testProduct := createTestProduct(t, db, "testProduct", 20, "This is a test product.", false, testSellerID, testCategoryID)

	createTestCartItem(t, db, testCustomerID, testProduct.ID)

	err := repo.DeleteCartItem(context.Background(), testCustomerID, testProduct.ID)
	if err != nil {
		t.Fatalf("delete cart item %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1 AND product_id = $2", testCustomerID, testProduct.ID); n != 0 {
		t.Fatalf("expected cart item (%d, %d) to be deleted, found %d rows", testCustomerID, testProduct.ID, n)
	}
}

func TestPostgresRepositoryDeleteCartItems(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	testUserTypeID := createTestUserType(t, db, "test")
	testSellerID := createTestUserAccount(t, db, "testSeller", "testPasswordHash", testUserTypeID)
	testCustomerID := createTestUserAccount(t, db, "testCustomer", "testPasswordHash", testUserTypeID)
	testCategoryID := createTestCategory(t, db, "testCategory")

	testFirstProduct := createTestProduct(t, db, "testFirstProduct", 20, "This is the first test product.", false, testSellerID, testCategoryID)
	testSecondProduct := createTestProduct(t, db, "testSecondProduct", 30, "This is the second test product.", false, testSellerID, testCategoryID)

	createTestCartItem(t, db, testCustomerID, testFirstProduct.ID)
	createTestCartItem(t, db, testCustomerID, testSecondProduct.ID)

	err := repo.DeleteCartItems(context.Background(), testCustomerID)
	if err != nil {
		t.Fatalf("delete cart items %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1", testCustomerID); n != 0 {
		t.Fatalf("expected cart items for customer %d to be deleted, found %d rows", testCustomerID, n)
	}
}

// endregion

// region Helper

func testDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("STOREFRONT_DATABASE_URL") // export STOREFRONT_DATABASE_URL="postgres://storefront:storefront@localhost:5432/storefrontapp?sslmode=disable"
	if dsn == "" {
		t.Fatal("TODO_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to database: %v", err)
	}

	return db
}

// Cleanup rows runs DELETE commands when the tests finish, whether they pass or fail
func cleanupRows(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	t.Cleanup(func() {
		if _, err := db.ExecContext(context.Background(), query, args...); err != nil {
			t.Errorf("clean up test data (%s): %v", query, err)
		}
	})
}

func countRows(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()

	var n int
	if err := db.QueryRowContext(context.Background(), query, args...).Scan(&n); err != nil {
		t.Fatalf("count rows: %v", err)
	}

	return n
}

func createTestUserType(t *testing.T, db *sql.DB, userType string) int64 {
	t.Helper()

	var id int64
	err := db.QueryRowContext(
		context.Background(),
		"INSERT INTO user_type (type) VALUES ($1) RETURNING (id)",
		userType,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test user type: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM user_type WHERE id = $1", id)

	return id
}

func createTestUserAccount(t *testing.T, db *sql.DB, username string, passwordHash string, userTypeID int64) int64 {
	t.Helper()

	var id int64
	err := db.QueryRowContext(
		context.Background(),
		"INSERT INTO user_account (username, password_hash, user_type_id) VALUES ($1, $2, $3) RETURNING (id)",
		username,
		passwordHash,
		userTypeID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test user account: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM user_account WHERE id = $1", id)

	return id
}

func createTestCategory(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()

	var id int64
	err := db.QueryRowContext(
		context.Background(),
		"INSERT INTO category (name) VALUES ($1) RETURNING (id)",
		name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test category: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM category WHERE id = $1", id)

	return id
}

func createTestProduct(t *testing.T, db *sql.DB, name string, price float64, description string, sold bool, sellerID int64, categoryID int64) domain.Product {
	t.Helper()

	now := time.Now()

	var id int64
	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO product (name, price, description, listed_at, updated_at, sold, seller_id, category_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING (id)`,
		name,
		price,
		description,
		now,
		now,
		sold,
		sellerID,
		categoryID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test product: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM product WHERE id = $1", id)

	return domain.Product{
		ID:          id,
		Name:        name,
		Price:       price,
		Description: description,
		ListedAt:    now,
		UpdatedAt:   now,
		Sold:        sold,
		SellerID:    sellerID,
		CategoryID:  categoryID,
	}
}

func createTestCartItem(t *testing.T, db *sql.DB, customerID int64, productID int64) {
	t.Helper()

	_, err := db.ExecContext(
		context.Background(),
		"INSERT INTO cart_item (customer_id, product_id) VALUES ($1, $2)",
		customerID,
		productID,
	)
	if err != nil {
		t.Fatalf("insert test cart item: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM cart_item WHERE customer_id = $1 AND product_id = $2", customerID, productID)
}

// endregion
