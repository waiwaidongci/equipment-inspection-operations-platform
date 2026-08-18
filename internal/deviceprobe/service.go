package deviceprobe

import (
	"context"
	"errors"
	"fmt"

	"github.com/example/inspection-platform/internal/domain"
)

var ErrDeviceMissing = errors.New("device missing")

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (service *Service) Lookup(ctx context.Context, code string) (string, error) {
	name, err := service.repository.Find(ctx, code)
	if err == nil {
		return name, nil
	}
	if errors.Is(err, domain.ErrNotFound) {
		return "", fmt.Errorf("%w: %s", ErrDeviceMissing, code)
	}
	return "", err
}
