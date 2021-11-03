package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/ioutils"
	"site/internal/services"
	"site/internal/store"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	Repository store.UserRepository
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := new(api.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	userResult, err := uh.Repository.Create(r.Context(), user)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	services.Sanitize(user)
	ioutils.Respond(w, r, http.StatusCreated, userResult)
}

func (uh *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	users, err := uh.Repository.All(r.Context(), &api.Pagination{})
	if err != nil {
		log.Println(err)
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, users)
}

func (uh *UserHandler) ByHandle(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	user, err := uh.Repository.ByHandle(r.Context(), &api.UserRequestHandle{Handle: handle})
	if err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, user)
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := new(api.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	userResult, err := uh.Repository.Update(r.Context(), user)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusAccepted, userResult)
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	if _, err := uh.Repository.Delete(r.Context(), &api.UserRequestHandle{Handle: handle}); err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, "")
}