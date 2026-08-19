package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type BatchStore struct {
	db *sql.DB
}

func NewBatchStore(db *sql.DB) *BatchStore {
	return &BatchStore{db: db}
}

func (s *BatchStore) Create(ctx context.Context, incubatorID int64, eggCount int, species string, startDate, expectedHatch time.Time, status string) (int64, error) {
	if incubatorID <= 0 {
		return 0, &model.ValidationError{Field: "incubator_id", Message: "incubator_id is required"}
	}
	if eggCount <= 0 {
		return 0, &model.ValidationError{Field: "egg_count", Message: "egg_count must be positive"}
	}
	if !model.ValidBatchStatus(status) {
		return 0, &model.ValidationError{Field: "status", Message: "invalid batch status"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO batches(incubator_id, egg_count, species, start_date, expected_hatch_date, status) VALUES(?,?,?,?,?,?)`,
		incubatorID, eggCount, species, startDate, expectedHatch, status)
	if err != nil {
		return 0, fmt.Errorf("create batch: %w", err)
	}
	return res.LastInsertId()
}

func (s *BatchStore) GetByID(ctx context.Context, id int64) (*model.Batch, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, egg_count, species, start_date, expected_hatch_date, status, created_at FROM batches WHERE id = ?`, id)
	var b model.Batch
	err := row.Scan(&b.ID, &b.IncubatorID, &b.EggCount, &b.Species, &b.StartDate, &b.ExpectedHatchDate, &b.Status, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrBatchNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get batch %d: %w", id, err)
	}
	return &b, nil
}

func (s *BatchStore) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Batch, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, egg_count, species, start_date, expected_hatch_date, status, created_at FROM batches WHERE incubator_id = ? ORDER BY id`, incubatorID)
	if err != nil {
		return nil, fmt.Errorf("list batches by incubator: %w", err)
	}
	defer rows.Close()
	var out []*model.Batch
	for rows.Next() {
		var b model.Batch
		if err := rows.Scan(&b.ID, &b.IncubatorID, &b.EggCount, &b.Species, &b.StartDate, &b.ExpectedHatchDate, &b.Status, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}
		out = append(out, &b)
	}
	return out, rows.Err()
}

func (s *BatchStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	if !model.ValidBatchStatus(status) {
		return &model.ValidationError{Field: "status", Message: "invalid batch status"}
	}
	res, err := s.db.ExecContext(ctx, `UPDATE batches SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update batch status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrBatchNotFound
	}
	return nil
}

func (s *BatchStore) ListByStatus(ctx context.Context, status string) ([]*model.Batch, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, egg_count, species, start_date, expected_hatch_date, status, created_at FROM batches WHERE status = ? ORDER BY id`, status)
	if err != nil {
		return nil, fmt.Errorf("list batches by status: %w", err)
	}
	defer rows.Close()
	var out []*model.Batch
	for rows.Next() {
		var b model.Batch
		if err := rows.Scan(&b.ID, &b.IncubatorID, &b.EggCount, &b.Species, &b.StartDate, &b.ExpectedHatchDate, &b.Status, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan batch: %w", err)
		}
		out = append(out, &b)
	}
	return out, rows.Err()
}

func (s *BatchStore) GetActiveBatch(ctx context.Context, incubatorID int64) (*model.Batch, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, egg_count, species, start_date, expected_hatch_date, status, created_at FROM batches WHERE incubator_id = ? AND status IN ('incubating','hatching') ORDER BY id DESC LIMIT 1`, incubatorID)
	var b model.Batch
	err := row.Scan(&b.ID, &b.IncubatorID, &b.EggCount, &b.Species, &b.StartDate, &b.ExpectedHatchDate, &b.Status, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrBatchNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get active batch: %w", err)
	}
	return &b, nil
}

func (s *BatchStore) CountByStatus(ctx context.Context, status string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM batches WHERE status = ?`, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count batches by status: %w", err)
	}
	return count, nil
}
