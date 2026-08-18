package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/auth"
	"github.com/example/inspection-platform/internal/platform/web"
	repairapp "github.com/example/inspection-platform/internal/repair/application"
)

type Handler struct {
	service *repairapp.Service
}

func NewHandler(service *repairapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/repairs", Handler: h.list},
		{Method: http.MethodPost, Path: "/api/repairs", Handler: h.create},
		{Method: http.MethodGet, Path: "/api/repairs/{id}", Handler: h.get},
		{Method: http.MethodPut, Path: "/api/repairs/{id}", Handler: h.update},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := domain.RepairFilter{
		DeviceID: web.QueryInt64(r, "deviceId"),
		Keyword:  web.Query(r, "keyword"),
		Page:     web.QueryInt(r, "page"),
		PageSize: web.QueryInt(r, "pageSize"),
	}
	page, err := h.service.List(r.Context(), filter)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		web.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}
	creatorID, _ := strconv.ParseInt(claims.Sub, 10, 64)
	var input repairapp.RepairInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	repair, err := h.service.Create(r.Context(), creatorID, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusCreated, repair)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	repair, err := h.service.Get(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, repair)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input repairapp.RepairInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	repair, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, repair)
}
