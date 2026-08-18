package http

import (
	"net/http"

	"github.com/example/inspection-platform/internal/platform/web"
	recordapp "github.com/example/inspection-platform/internal/record/application"
)

type Handler struct {
	service *recordapp.Service
}

func NewHandler(service *recordapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() []web.Route {
	return []web.Route{
		{Method: http.MethodGet, Path: "/api/records", Handler: h.list},
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	deviceID := web.QueryInt64(r, "deviceId")
	if deviceID <= 0 {
		web.Error(w, http.StatusBadRequest, errDeviceRequired())
		return
	}
	records, err := h.service.ListByDevice(r.Context(), deviceID, web.QueryInt(r, "limit"), web.QueryInt(r, "offset"))
	if err != nil {
		web.Error(w, web.StatusForError(err), err)
		return
	}
	web.JSON(w, http.StatusOK, records)
}

func errDeviceRequired() error {
	return errBadRequest("deviceId is required")
}

type badRequest string

func (b badRequest) Error() string { return string(b) }

func errBadRequest(msg string) error { return badRequest(msg) }
