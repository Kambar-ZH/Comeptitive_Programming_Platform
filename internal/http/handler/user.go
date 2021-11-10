package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/ioutils"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	service services.UserService
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := new(api.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	userResult, err := uh.service.Create(r.Context(), user)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	services.Sanitize(user)
	ioutils.Respond(w, r, http.StatusCreated, userResult)
}

func (uh *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
	}
	users, err := uh.service.All(r.Context(), int32(page))
	if err != nil {
		log.Println(err)
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, users)
}

func (uh *UserHandler) ByHandle(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	user, err := uh.service.ByHandle(r.Context(), handle)
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

	userResult, err := uh.service.Update(r.Context(), user)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusAccepted, userResult)
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	if _, err := uh.service.Delete(r.Context(), &api.UserRequestHandle{Handle: handle}); err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, "")
}
