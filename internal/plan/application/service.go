package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
)

type Service struct {
	plans   domain.PlanRepository
	devices domain.DeviceRepository
	tasks   domain.TaskRepository
}

func NewService(plans domain.PlanRepository, devices domain.DeviceRepository, tasks domain.TaskRepository) *Service {
	return &Service{plans: plans, devices: devices, tasks: tasks}
}

type PlanItemInput struct {
	Name          string `json:"name"`
	ItemType      string `json:"itemType"`
	StandardRange string `json:"standardRange"`
	Required      bool   `json:"required"`
	Description   string `json:"description"`
}

type PlanInput struct {
	Name      string          `json:"name"`
	ScopeType string          `json:"scopeType"`
	DeviceID  *int64          `json:"deviceId,omitempty"`
	Category  string          `json:"category"`
	Period    string          `json:"period"`
	StartDate string          `json:"startDate"`
	EndDate   string          `json:"endDate"`
	Owner     string          `json:"owner"`
	Status    string          `json:"status"`
	Items     []PlanItemInput `json:"items"`
}

func (s *Service) Create(ctx context.Context, input PlanInput) (*domain.InspectionPlan, error) {
	plan, items, err := s.buildPlan(0, input)
	if err != nil {
		return nil, err
	}
	if err := s.plans.Create(ctx, plan); err != nil {
		return nil, err
	}
	if err := s.plans.ReplaceItems(ctx, plan.ID, items); err != nil {
		return nil, err
	}
	plan.Items = items
	if plan.Status == "enabled" {
		_, _ = s.generateForPlan(ctx, plan, time.Now().UTC().AddDate(0, 0, 14))
	}
	return plan, nil
}

func (s *Service) Update(ctx context.Context, id int64, input PlanInput) (*domain.InspectionPlan, error) {
	existing, err := s.plans.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	plan, items, err := s.buildPlan(existing.ID, input)
	if err != nil {
		return nil, err
	}
	plan.CreatedAt = existing.CreatedAt
	if err := s.plans.Update(ctx, plan); err != nil {
		return nil, err
	}
	if err := s.plans.ReplaceItems(ctx, plan.ID, items); err != nil {
		return nil, err
	}
	plan.Items = items
	if plan.Status == "enabled" {
		_, _ = s.generateForPlan(ctx, plan, time.Now().UTC().AddDate(0, 0, 14))
	}
	return plan, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.InspectionPlan, error) {
	plan, err := s.plans.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	items, err := s.plans.ListItems(ctx, id)
	if err != nil {
		return nil, err
	}
	plan.Items = items
	return plan, nil
}

func (s *Service) List(ctx context.Context, filter domain.PlanFilter) (domain.Page, error) {
	items, total, err := s.plans.List(ctx, filter)
	if err != nil {
		return domain.Page{}, err
	}
	return newPage(items, total, filter.Page, filter.PageSize), nil
}

func (s *Service) SetStatus(ctx context.Context, id int64, status string) (*domain.InspectionPlan, error) {
	plan, err := s.plans.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	switch status {
	case "enabled", "paused":
		if plan.Status == "terminated" {
			return nil, fmt.Errorf("%w: terminated plan cannot be restarted", domain.ErrForbidden)
		}
	case "terminated":
	default:
		return nil, fmt.Errorf("%w: unsupported plan status", domain.ErrInvalid)
	}
	if err := s.plans.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	if status == "enabled" {
		_, _ = s.generateForPlan(ctx, plan, time.Now().UTC().AddDate(0, 0, 14))
	}
	plan.Status = status
	return plan, nil
}

func (s *Service) GenerateDue(ctx context.Context) (int, error) {
	horizon := time.Now().UTC().AddDate(0, 0, 14)
	plans, err := s.plans.ListEnabled(ctx)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, plan := range plans {
		items, err := s.plans.ListItems(ctx, plan.ID)
		if err != nil {
			return total, err
		}
		plan.Items = items
		created, err := s.generateForPlan(ctx, &plan, horizon)
		if err != nil {
			return total, err
		}
		total += created
	}
	return total, nil
}

