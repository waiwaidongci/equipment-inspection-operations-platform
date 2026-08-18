package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteDeviceRepository struct {
	db *sql.DB
}

func NewSQLiteDeviceRepository(db *sql.DB) *SQLiteDeviceRepository {
	return &SQLiteDeviceRepository{db: db}
}

func (r *SQLiteDeviceRepository) Create(ctx context.Context, device *domain.Device) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO devices (
			code, name, category, model, serial_number, manufacturer, location,
			install_date, warranty_end, owner, status, remark, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		device.Code, device.Name, device.Category, device.Model, device.SerialNumber,
		device.Manufacturer, device.Location, device.InstallDate, device.WarrantyEnd,
		device.Owner, device.Status, device.Remark,
		database.FormatTime(device.CreatedAt), database.FormatTime(device.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		device.ID = id
	}
	return nil
}

func (r *SQLiteDeviceRepository) Update(ctx context.Context, device *domain.Device) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET
			code = ?, name = ?, category = ?, model = ?, serial_number = ?, manufacturer = ?,
			location = ?, install_date = ?, warranty_end = ?, owner = ?, status = ?, remark = ?,
			updated_at = ?
		WHERE id = ?`,
		device.Code, device.Name, device.Category, device.Model, device.SerialNumber,
		device.Manufacturer, device.Location, device.InstallDate, device.WarrantyEnd,
		device.Owner, device.Status, device.Remark, database.FormatTime(device.UpdatedAt), device.ID,
	)
	if err != nil {
		return fmt.Errorf("update device: %w", err)
	}
	return nil
}

func (r *SQLiteDeviceRepository) FindByID(ctx context.Context, id int64) (*domain.Device, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, category, model, serial_number, manufacturer, location,
			install_date, warranty_end, owner, status, remark, created_at, updated_at
		FROM devices WHERE id = ?`, id)
	return scanDevice(row)
}

func (r *SQLiteDeviceRepository) FindByCode(ctx context.Context, code string) (*domain.Device, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, code, name, category, model, serial_number, manufacturer, location,
			install_date, warranty_end, owner, status, remark, created_at, updated_at
		FROM devices WHERE code = ?`, code)
	return scanDevice(row)
}

func (r *SQLiteDeviceRepository) List(ctx context.Context, filter domain.DeviceFilter) ([]domain.Device, int64, error) {
	where := []string{"1 = 1"}
	var args []any
	if filter.Query != "" {
		where = append(where, "(code LIKE ? OR name LIKE ? OR serial_number LIKE ?)")
		like := "%" + filter.Query + "%"
		args = append(args, like, like, like)
	}
	if filter.Category != "" {
		where = append(where, "category = ?")
		args = append(args, filter.Category)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Owner != "" {
		where = append(where, "owner = ?")
		args = append(args, filter.Owner)
	}
	if filter.Location != "" {
		where = append(where, "location = ?")
		args = append(args, filter.Location)
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count devices: %w", err)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	query := `
		SELECT id, code, name, category, model, serial_number, manufacturer, location,
			install_date, warranty_end, owner, status, remark, created_at, updated_at
		FROM devices WHERE ` + whereSQL + ` ORDER BY updated_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query devices: %w", err)
	}
	defer rows.Close()
	var devices []domain.Device
	for rows.Next() {
		device, err := scanDevice(rows)
		if err != nil {
			return nil, 0, err
		}
		devices = append(devices, *device)
	}
	return devices, total, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanDevice(row rowScanner) (*domain.Device, error) {
	var (
		device  domain.Device
		created string
		updated string
	)
	if err := row.Scan(
		&device.ID, &device.Code, &device.Name, &device.Category, &device.Model,
		&device.SerialNumber, &device.Manufacturer, &device.Location, &device.InstallDate,
		&device.WarrantyEnd, &device.Owner, &device.Status, &device.Remark, &created, &updated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan device: %w", err)
	}
	var err error
	device.CreatedAt, err = database.ParseTime(created)
	if err != nil {
		return nil, fmt.Errorf("parse device created_at: %w", err)
	}
	device.UpdatedAt, err = database.ParseTime(updated)
	if err != nil {
		return nil, fmt.Errorf("parse device updated_at: %w", err)
	}
	return &device, nil
}
