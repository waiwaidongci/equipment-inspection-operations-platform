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

type SQLitePlanRepository struct {
	db *sql.DB
}

func NewSQLitePlanRepository(db *sql.DB) *SQLitePlanRepository {
	return &SQLitePlanRepository{db: db}
}

func (r *SQLitePlanRepository) Create(ctx context.Context, plan *domain.InspectionPlan) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO inspection_plans (
			name, scope_type, device_id, category, period, start_date, end_date,
			owner, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		plan.Name, plan.ScopeType, database.Int64Pointer(plan.DeviceID), plan.Category,
		plan.Period, plan.StartDate, plan.EndDate, plan.Owner, plan.Status,
		database.FormatTime(plan.CreatedAt), database.FormatTime(plan.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert plan: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		plan.ID = id
	}
	return nil
}

func (r *SQLitePlanRepository) Update(ctx context.Context, plan *domain.InspectionPlan) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE inspection_plans SET
			name = ?, scope_type = ?, device_id = ?, category = ?, period = ?,
			start_date = ?, end_date = ?, owner = ?, status = ?, updated_at = ?
		WHERE id = ?`,
		plan.Name, plan.ScopeType, database.Int64Pointer(plan.DeviceID), plan.Category,
		plan.Period, plan.StartDate, plan.EndDate, plan.Owner, plan.Status,
		database.FormatTime(plan.UpdatedAt), plan.ID,
	)
	if err != nil {
		return fmt.Errorf("update plan: %w", err)
	}
	return nil
}

func (r *SQLitePlanRepository) FindByID(ctx context.Context, id int64) (*domain.InspectionPlan, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, scope_type, device_id, category, period, start_date, end_date,
			owner, status, created_at, updated_at
		FROM inspection_plans WHERE id = ?`, id)
	return scanPlan(row)
}

func (r *SQLitePlanRepository) List(ctx context.Context, filter domain.PlanFilter) ([]domain.InspectionPlan, int64, error) {
	where := []string{"1 = 1"}
	var args []any
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Keyword != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+filter.Keyword+"%")
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inspection_plans WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count plans: %w", err)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	query := `
		SELECT id, name, scope_type, device_id, category, period, start_date, end_date,
			owner, status, created_at, updated_at
		FROM inspection_plans WHERE ` + whereSQL + ` ORDER BY updated_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query plans: %w", err)
	}
	defer rows.Close()
	var plans []domain.InspectionPlan
	for rows.Next() {
		plan, err := scanPlan(rows)
		if err != nil {
			return nil, 0, err
		}
		plans = append(plans, *plan)
	}
	return plans, total, rows.Err()
}

func (r *SQLitePlanRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE inspection_plans SET status = ?, updated_at = ? WHERE id = ?`,
		status, database.FormatTime(now()), id)
	if err != nil {
		return fmt.Errorf("update plan status: %w", err)
	}
	return nil
}

func (r *SQLitePlanRepository) ListEnabled(ctx context.Context) ([]domain.InspectionPlan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, scope_type, device_id, category, period, start_date, end_date,
			owner, status, created_at, updated_at
		FROM inspection_plans WHERE status = 'enabled' ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query enabled plans: %w", err)
	}
	defer rows.Close()
	var plans []domain.InspectionPlan
	for rows.Next() {
		plan, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, *plan)
	}
	return plans, rows.Err()
}

func (r *SQLitePlanRepository) ReplaceItems(ctx context.Context, planID int64, items []domain.PlanCheckItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin replace items: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM plan_check_items WHERE plan_id = ?`, planID); err != nil {
		return fmt.Errorf("delete plan items: %w", err)
	}
	for i := range items {
		items[i].PlanID = planID
		items[i].SortOrder = i
		result, err := tx.ExecContext(ctx, `
			INSERT INTO plan_check_items (plan_id, name, item_type, standard_range, required, description, sort_order)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			planID, items[i].Name, items[i].ItemType, items[i].StandardRange,
			boolToInt(items[i].Required), items[i].Description, i,
		)
		if err != nil {
			return fmt.Errorf("insert plan item: %w", err)
		}
		if itemID, err := result.LastInsertId(); err == nil {
			items[i].ID = itemID
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit plan items: %w", err)
	}
	return nil
}

func (r *SQLitePlanRepository) ListItems(ctx context.Context, planID int64) ([]domain.PlanCheckItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, plan_id, name, item_type, standard_range, required, description, sort_order
		FROM plan_check_items WHERE plan_id = ? ORDER BY sort_order`, planID)
	if err != nil {
		return nil, fmt.Errorf("query plan items: %w", err)
	}
	defer rows.Close()
	var items []domain.PlanCheckItem
	for rows.Next() {
		var (
			item     domain.PlanCheckItem
			required int
		)
		if err := rows.Scan(&item.ID, &item.PlanID, &item.Name, &item.ItemType, &item.StandardRange,
			&required, &item.Description, &item.SortOrder); err != nil {
			return nil, fmt.Errorf("scan plan item: %w", err)
		}
		item.Required = required == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPlan(row rowScanner) (*domain.InspectionPlan, error) {
	var (
		plan     domain.InspectionPlan
		deviceID sql.NullInt64
		created  string
		updated  string
	)
	if err := row.Scan(
		&plan.ID, &plan.Name, &plan.ScopeType, &deviceID, &plan.Category, &plan.Period,
		&plan.StartDate, &plan.EndDate, &plan.Owner, &plan.Status, &created, &updated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan plan: %w", err)
	}
	plan.DeviceID = database.NullInt64(deviceID)
	var err error
	plan.CreatedAt, err = database.ParseTime(created)
	if err != nil {
		return nil, fmt.Errorf("parse plan created_at: %w", err)
	}
	plan.UpdatedAt, err = database.ParseTime(updated)
	if err != nil {
		return nil, fmt.Errorf("parse plan updated_at: %w", err)
	}
	return &plan, nil
}

func now() time.Time {
	return time.Now()
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