func (s *Service) buildPlan(id int64, input PlanInput) (*domain.InspectionPlan, []domain.PlanCheckItem, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.ScopeType = strings.TrimSpace(input.ScopeType)
	input.Category = strings.TrimSpace(input.Category)
	input.Period = strings.TrimSpace(input.Period)
	input.Owner = strings.TrimSpace(input.Owner)
	if input.Name == "" || input.Owner == "" {
		return nil, nil, fmt.Errorf("%w: name and owner are required", domain.ErrInvalid)
	}
	if input.ScopeType != "device" && input.ScopeType != "category" {
		return nil, nil, fmt.Errorf("%w: scopeType must be device or category", domain.ErrInvalid)
	}
	if input.ScopeType == "device" && (input.DeviceID == nil || *input.DeviceID <= 0) {
		return nil, nil, fmt.Errorf("%w: deviceId is required for device scope", domain.ErrInvalid)
	}
	if input.ScopeType == "category" && input.Category == "" {
		return nil, nil, fmt.Errorf("%w: category is required for category scope", domain.ErrInvalid)
	}
	if input.Period != "daily" && input.Period != "weekly" && input.Period != "monthly" {
		return nil, nil, fmt.Errorf("%w: period must be daily, weekly or monthly", domain.ErrInvalid)
	}
	if _, err := time.Parse("2006-01-02", input.StartDate); err != nil {
		return nil, nil, fmt.Errorf("%w: invalid startDate", domain.ErrInvalid)
	}
	if _, err := time.Parse("2006-01-02", input.EndDate); err != nil {
		return nil, nil, fmt.Errorf("%w: invalid endDate", domain.ErrInvalid)
	}
	if input.EndDate < input.StartDate {
		return nil, nil, fmt.Errorf("%w: endDate must not be before startDate", domain.ErrInvalid)
	}
	if len(input.Items) == 0 {
		return nil, nil, fmt.Errorf("%w: at least one check item is required", domain.ErrInvalid)
	}
	if input.Status == "" {
		input.Status = "enabled"
	}
	items := make([]domain.PlanCheckItem, 0, len(input.Items))
	for i, item := range input.Items {
		item.Name = strings.TrimSpace(item.Name)
		item.ItemType = strings.TrimSpace(item.ItemType)
		if item.Name == "" {
			return nil, nil, fmt.Errorf("%w: check item name is required", domain.ErrInvalid)
		}
		if item.ItemType == "" {
			item.ItemType = "text"
		}
		items = append(items, domain.PlanCheckItem{
			Name: item.Name, ItemType: item.ItemType, StandardRange: item.StandardRange,
			Required: item.Required, Description: item.Description, SortOrder: i,
		})
	}
	now := time.Now().UTC()
	plan := &domain.InspectionPlan{
		ID: id, Name: input.Name, ScopeType: input.ScopeType, DeviceID: input.DeviceID,
		Category: input.Category, Period: input.Period, StartDate: input.StartDate,
		EndDate: input.EndDate, Owner: input.Owner, Status: input.Status,
		CreatedAt: now, UpdatedAt: now,
	}
	return plan, items, nil
}

func (s *Service) generateForPlan(ctx context.Context, plan *domain.InspectionPlan, horizon time.Time) (int, error) {
	devices, err := s.scopeDevices(ctx, plan)
	if err != nil {
		return 0, err
	}
	start, _ := time.Parse("2006-01-02", plan.StartDate)
	end, _ := time.Parse("2006-01-02", plan.EndDate)
	start = start.UTC()
	end = end.UTC()
	now := time.Now().UTC()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if day.Before(start) {
		day = start
	}
	endLimit := end
	if endLimit.After(horizon) {
		endLimit = horizon
	}
	created := 0
	for !day.After(endLimit) {
		if day.After(end) {
			break
		}
		if matchesPeriod(plan.Period, day, start) {
			for _, device := range devices {
				if device.Status != "active" {
					continue
				}
				exists, err := s.tasks.ExistsByPlanDeviceDate(ctx, plan.ID, device.ID, day)
				if err != nil {
					return created, err
				}
				if exists {
					continue
				}
				results := make([]domain.TaskCheckResult, 0, len(plan.Items))
				for _, item := range plan.Items {
					results = append(results, domain.TaskCheckResult{
						ItemID: item.ID, Name: item.Name, ItemType: item.ItemType,
						StandardRange: item.StandardRange, Required: item.Required, Description: item.Description,
					})
				}
				plannedAt := time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, time.UTC)
				now := time.Now().UTC()
				task := domain.InspectionTask{
					PlanID: plan.ID, DeviceID: device.ID, PlannedAt: plannedAt,
					Status: "pending", Results: results, CreatedAt: now, UpdatedAt: now,
				}
				if err := s.tasks.Create(ctx, &task); err != nil {
					return created, err
				}
				created++
			}
		}
		day = day.AddDate(0, 0, 1)
	}
	return created, nil
}

func (s *Service) scopeDevices(ctx context.Context, plan *domain.InspectionPlan) ([]domain.Device, error) {
	if plan.ScopeType == "device" {
		device, err := s.devices.FindByID(ctx, *plan.DeviceID)
		if err != nil {
			return nil, err
		}
		return []domain.Device{*device}, nil
	}
	devices, _, err := s.devices.List(ctx, domain.DeviceFilter{Category: plan.Category, Status: "active", Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func matchesPeriod(period string, date, start time.Time) bool {
	switch period {
	case "daily":
		return true
	case "weekly":
		return date.Weekday() == start.Weekday()
	case "monthly":
		return date.Day() == start.Day()
	default:
		return false
	}
}

func newPage(items any, total int64, page, pageSize int) domain.Page {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	return domain.Page{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}
}

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
