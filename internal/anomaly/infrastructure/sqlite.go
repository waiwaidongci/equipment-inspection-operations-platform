package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteAnomalyRepository struct {
	db *sql.DB
}

func NewSQLiteAnomalyRepository(db *sql.DB) *SQLiteAnomalyRepository {
	return &SQLiteAnomalyRepository{db: db}
}

func (r *SQLiteAnomalyRepository) Create(ctx context.Context, anomaly *domain.Anomaly) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO anomalies (
			device_id, task_id, discoverer_id, discovered_at, severity, description, status,
			assignee_id, progress, cause_analysis, close_note, verification_result, closed_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		anomaly.DeviceID, database.Int64Pointer(anomaly.TaskID), anomaly.DiscovererID,
		database.FormatTime(anomaly.DiscoveredAt), anomaly.Severity, anomaly.Description, anomaly.Status,
		database.Int64Pointer(anomaly.AssigneeID), anomaly.Progress, anomaly.CauseAnalysis,
		anomaly.CloseNote, anomaly.VerificationResult, nullTime(anomaly.ClosedAt),
		database.FormatTime(anomaly.CreatedAt), database.FormatTime(anomaly.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert anomaly: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		anomaly.ID = id
	}
	return nil
}

func (r *SQLiteAnomalyRepository) Update(ctx context.Context, anomaly *domain.Anomaly) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE anomalies SET
			device_id = ?, task_id = ?, discoverer_id = ?, discovered_at = ?, severity = ?,
			description = ?, status = ?, assignee_id = ?, progress = ?, cause_analysis = ?,
			close_note = ?, verification_result = ?, closed_at = ?, updated_at = ?
		WHERE id = ?`,
		anomaly.DeviceID, database.Int64Pointer(anomaly.TaskID), anomaly.DiscovererID,
		database.FormatTime(anomaly.DiscoveredAt), anomaly.Severity, anomaly.Description,
		anomaly.Status, database.Int64Pointer(anomaly.AssigneeID), anomaly.Progress,
		anomaly.CauseAnalysis, anomaly.CloseNote, anomaly.VerificationResult,
		nullTime(anomaly.ClosedAt), database.FormatTime(anomaly.UpdatedAt), anomaly.ID,
	)
	if err != nil {
		return fmt.Errorf("update anomaly: %w", err)
	}
	return nil
}

func (r *SQLiteAnomalyRepository) FindByID(ctx context.Context, id int64) (*domain.Anomaly, error) {
	row := r.db.QueryRowContext(ctx, anomalySelect+" WHERE a.id = ?", id)
	return scanAnomaly(row)
}

func (r *SQLiteAnomalyRepository) List(ctx context.Context, filter domain.AnomalyFilter) ([]domain.Anomaly, int64, error) {
	where := []string{"1 = 1"}
	var args []any
	if filter.Status != "" {
		where = append(where, "a.status = ?")
		args = append(args, filter.Status)
	}
	if filter.DeviceID > 0 {
		where = append(where, "a.device_id = ?")
		args = append(args, filter.DeviceID)
	}
	if filter.Severity != "" {
		where = append(where, "a.severity = ?")
		args = append(args, filter.Severity)
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM anomalies a WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count anomalies: %w", err)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	query := anomalySelect + " WHERE " + whereSQL + " ORDER BY a.updated_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query anomalies: %w", err)
	}
	defer rows.Close()
	var anomalies []domain.Anomaly
	for rows.Next() {
		anomaly, err := scanAnomaly(rows)
		if err != nil {
			return nil, 0, err
		}
		anomalies = append(anomalies, *anomaly)
	}
	return anomalies, total, rows.Err()
}

func (r *SQLiteAnomalyRepository) ListByDevice(ctx context.Context, deviceID int64) ([]domain.Anomaly, error) {
	rows, err := r.db.QueryContext(ctx, anomalySelect+" WHERE a.device_id = ? ORDER BY a.updated_at DESC", deviceID)
	if err != nil {
		return nil, fmt.Errorf("query device anomalies: %w", err)
	}
	defer rows.Close()
	var anomalies []domain.Anomaly
	for rows.Next() {
		anomaly, err := scanAnomaly(rows)
		if err != nil {
			return nil, err
		}
		anomalies = append(anomalies, *anomaly)
	}
	return anomalies, rows.Err()
}

const anomalySelect = `
	SELECT a.id, a.device_id, a.task_id, a.discoverer_id, a.discovered_at, a.severity,
		a.description, a.status, a.assignee_id, a.progress, a.cause_analysis, a.close_note,
		a.verification_result, a.closed_at, a.created_at, a.updated_at
	FROM anomalies a`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAnomaly(row rowScanner) (*domain.Anomaly, error) {
	var (
		anomaly    domain.Anomaly
		taskID     sql.NullInt64
		assigneeID sql.NullInt64
		discovered string
		closedAt   sql.NullString
		createdAt  string
		updatedAt  string
	)
	if err := row.Scan(
		&anomaly.ID, &anomaly.DeviceID, &taskID, &anomaly.DiscovererID, &discovered,
		&anomaly.Severity, &anomaly.Description, &anomaly.Status, &assigneeID,
		&anomaly.Progress, &anomaly.CauseAnalysis, &anomaly.CloseNote,
		&anomaly.VerificationResult, &closedAt, &createdAt, &updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan anomaly: %w", err)
	}
	var err error
	anomaly.TaskID = database.NullInt64(taskID)
	anomaly.AssigneeID = database.NullInt64(assigneeID)
	anomaly.DiscoveredAt, err = database.ParseTime(discovered)
	if err != nil {
		return nil, fmt.Errorf("parse anomaly discovered_at: %w", err)
	}
	anomaly.ClosedAt, err = database.NullTime(closedAt)
	if err != nil {
		return nil, fmt.Errorf("parse anomaly closed_at: %w", err)
	}
	anomaly.CreatedAt, err = database.ParseTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse anomaly created_at: %w", err)
	}
	anomaly.UpdatedAt, err = database.ParseTime(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse anomaly updated_at: %w", err)
	}
	return &anomaly, nil
}

func nullTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return database.FormatTime(*value)
}
