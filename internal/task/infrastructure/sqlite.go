package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

func (r *SQLiteTaskRepository) Create(ctx context.Context, task *domain.InspectionTask) error {
	return r.createOne(ctx, task)
}

func (r *SQLiteTaskRepository) CreateBatch(ctx context.Context, tasks []domain.InspectionTask) error {
	if len(tasks) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin task batch: %w", err)
	}
	defer tx.Rollback()
	for i := range tasks {
		if err := r.createOneTx(ctx, tx, &tasks[i]); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit task batch: %w", err)
	}
	return nil
}

func (r *SQLiteTaskRepository) createOne(ctx context.Context, task *domain.InspectionTask) error {
	return r.createOneTx(ctx, r.db, task)
}

func (r *SQLiteTaskRepository) createOneTx(ctx context.Context, executor execContext, task *domain.InspectionTask) error {
	results := task.Results
	if results == nil {
		results = []domain.TaskCheckResult{}
	}
	raw, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("marshal task results: %w", err)
	}
	var executedAt any
	if task.ExecutedAt != nil {
		executedAt = database.FormatTime(*task.ExecutedAt)
	}
	result, err := executor.ExecContext(ctx, `
		INSERT INTO inspection_tasks (
			plan_id, device_id, planned_at, status, actual_executor_id, executed_at,
			results, remark, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(plan_id, device_id, planned_at) DO NOTHING`,
		task.PlanID, task.DeviceID, database.FormatTime(task.PlannedAt), task.Status,
		database.Int64Pointer(task.ActualExecutorID), executedAt, string(raw), task.Remark,
		database.FormatTime(task.CreatedAt), database.FormatTime(task.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		task.ID = id
	}
	return nil
}

type execContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (r *SQLiteTaskRepository) FindByID(ctx context.Context, id int64) (*domain.InspectionTask, error) {
	row := r.db.QueryRowContext(ctx, taskSelect+" WHERE t.id = ?", id)
	return scanTask(row)
}

func (r *SQLiteTaskRepository) List(ctx context.Context, filter domain.TaskFilter) ([]domain.InspectionTask, int64, error) {
	where := []string{"1 = 1"}
	var args []any
	if filter.Status != "" {
		where = append(where, "t.status = ?")
		args = append(args, filter.Status)
	}
	if filter.DeviceID > 0 {
		where = append(where, "t.device_id = ?")
		args = append(args, filter.DeviceID)
	}
	if filter.ExecutorID > 0 {
		where = append(where, "t.actual_executor_id = ?")
		args = append(args, filter.ExecutorID)
	}
	if !filter.From.IsZero() {
		where = append(where, "t.planned_at >= ?")
		args = append(args, database.FormatTime(filter.From))
	}
	if !filter.To.IsZero() {
		where = append(where, "t.planned_at <= ?")
		args = append(args, database.FormatTime(filter.To))
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inspection_tasks t WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	query := taskSelect + " WHERE " + whereSQL + " ORDER BY t.planned_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()
	var tasks []domain.InspectionTask
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, total, rows.Err()
}

func (r *SQLiteTaskRepository) Update(ctx context.Context, task *domain.InspectionTask) error {
	raw, err := json.Marshal(task.Results)
	if err != nil {
		return fmt.Errorf("marshal task results: %w", err)
	}
	var executedAt any
	if task.ExecutedAt != nil {
		executedAt = database.FormatTime(*task.ExecutedAt)
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE inspection_tasks SET
			status = ?, actual_executor_id = ?, executed_at = ?, results = ?, remark = ?, updated_at = ?
		WHERE id = ?`,
		task.Status, database.Int64Pointer(task.ActualExecutorID), executedAt, string(raw),
		task.Remark, database.FormatTime(task.UpdatedAt), task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (r *SQLiteTaskRepository) ExistsByPlanDeviceDate(ctx context.Context, planID, deviceID int64, plannedAt time.Time) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inspection_tasks WHERE plan_id = ? AND device_id = ? AND planned_at = ?`,
		planID, deviceID, database.FormatTime(plannedAt),
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check duplicate task: %w", err)
	}
	return exists > 0, nil
}

func (r *SQLiteTaskRepository) ListPendingDue(ctx context.Context, before time.Time) ([]domain.InspectionTask, error) {
	rows, err := r.db.QueryContext(ctx, taskSelect+" WHERE t.status IN ('pending','draft') AND t.planned_at < ? ORDER BY t.planned_at",
		database.FormatTime(before))
	if err != nil {
		return nil, fmt.Errorf("query due tasks: %w", err)
	}
	defer rows.Close()
	var tasks []domain.InspectionTask
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, rows.Err()
}

const taskSelect = `
	SELECT t.id, t.plan_id, t.device_id, t.planned_at, t.status, t.actual_executor_id,
		t.executed_at, t.results, t.remark, t.created_at, t.updated_at
	FROM inspection_tasks t`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTask(row rowScanner) (*domain.InspectionTask, error) {
	var (
		task       domain.InspectionTask
		plannedAt  string
		executorID sql.NullInt64
		executedAt sql.NullString
		resultsRaw string
		createdAt  string
		updatedAt  string
	)
	if err := row.Scan(
		&task.ID, &task.PlanID, &task.DeviceID, &plannedAt, &task.Status, &executorID,
		&executedAt, &resultsRaw, &task.Remark, &createdAt, &updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan task: %w", err)
	}
	var err error
	task.PlannedAt, err = database.ParseTime(plannedAt)
	if err != nil {
		return nil, fmt.Errorf("parse task planned_at: %w", err)
	}
	task.ActualExecutorID = database.NullInt64(executorID)
	task.ExecutedAt, err = database.NullTime(executedAt)
	if err != nil {
		return nil, fmt.Errorf("parse task executed_at: %w", err)
	}
	task.CreatedAt, err = database.ParseTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse task created_at: %w", err)
	}
	task.UpdatedAt, err = database.ParseTime(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse task updated_at: %w", err)
	}
	if err := json.Unmarshal([]byte(resultsRaw), &task.Results); err != nil {
		return nil, fmt.Errorf("unmarshal task results: %w", err)
	}
	if task.Results == nil {
		task.Results = []domain.TaskCheckResult{}
	}
	return &task, nil
}
