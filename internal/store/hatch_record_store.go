package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type HatchRecordStore struct {
	db    *sql.DB
	mu    sync.RWMutex
	cache map[int64][]*model.HatchRecord
}

func NewHatchRecordStore(db *sql.DB) *HatchRecordStore {
	return &HatchRecordStore{db: db, cache: make(map[int64][]*model.HatchRecord)}
}

func (s *HatchRecordStore) Create(ctx context.Context, r *model.HatchRecord) (int64, error) {
	if r.BatchID <= 0 {
		return 0, &model.ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}
	if r.HatchedCount < 0 || r.HealthyCount < 0 || r.WeakCount < 0 || r.DeadCount < 0 {
		return 0, &model.ValidationError{Field: "counts", Message: "counts must be non-negative"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO hatch_records(batch_id, hatched_count, healthy_count, weak_count, dead_count, hatch_date, graded_by) VALUES(?,?,?,?,?,?,?)`,
		r.BatchID, r.HatchedCount, r.HealthyCount, r.WeakCount, r.DeadCount, r.HatchDate, r.GradedBy)
	if err != nil {
		return 0, fmt.Errorf("create hatch record: %w", err)
	}
	return res.LastInsertId()
}

func (s *HatchRecordStore) GetByID(ctx context.Context, id int64) (*model.HatchRecord, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, batch_id, hatched_count, healthy_count, weak_count, dead_count, hatch_date, graded_by, created_at FROM hatch_records WHERE id = ?`, id)
	var r model.HatchRecord
	err := row.Scan(&r.ID, &r.BatchID, &r.HatchedCount, &r.HealthyCount, &r.WeakCount, &r.DeadCount, &r.HatchDate, &r.GradedBy, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrHatchRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get hatch record %d: %w", id, err)
	}
	return &r, nil
}

func (s *HatchRecordStore) ListByBatch(ctx context.Context, batchID int64) ([]*model.HatchRecord, error) {
	s.mu.RLock()
	if cached, ok := s.cache[batchID]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, batch_id, hatched_count, healthy_count, weak_count, dead_count, hatch_date, graded_by, created_at FROM hatch_records WHERE batch_id = ? ORDER BY hatch_date`, batchID)
	if err != nil {
		return nil, fmt.Errorf("list hatch records by batch: %w", err)
	}
	defer rows.Close()
	var out []*model.HatchRecord
	for rows.Next() {
		var r model.HatchRecord
		if err := rows.Scan(&r.ID, &r.BatchID, &r.HatchedCount, &r.HealthyCount, &r.WeakCount, &r.DeadCount, &r.HatchDate, &r.GradedBy, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan hatch record: %w", err)
		}
		out = append(out, &r)
	}
	s.mu.Lock()
	s.cache[batchID] = out
	s.mu.Unlock()
	return out, rows.Err()
}

func (s *HatchRecordStore) ListByDateRange(ctx context.Context, from, to time.Time) ([]*model.HatchRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, batch_id, hatched_count, healthy_count, weak_count, dead_count, hatch_date, graded_by, created_at FROM hatch_records WHERE hatch_date >= ? AND hatch_date <= ? ORDER BY hatch_date`,
		from, to)
	if err != nil {
		return nil, fmt.Errorf("list hatch records by date range: %w", err)
	}
	defer rows.Close()
	var out []*model.HatchRecord
	for rows.Next() {
		var r model.HatchRecord
		if err := rows.Scan(&r.ID, &r.BatchID, &r.HatchedCount, &r.HealthyCount, &r.WeakCount, &r.DeadCount, &r.HatchDate, &r.GradedBy, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan hatch record: %w", err)
		}
		out = append(out, &r)
	}
	return out, rows.Err()
}

func (s *HatchRecordStore) SummaryByBatch(ctx context.Context, batchID int64) (totalHatched, healthy, weak, dead int, err error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(hatched_count),0), COALESCE(SUM(healthy_count),0), COALESCE(SUM(weak_count),0), COALESCE(SUM(dead_count),0) FROM hatch_records WHERE batch_id = ?`,
		batchID)
	err = row.Scan(&totalHatched, &healthy, &weak, &dead)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("summary by batch: %w", err)
	}
	return totalHatched, healthy, weak, dead, nil
}

func (s *HatchRecordStore) CountByBatch(ctx context.Context, batchID int64) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM hatch_records WHERE batch_id = ?`, batchID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count hatch records: %w", err)
	}
	return count, nil
}
