package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"site/internal/dto"
	"site/internal/services"
)

type UserFriendHandler struct {
	service services.UserFriendsService
}

func NewUserFriendHandler(opts ...UserFriendHandlerOption) *UserFriendHandler {
	uh := &UserFriendHandler{}
	for _, v := range opts {
		v(uh)
	}
	return uh
}

func (ufh *UserFriendHandler) Create(w http.ResponseWriter, r *http.Request) {
	friendHandle := chi.URLParam(r, "handle")

	err := ufh.service.Create(r.Context(), &dto.UserFriendCreateRequest{
		FriendHandle: friendHandle,
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (ufh *UserFriendHandler) Delete(w http.ResponseWriter, r *http.Request) {
	friendHandle := chi.URLParam(r, "handle")

	if err := ufh.service.Delete(r.Context(), &dto.UserFriendDeleteRequest{
		FriendHandle: friendHandle,
	}); err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, nil)
}
