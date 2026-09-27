package user

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		email string,
		passwordHash string,
		firstName string,
		lastName string,
	) (*User, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*User, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*User, error)
}

type Service struct {
	repo UserRepository
}

type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(
	ctx context.Context,
	input RegisterInput,
) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	if email == "" {
		return nil, ErrEmailRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(input.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)

	return s.repo.Create(
		ctx,
		email,
		string(passwordHash),
		firstName,
		lastName,
	)
}
