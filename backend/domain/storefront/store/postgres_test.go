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

func TestPostgresRepositoryCreateUserAccount(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	userAccount, err := domain.NewUserAccount("  johnsmith123 ", " abc  ", 1)
	if err != nil {
		t.Fatalf("create domain user account %v", err)
	}

	created, err := repo.CreateUserAccount(context.Background(), userAccount)
	if err != nil {
		t.Fatalf("create user account %v", err)
	}

	if created.ID == 0 {
		t.Fatalf("expected database to generate id")
	}

	if created.Username != "johnsmith123" {
		t.Fatalf("expected username johnsmith123, got %s", created.Username)
	}

	if created.PasswordHash != "abc" {
		t.Fatalf("expected password hash abc, got %s", created.PasswordHash)
	}

	if created.UserTypeID != 1 {
		t.Fatalf("expected user type 1, got %v", created.UserTypeID)
	}

	_, err = db.ExecContext(
		context.Background(),
		"DELETE FROM user_account WHERE id = $1",
		created.ID,
	)

	if err != nil {
		t.Fatalf("clean up test user account: %v", err)
	}
}

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
