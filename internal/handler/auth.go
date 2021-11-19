package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/services"

	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	service      services.AuthService
	sessionStore sessions.Store
}

func NewAuthHandler(opts ...AuthHandlerOption) *AuthHandler {
	ah := &AuthHandler{}
	for _, v := range opts {
		v(ah)
	}
	return ah
}

func (a AuthHandler) CreateSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &dto.Cridentials{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := a.service.ByEmail(r.Context(), req)
		if err != nil {
			Error(w, r, http.StatusUnauthorized, err)
			return
		}

		session, err := a.sessionStore.Get(r, middleware.SessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_handle"] = user.Handle
		if err := a.sessionStore.Save(r, w, session); err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		Respond(w, r, http.StatusOK, nil)
	}
}

func (a AuthHandler) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sessionStore.Get(r, middleware.SessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		handle, ok := session.Values["user_handle"]
		if !ok {
			Error(w, r, http.StatusUnauthorized, middleware.ErrNotAuthenticated)
			return
		}
		// TODO: CHECK THAT HANDLE IS STRING
		user, err := a.service.ByHandle(r.Context(), handle.(string))
		if err != nil {
			Error(w, r, http.StatusUnauthorized, middleware.ErrNotAuthenticated)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), middleware.CtxKeyUser, user)))
	})
}

func (a AuthHandler) HandleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Respond(w, r, http.StatusOK, middleware.UserFromCtx(r.Context()))
	}
}
