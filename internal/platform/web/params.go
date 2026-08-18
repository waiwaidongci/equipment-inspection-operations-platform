package web

import (
	"net/http"
	"strconv"
	"time"
)

func PathID(r *http.Request, name string) (int64, error) {
	value := r.PathValue(name)
	return strconv.ParseInt(value, 10, 64)
}

func Query(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func QueryInt(r *http.Request, name string) int {
	value, _ := strconv.Atoi(r.URL.Query().Get(name))
	return value
}

func QueryInt64(r *http.Request, name string) int64 {
	value, _ := strconv.ParseInt(r.URL.Query().Get(name), 10, 64)
	return value
}

func QueryTime(r *http.Request, name string) time.Time {
	value := r.URL.Query().Get(name)
	if value == "" {
		return time.Time{}
	}
	t, _ := time.Parse("2006-01-02", value)
	return t
}
