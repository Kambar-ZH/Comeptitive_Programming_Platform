package handler

import (
	"encoding/json"
	"net/http"
	"site/internal/datastruct"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(opts ...UserHandlerOption) *UserHandler {
	uh := &UserHandler{}
	for _, v := range opts {
		v(uh)
	}
	return uh
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := new(datastruct.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	err := uh.service.Create(r.Context(), user)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (uh *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	query := &datastruct.UserQuery{Page: 1}
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil {
			Error(w, r, http.StatusBadRequest, err)
			return
		}
		query.Page = int32(pageNum)
	}
	query.Filter = r.URL.Query().Get("filter")

	users, err := uh.service.All(r.Context(), query)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}

func (uh *UserHandler) ByHandle(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	user, err := uh.service.ByHandle(r.Context(), handle)
	if err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, user)
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := new(datastruct.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	err := uh.service.Update(r.Context(), user)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusAccepted, nil)
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	if err := uh.service.Delete(r.Context(), handle); err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, "")
}
