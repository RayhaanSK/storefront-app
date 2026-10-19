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

// endregion

// region Cart Item

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

// endregion
