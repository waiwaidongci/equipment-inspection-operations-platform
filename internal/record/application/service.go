package application

import (
	"context"

	"github.com/example/inspection-platform/internal/domain"
)

type Service struct {
	records domain.RecordRepository
}

func NewService(records domain.RecordRepository) *Service {
	return &Service{records: records}
}

func (s *Service) ListByDevice(ctx context.Context, deviceID int64, limit, offset int) ([]domain.InspectionRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return s.records.ListByDevice(ctx, deviceID, limit, offset)
}
