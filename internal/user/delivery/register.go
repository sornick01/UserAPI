package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sornick01/UserAPI/internal/user"
	"net/http"
	"time"
)

func RegisterEndpoints(r *chi.Mux, uc user.UseCase) {
	h := NewHandlers(uc)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", h.SearchUsersHandler)
				r.Post("/", h.CreateUserHandler)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.GetUserHandler)
					r.Patch("/", h.UpdateUserHandler)
					r.Delete("/", h.DeleteUserHandler)
				})
			})
		})
	})
}
