package http

import (
	"site/internal/http/handler"

	"github.com/go-chi/chi"
)

func (s *Server) PrepareUserRoutes(r chi.Router, h *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.All)
		r.Get("/{handle}", h.ByHandle)
		r.Put("/", h.Update)
		r.Delete("/{handle}", h.Delete)
	})
}

func (s *Server) PrepareSubmissionRoutes(r chi.Router, h *handler.SubmissionHandler) {
	r.Route("/contests", func(r chi.Router) {
		r.Post("/{contestId}/submissions", h.Create)
		r.Get("/{contestId}/submissions", h.All)
		r.Get("/{contestId}/submissions/{id}", h.ById)
		r.Put("/{contestId}/submissions", h.Update)
		r.Delete("/{contestId}/submissions/{id}", h.Delete)
	})
}

func (s *Server) PrepareGuestWebPageRoutes(r chi.Router) {
	r.HandleFunc("/", s.HomePage())
	r.HandleFunc("/register", s.RegisterPage())
	r.HandleFunc("/login", s.LoginPage())
}

func (s *Server) PrepareAuthenticatedWebPageRoutes(r chi.Router) {
	r.HandleFunc("/upload", s.UploadPage())
	r.HandleFunc("/ratings", s.RatingsPage())
}
