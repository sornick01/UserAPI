package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sornick01/UserAPI/internal/user"
	"net/http"
)

type Handlers struct {
	useCase user.UseCase
}

type CreateRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateRequest) Bind(r *http.Request) error { return nil }

type UpdateRequest struct {
	DisplayName string `json:"display_name"`
}

func (u *UpdateRequest) Bind(r *http.Request) error { return nil }

func NewHandlers(useCase user.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req := CreateRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := h.useCase.CreateUser(r.Context(), req.DisplayName, req.Email)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServ(err))
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (h *Handlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.useCase.DeleteUser(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	u, err := h.useCase.GetUser(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		_ = render.Render(w, r, ErrInternalServ(err))
		return
	}

	render.JSON(w, r, u)
}

func (h *Handlers) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.useCase.SearchUser(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrInternalServ(err))
		return
	}

	render.JSON(w, r, users)
}

func (h *Handlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	req := UpdateRequest{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")
	err := h.useCase.UpdateUser(r.Context(), id, req.DisplayName)
	if err != nil {
		if err == user.ErrUserNotFound {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, ErrInternalServ(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServ(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal error",
		ErrorText:      err.Error(),
	}
}
