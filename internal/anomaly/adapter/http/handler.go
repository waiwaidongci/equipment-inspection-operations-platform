package http

import (
	"net/http"

	anomalyapp "github.com/example/inspection-platform/internal/anomaly/application"
	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/web"
)

type Handler struct {
	service *anomalyapp.Service
}

func NewHandler(service *anomalyapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/anomalies", Handler: h.list},
		{Method: http.MethodGet, Path: "/api/anomalies/{id}", Handler: h.get},
		{Method: http.MethodPost, Path: "/api/anomalies/{id}/assign", Handler: h.assign},
		{Method: http.MethodPost, Path: "/api/anomalies/{id}/progress", Handler: h.progress},
		{Method: http.MethodPost, Path: "/api/anomalies/{id}/close", Handler: h.close},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := domain.AnomalyFilter{
		Status:   web.Query(r, "status"),
		DeviceID: web.QueryInt64(r, "deviceId"),
		Severity: web.Query(r, "severity"),
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

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	anomaly, err := h.service.Get(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, anomaly)
}

func (h *Handler) assign(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input struct {
		AssigneeID int64 `json:"assigneeId"`
	}
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	anomaly, err := h.service.Assign(r.Context(), id, input.AssigneeID)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, anomaly)
}

func (h *Handler) progress(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input struct {
		Progress string `json:"progress"`
	}
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	anomaly, err := h.service.Progress(r.Context(), id, input.Progress)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, anomaly)
}

func (h *Handler) close(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input anomalyapp.CloseInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	anomaly, err := h.service.Close(r.Context(), id, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, anomaly)
}
