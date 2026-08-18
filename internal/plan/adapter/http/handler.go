package http

import (
	"net/http"

	"github.com/example/inspection-platform/internal/domain"
	planapp "github.com/example/inspection-platform/internal/plan/application"
	"github.com/example/inspection-platform/internal/platform/web"
)

type Handler struct {
	service *planapp.Service
}

func NewHandler(service *planapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/plans", Handler: h.list},
		{Method: http.MethodPost, Path: "/api/plans", Handler: h.create},
		{Method: http.MethodGet, Path: "/api/plans/{id}", Handler: h.get},
		{Method: http.MethodPut, Path: "/api/plans/{id}", Handler: h.update},
		{Method: http.MethodPost, Path: "/api/plans/{id}/status", Handler: h.status},
		{Method: http.MethodPost, Path: "/api/plans/generate", Handler: h.generate},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := domain.PlanFilter{
		Status:   web.Query(r, "status"),
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
	var input planapp.PlanInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	plan, err := h.service.Create(r.Context(), input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusCreated, plan)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	plan, err := h.service.Get(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, plan)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input planapp.PlanInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	plan, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, plan)
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input struct {
		Status string `json:"status"`
	}
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	plan, err := h.service.SetStatus(r.Context(), id, input.Status)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, plan)
}

func (h *Handler) generate(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.GenerateDue(r.Context())
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]int{"created": count})
}
