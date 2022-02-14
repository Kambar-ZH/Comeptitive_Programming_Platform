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
	ps2 := services.NewParticipantRepository(services.ParticipantServiceWithStore(s.store))
	ufs2 := services.NewUserFriendService(services.UserFriendsServiceWithStore(s.store))

	uh := handler.NewUserHandler(handler.WithUserService(us))
	sh := handler.NewSubmissionHandler(handler.WithSubmissionService(ss))
	ufh := handler.NewUploadFileHandler(handler.WithUploadFileService(ufs))
	ah := handler.NewAuthHandler(handler.WithAuthService(as), handler.WithSessionStore(s.sessionStore))
	ch := handler.NewContestHandler(handler.WithContestService(cs))
	ph := handler.NewProblemHandler(handler.WithProblemService(ps))
	ph2 := handler.NewParticipantHandler(handler.WithParticipantService(ps2))
	ufh2 := handler.NewUserFriendHandler(handler.WithUserFriendService(ufs2))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.Create)
		r.With(middleware.Paginate).Get("/", uh.FindAll)
		r.Put("/", uh.Update)
		r.Route("/{handle}", func(r chi.Router) {
			r.Get("/", uh.GetByHandle)
			r.Delete("/", uh.Delete)
			r.Route("/friends", func(r chi.Router) {
				r.Use(ah.AuthenticateUser)
				r.With(middleware.Paginate).Get("/", uh.FindFriends)
				r.Post("/", ufh2.Create)
				r.Delete("/", ufh2.Delete)
			})
		})
	})

	r.Route("/profile", func(r chi.Router) {
		r.Use(ah.AuthenticateUser)
		r.HandleFunc("/", ah.HandleProfile())
	})

	r.Route("/problemset", func(r chi.Router) {
		r.Get("/", ph.Problemset)
	})

	r.HandleFunc("/sessions", ah.CreateSession())
	r.Route("/contests", func(r chi.Router) {
		r.Get("/", ch.FindAll)
		r.Post("/", ch.Create)
		r.With(middleware.Paginate).Get("/", ch.FindAll)
		r.Put("/", ch.Update)
		r.Delete("/{id}", ch.Delete)
		r.Route("/{contestId}", func(r chi.Router) {
			r.Use(ah.AuthenticateUser)
			r.Post("/register", ph2.Register)
			r.Get("/", ch.GetById)
			s.PrepareSubmissionRoutes(r.With(ah.AuthenticateUser), sh)
			s.PrepareParticipantsRoutes(r.With(ah.AuthenticateUser), ph2)
			s.PrepareProblemRoutes(r, ph)
			r.HandleFunc("/upload", ufh.UploadFile())
		})
	})
	return r
}

func (s *Server) PrepareProblemRoutes(r chi.Router, h *handler.ProblemHandler) {
	r.Route("/problems", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.FindAll)
		r.Get("/{id}", h.GetById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (s *Server) PrepareSubmissionRoutes(r chi.Router, h *handler.SubmissionHandler) {
	r.Route("/submissions", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.FindAll)
		r.Get("/{id}", h.GetById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (s *Server) PrepareContestRoutes(r chi.Router, h *handler.ContestHander) {
	r.Route("/contests", func(r chi.Router) {
		r.Post("/", h.Create)
		r.With(middleware.Paginate).Get("/", h.FindAll)
		r.Get("/{id}", h.GetById)
		r.Put("/", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (s *Server) PrepareParticipantsRoutes(r chi.Router, h *handler.ParticipantHandler) {
	r.Route("/standings", func(r chi.Router) {
		r.With(middleware.Paginate).Get("/", h.FindAll)
		r.With(middleware.Paginate).Get("/friends", h.FindFriends)
	})
}
