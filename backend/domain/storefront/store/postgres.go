package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "storefrontapp/domain/storefront"
	storefrontservice "storefrontapp/services/storefront"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a PostgreSQL STOREFRONT repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// region User Type

func (r *PostgresRepository) ListUserTypes(ctx context.Context) ([]domain.UserType, error) {
	const query = `
		SELECT (user_type.id, user_type.type)
		FROM user_type;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list user types: %w", err)
	}
	defer rows.Close()

	var userTypes []domain.UserType

	for rows.Next() {
		var userType domain.UserType

		if err := rows.Scan(
			&userType.ID,
			&userType.Type,
		); err != nil {
			return nil, fmt.Errorf("scan user type: %w", err)
		}

		userTypes = append(userTypes, userType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user types: %w", err)
	}

	return userTypes, nil
}

// endregion

// region User Account

func (r *PostgresRepository) CreateUserAccount(ctx context.Context, userAccount domain.UserAccount) (domain.UserAccount, error) {
	const query = `
		INSERT INTO user_account (username, password_hash, user_type_id)
		VALUES ($1, $2, $3)
		RETURNING id, username, password_hash, user_type_id;
	`

	var created domain.UserAccount

	err := r.db.QueryRowContext(
		ctx,
		query,
		userAccount.Username,
		userAccount.PasswordHash,
		userAccount.UserTypeID,
	).Scan(
		&created.ID,
		&created.Username,
		&created.PasswordHash,
		&created.UserTypeID,
	)

	if err != nil {
		return domain.UserAccount{}, fmt.Errorf("create user account: %w", err)
	}

	return created, nil
}

func (r *PostgresRepository) ListUserAccounts(ctx context.Context) ([]domain.UserAccount, error) {
	const query = `
		SELECT user_account.id, username, password_hash, user_type_id, type
		FROM user_account
		INNER JOIN user_type ON user_type_id = user_type.id
		ORDER BY user_account.id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list user accounts: %w", err)
	}
	defer rows.Close()

	var userAccounts []domain.UserAccount

	for rows.Next() {
		var userAccount domain.UserAccount

		if err := rows.Scan(
			&userAccount.ID,
			&userAccount.Username,
			&userAccount.PasswordHash,
			&userAccount.UserTypeID,
		); err != nil {
			return nil, fmt.Errorf("scan user account: %w", err)
		}

		userAccounts = append(userAccounts, userAccount)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user accounts: %w", err)
	}

	return userAccounts, nil
}

func (r *PostgresRepository) GetUserAccount(ctx context.Context, id int64) (domain.UserAccount, error) {
	const query = `
		SELECT user_account.id, username, password_hash, user_type_id, type
		FROM user_account
		INNER JOIN user_type ON user_type_id = user_type.id
		WHERE user_account.id = $1;
	`

	var userAccount domain.UserAccount

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&userAccount.ID,
		&userAccount.Username,
		&userAccount.PasswordHash,
		&userAccount.UserTypeID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UserAccount{}, storefrontservice.ErrUserAccountNotFound
		}

		return domain.UserAccount{}, fmt.Errorf("get user account %w", err)
	}

	return userAccount, nil
}

func (r *PostgresRepository) DeleteUserAccount(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM user_account
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user account: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted user account: %w", err)
	}

	if count == 0 {
		return storefrontservice.ErrUserAccountNotFound
	}

	return nil
}

// endregion

// region Product

func (r *PostgresRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	const query = `
		INSERT INTO product (name, price, description, listed_at, updated_at, seller_id, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, price, description, listed_at, updated_at, sold, seller_id, category_id;
	`

	var created domain.Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Description,
		product.ListedAt,
		product.UpdatedAt,
		product.SellerID,
		product.CategoryID,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Price,
		&created.Description,
		&created.ListedAt,
		&created.UpdatedAt,
		&created.Sold,
		&created.SellerID,
		&created.CategoryID,
	)

	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}

	return created, nil
}

