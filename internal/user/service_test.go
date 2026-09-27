package user

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	createEmail        string
	createPasswordHash string
	createFirstName    string
	createLastName     string
	createErr          error
}

func (f *fakeUserRepository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
	firstName string,
	lastName string,
) (*User, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}

	f.createEmail = email
	f.createPasswordHash = passwordHash
	f.createFirstName = firstName
	f.createLastName = lastName

	return &User{
		ID:           1,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
	}, nil
}

func (f *fakeUserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*User, error) {
	return nil, ErrNotFound
}

func (f *fakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	return nil, ErrNotFound
}

func TestServiceRegister(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewService(repo)

	input := RegisterInput{
		Email:     "  TEST@Example.COM  ",
		Password:  "supersecret123",
		FirstName: "  Test  ",
		LastName:  "  User  ",
	}

	createdUser, err := service.Register(
		context.Background(),
		input,
	)
	if err != nil {
		t.Fatal(err)
	}

	if createdUser.ID != 1 {
		t.Errorf(
			"expected user ID %d, got %d",
			1,
			createdUser.ID,
		)
	}

	if repo.createEmail != "test@example.com" {
		t.Errorf(
			"expected normalized email %q, got %q",
			"test@example.com",
			repo.createEmail,
		)
	}

	if repo.createFirstName != "Test" {
		t.Errorf(
			"expected first name %q, got %q",
			"Test",
			repo.createFirstName,
		)
	}

	if repo.createLastName != "User" {
		t.Errorf(
			"expected last name %q, got %q",
			"User",
			repo.createLastName,
		)
	}

	if repo.createPasswordHash == input.Password {
		t.Error("expected password to be hashed")
	}

	if repo.createPasswordHash == "" {
		t.Error("expected password hash, got empty string")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(repo.createPasswordHash),
		[]byte(input.Password),
	)
	if err != nil {
		t.Errorf("password hash does not match password: %v", err)
	}
}

func TestServiceRegisterValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    RegisterInput
		expected error
	}{
		{
			name: "missing email",
			input: RegisterInput{
				Email:    "",
				Password: "supersecret123",
			},
			expected: ErrEmailRequired,
		},
		{
			name: "email containing only whitespace",
			input: RegisterInput{
				Email:    "   ",
				Password: "supersecret123",
			},
			expected: ErrEmailRequired,
		},
		{
			name: "missing password",
			input: RegisterInput{
				Email:    "test@example.com",
				Password: "",
			},
			expected: ErrPasswordRequired,
		},
		{
			name: "password too short",
			input: RegisterInput{
				Email:    "test@example.com",
				Password: "short",
			},
			expected: ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepository{}
			service := NewService(repo)

			_, err := service.Register(
				context.Background(),
				tt.input,
			)

			if !errors.Is(err, tt.expected) {
				t.Errorf(
					"expected error %v, got %v",
					tt.expected,
					err,
				)
			}
		})
	}
}

func TestServiceRegisterDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepository{
		createErr: ErrEmailAlreadyExists,
	}

	service := NewService(repo)

	_, err := service.Register(
		context.Background(),
		RegisterInput{
			Email:     "test@example.com",
			Password:  "supersecret123",
			FirstName: "Test",
			LastName:  "User",
		},
	)

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Errorf(
			"expected ErrEmailAlreadyExists, got %v",
			err,
		)
	}
}
