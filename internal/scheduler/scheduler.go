package scheduler

import (
	"context"
	"log/slog"
	"time"

	planapp "github.com/example/inspection-platform/internal/plan/application"
)

type Scheduler struct {
	service  *planapp.Service
	interval time.Duration
	done     chan struct{}
}

func New(service *planapp.Service, interval time.Duration) *Scheduler {
	return &Scheduler{service: service, interval: interval, done: make(chan struct{})}
}

func (s *Scheduler) Start(ctx context.Context) {
	if s.interval <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			case <-ticker.C:
				count, err := s.service.GenerateDue(ctx)
				if err != nil {
					slog.Error("scheduler generation failed", "error", err)
					continue
				}
				if count > 0 {
					slog.Info("scheduler generated tasks", "count", count)
				}
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.done)
}
