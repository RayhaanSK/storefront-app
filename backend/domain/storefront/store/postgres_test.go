package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	domain "storefrontapp/domain/storefront"
	storefrontservice "storefrontapp/services/storefront"
)

// These tests run against a real PostgreSQL database. Set STOREFRONT_DATABASE_URL first, e.g.
//   export STOREFRONT_DATABASE_URL="postgres://storefront:storefront@localhost:5432/storefrontapp?sslmode=disable"
//
// The user_type table is assumed to be seeded (id 1 must exist). Every other row a test needs is
// created by the test itself and removed again with t.Cleanup, so tests don't depend on each other.

const (
	testUserTypeID = 1

	// SERIAL / IDENTITY columns never generate negative ids, so this id never exists.
	missingID int64 = -1

	pgForeignKeyViolation = "23503"
	pgUniqueViolation     = "23505"
)

// region User Type

func TestPostgresRepositoryListUserTypes(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	userTypes, err := repo.ListUserTypes(context.Background())
	if err != nil {
		t.Fatalf("list user types: %v", err)
	}

	if len(userTypes) == 0 {
		t.Fatalf("expected seeded user types, got none")
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

	username := uniqueName("user")

	userAccount, err := domain.NewUserAccount("  "+username+" ", " abc  ", testUserTypeID)
	if err != nil {
		t.Fatalf("create domain user account: %v", err)
	}

	created, err := repo.CreateUserAccount(context.Background(), userAccount)
	if err != nil {
		t.Fatalf("create user account: %v", err)
	}
	cleanupUserAccount(t, db, created.ID)

	if created.ID == 0 {
		t.Fatalf("expected database to generate id")
	}

	if created.Username != username {
		t.Fatalf("expected username %s, got %s", username, created.Username)
	}

	if created.PasswordHash != "abc" {
		t.Fatalf("expected password hash abc, got %s", created.PasswordHash)
	}

	if created.UserTypeID != testUserTypeID {
		t.Fatalf("expected user type %d, got %v", testUserTypeID, created.UserTypeID)
	}
}

func TestPostgresRepositoryCreateUserAccountDuplicateUsername(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	existing := createTestUserAccount(t, db)

	duplicate := domain.UserAccount{
		Username:     existing.Username,
		PasswordHash: "different-hash",
		UserTypeID:   testUserTypeID,
	}

	created, err := repo.CreateUserAccount(context.Background(), duplicate)
	if err == nil {
		cleanupUserAccount(t, db, created.ID)
		t.Fatalf("expected error creating duplicate username %s", existing.Username)
	}

	requirePgErrorCode(t, err, pgUniqueViolation)
}

func TestPostgresRepositoryCreateUserAccountMissingUserType(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	userAccount := domain.UserAccount{
		Username:     uniqueName("user"),
		PasswordHash: "hash",
		UserTypeID:   missingID,
	}

	created, err := repo.CreateUserAccount(context.Background(), userAccount)
	if err == nil {
		cleanupUserAccount(t, db, created.ID)
		t.Fatalf("expected error creating user account with missing user type")
	}

	requirePgErrorCode(t, err, pgForeignKeyViolation)
}

func TestPostgresRepositoryListUserAccounts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	first := createTestUserAccount(t, db)
	second := createTestUserAccount(t, db)

	userAccounts, err := repo.ListUserAccounts(context.Background())
	if err != nil {
		t.Fatalf("list user accounts: %v", err)
	}

	for _, want := range []domain.UserAccount{first, second} {
		got, ok := findUserAccount(userAccounts, want.ID)
		if !ok {
			t.Fatalf("expected user account %d in list", want.ID)
		}
		assertUserAccount(t, want, got)
	}

	for i := 1; i < len(userAccounts); i++ {
		if userAccounts[i-1].ID > userAccounts[i].ID {
			t.Fatalf("expected user accounts ordered by id, got %d before %d", userAccounts[i-1].ID, userAccounts[i].ID)
		}
	}
}

func TestPostgresRepositoryGetUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	want := createTestUserAccount(t, db)

	got, err := repo.GetUserAccount(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("get user account: %v", err)
	}

	assertUserAccount(t, want, got)
}

func TestPostgresRepositoryGetUserAccountNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	_, err := repo.GetUserAccount(context.Background(), missingID)
	if !errors.Is(err, storefrontservice.ErrUserAccountNotFound) {
		t.Fatalf("expected ErrUserAccountNotFound, got %v", err)
	}
}

func TestPostgresRepositoryDeleteUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	userAccount := createTestUserAccount(t, db)

	if err := repo.DeleteUserAccount(context.Background(), userAccount.ID); err != nil {
		t.Fatalf("delete user account: %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM user_account WHERE id = $1", userAccount.ID); n != 0 {
		t.Fatalf("expected user account %d to be deleted, found %d rows", userAccount.ID, n)
	}
}

func TestPostgresRepositoryDeleteUserAccountNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	err := repo.DeleteUserAccount(context.Background(), missingID)
	if !errors.Is(err, storefrontservice.ErrUserAccountNotFound) {
		t.Fatalf("expected ErrUserAccountNotFound, got %v", err)
	}
}

// endregion

// region Product

func TestPostgresRepositoryCreateProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	now := testNow()

	product := domain.Product{
		Name:        uniqueName("product"),
		Price:       20,
		Description: "A test product",
		ListedAt:    now,
		UpdatedAt:   now,
		SellerID:    seller.ID,
		CategoryID:  categoryID,
	}

	created, err := repo.CreateProduct(context.Background(), product)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}
	cleanupProduct(t, db, created.ID)

	if created.ID == 0 {
		t.Fatalf("expected database to generate id")
	}
	if created.Name != product.Name {
		t.Fatalf("expected name %s, got %s", product.Name, created.Name)
	}
	if created.Price != product.Price {
		t.Fatalf("expected price %v, got %v", product.Price, created.Price)
	}
	if created.Description != product.Description {
		t.Fatalf("expected description %q, got %q", product.Description, created.Description)
	}
	if !created.ListedAt.Equal(now) {
		t.Fatalf("expected listed at %v, got %v", now, created.ListedAt)
	}
	if !created.UpdatedAt.Equal(now) {
		t.Fatalf("expected updated at %v, got %v", now, created.UpdatedAt)
	}
	if created.Sold {
		t.Fatalf("expected new product to default to unsold")
	}
	if created.SellerID != seller.ID {
		t.Fatalf("expected seller id %d, got %d", seller.ID, created.SellerID)
	}
	if created.CategoryID != categoryID {
		t.Fatalf("expected category id %d, got %d", categoryID, created.CategoryID)
	}
}

func TestPostgresRepositoryCreateProductMissingSeller(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	categoryID, _ := createTestCategory(t, db)
	now := testNow()

	product := domain.Product{
		Name:       uniqueName("product"),
		Price:      20,
		ListedAt:   now,
		UpdatedAt:  now,
		SellerID:   missingID,
		CategoryID: categoryID,
	}

	created, err := repo.CreateProduct(context.Background(), product)
	if err == nil {
		cleanupProduct(t, db, created.ID)
		t.Fatalf("expected error creating product with missing seller")
	}

	requirePgErrorCode(t, err, pgForeignKeyViolation)
}

func TestPostgresRepositoryCreateProductMissingCategory(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	now := testNow()

	product := domain.Product{
		Name:       uniqueName("product"),
		Price:      20,
		ListedAt:   now,
		UpdatedAt:  now,
		SellerID:   seller.ID,
		CategoryID: missingID,
	}

	created, err := repo.CreateProduct(context.Background(), product)
	if err == nil {
		cleanupProduct(t, db, created.ID)
		t.Fatalf("expected error creating product with missing category")
	}

	requirePgErrorCode(t, err, pgForeignKeyViolation)
}

