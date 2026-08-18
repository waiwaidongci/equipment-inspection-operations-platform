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
	tasks     domain.TaskRepository
	plans     domain.PlanRepository
	records   domain.RecordRepository
	anomalies domain.AnomalyRepository
}

func NewService(tasks domain.TaskRepository, plans domain.PlanRepository, records domain.RecordRepository, anomalies domain.AnomalyRepository) *Service {
	return &Service{tasks: tasks, plans: plans, records: records, anomalies: anomalies}
}

type SubmitResultInput struct {
	ItemID   int64  `json:"itemId"`
	Value    string `json:"value"`
	Passed   bool   `json:"passed"`
	Abnormal bool   `json:"abnormal"`
	Remark   string `json:"remark"`
}

type SubmitInput struct {
	Results            []SubmitResultInput `json:"results"`
	Remark             string              `json:"remark"`
	Abnormal           bool                `json:"abnormal"`
	Severity           string              `json:"severity"`
	AnomalyDescription string              `json:"anomalyDescription"`
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.InspectionTask, error) {
	return s.tasks.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.TaskFilter) (domain.Page, error) {
	items, total, err := s.tasks.List(ctx, filter)
	if err != nil {
		return domain.Page{}, err
	}
	return newPage(items, total, filter.Page, filter.PageSize), nil
}

func (s *Service) Submit(ctx context.Context, id, executorID int64, input SubmitInput) (*domain.InspectionTask, error) {
	task, err := s.tasks.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.Status == "completed" || task.Status == "abnormal" {
		return nil, fmt.Errorf("%w: completed task cannot be submitted again", domain.ErrAlreadyDone)
	}
	results, err := mergeResults(task.Results, input.Results)
	if err != nil {
		return nil, err
	}
	abnormal := input.Abnormal
	for _, result := range results {
		if result.Required && strings.TrimSpace(result.Value) == "" {
			return nil, fmt.Errorf("%w: required check item %q must have a value", domain.ErrInvalid, result.Name)
		}
		if result.Abnormal {
			abnormal = true
		}
	}
	now := time.Now().UTC()
	task.Results = results
	task.Status = "completed"
	if abnormal {
		task.Status = "abnormal"
	}
	task.ActualExecutorID = &executorID
	task.ExecutedAt = &now
	task.Remark = input.Remark
	task.UpdatedAt = now
	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, err
	}
	record := &domain.InspectionRecord{
		TaskID: task.ID, DeviceID: task.DeviceID, PlanID: task.PlanID,
		ExecutorID: executorID, ExecutedAt: now, Status: task.Status,
		ResultSummary: resultSummary(task.Status, input.Remark), CreatedAt: now,
	}
	if err := s.records.Create(ctx, record); err != nil {
		return nil, err
	}
	if abnormal {
		severity := input.Severity
		if severity == "" {
			severity = "medium"
		}
		description := strings.TrimSpace(input.AnomalyDescription)
		if description == "" {
			description = "巡检发现异常，请及时处理。"
		}
		anomaly := &domain.Anomaly{
			DeviceID: task.DeviceID, TaskID: &task.ID, DiscovererID: executorID,
			DiscoveredAt: now, Severity: severity, Description: description,
			Status: "open", CreatedAt: now, UpdatedAt: now,
		}
		if err := s.anomalies.Create(ctx, anomaly); err != nil {
			return nil, err
		}
	}
	return task, nil
}

func (s *Service) SaveDraft(ctx context.Context, id, executorID int64, input SubmitInput) (*domain.InspectionTask, error) {
	task, err := s.tasks.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.Status == "completed" || task.Status == "abnormal" {
		return nil, fmt.Errorf("%w: completed task cannot be changed", domain.ErrAlreadyDone)
	}
	results, err := mergeResults(task.Results, input.Results)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	task.Results = results
	task.Status = "draft"
	task.ActualExecutorID = &executorID
	task.ExecutedAt = nil
	task.Remark = input.Remark
	task.UpdatedAt = now
	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func mergeResults(base []domain.TaskCheckResult, input []SubmitResultInput) ([]domain.TaskCheckResult, error) {
	byID := make(map[int64]SubmitResultInput, len(input))
	for _, item := range input {
		byID[item.ItemID] = item
	}
	merged := make([]domain.TaskCheckResult, 0, len(base))
	for _, current := range base {
		item, ok := byID[current.ItemID]
		if !ok {
			merged = append(merged, current)
			continue
		}
		current.Value = item.Value
		current.Passed = item.Passed
		current.Abnormal = item.Abnormal
		current.Remark = item.Remark
		merged = append(merged, current)
	}
	return merged, nil
}

func resultSummary(status, remark string) string {
	if status == "abnormal" {
		if strings.TrimSpace(remark) != "" {
			return "发现异常：" + remark
		}
		return "发现异常"
	}
	if strings.TrimSpace(remark) != "" {
		return "正常完成：" + remark
	}
	return "正常完成"
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

func IsAlreadyDone(err error) bool {
	return errors.Is(err, domain.ErrAlreadyDone)
}
