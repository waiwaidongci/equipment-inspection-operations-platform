package api

import (
	"database/sql"
	"net/http"

	anomalyhttp "github.com/example/inspection-platform/internal/anomaly/adapter/http"
	devicehttp "github.com/example/inspection-platform/internal/device/adapter/http"
	planhttp "github.com/example/inspection-platform/internal/plan/adapter/http"
	"github.com/example/inspection-platform/internal/platform/middleware"
	"github.com/example/inspection-platform/internal/platform/web"
	recordhttp "github.com/example/inspection-platform/internal/record/adapter/http"
	repairhttp "github.com/example/inspection-platform/internal/repair/adapter/http"
	taskhttp "github.com/example/inspection-platform/internal/task/adapter/http"
	userhttp "github.com/example/inspection-platform/internal/user/adapter/http"
)

type Dependencies struct {
	DB        *sql.DB
	Metrics   *middleware.Metrics
	JWTSecret string
	Users     *userhttp.Handler
	Devices   *devicehttp.Handler
	Plans     *planhttp.Handler
	Tasks     *taskhttp.Handler
	Records   *recordhttp.Handler
	Anomalies *anomalyhttp.Handler
	Repairs   *repairhttp.Handler
}

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler(deps.DB))
	mux.HandleFunc("GET /readyz", readyHandler(deps.DB))
	mux.Handle("GET /metrics", deps.Metrics.Handler())

	for _, routes := range [][]web.Route{
		deps.Users.Routes(),
		deps.Devices.Routes(),
		deps.Plans.Routes(),
		deps.Tasks.Routes(),
		deps.Records.Routes(),
		deps.Anomalies.Routes(),
		deps.Repairs.Routes(),
	} {
		for _, route := range routes {
			mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
		}
	}

	publicPaths := map[string]bool{
		"/healthz":        true,
		"/readyz":         true,
		"/metrics":        true,
		"/api/auth/login": true,
	}
	return middleware.Chain(
		mux,
		deps.Metrics.Middleware,
		middleware.RequestID,
		middleware.Logging,
		middleware.Recover,
		middleware.CORS,
		middleware.Auth(deps.JWTSecret, publicPaths),
	)
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			web.Error(w, http.StatusServiceUnavailable, err)
			return
		}
		web.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func readyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			web.Error(w, http.StatusServiceUnavailable, err)
			return
		}
		web.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
