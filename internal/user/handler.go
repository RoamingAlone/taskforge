package user

import (
	"context"
	"errors"
	"net/http"

	"html/template"
)

type UserService interface {
	Register(
		ctx context.Context,
		input RegisterInput,
	) (*User, error)
}

type Handler struct {
	service          UserService
	registerTemplate *template.Template
}

func NewHandler(
	service UserService,
	registerTemplate *template.Template,
) *Handler {
	return &Handler{
		service:          service,
		registerTemplate: registerTemplate,
	}
}

func (h *Handler) RegisterForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := h.registerTemplate.ExecuteTemplate(
		w,
		"base",
		nil,
	)

	if err != nil {
		http.Error(
			w,
			"failed to render page",
			http.StatusInternalServerError,
		)
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