func (r *PostgresRepository) ListUnsoldProducts(ctx context.Context) ([]domain.Product, error) {
	const query = `
		SELECT (product.id, product.name, product.price, product.seller_id, user_account.username, product.category_id, category.name)
		FROM product
		INNER JOIN user_account ON seller_id = user_account.id
		INNER JOIN category ON category_id = category.id
		WHERE sold = FALSE;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list unsold products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product

	if rows.Next() {
		var product domain.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.SellerID,
			&product.SellerUsername,
			&product.CategoryID,
			&product.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("scan unsold products: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate unsold products: %w", err)
	}

	return products, nil
}

func (r *PostgresRepository) ListOwnProducts(ctx context.Context, sellerId int64) ([]domain.Product, error) {
	const query = `
		SELECT (product.id, product.name, product.price, product.description, product.listed_at, product.updated_at, product.sold, product.category_id, category.name)
		FROM product
		INNER JOIN category ON product.category_id = category.id;
		WHERE product.seller_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list own products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product

	if rows.Next() {
		var product domain.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.ListedAt,
			&product.UpdatedAt,
			&product.Sold,
			&product.CategoryID,
			&product.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("scan own products: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate own products: %w", err)
	}

	return products, nil
}

func (r *PostgresRepository) ListProductsInCart(ctx context.Context, customerId int64) ([]domain.Product, error) {
	const query = `
		SELECT (product.id, product.name, product.price, product.seller_id, user_account.username, product.category_id, category.name)
		FROM cart_item
		INNER JOIN product ON cart_item.product_id = product.id
		WHERE cart_item.customer_id = $1;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list cart products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product

	if rows.Next() {
		var product domain.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.SellerID,
			&product.SellerUsername,
			&product.CategoryID,
			&product.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("scan cart products: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cart products: %w", err)
	}

	return products, nil
}

func (r *PostgresRepository) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	const query = `
		SELECT (product.id, product.name, product.price, product.description, product.listed_at, product.seller_id, user_account.username, product.category_id, category.name)
		FROM product
		WHERE product.id = $1
		INNER JOIN category ON product.category_id = category.id;
	`

	var product domain.Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.ListedAt,
		&product.SellerID,
		&product.SellerUsername,
		&product.CategoryID,
		&product.CategoryName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, storefrontservice.ErrProductNotFound
		}

		return domain.Product{}, fmt.Errorf("get product %w", err)
	}

	return product, nil
}

func (r *PostgresRepository) UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	const query = `
		UPDATE product
		SET product.name = $1,
			product.price = $2,
			product.description = $3,
			product.updated_at = $4,
			product.category_id = $5
		FROM category
		WHERE product.id = $6
		  AND product.category_id = category.id
		RETURNING product.id, product.name, product.price, product.description, product.listed_at, product.updated_at, product.sold, product.category_id, category.name;
	`

	var updatedProduct domain.Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Description,
		product.UpdatedAt,
		product.CategoryID,
		product.ID,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Description,
		&updatedProduct.ListedAt,
		&updatedProduct.UpdatedAt,
		&updatedProduct.Sold,
		&updatedProduct.CategoryID,
		&updatedProduct.CategoryName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, storefrontservice.ErrProductNotFound
		}

		return domain.Product{}, fmt.Errorf("update product %w", err)
	}

	return product, nil
}

func (r *PostgresRepository) DeleteProduct(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM product
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted product: %w", err)
	}

	if count == 0 {
		return storefrontservice.ErrProductNotFound
	}

	return nil
}

// endregion

// region Cart Item

func (r *PostgresRepository) CreateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {
	const query = `
		INSERT INTO cart_item (customer_id, product_id)
		VALUES ($1, $2)
		RETURNING customer_id, product_id;
	`

	var created domain.CartItem

	err := r.db.QueryRowContext(
		ctx,
		query,
		cartItem.CustomerID,
		cartItem.ProductID,
	).Scan(
		&created.CustomerID,
		&created.ProductID,
	)

	if err != nil {
		return domain.CartItem{}, fmt.Errorf("create cart item: %w", err)
	}

	return created, nil
}

func (r *PostgresRepository) DeleteCartItem(ctx context.Context, customerId int64, productId int64) error {
	const query = `
		DELETE FROM cart_item
		WHERE customer_id = $1
		  AND product_id = $2;
	`

	result, err := r.db.ExecContext(ctx, query, customerId, productId)
	if err != nil {
		return fmt.Errorf("delete cart item: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted cart item: %w", err)
	}

	if count == 0 {
		return storefrontservice.ErrCartItemNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteCartItems(ctx context.Context, customerId int64) error {
	const query = `
		DELETE FROM cart_item
		WHERE customer_id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, customerId)
	if err != nil {
		return fmt.Errorf("delete cart items: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted cart items: %w", err)
	}

	if count == 0 {
		return storefrontservice.ErrCartItemNotFound
	}

	return nil
}

// endregion
