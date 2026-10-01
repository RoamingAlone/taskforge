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

type RegisterPageData struct {
	Error     string
	Email     string
	FirstName string
	LastName  string
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
	h.renderRegister(
		w,
		http.StatusOK,
		RegisterPageData{},
	)
}

func (h *Handler) renderRegister(
	w http.ResponseWriter,
	status int,
	data RegisterPageData,
) {
	w.WriteHeader(status)

	err := h.registerTemplate.ExecuteTemplate(
		w,
		"base",
		data,
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

	pageData := RegisterPageData{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	_, err = h.service.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailRequired),
			errors.Is(err, ErrPasswordRequired),
			errors.Is(err, ErrPasswordTooShort):

			pageData.Error = err.Error()

			h.renderRegister(
				w,
				http.StatusBadRequest,
				pageData,
			)

		case errors.Is(err, ErrEmailAlreadyExists):
			pageData.Error = err.Error()

			h.renderRegister(
				w,
				http.StatusConflict,
				pageData,
			)

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
