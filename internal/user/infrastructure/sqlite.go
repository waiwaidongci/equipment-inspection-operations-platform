package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/database"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, display_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.PasswordHash, user.DisplayName,
		database.FormatTime(user.CreatedAt), database.FormatTime(user.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	if id, err := result.LastInsertId(); err == nil {
		user.ID = id
	}
	return nil
}

func (r *SQLiteUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, display_name, created_at, updated_at
		FROM users WHERE username = ?`, username)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, password_hash, display_name, created_at, updated_at
		FROM users WHERE id = ?`, id)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, password_hash, display_name, created_at, updated_at
		FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (*domain.User, error) {
	var (
		user    domain.User
		created string
		updated string
	)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	var err error
	user.CreatedAt, err = database.ParseTime(created)
	if err != nil {
		return nil, fmt.Errorf("parse user created_at: %w", err)
	}
	user.UpdatedAt, err = database.ParseTime(updated)
	if err != nil {
		return nil, fmt.Errorf("parse user updated_at: %w", err)
	}
	return &user, nil
}
