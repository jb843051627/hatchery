package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type MaintenanceStore struct {
	db *sql.DB
}

func NewMaintenanceStore(db *sql.DB) *MaintenanceStore {
	return &MaintenanceStore{db: db}
}

func (s *MaintenanceStore) Create(ctx context.Context, incubatorID int64, description string, scheduledFor time.Time) (int64, error) {
	if incubatorID <= 0 {
		return 0, &model.ValidationError{Field: "incubator_id", Message: "incubator_id is required"}
	}
	if description == "" {
		return 0, &model.ValidationError{Field: "description", Message: "description is required"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO maintenance_tasks(incubator_id, description, status, scheduled_for) VALUES(?,?,?,?)`,
		incubatorID, description, model.MaintenanceStatusPending, scheduledFor)
	if err != nil {
		return 0, fmt.Errorf("create maintenance task: %w", err)
	}
	return res.LastInsertId()
}

func (s *MaintenanceStore) GetByID(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, description, status, scheduled_for, completed_at, created_at FROM maintenance_tasks WHERE id = ?`, id)
	var m model.MaintenanceTask
	var completedAt sql.NullTime
	err := row.Scan(&m.ID, &m.IncubatorID, &m.Description, &m.Status, &m.ScheduledFor, &completedAt, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get maintenance task %d: %w", id, err)
	}
	if completedAt.Valid {
		m.CompletedAt = &completedAt.Time
	}
	return &m, nil
}

func (s *MaintenanceStore) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.MaintenanceTask, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, description, status, scheduled_for, completed_at, created_at FROM maintenance_tasks WHERE incubator_id = ? ORDER BY scheduled_for`, incubatorID)
	if err != nil {
		return nil, fmt.Errorf("list maintenance by incubator: %w", err)
	}
	defer rows.Close()
	var out []*model.MaintenanceTask
	for rows.Next() {
		var m model.MaintenanceTask
		var completedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.IncubatorID, &m.Description, &m.Status, &m.ScheduledFor, &completedAt, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maintenance: %w", err)
		}
		if completedAt.Valid {
			m.CompletedAt = &completedAt.Time
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (s *MaintenanceStore) ListPending(ctx context.Context) ([]*model.MaintenanceTask, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, description, status, scheduled_for, completed_at, created_at FROM maintenance_tasks WHERE status = 'pending' ORDER BY scheduled_for`)
	if err != nil {
		return nil, fmt.Errorf("list pending maintenance: %w", err)
	}
	defer rows.Close()
	var out []*model.MaintenanceTask
	for rows.Next() {
		var m model.MaintenanceTask
		var completedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.IncubatorID, &m.Description, &m.Status, &m.ScheduledFor, &completedAt, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maintenance: %w", err)
		}
		if completedAt.Valid {
			m.CompletedAt = &completedAt.Time
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (s *MaintenanceStore) Complete(ctx context.Context, id int64, completedAt time.Time) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE maintenance_tasks SET status = 'completed', completed_at = ? WHERE id = ? AND status = 'pending'`,
		completedAt, id)
	if err != nil {
		return fmt.Errorf("complete maintenance: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrMaintenanceNotFound
	}
	return nil
}

func (s *MaintenanceStore) BatchCreate(ctx context.Context, tasks []*model.MaintenanceTask) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO maintenance_tasks(incubator_id, description, status, scheduled_for) VALUES(?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()
	for _, t := range tasks {
		if _, err := stmt.ExecContext(ctx, t.IncubatorID, t.Description, model.MaintenanceStatusPending, t.ScheduledFor); err != nil {
			return fmt.Errorf("batch insert maintenance: %w", err)
		}
	}
	return tx.Commit()
}
