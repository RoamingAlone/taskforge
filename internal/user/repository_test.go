package user

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepositoryUserLifecycle(t *testing.T) {
	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://taskforge:password@localhost:5432/taskforge",
	)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)

	_, err = db.Exec(
		ctx,
		"DELETE FROM users WHERE email = $1",
		"test@example.com",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			context.Background(),
			"DELETE FROM users WHERE email = $1",
			"test@example.com",
		)
		if err != nil {
			t.Errorf("cleanup test user: %v", err)
		}

		db.Close()
	})

	createdUser, err := repo.Create(
		ctx,
		"test@example.com",
		"fake-password-hash",
		"Test",
		"User",
	)
	if err != nil {
		t.Fatal(err)
	}

	if createdUser.ID == 0 {
		t.Error("expected user ID to be generated")
	}

	if createdUser.Email != "test@example.com" {
		t.Errorf(
			"expected email %q, got %q",
			"test@example.com",
			createdUser.Email,
		)
	}

	_, err = repo.Create(
		ctx,
		"test@example.com",
		"another-password-hash",
		"Another",
		"User",
	)

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Errorf(
			"expected ErrEmailAlreadyExists, got %v",
			err,
		)
	}

	foundUser, err := repo.GetByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf(
			"expected user ID %d, got %d",
			createdUser.ID,
			foundUser.ID,
		)
	}

	if foundUser.Email != createdUser.Email {
		t.Errorf(
			"expected email %q, got %q",
			createdUser.Email,
			foundUser.Email,
		)
	}

	foundByEmail, err := repo.GetByEmail(ctx, createdUser.Email)
	if err != nil {
		t.Fatal(err)
	}

	if foundByEmail.ID != createdUser.ID {
		t.Errorf(
			"expected user ID %d, got %d",
			createdUser.ID,
			foundByEmail.ID,
		)
	}

	if foundByEmail.Email != createdUser.Email {
		t.Errorf(
			"expected email %q, got %q",
			createdUser.Email,
			foundByEmail.Email,
		)
	}

	_, err = repo.GetByID(ctx, -1)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}