func TestPostgresRepositoryListUnsoldProducts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	categoryID, categoryName := createTestCategory(t, db)

	// Two unsold products, so the test fails if only the first row is read.
	unsoldA := createTestProduct(t, db, seller.ID, categoryID, false)
	unsoldB := createTestProduct(t, db, seller.ID, categoryID, false)
	sold := createTestProduct(t, db, seller.ID, categoryID, true)

	products, err := repo.ListUnsoldProducts(context.Background())
	if err != nil {
		t.Fatalf("list unsold products: %v", err)
	}

	for _, want := range []domain.Product{unsoldA, unsoldB} {
		got, ok := findProduct(products, want.ID)
		if !ok {
			t.Fatalf("expected unsold product %d in list", want.ID)
		}

		if got.Name != want.Name {
			t.Fatalf("expected name %s, got %s", want.Name, got.Name)
		}
		if got.Price != want.Price {
			t.Fatalf("expected price %v, got %v", want.Price, got.Price)
		}
		if got.SellerID != seller.ID {
			t.Fatalf("expected seller id %d, got %d", seller.ID, got.SellerID)
		}
		if got.SellerUsername != seller.Username {
			t.Fatalf("expected seller username %s, got %s", seller.Username, got.SellerUsername)
		}
		if got.CategoryID != categoryID {
			t.Fatalf("expected category id %d, got %d", categoryID, got.CategoryID)
		}
		if got.CategoryName != categoryName {
			t.Fatalf("expected category name %s, got %s", categoryName, got.CategoryName)
		}
	}

	if _, ok := findProduct(products, sold.ID); ok {
		t.Fatalf("expected sold product %d to be excluded", sold.ID)
	}
}

func TestPostgresRepositoryListOwnProducts(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	otherSeller := createTestUserAccount(t, db)
	categoryID, categoryName := createTestCategory(t, db)

	ownA := createTestProduct(t, db, seller.ID, categoryID, false)
	ownB := createTestProduct(t, db, seller.ID, categoryID, true) // sold products still belong to the seller
	createTestProduct(t, db, otherSeller.ID, categoryID, false)

	products, err := repo.ListOwnProducts(context.Background(), seller.ID)
	if err != nil {
		t.Fatalf("list own products: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products for seller %d, got %d", seller.ID, len(products))
	}

	for _, want := range []domain.Product{ownA, ownB} {
		got, ok := findProduct(products, want.ID)
		if !ok {
			t.Fatalf("expected product %d in seller's list", want.ID)
		}

		if got.Name != want.Name {
			t.Fatalf("expected name %s, got %s", want.Name, got.Name)
		}
		if got.Price != want.Price {
			t.Fatalf("expected price %v, got %v", want.Price, got.Price)
		}
		if got.Description != want.Description {
			t.Fatalf("expected description %q, got %q", want.Description, got.Description)
		}
		if !got.ListedAt.Equal(want.ListedAt) {
			t.Fatalf("expected listed at %v, got %v", want.ListedAt, got.ListedAt)
		}
		if !got.UpdatedAt.Equal(want.UpdatedAt) {
			t.Fatalf("expected updated at %v, got %v", want.UpdatedAt, got.UpdatedAt)
		}
		if got.Sold != want.Sold {
			t.Fatalf("expected sold %v, got %v", want.Sold, got.Sold)
		}
		if got.CategoryID != categoryID {
			t.Fatalf("expected category id %d, got %d", categoryID, got.CategoryID)
		}
		if got.CategoryName != categoryName {
			t.Fatalf("expected category name %s, got %s", categoryName, got.CategoryName)
		}
	}
}

func TestPostgresRepositoryListOwnProductsNone(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)

	products, err := repo.ListOwnProducts(context.Background(), seller.ID)
	if err != nil {
		t.Fatalf("list own products: %v", err)
	}

	if len(products) != 0 {
		t.Fatalf("expected no products, got %d", len(products))
	}
}

