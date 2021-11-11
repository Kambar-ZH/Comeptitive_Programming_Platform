package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"site/internal/http/ioutils"
)

type ctxKey int8

const (
	sessionName        = "gopherschool"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = fmt.Errorf("incorrect email or password")
	errNotAuthenticated         = fmt.Errorf("not authenticated")
)

func (s *Server) HandleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			ioutils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.Users().ByEmail(s.ctx, req.Email)
		// if err != nil || !services.ComparePassword(user, req.Password) {
		// 	ioutils.Error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		// 	return
		// }

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_handle"] = user.Handle
		if err := s.sessionStore.Save(r, w, session); err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		ioutils.Respond(w, r, http.StatusOK, nil)
	}
}

func (s *Server) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		handle, ok := session.Values["user_handle"]
		if !ok {
			ioutils.Error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		// TODO: CHECK THAT HANDLE IS STRING
		user, err := s.store.Users().ByHandle(s.ctx, handle.(string))
		if err != nil {
			ioutils.Error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	})
}
