package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context) ([]User, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, device *Device) error
	Update(ctx context.Context, device *Device) error
	FindByID(ctx context.Context, id int64) (*Device, error)
	FindByCode(ctx context.Context, code string) (*Device, error)
	List(ctx context.Context, filter DeviceFilter) ([]Device, int64, error)
}

type PlanRepository interface {
	Create(ctx context.Context, plan *InspectionPlan) error
	Update(ctx context.Context, plan *InspectionPlan) error
	FindByID(ctx context.Context, id int64) (*InspectionPlan, error)
	List(ctx context.Context, filter PlanFilter) ([]InspectionPlan, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	ListEnabled(ctx context.Context) ([]InspectionPlan, error)
	ReplaceItems(ctx context.Context, planID int64, items []PlanCheckItem) error
	ListItems(ctx context.Context, planID int64) ([]PlanCheckItem, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *InspectionTask) error
	CreateBatch(ctx context.Context, tasks []InspectionTask) error
	FindByID(ctx context.Context, id int64) (*InspectionTask, error)
	List(ctx context.Context, filter TaskFilter) ([]InspectionTask, int64, error)
	Update(ctx context.Context, task *InspectionTask) error
	ExistsByPlanDeviceDate(ctx context.Context, planID, deviceID int64, plannedAt time.Time) (bool, error)
	ListPendingDue(ctx context.Context, before time.Time) ([]InspectionTask, error)
}

type RecordRepository interface {
	Create(ctx context.Context, record *InspectionRecord) error
	ListByDevice(ctx context.Context, deviceID int64, limit, offset int) ([]InspectionRecord, error)
	ListByTask(ctx context.Context, taskID int64) ([]InspectionRecord, error)
}

type AnomalyRepository interface {
	Create(ctx context.Context, anomaly *Anomaly) error
	Update(ctx context.Context, anomaly *Anomaly) error
	FindByID(ctx context.Context, id int64) (*Anomaly, error)
	List(ctx context.Context, filter AnomalyFilter) ([]Anomaly, int64, error)
	ListByDevice(ctx context.Context, deviceID int64) ([]Anomaly, error)
}

type RepairRepository interface {
	Create(ctx context.Context, repair *Repair) error
	Update(ctx context.Context, repair *Repair) error
	FindByID(ctx context.Context, id int64) (*Repair, error)
	List(ctx context.Context, filter RepairFilter) ([]Repair, int64, error)
	ListByDevice(ctx context.Context, deviceID int64) ([]Repair, error)
}
