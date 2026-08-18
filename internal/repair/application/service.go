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
	repairs   domain.RepairRepository
	anomalies domain.AnomalyRepository
}

func NewService(repairs domain.RepairRepository, anomalies domain.AnomalyRepository) *Service {
	return &Service{repairs: repairs, anomalies: anomalies}
}

type RepairInput struct {
	DeviceID   int64   `json:"deviceId"`
	AnomalyID  *int64  `json:"anomalyId,omitempty"`
	RepairType string  `json:"repairType"`
	Content    string  `json:"content"`
	Vendor     string  `json:"vendor"`
	StartedAt  string  `json:"startedAt"`
	EndedAt    string  `json:"endedAt"`
	Cost       float64 `json:"cost"`
	Result     string  `json:"result"`
	Remark     string  `json:"remark"`
}

func (s *Service) Create(ctx context.Context, creatorID int64, input RepairInput) (*domain.Repair, error) {
	if err := validateRepairInput(&input); err != nil {
		return nil, err
	}
	if input.AnomalyID != nil {
		anomaly, err := s.anomalies.FindByID(ctx, *input.AnomalyID)
		if err != nil {
			return nil, err
		}
		if anomaly.Status != "closed" {
			return nil, fmt.Errorf("%w: only closed anomaly can be linked to a repair", domain.ErrForbidden)
		}
		if anomaly.DeviceID != input.DeviceID {
			return nil, fmt.Errorf("%w: anomaly does not belong to the selected device", domain.ErrInvalid)
		}
	}
	now := time.Now().UTC()
	repair := &domain.Repair{
		DeviceID: input.DeviceID, AnomalyID: input.AnomalyID, RepairType: input.RepairType,
		Content: input.Content, Vendor: input.Vendor, StartedAt: input.StartedAt,
		EndedAt: input.EndedAt, Cost: input.Cost, Result: input.Result, Remark: input.Remark,
		CreatedBy: creatorID, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repairs.Create(ctx, repair); err != nil {
		return nil, err
	}
	return repair, nil
}

func (s *Service) Update(ctx context.Context, id int64, input RepairInput) (*domain.Repair, error) {
	repair, err := s.repairs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := validateRepairInput(&input); err != nil {
		return nil, err
	}
	if input.AnomalyID != nil && repair.AnomalyID != nil && *repair.AnomalyID != *input.AnomalyID {
		return nil, fmt.Errorf("%w: linked anomaly cannot be changed after creation", domain.ErrForbidden)
	}
	repair.DeviceID, repair.AnomalyID, repair.RepairType = input.DeviceID, input.AnomalyID, input.RepairType
	repair.Content, repair.Vendor, repair.StartedAt = input.Content, input.Vendor, input.StartedAt
	repair.EndedAt, repair.Cost, repair.Result, repair.Remark = input.EndedAt, input.Cost, input.Result, input.Remark
	repair.UpdatedAt = time.Now().UTC()
	if err := s.repairs.Update(ctx, repair); err != nil {
		return nil, err
	}
	return repair, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.Repair, error) {
	return s.repairs.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.RepairFilter) (domain.Page, error) {
	items, total, err := s.repairs.List(ctx, filter)
	if err != nil {
		return domain.Page{}, err
	}
	return newPage(items, total, filter.Page, filter.PageSize), nil
}

func validateRepairInput(input *RepairInput) error {
	input.RepairType = strings.TrimSpace(input.RepairType)
	input.Content = strings.TrimSpace(input.Content)
	input.Result = strings.TrimSpace(input.Result)
	if input.DeviceID <= 0 || input.RepairType == "" || input.Content == "" || input.Result == "" {
		return fmt.Errorf("%w: deviceId, repairType, content and result are required", domain.ErrInvalid)
	}
	if _, err := time.Parse("2006-01-02", input.StartedAt); err != nil {
		return fmt.Errorf("%w: invalid startedAt", domain.ErrInvalid)
	}
	if _, err := time.Parse("2006-01-02", input.EndedAt); err != nil {
		return fmt.Errorf("%w: invalid endedAt", domain.ErrInvalid)
	}
	if input.EndedAt < input.StartedAt {
		return fmt.Errorf("%w: endedAt must not be before startedAt", domain.ErrInvalid)
	}
	return nil
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

func IsNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
