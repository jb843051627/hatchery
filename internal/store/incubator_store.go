package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type IncubatorStore struct {
	db *sql.DB
}

func NewIncubatorStore(db *sql.DB) *IncubatorStore {
	return &IncubatorStore{db: db}
}

func (s *IncubatorStore) Create(ctx context.Context, name, location string, capacity int, status string, installedAt time.Time) (int64, error) {
	if name == "" {
		return 0, &model.ValidationError{Field: "name", Message: "name is required"}
	}
	if capacity <= 0 {
		return 0, &model.ValidationError{Field: "capacity", Message: "capacity must be positive"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO incubators(name, capacity, status, location, installed_at) VALUES(?,?,?,?,?)`,
		name, capacity, status, location, installedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("create incubator: %w", err)
	}
	return res.LastInsertId()
}

func (s *IncubatorStore) GetByID(ctx context.Context, id int64) (*model.Incubator, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, capacity, status, location, installed_at, created_at FROM incubators WHERE id = ?`, id)
	var inc model.Incubator
	err := row.Scan(&inc.ID, &inc.Name, &inc.Capacity, &inc.Status, &inc.Location, &inc.InstalledAt, &inc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrIncubatorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get incubator %d: %w", id, err)
	}
	return &inc, nil
}

func (s *IncubatorStore) ListAll(ctx context.Context) ([]*model.Incubator, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, capacity, status, location, installed_at, created_at FROM incubators ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list incubators: %w", err)
	}
	defer rows.Close()
	var out []*model.Incubator
	for rows.Next() {
		var inc model.Incubator
		if err := rows.Scan(&inc.ID, &inc.Name, &inc.Capacity, &inc.Status, &inc.Location, &inc.InstalledAt, &inc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan incubator: %w", err)
		}
		out = append(out, &inc)
	}
	return out, rows.Err()
}

func (s *IncubatorStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	if status != model.IncubatorStatusActive && status != model.IncubatorStatusInactive && status != model.IncubatorStatusMaintenance {
		return &model.ValidationError{Field: "status", Message: "invalid status"}
	}
	res, err := s.db.ExecContext(ctx, `UPDATE incubators SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update incubator status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrIncubatorNotFound
	}
	return nil
}

func (s *IncubatorStore) Count(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM incubators`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count incubators: %w", err)
	}
	return count, nil
}

func (s *IncubatorStore) ListByStatus(ctx context.Context, status string) ([]*model.Incubator, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, capacity, status, location, installed_at, created_at FROM incubators WHERE status = ? ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list incubators by status: %w", err)
	}
	defer rows.Close()
	var out []*model.Incubator
	for rows.Next() {
		var inc model.Incubator
		if err := rows.Scan(&inc.ID, &inc.Name, &inc.Capacity, &inc.Status, &inc.Location, &inc.InstalledAt, &inc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan incubator: %w", err)
		}
		out = append(out, &inc)
	}
	return out, rows.Err()
}
