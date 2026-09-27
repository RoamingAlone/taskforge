package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
	firstName string,
	lastName string,
) (*User, error) {
	query := `
		INSERT INTO users (
			email,
			password_hash,
			first_name,
			last_name
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			email,
			password_hash,
			first_name,
			last_name,
			created_at,
			updated_at
	`

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
		passwordHash,
		firstName,
		lastName,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id int64,
) (*User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
