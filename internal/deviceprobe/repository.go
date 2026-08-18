package deviceprobe

import (
	"context"
	"fmt"

	"github.com/example/inspection-platform/internal/domain"
)

type Backend interface {
	FindDevice(context.Context, string) (string, error)
}

type Repository struct{ backend Backend }

func NewRepository(backend Backend) *Repository { return &Repository{backend: backend} }

func (repository *Repository) Find(ctx context.Context, code string) (string, error) {
	name, err := repository.backend.FindDevice(ctx, code)
	if err != nil {
		return "", fmt.Errorf("device lookup failed: %w", err)
	}
	if name == "" {
		return "", fmt.Errorf("device lookup failed: %w", domain.ErrNotFound)
	}
	return name, nil
}
