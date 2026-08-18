package http

import (
	"net/http"

	"github.com/example/inspection-platform/internal/device/application"
	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/web"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/devices", Handler: h.list},
		{Method: http.MethodPost, Path: "/api/devices", Handler: h.create},
		{Method: http.MethodGet, Path: "/api/devices/{id}", Handler: h.get},
		{Method: http.MethodPut, Path: "/api/devices/{id}", Handler: h.update},
		{Method: http.MethodGet, Path: "/api/devices/{id}/history", Handler: h.history},
		{Method: http.MethodGet, Path: "/api/dictionaries", Handler: h.dictionaries},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := domain.DeviceFilter{
		Query:    web.Query(r, "query"),
		Category: web.Query(r, "category"),
		Status:   web.Query(r, "status"),
		Owner:    web.Query(r, "owner"),
		Location: web.Query(r, "location"),
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
	var input application.DeviceInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	device, err := h.service.Create(r.Context(), input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusCreated, device)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	device, err := h.service.Get(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, device)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	var input application.DeviceInput
	if err := web.Decode(r, &input); err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	device, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, device)
}

func (h *Handler) history(w http.ResponseWriter, r *http.Request) {
	id, err := web.PathID(r, "id")
	if err != nil {
		web.Error(w, http.StatusBadRequest, err)
		return
	}
	history, err := h.service.History(r.Context(), id)
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, history)
}

func (h *Handler) dictionaries(w http.ResponseWriter, r *http.Request) {
	web.JSON(w, http.StatusOK, map[string][]string{
		"categories": {"电气设备", "机械设备", "消防设备", "暖通设备", "电梯设备", "其他"},
		"locations":  {"生产车间", "动力站房", "办公楼", "仓库", "室外区域"},
		"statuses":   {"active", "inactive"},
	})
}
