package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/hatchery/internal/model"
)

type AlertStore struct {
	db *sql.DB
}

func NewAlertStore(db *sql.DB) *AlertStore {
	return &AlertStore{db: db}
}

func (s *AlertStore) Create(ctx context.Context, a *model.Alert) (int64, error) {
	if a.IncubatorID <= 0 {
		return 0, &model.ValidationError{Field: "incubator_id", Message: "incubator_id is required"}
	}
	if a.Message == "" {
		return 0, &model.ValidationError{Field: "message", Message: "message is required"}
	}
	if !model.ValidAlertLevel(a.Level) {
		return 0, &model.ValidationError{Field: "level", Message: "invalid alert level"}
	}
	if a.Status == "" {
		a.Status = model.AlertStatusActive
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO alerts(incubator_id, batch_id, level, message, status) VALUES(?,?,?,?,?)`,
		a.IncubatorID, a.BatchID, a.Level, a.Message, a.Status)
	if err != nil {
		return 0, fmt.Errorf("create alert: %w", err)
	}
	return res.LastInsertId()
}

func (s *AlertStore) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, batch_id, level, message, status, created_at FROM alerts WHERE id = ?`, id)
	var a model.Alert
	err := row.Scan(&a.ID, &a.IncubatorID, &a.BatchID, &a.Level, &a.Message, &a.Status, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrAlertNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get alert %d: %w", id, err)
	}
	return &a, nil
}

func (s *AlertStore) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, batch_id, level, message, status, created_at FROM alerts WHERE incubator_id = ? ORDER BY created_at DESC`, incubatorID)
	if err != nil {
		return nil, fmt.Errorf("list alerts by incubator: %w", err)
	}
	defer rows.Close()
	var out []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.IncubatorID, &a.BatchID, &a.Level, &a.Message, &a.Status, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan alert: %w", err)
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

func (s *AlertStore) ListByLevel(ctx context.Context, level string) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, batch_id, level, message, status, created_at FROM alerts WHERE level = ? AND status = 'active' ORDER BY created_at DESC`, level)
	if err != nil {
		return nil, fmt.Errorf("list alerts by level: %w", err)
	}
	defer rows.Close()
	var out []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.IncubatorID, &a.BatchID, &a.Level, &a.Message, &a.Status, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan alert: %w", err)
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

func (s *AlertStore) Acknowledge(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE alerts SET status = 'acknowledged' WHERE id = ? AND status = 'active'`, id)
	if err != nil {
		return fmt.Errorf("acknowledge alert: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrAlertNotFound
	}
	return nil
}

func (s *AlertStore) CountActive(ctx context.Context, incubatorID int64) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE incubator_id = ? AND status = 'active'`, incubatorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count active alerts: %w", err)
	}
	return count, nil
}
