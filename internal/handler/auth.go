package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"site/internal/consts"
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

func (a AuthHandler) UpdateLocale() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			LanguageCode int `json:"language_code"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			return
		}

		session, err := a.sessionStore.Get(r, middleware.SessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["language_code"] = req.LanguageCode
		fmt.Println(req.LanguageCode)
		if err := a.sessionStore.Save(r, w, session); err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		Respond(w, r, http.StatusOK, nil)
	}
}

func (a AuthHandler) VerifyLocale(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sessionStore.Get(r, middleware.SessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		languageCode, ok := session.Values["language_code"]
		if !ok {
			languageCode = consts.EN
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), middleware.CtxKeyLanguageCode, languageCode)))
	})
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
		user, err := a.service.ByHandle(r.Context(), handle.(string))
		if err != nil {
			Error(w, r, http.StatusUnauthorized, middleware.ErrNotAuthenticated)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), middleware.CtxKeyUser, user)))
	})
}

func (a AuthHandler) HandleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromCtx(r.Context())
		Respond(w, r, http.StatusOK, user)
	}
}
