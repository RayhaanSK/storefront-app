package storefront

import (
	"testing"
	"time"
)

// region UserType

func TestNewUserTypeRejectsBlankTitle(t *testing.T) {
	_, err := NewUserType("   ")

	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNewUserTypeTrimsTitle(t *testing.T) {
	userType, err := NewUserType("    customer  ")

	if err != nil {
		t.Fatalf("exepcted no error, got %v", err)
	}

	if userType.Type != "customer" {
		t.Fatalf("expected Type customer, got %s", userType.Type)
	}
}

// endregion

// region UserAccount

func TestNewUserAccountRejectsBlankTitle(t *testing.T) {
	_, err := NewUserAccount("   ", "abc", 1)
	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}

	_, err = NewUserAccount("abc", "   ", 1)
	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNewUserAccountTrimsTitle(t *testing.T) {
	userAccount, err := NewUserAccount("    john123  ", "   abc ", 1)
	if err != nil {
		t.Fatalf("exepcted no error, got %v", err)
	}

	if userAccount.Username != "john123" {
		t.Fatalf("expected Username john123, got %s", userAccount.Username)
	}

	if userAccount.PasswordHash != "abc" {
		t.Fatalf("expected PasswordHash abc, got %s", userAccount.PasswordHash)
	}
}

// endregion

// region Category

func TestNewCategoryRejectsBlankTitle(t *testing.T) {
	_, err := NewCategory("   ")
	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNewCategoryTrimsTitle(t *testing.T) {
	category, err := NewCategory("   books ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.Name != "books" {
		t.Fatalf("expected Name books, got %s", category.Name)
	}
}

// endregion

// region Product

func TestNewProductRejectsBlankTitle(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	_, err := NewProduct("   ", 123, "", now, 1, 1)
	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNewProductTrimsTitle(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	product, err := NewProduct("  chair ", 123, "", now, 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Name != "chair" {
		t.Fatalf("expected Name chair, got %s", product.Name)
	}
}

func TestUpdateProductRejectsBlankTitle(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	product, _ := NewProduct("chair", 123, "", now, 1, 1)

	err := product.UpdateProduct("   ", 123, "", now, 1)
	if err != ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestUpdateProductTrimsTitle(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	product, _ := NewProduct("chair", 123, "", now, 1, 1)

	product.UpdateProduct("  table ", 123, "", now, 1)
	if product.Name != "table" {
		t.Fatalf("expected Name table, got %v", product.Name)
	}
}

// endregion
