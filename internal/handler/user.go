package handler

import (
	"encoding/json"
	"net/http"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/services"

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

func (uh *UserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserFindAllRequest{
		Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
		Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
	}

	users, err := uh.service.FindAll(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}

func (uh *UserHandler) FindFriends(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserFindFriendsRequest{
		Page: r.Context().Value(middleware.CtxKeyPage).(int32),
	}

	users, err := uh.service.FindFriends(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}

func (uh *UserHandler) GetByHandle(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")

	user, err := uh.service.GetByHandle(r.Context(), handle)
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
	Respond(w, r, http.StatusOK, nil)
}
