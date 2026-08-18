package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteRecordRepository struct {
	db *sql.DB
}

func NewSQLiteRecordRepository(db *sql.DB) *SQLiteRecordRepository {
	return &SQLiteRecordRepository{db: db}
}

func (r *SQLiteRecordRepository) Create(ctx context.Context, record *domain.InspectionRecord) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO inspection_records (
			task_id, device_id, plan_id, executor_id, executed_at, status, result_summary, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		record.TaskID, record.DeviceID, record.PlanID, record.ExecutorID,
		database.FormatTime(record.ExecutedAt), record.Status, record.ResultSummary,
		database.FormatTime(record.CreatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert inspection record: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		record.ID = id
	}
	return nil
}

func (r *SQLiteRecordRepository) ListByDevice(ctx context.Context, deviceID int64, limit, offset int) ([]domain.InspectionRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, task_id, device_id, plan_id, executor_id, executed_at, status, result_summary, created_at
		FROM inspection_records WHERE device_id = ? ORDER BY executed_at DESC LIMIT ? OFFSET ?`,
		deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query device records: %w", err)
	}
	defer rows.Close()
	var records []domain.InspectionRecord
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}
	return records, rows.Err()
}

func (r *SQLiteRecordRepository) ListByTask(ctx context.Context, taskID int64) ([]domain.InspectionRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, task_id, device_id, plan_id, executor_id, executed_at, status, result_summary, created_at
		FROM inspection_records WHERE task_id = ? ORDER BY executed_at DESC`, taskID)
	if err != nil {
		return nil, fmt.Errorf("query task records: %w", err)
	}
	defer rows.Close()
	var records []domain.InspectionRecord
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}
	return records, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRecord(row rowScanner) (*domain.InspectionRecord, error) {
	var (
		record     domain.InspectionRecord
		executedAt string
		createdAt  string
	)
	if err := row.Scan(
		&record.ID, &record.TaskID, &record.DeviceID, &record.PlanID, &record.ExecutorID,
		&executedAt, &record.Status, &record.ResultSummary, &createdAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan record: %w", err)
	}
	var err error
	record.ExecutedAt, err = database.ParseTime(executedAt)
	if err != nil {
		return nil, fmt.Errorf("parse record executed_at: %w", err)
	}
	record.CreatedAt, err = database.ParseTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse record created_at: %w", err)
	}
	return &record, nil
}
