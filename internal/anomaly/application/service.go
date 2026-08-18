package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
)

type Service struct {
	anomalies domain.AnomalyRepository
}

func NewService(anomalies domain.AnomalyRepository) *Service {
	return &Service{anomalies: anomalies}
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.Anomaly, error) {
	return s.anomalies.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.AnomalyFilter) (domain.Page, error) {
	items, total, err := s.anomalies.List(ctx, filter)
	if err != nil {
		return domain.Page{}, err
	}
	return newPage(items, total, filter.Page, filter.PageSize), nil
}

func (s *Service) Assign(ctx context.Context, id, assigneeID int64) (*domain.Anomaly, error) {
	anomaly, err := s.anomalies.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if anomaly.Status == "closed" {
		return nil, fmt.Errorf("%w: closed anomaly cannot be assigned", domain.ErrForbidden)
	}
	if assigneeID <= 0 {
		return nil, fmt.Errorf("%w: assigneeId is required", domain.ErrInvalid)
	}
	anomaly.AssigneeID = &assigneeID
	anomaly.Status = "assigned"
	anomaly.UpdatedAt = time.Now().UTC()
	if err := s.anomalies.Update(ctx, anomaly); err != nil {
		return nil, err
	}
	return anomaly, nil
}

func (s *Service) Progress(ctx context.Context, id int64, progress string) (*domain.Anomaly, error) {
	anomaly, err := s.anomalies.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if anomaly.Status == "closed" {
		return nil, fmt.Errorf("%w: closed anomaly cannot be updated", domain.ErrForbidden)
	}
	progress = strings.TrimSpace(progress)
	if progress == "" {
		return nil, fmt.Errorf("%w: progress is required", domain.ErrInvalid)
	}
	anomaly.Progress = progress
	anomaly.Status = "in_progress"
	anomaly.UpdatedAt = time.Now().UTC()
	if err := s.anomalies.Update(ctx, anomaly); err != nil {
		return nil, err
	}
	return anomaly, nil
}

type CloseInput struct {
	CauseAnalysis      string `json:"causeAnalysis"`
	CloseNote          string `json:"closeNote"`
	VerificationResult string `json:"verificationResult"`
}

func (s *Service) Close(ctx context.Context, id int64, input CloseInput) (*domain.Anomaly, error) {
	anomaly, err := s.anomalies.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if anomaly.Status == "closed" {
		return nil, fmt.Errorf("%w: anomaly is already closed", domain.ErrAlreadyDone)
	}
	if strings.TrimSpace(input.CauseAnalysis) == "" || strings.TrimSpace(input.CloseNote) == "" || strings.TrimSpace(input.VerificationResult) == "" {
		return nil, fmt.Errorf("%w: causeAnalysis, closeNote and verificationResult are required", domain.ErrInvalid)
	}
	now := time.Now().UTC()
	anomaly.CauseAnalysis = input.CauseAnalysis
	anomaly.CloseNote = input.CloseNote
	anomaly.VerificationResult = input.VerificationResult
	anomaly.Status = "closed"
	anomaly.ClosedAt = &now
	anomaly.UpdatedAt = now
	if err := s.anomalies.Update(ctx, anomaly); err != nil {
		return nil, err
	}
	return anomaly, nil
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
