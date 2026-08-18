package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/auth"
	"github.com/example/inspection-platform/internal/platform/web"
	taskapp "github.com/example/inspection-platform/internal/task/application"
)

type Handler struct {
	service *taskapp.Service
}

func NewHandler(service *taskapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/tasks", Handler: h.list},
		{Method: http.MethodGet, Path: "/api/tasks/{id}", Handler: h.get},
		{Method: http.MethodPost, Path: "/api/tasks/{id}/submit", Handler: h.submit},
		{Method: http.MethodPost, Path: "/api/tasks/{id}/draft", Handler: h.draft},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := domain.TaskFilter{
		Status:     web.Query(r, "status"),
		DeviceID:   web.QueryInt64(r, "deviceId"),
		ExecutorID: web.QueryInt64(r, "executorId"),
		From:       web.QueryTime(r, "from"),
		To:         web.QueryTime(r, "to"),
		Page:       web.QueryInt(r, "page"),
		PageSize:   web.QueryInt(r, "pageSize"),
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
	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, task)
}

func (h *Handler) submit(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		web.Error(w, http.StatusUnauthorized, errMissingClaims())
		return
	}
	executorID, _ := strconv.ParseInt(claims.Sub, 10, 64)
	var input taskapp.SubmitInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	task, err := h.service.Submit(r.Context(), id, executorID, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, task)
}

func (h *Handler) draft(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		web.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}
	executorID, _ := strconv.ParseInt(claims.Sub, 10, 64)
	var input taskapp.SubmitInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	task, err := h.service.SaveDraft(r.Context(), id, executorID, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, task)
}

func errMissingClaims() error {
	return errors.New("missing auth claims")
}
