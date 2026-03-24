package repository

import (
	"context"
	"time"
	"todos/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at, updated_at"

	err := pool.QueryRow(
		ctx,
		query,
		user.Email, user.PasswordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, err
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User = &models.User{}

	var query string = "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1"
	err := pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(pool *pgxpool.Pool, id int) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User = &models.User{}

	var query string = "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1"
	err := pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