func TestPostgresRepositoryListProductsInCart(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	customer := createTestUserAccount(t, db)
	otherCustomer := createTestUserAccount(t, db)
	categoryID, categoryName := createTestCategory(t, db)

	inCartA := createTestProduct(t, db, seller.ID, categoryID, false)
	inCartB := createTestProduct(t, db, seller.ID, categoryID, false)
	inOtherCart := createTestProduct(t, db, seller.ID, categoryID, false)

	createTestCartItem(t, db, customer.ID, inCartA.ID)
	createTestCartItem(t, db, customer.ID, inCartB.ID)
	createTestCartItem(t, db, otherCustomer.ID, inOtherCart.ID)

	products, err := repo.ListProductsInCart(context.Background(), customer.ID)
	if err != nil {
		t.Fatalf("list products in cart: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products in cart, got %d", len(products))
	}

	for _, want := range []domain.Product{inCartA, inCartB} {
		got, ok := findProduct(products, want.ID)
		if !ok {
			t.Fatalf("expected product %d in cart", want.ID)
		}

		if got.Name != want.Name {
			t.Fatalf("expected name %s, got %s", want.Name, got.Name)
		}
		if got.Price != want.Price {
			t.Fatalf("expected price %v, got %v", want.Price, got.Price)
		}
		if got.SellerID != seller.ID {
			t.Fatalf("expected seller id %d, got %d", seller.ID, got.SellerID)
		}
		if got.SellerUsername != seller.Username {
			t.Fatalf("expected seller username %s, got %s", seller.Username, got.SellerUsername)
		}
		if got.CategoryID != categoryID {
			t.Fatalf("expected category id %d, got %d", categoryID, got.CategoryID)
		}
		if got.CategoryName != categoryName {
			t.Fatalf("expected category name %s, got %s", categoryName, got.CategoryName)
		}
	}
}

func TestPostgresRepositoryListProductsInCartEmpty(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	customer := createTestUserAccount(t, db)

	products, err := repo.ListProductsInCart(context.Background(), customer.ID)
	if err != nil {
		t.Fatalf("list products in cart: %v", err)
	}

	if len(products) != 0 {
		t.Fatalf("expected empty cart, got %d products", len(products))
	}
}

func TestPostgresRepositoryGetProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	categoryID, categoryName := createTestCategory(t, db)
	want := createTestProduct(t, db, seller.ID, categoryID, false)

	got, err := repo.GetProduct(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("get product: %v", err)
	}

	if got.ID != want.ID {
		t.Fatalf("expected id %d, got %d", want.ID, got.ID)
	}
	if got.Name != want.Name {
		t.Fatalf("expected name %s, got %s", want.Name, got.Name)
	}
	if got.Price != want.Price {
		t.Fatalf("expected price %v, got %v", want.Price, got.Price)
	}
	if got.Description != want.Description {
		t.Fatalf("expected description %q, got %q", want.Description, got.Description)
	}
	if !got.ListedAt.Equal(want.ListedAt) {
		t.Fatalf("expected listed at %v, got %v", want.ListedAt, got.ListedAt)
	}
	if got.SellerID != seller.ID {
		t.Fatalf("expected seller id %d, got %d", seller.ID, got.SellerID)
	}
	if got.SellerUsername != seller.Username {
		t.Fatalf("expected seller username %s, got %s", seller.Username, got.SellerUsername)
	}
	if got.CategoryID != categoryID {
		t.Fatalf("expected category id %d, got %d", categoryID, got.CategoryID)
	}
	if got.CategoryName != categoryName {
		t.Fatalf("expected category name %s, got %s", categoryName, got.CategoryName)
	}
}

func TestPostgresRepositoryGetProductNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	_, err := repo.GetProduct(context.Background(), missingID)
	if !errors.Is(err, storefrontservice.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestPostgresRepositoryUpdateProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	oldCategoryID, _ := createTestCategory(t, db)
	newCategoryID, newCategoryName := createTestCategory(t, db)
	original := createTestProduct(t, db, seller.ID, oldCategoryID, false)

	updatedAt := original.UpdatedAt.Add(time.Hour)

	// Only the fields the service would send: the repository must fill in the rest from the database.
	changes := domain.Product{
		ID:          original.ID,
		Name:        original.Name + "_updated",
		Price:       35,
		Description: "An updated description",
		UpdatedAt:   updatedAt,
		CategoryID:  newCategoryID,
	}

	updated, err := repo.UpdateProduct(context.Background(), changes)
	if err != nil {
		t.Fatalf("update product: %v", err)
	}

	if updated.ID != original.ID {
		t.Fatalf("expected id %d, got %d", original.ID, updated.ID)
	}
	if updated.Name != changes.Name {
		t.Fatalf("expected name %s, got %s", changes.Name, updated.Name)
	}
	if updated.Price != changes.Price {
		t.Fatalf("expected price %v, got %v", changes.Price, updated.Price)
	}
	if updated.Description != changes.Description {
		t.Fatalf("expected description %q, got %q", changes.Description, updated.Description)
	}
	if !updated.ListedAt.Equal(original.ListedAt) {
		t.Fatalf("expected listed at to stay %v, got %v", original.ListedAt, updated.ListedAt)
	}
	if !updated.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("expected updated at %v, got %v", updatedAt, updated.UpdatedAt)
	}
	if updated.Sold != original.Sold {
		t.Fatalf("expected sold to stay %v, got %v", original.Sold, updated.Sold)
	}
	if updated.CategoryID != newCategoryID {
		t.Fatalf("expected category id %d, got %d", newCategoryID, updated.CategoryID)
	}
	if updated.CategoryName != newCategoryName {
		t.Fatalf("expected category name %s, got %s", newCategoryName, updated.CategoryName)
	}

	// Check the change was actually persisted, not just echoed back.
	var storedName string
	var storedCategoryID int64
	err = db.QueryRowContext(
		context.Background(),
		"SELECT name, category_id FROM product WHERE id = $1",
		original.ID,
	).Scan(&storedName, &storedCategoryID)
	if err != nil {
		t.Fatalf("read updated product: %v", err)
	}
	if storedName != changes.Name || storedCategoryID != newCategoryID {
		t.Fatalf("expected stored name %s and category %d, got %s and %d", changes.Name, newCategoryID, storedName, storedCategoryID)
	}
}

func TestPostgresRepositoryUpdateProductNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	categoryID, _ := createTestCategory(t, db)

	changes := domain.Product{
		ID:         missingID,
		Name:       "does not matter",
		Price:      20,
		UpdatedAt:  testNow(),
		CategoryID: categoryID,
	}

	_, err := repo.UpdateProduct(context.Background(), changes)
	if !errors.Is(err, storefrontservice.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestPostgresRepositoryDeleteProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	product := createTestProduct(t, db, seller.ID, categoryID, false)

	if err := repo.DeleteProduct(context.Background(), product.ID); err != nil {
		t.Fatalf("delete product: %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM product WHERE id = $1", product.ID); n != 0 {
		t.Fatalf("expected product %d to be deleted, found %d rows", product.ID, n)
	}
}

func TestPostgresRepositoryDeleteProductNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	err := repo.DeleteProduct(context.Background(), missingID)
	if !errors.Is(err, storefrontservice.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

// endregion

// region Cart Item

func TestPostgresRepositoryCreateCartItem(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	customer := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	product := createTestProduct(t, db, seller.ID, categoryID, false)

	cartItem := domain.CartItem{CustomerID: customer.ID, ProductID: product.ID}

	created, err := repo.CreateCartItem(context.Background(), cartItem)
	if err != nil {
		t.Fatalf("create cart item: %v", err)
	}
	cleanupCartItem(t, db, customer.ID, product.ID)

	if created.CustomerID != customer.ID {
		t.Fatalf("expected customer id %d, got %d", customer.ID, created.CustomerID)
	}
	if created.ProductID != product.ID {
		t.Fatalf("expected product id %d, got %d", product.ID, created.ProductID)
	}

	n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1 AND product_id = $2", customer.ID, product.ID)
	if n != 1 {
		t.Fatalf("expected 1 stored cart item, found %d", n)
	}
}

func TestPostgresRepositoryCreateCartItemMissingCustomer(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	product := createTestProduct(t, db, seller.ID, categoryID, false)

	_, err := repo.CreateCartItem(context.Background(), domain.CartItem{CustomerID: missingID, ProductID: product.ID})
	if err == nil {
		cleanupCartItem(t, db, missingID, product.ID)
		t.Fatalf("expected error creating cart item with missing customer")
	}

	requirePgErrorCode(t, err, pgForeignKeyViolation)
}

func TestPostgresRepositoryCreateCartItemMissingProduct(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	customer := createTestUserAccount(t, db)

	_, err := repo.CreateCartItem(context.Background(), domain.CartItem{CustomerID: customer.ID, ProductID: missingID})
	if err == nil {
		cleanupCartItem(t, db, customer.ID, missingID)
		t.Fatalf("expected error creating cart item with missing product")
	}

	requirePgErrorCode(t, err, pgForeignKeyViolation)
}

func TestPostgresRepositoryDeleteCartItem(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	customer := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	toDelete := createTestProduct(t, db, seller.ID, categoryID, false)
	toKeep := createTestProduct(t, db, seller.ID, categoryID, false)

	createTestCartItem(t, db, customer.ID, toDelete.ID)
	createTestCartItem(t, db, customer.ID, toKeep.ID)

	if err := repo.DeleteCartItem(context.Background(), customer.ID, toDelete.ID); err != nil {
		t.Fatalf("delete cart item: %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1 AND product_id = $2", customer.ID, toDelete.ID); n != 0 {
		t.Fatalf("expected cart item to be deleted, found %d rows", n)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1 AND product_id = $2", customer.ID, toKeep.ID); n != 1 {
		t.Fatalf("expected other cart item to remain, found %d rows", n)
	}
}

func TestPostgresRepositoryDeleteCartItemNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	customer := createTestUserAccount(t, db)

	err := repo.DeleteCartItem(context.Background(), customer.ID, missingID)
	if !errors.Is(err, storefrontservice.ErrCartItemNotFound) {
		t.Fatalf("expected ErrCartItemNotFound, got %v", err)
	}
}

func TestPostgresRepositoryDeleteCartItems(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	seller := createTestUserAccount(t, db)
	customer := createTestUserAccount(t, db)
	otherCustomer := createTestUserAccount(t, db)
	categoryID, _ := createTestCategory(t, db)
	productA := createTestProduct(t, db, seller.ID, categoryID, false)
	productB := createTestProduct(t, db, seller.ID, categoryID, false)

	createTestCartItem(t, db, customer.ID, productA.ID)
	createTestCartItem(t, db, customer.ID, productB.ID)
	createTestCartItem(t, db, otherCustomer.ID, productA.ID)

	if err := repo.DeleteCartItems(context.Background(), customer.ID); err != nil {
		t.Fatalf("delete cart items: %v", err)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1", customer.ID); n != 0 {
		t.Fatalf("expected customer's cart to be empty, found %d rows", n)
	}

	if n := countRows(t, db, "SELECT COUNT(*) FROM cart_item WHERE customer_id = $1", otherCustomer.ID); n != 1 {
		t.Fatalf("expected other customer's cart to be untouched, found %d rows", n)
	}
}

func TestPostgresRepositoryDeleteCartItemsNotFound(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	customer := createTestUserAccount(t, db)

	err := repo.DeleteCartItems(context.Background(), customer.ID)
	if !errors.Is(err, storefrontservice.ErrCartItemNotFound) {
		t.Fatalf("expected ErrCartItemNotFound, got %v", err)
	}
}

// endregion

// region Helpers

func testDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("STOREFRONT_DATABASE_URL")
	if dsn == "" {
		t.Fatal("STOREFRONT_DATABASE_URL is not set")
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

var nameCounter atomic.Int64

// uniqueName returns a name that won't clash with rows from other tests or earlier runs.
func uniqueName(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), nameCounter.Add(1))
}

// testNow returns the current time at PostgreSQL's microsecond precision so it compares equal after a round trip.
func testNow() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}

// Fixtures are inserted with plain SQL rather than through the repository,
// so a bug in one repository function can't break the setup of another function's test.
// t.Cleanup runs in reverse order, so dependent rows are removed before the rows they reference.

