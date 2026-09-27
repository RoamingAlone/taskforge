package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type fakeUserService struct {
	registerInput RegisterInput
	registerErr   error
}

func (f *fakeUserService) Register(
	ctx context.Context,
	input RegisterInput,
) (*User, error) {
	f.registerInput = input

	if f.registerErr != nil {
		return nil, f.registerErr
	}

	return &User{
		ID:        1,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}, nil
}

func TestHandlerRegisterSubmit(t *testing.T) {
	service := &fakeUserService{}
	handler := NewHandler(service)

	form := url.Values{}
	form.Set("email", "test@example.com")
	form.Set("password", "supersecret123")
	form.Set("first_name", "Test")
	form.Set("last_name", "User")

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	recorder := httptest.NewRecorder()

	handler.RegisterSubmit(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusSeeOther {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusSeeOther,
			response.StatusCode,
		)
	}

	if location := response.Header.Get("Location"); location != "/" {
		t.Errorf(
			"expected redirect to /, got %q",
			location,
		)
	}

	if service.registerInput.Email != "test@example.com" {
		t.Errorf(
			"expected email %q, got %q",
			"test@example.com",
			service.registerInput.Email,
		)
	}
}

func TestHandlerRegisterSubmitErrors(t *testing.T) {
	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
	}{
		{
			name:           "email required",
			serviceErr:     ErrEmailRequired,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "password required",
			serviceErr:     ErrPasswordRequired,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "password too short",
			serviceErr:     ErrPasswordTooShort,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "email already exists",
			serviceErr:     ErrEmailAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "unexpected error",
			serviceErr:     errors.New("something went wrong"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeUserService{
				registerErr: tt.serviceErr,
			}

			handler := NewHandler(service)

			form := url.Values{}
			form.Set("email", "test@example.com")
			form.Set("password", "supersecret123")
			form.Set("first_name", "Test")
			form.Set("last_name", "User")

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				strings.NewReader(form.Encode()),
			)

			req.Header.Set(
				"Content-Type",
				"application/x-www-form-urlencoded",
			)

			recorder := httptest.NewRecorder()

			handler.RegisterSubmit(recorder, req)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatus,
					response.StatusCode,
				)
			}
		})
	}
}
