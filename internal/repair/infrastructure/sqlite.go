package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteRepairRepository struct {
	db *sql.DB
}

func NewSQLiteRepairRepository(db *sql.DB) *SQLiteRepairRepository {
	return &SQLiteRepairRepository{db: db}
}

func (r *SQLiteRepairRepository) Create(ctx context.Context, repair *domain.Repair) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO repairs (
			device_id, anomaly_id, repair_type, content, vendor, started_at, ended_at,
			cost, result, remark, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		repair.DeviceID, database.Int64Pointer(repair.AnomalyID), repair.RepairType,
		repair.Content, repair.Vendor, repair.StartedAt, repair.EndedAt, repair.Cost,
		repair.Result, repair.Remark, repair.CreatedBy,
		database.FormatTime(repair.CreatedAt), database.FormatTime(repair.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert repair: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		repair.ID = id
	}
	return nil
}

func (r *SQLiteRepairRepository) Update(ctx context.Context, repair *domain.Repair) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE repairs SET
			device_id = ?, anomaly_id = ?, repair_type = ?, content = ?, vendor = ?,
			started_at = ?, ended_at = ?, cost = ?, result = ?, remark = ?, updated_at = ?
		WHERE id = ?`,
		repair.DeviceID, database.Int64Pointer(repair.AnomalyID), repair.RepairType,
		repair.Content, repair.Vendor, repair.StartedAt, repair.EndedAt, repair.Cost,
		repair.Result, repair.Remark, database.FormatTime(repair.UpdatedAt), repair.ID,
	)
	if err != nil {
		return fmt.Errorf("update repair: %w", err)
	}
	return nil
}

func (r *SQLiteRepairRepository) FindByID(ctx context.Context, id int64) (*domain.Repair, error) {
	row := r.db.QueryRowContext(ctx, repairSelect+" WHERE r.id = ?", id)
	return scanRepair(row)
}

func (r *SQLiteRepairRepository) List(ctx context.Context, filter domain.RepairFilter) ([]domain.Repair, int64, error) {
	where := []string{"1 = 1"}
	var args []any
	if filter.DeviceID > 0 {
		where = append(where, "r.device_id = ?")
		args = append(args, filter.DeviceID)
	}
	if filter.Keyword != "" {
		where = append(where, "(r.content LIKE ? OR r.vendor LIKE ? OR r.result LIKE ?)")
		like := "%" + filter.Keyword + "%"
		args = append(args, like, like, like)
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM repairs r WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count repairs: %w", err)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	query := repairSelect + " WHERE " + whereSQL + " ORDER BY r.updated_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query repairs: %w", err)
	}
	defer rows.Close()
	var repairs []domain.Repair
	for rows.Next() {
		repair, err := scanRepair(rows)
		if err != nil {
			return nil, 0, err
		}
		repairs = append(repairs, *repair)
	}
	return repairs, total, rows.Err()
}

func (r *SQLiteRepairRepository) ListByDevice(ctx context.Context, deviceID int64) ([]domain.Repair, error) {
	rows, err := r.db.QueryContext(ctx, repairSelect+" WHERE r.device_id = ? ORDER BY r.updated_at DESC", deviceID)
	if err != nil {
		return nil, fmt.Errorf("query device repairs: %w", err)
	}
	defer rows.Close()
	var repairs []domain.Repair
	for rows.Next() {
		repair, err := scanRepair(rows)
		if err != nil {
			return nil, err
		}
		repairs = append(repairs, *repair)
	}
	return repairs, rows.Err()
}

const repairSelect = `
	SELECT r.id, r.device_id, r.anomaly_id, r.repair_type, r.content, r.vendor,
		r.started_at, r.ended_at, r.cost, r.result, r.remark, r.created_by,
		r.created_at, r.updated_at
	FROM repairs r`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRepair(row rowScanner) (*domain.Repair, error) {
	var (
		repair    domain.Repair
		anomalyID sql.NullInt64
		createdAt string
		updatedAt string
	)
	if err := row.Scan(
		&repair.ID, &repair.DeviceID, &anomalyID, &repair.RepairType, &repair.Content,
		&repair.Vendor, &repair.StartedAt, &repair.EndedAt, &repair.Cost, &repair.Result,
		&repair.Remark, &repair.CreatedBy, &createdAt, &updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan repair: %w", err)
	}
	repair.AnomalyID = database.NullInt64(anomalyID)
	var err error
	repair.CreatedAt, err = database.ParseTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse repair created_at: %w", err)
	}
	repair.UpdatedAt, err = database.ParseTime(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse repair updated_at: %w", err)
	}
	return &repair, nil
}