func createTestUserAccount(t *testing.T, db *sql.DB) domain.UserAccount {
	t.Helper()

	userAccount := domain.UserAccount{
		Username:     uniqueName("user"),
		PasswordHash: "test-hash",
		UserTypeID:   testUserTypeID,
	}

	err := db.QueryRowContext(
		context.Background(),
		"INSERT INTO user_account (username, password_hash, user_type_id) VALUES ($1, $2, $3) RETURNING id",
		userAccount.Username,
		userAccount.PasswordHash,
		userAccount.UserTypeID,
	).Scan(&userAccount.ID)
	if err != nil {
		t.Fatalf("insert test user account: %v", err)
	}

	cleanupUserAccount(t, db, userAccount.ID)

	return userAccount
}

func createTestCategory(t *testing.T, db *sql.DB) (int64, string) {
	t.Helper()

	name := uniqueName("category")

	var id int64
	err := db.QueryRowContext(
		context.Background(),
		"INSERT INTO category (name) VALUES ($1) RETURNING id",
		name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test category: %v", err)
	}

	cleanupRows(t, db, "DELETE FROM category WHERE id = $1", id)

	return id, name
}

func createTestProduct(t *testing.T, db *sql.DB, sellerID int64, categoryID int64, sold bool) domain.Product {
	t.Helper()

	now := testNow()

	product := domain.Product{
		Name:        uniqueName("product"),
		Price:       20,
		Description: "A test product",
		ListedAt:    now,
		UpdatedAt:   now,
		Sold:        sold,
		SellerID:    sellerID,
		CategoryID:  categoryID,
	}

	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO product (name, price, description, listed_at, updated_at, sold, seller_id, category_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		product.Name,
		product.Price,
		product.Description,
		product.ListedAt,
		product.UpdatedAt,
		product.Sold,
		product.SellerID,
		product.CategoryID,
	).Scan(&product.ID)
	if err != nil {
		t.Fatalf("insert test product: %v", err)
	}

	cleanupProduct(t, db, product.ID)

	return product
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

	cleanupCartItem(t, db, customerID, productID)
}

func cleanupUserAccount(t *testing.T, db *sql.DB, id int64) {
	cleanupRows(t, db, "DELETE FROM user_account WHERE id = $1", id)
}

func cleanupProduct(t *testing.T, db *sql.DB, id int64) {
	cleanupRows(t, db, "DELETE FROM product WHERE id = $1", id)
}

func cleanupCartItem(t *testing.T, db *sql.DB, customerID int64, productID int64) {
	cleanupRows(t, db, "DELETE FROM cart_item WHERE customer_id = $1 AND product_id = $2", customerID, productID)
}

// cleanupRows registers a DELETE to run when the test finishes, whether it passed or failed.
// Deleting a row that is already gone (e.g. the test deleted it) is not an error.
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

// requirePgErrorCode checks that err wraps a PostgreSQL error with the given SQLSTATE code,
// so a test for e.g. a foreign key violation doesn't pass because of an unrelated SQL error.
func requirePgErrorCode(t *testing.T, err error, code string) {
	t.Helper()

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected PostgreSQL error %s, got %v", code, err)
	}

	if pgErr.Code != code {
		t.Fatalf("expected PostgreSQL error %s, got %s: %s", code, pgErr.Code, pgErr.Message)
	}
}

func assertUserAccount(t *testing.T, want domain.UserAccount, got domain.UserAccount) {
	t.Helper()

	if got.ID != want.ID {
		t.Fatalf("expected id %d, got %d", want.ID, got.ID)
	}
	if got.Username != want.Username {
		t.Fatalf("expected username %s, got %s", want.Username, got.Username)
	}
	if got.PasswordHash != want.PasswordHash {
		t.Fatalf("expected password hash %s, got %s", want.PasswordHash, got.PasswordHash)
	}
	if got.UserTypeID != want.UserTypeID {
		t.Fatalf("expected user type id %v, got %v", want.UserTypeID, got.UserTypeID)
	}
}

func findUserAccount(userAccounts []domain.UserAccount, id int64) (domain.UserAccount, bool) {
	for _, userAccount := range userAccounts {
		if userAccount.ID == id {
			return userAccount, true
		}
	}
	return domain.UserAccount{}, false
}

func findProduct(products []domain.Product, id int64) (domain.Product, bool) {
	for _, product := range products {
		if product.ID == id {
			return product, true
		}
	}
	return domain.Product{}, false
}

// endregion
