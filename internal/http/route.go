package http

import (
	"site/internal/handler"
	"site/internal/middleware"

	"github.com/go-chi/chi"
)

func (s *Server) PrepareUserRoutes(r chi.Router, h *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.All)
		r.Get("/{handle}", h.ByHandle)
		r.Put("/", h.Update)
		r.Delete("/{handle}", h.Delete)
	})
}

func (s *Server) PrepareSubmissionRoutes(r chi.Router, h *handler.SubmissionHandler) {
	r.Route("/submissions", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.All)
		r.Get("/{id}", h.ById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (s *Server) PrepareUploadFileRoutes(r chi.Router, h *handler.UploadFileHandler) {
	r.HandleFunc("/upload", h.UploadFile())
}
