package deviceprobe

import (
	"context"
	"testing"

	"github.com/example/inspection-platform/internal/domain"
)

type missingBackend struct{}

func (missingBackend) FindDevice(context.Context, string) (string, error) {
	return "", domain.ErrNotFound
}

func TestMissingDeviceReturnsNotFoundWithoutRetry(t *testing.T) {
	service := NewService(NewRepository(missingBackend{}))
	attempts := 0
	response := Request(context.Background(), service, "EQ-404", &attempts)
	if response.Status != 404 {
		t.Fatalf("status = %d, want 404", response.Status)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
