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

// NewPostgresRepository constructs a PostgreSQL TODO repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

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
