package user

import (
	"context"
	"errors"
	"net/http"
)

type UserService interface {
	Register(
		ctx context.Context,
		input RegisterInput,
	) (*User, error)
}

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Register | TaskForge</title>
		</head>
		<body>
			<h1>Create an account</h1>

			<form method="POST" action="/register">
				<label>
					Email
					<input type="email" name="email" required>
				</label>

				<br>

				<label>
					First Name
					<input type="text" name="first_name">
				</label>

				<br>

				<label>
					Last Name
					<input type="text" name="last_name">
				</label>

				<br>

				<label>
					Password
					<input type="password" name="password" required>
				</label>

				<br>

				<button type="submit">Register</button>
			</form>
		</body>
		</html>
	`))

	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterSubmit(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	input := RegisterInput{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
	}

	_, err = h.service.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailRequired):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrPasswordRequired):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrPasswordTooShort):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)

		default:
			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
		}

		return
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}
