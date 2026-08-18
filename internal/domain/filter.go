package domain

import "time"

type DeviceFilter struct {
	Query    string
	Category string
	Status   string
	Owner    string
	Location string
	Page     int
	PageSize int
}

type PlanFilter struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type TaskFilter struct {
	Status     string
	DeviceID   int64
	ExecutorID int64
	From       time.Time
	To         time.Time
	Page       int
	PageSize   int
}

type AnomalyFilter struct {
	Status   string
	DeviceID int64
	Severity string
	Page     int
	PageSize int
}

type RepairFilter struct {
	DeviceID int64
	Keyword  string
	Page     int
	PageSize int
}

type Page struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int64 `json:"totalPages"`
}
