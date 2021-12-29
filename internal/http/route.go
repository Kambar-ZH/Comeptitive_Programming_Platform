package http

import (
	"site/internal/handler"
	"site/internal/middleware"
	"site/internal/services"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
)

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	us := services.NewUserService(services.UserServiceWithStore(s.store))
	ss := services.NewSubmissionService(services.SubmissionServiceWithStore(s.store), services.SubmissionServiceWithBroker(s.broker), services.SubmissionServiceWithCache(s.cache))
	ufs := services.NewUploadFileService(services.UploadFileServiceWithStore(s.store))
	as := services.NewAuthService(services.AuthServiceWithStore(s.store))
	cs := services.NewContestService(services.ContestServiceWithStore(s.store))
	ps := services.NewProblemService(services.ProblemServiceWithStore(s.store))

	uh := handler.NewUserHandler(handler.WithUserService(us))
	sh := handler.NewSubmissionHandler(handler.WithSubmissionService(ss))
	ufh := handler.NewUploadFileHandler(handler.WithUploadFileService(ufs))
	ah := handler.NewAuthHandler(handler.WithAuthService(as), handler.WithSessionStore(s.sessionStore))
	ch := handler.NewContestHandler(handler.WithContestService(cs))
	ph := handler.NewProblemHandler(handler.WithProblemService(ps))

	s.PrepareUserRoutes(r, uh)

	r.Route("/profile", func(r chi.Router) {
		r.Use(ah.AuthenticateUser)
		r.HandleFunc("/", ah.HandleWhoami())
	})

	r.Route("/problemset", func(r chi.Router) {
		r.Get("/", ph.Problemset)
	})

	r.HandleFunc("/sessions", ah.CreateSession())
	r.Route("/contests", func(r chi.Router) {
		r.Use(ah.AuthenticateUser)
		r.Post("/", ch.Create)
		r.With(middleware.Paginate).Get("/", ch.All)
		r.Put("/", ch.Update)
		r.Delete("/{id}", ch.Delete)
		r.Route("/{contestId}", func(r chi.Router) {
			r.Get("/", ch.ById)
			s.PrepareSubmissionRoutes(r, sh)
			s.PrepareProblemRoutes(r, ph)
			r.HandleFunc("/upload", ufh.UploadFile())
		})
	})
	return r
}

func (s *Server) PrepareProblemRoutes(r chi.Router, h *handler.ProblemHandler) {
	r.Route("/problems", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.All)
		r.Get("/{id}", h.ById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

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

func (s *Server) PrepareContestRoutes(r chi.Router, h *handler.ContestHander) {
	r.Route("/contests", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.All)
		r.Get("/{id}", h.ById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}