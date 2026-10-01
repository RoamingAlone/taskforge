package user

import (
	"context"
	"errors"
	"html/template"
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

func newTestHandler(
	t *testing.T,
	service UserService,
) *Handler {
	t.Helper()

	tmpl, err := template.New("register").Parse(`
		{{ define "base" }}
			<!DOCTYPE html>
			<html>
			<body>
				<h1>Register</h1>

				{{ if .Error }}
					<div>{{ .Error }}</div>
				{{ end }}

				<input name="email" value="{{ .Email }}">
				<input name="first_name" value="{{ .FirstName }}">
				<input name="last_name" value="{{ .LastName }}">
				<input name="password" type="password">
			</body>
			</html>
		{{ end }}
	`)

	if err != nil {
		t.Fatalf("parse test template: %v", err)
	}

	return NewHandler(service, tmpl)
}

func TestHandlerRegisterSubmit(t *testing.T) {
	service := &fakeUserService{}
	handler := newTestHandler(t, service)

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

			handler := newTestHandler(t, service)

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

			if errors.Is(tt.serviceErr, ErrEmailAlreadyExists) {
				body := recorder.Body.String()

				if !strings.Contains(body, "email already exists") {
					t.Errorf(
						"expected response body to contain duplicate email error",
					)
				}

				if !strings.Contains(body, "test@example.com") {
					t.Errorf(
						"expected response body to preserve email",
					)
				}

				if !strings.Contains(body, "Test") {
					t.Errorf(
						"expected response body to preserve first name",
					)
				}

				if !strings.Contains(body, "User") {
					t.Errorf(
						"expected response body to preserve last name",
					)
				}

				if strings.Contains(body, "supersecret123") {
					t.Errorf(
						"expected response body not to contain password",
					)

				}
			}
		})
	}
}

func TestHandlerRegisterForm(t *testing.T) {
	service := &fakeUserService{}
	handler := newTestHandler(t, service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/register",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.RegisterForm(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	if !strings.Contains(recorder.Body.String(), "<h1>Register</h1>") {
		t.Errorf(
			"expected response body to contain registration heading",
		)
	}
}
