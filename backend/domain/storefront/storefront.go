package storefront

import (
	"errors"
	"strings"
	"time"
)

var ErrTitleRequired = errors.New("title is required; empty entries are not allowed.")

// region UserType

type UserType struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func NewUserType(typeName string) (UserType, error) {
	typeName = strings.TrimSpace(typeName)

	if typeName == "" {
		return UserType{}, ErrTitleRequired
	}

	return UserType{
		Type: typeName,
	}, nil
}

// endregion

// region UserAccount

type UserAccount struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	UserTypeID   int64  `json:"userTypeId"`
}

func NewUserAccount(username string, passwordHash string, userTypeId int64) (UserAccount, error) {
	username = strings.TrimSpace(username)

	if username == "" {
		return UserAccount{}, ErrTitleRequired
	}

	passwordHash = strings.TrimSpace(passwordHash)

	if passwordHash == "" {
		return UserAccount{}, ErrTitleRequired
	}

	return UserAccount{
		Username:     username,
		PasswordHash: passwordHash,
		UserTypeID:   userTypeId,
	}, nil
}

// endregion

// region Category

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewCategory(name string) (Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Category{}, ErrTitleRequired
	}

	return Category{
		Name: name,
	}, nil
}

// endregion

// region Product

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	ListedAt    time.Time `json:"listedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Sold        bool      `json:"sold"`
	SellerID    int64     `json:"sellerId"`
	CategoryID  int64     `json:"categoryID"`
}

func NewProduct(name string, price float64, description string, now time.Time, sellerId int64, categoryId int64) (Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Product{}, ErrTitleRequired
	}

	return Product{
		Name:        name,
		Price:       price,
		Description: description,
		ListedAt:    now,
		UpdatedAt:   now,
		SellerID:    sellerId,
		CategoryID:  categoryId,
	}, nil
}

func (p *Product) UpdateProduct(name string, price float64, description string, now time.Time, categoryId int64) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrTitleRequired
	}

	p.Name = name
	p.Price = price
	p.Description = description
	p.UpdatedAt = now
	p.CategoryID = categoryId

	return nil
}

// endregion

// region CartItem

type CartItem struct {
	CustomerID int64 `json:"customerId"`
	ProductID  int64 `json:"productId"`
}

func NewCartItem(customerId int64, productId int64) (CartItem, error) {
	return CartItem{
		CustomerID: customerId,
		ProductID:  productId,
	}, nil
}

// endregion
