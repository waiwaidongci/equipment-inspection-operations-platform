package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/example/inspection-platform/internal/platform/auth"
	"github.com/example/inspection-platform/internal/platform/web"
	userapp "github.com/example/inspection-platform/internal/user/application"
)

type Handler struct {
	service *userapp.Service
}

func NewHandler(service *userapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodPost, Path: "/api/auth/login", Handler: h.login},
		{Method: http.MethodGet, Path: "/api/auth/me", Handler: h.me},
		{Method: http.MethodGet, Path: "/api/users", Handler: h.list},
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input userapp.LoginInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.service.Login(r.Context(), input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, result)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		web.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}
	id, _ := strconv.ParseInt(claims.Sub, 10, 64)
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, user)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, users)
}
