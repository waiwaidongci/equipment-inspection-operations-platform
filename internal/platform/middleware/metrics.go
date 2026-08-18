package middleware

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
)

type Metrics struct {
	total    atomic.Int64
	byStatus sync.Map
	requests atomic.Int64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) observe(status int) {
	m.total.Add(1)
	m.requests.Add(1)
	key := fmt.Sprintf("%d", status)
	value, _ := m.byStatus.LoadOrStore(key, new(atomic.Int64))
	value.(*atomic.Int64).Add(1)
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		m.observe(recorder.status)
	})
}

func (m *Metrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = fmt.Fprintf(w, "# HELP inspection_http_requests_total Total HTTP requests.\n")
		_, _ = fmt.Fprintf(w, "# TYPE inspection_http_requests_total counter\n")
		_, _ = fmt.Fprintf(w, "inspection_http_requests_total %d\n", m.total.Load())
		_, _ = fmt.Fprintf(w, "# HELP inspection_http_responses_total HTTP responses by status code.\n")
		_, _ = fmt.Fprintf(w, "# TYPE inspection_http_responses_total counter\n")
		keys := make([]string, 0)
		m.byStatus.Range(func(key, _ any) bool {
			keys = append(keys, key.(string))
			return true
		})
		sort.Strings(keys)
		for _, key := range keys {
			value, _ := m.byStatus.Load(key)
			_, _ = fmt.Fprintf(w, "inspection_http_responses_total{code=%q} %d\n", key, value.(*atomic.Int64).Load())
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
