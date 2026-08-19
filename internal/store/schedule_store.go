package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type ScheduleStore struct {
	db *sql.DB
}

func NewScheduleStore(db *sql.DB) *ScheduleStore {
	return &ScheduleStore{db: db}
}

func (s *ScheduleStore) Create(ctx context.Context, sc *model.Schedule) (int64, error) {
	if sc.IncubatorID <= 0 {
		return 0, &model.ValidationError{Field: "incubator_id", Message: "incubator_id is required"}
	}
	if !model.ValidSchedulePhase(sc.Phase) {
		return 0, &model.ValidationError{Field: "phase", Message: "invalid phase"}
	}
	if sc.EndTime.Before(sc.StartTime) {
		return 0, &model.ValidationError{Field: "end_time", Message: "end_time must be after start_time"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO schedules(incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status) VALUES(?,?,?,?,?,?,?,?)`,
		sc.IncubatorID, sc.BatchID, sc.Phase, sc.TargetTemp, sc.TargetHumidity, sc.StartTime, sc.EndTime, sc.Status)
	if err != nil {
		return 0, fmt.Errorf("create schedule: %w", err)
	}
	return res.LastInsertId()
}

func (s *ScheduleStore) GetByID(ctx context.Context, id int64) (*model.Schedule, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status, created_at FROM schedules WHERE id = ?`, id)
	var sc model.Schedule
	err := row.Scan(&sc.ID, &sc.IncubatorID, &sc.BatchID, &sc.Phase, &sc.TargetTemp, &sc.TargetHumidity, &sc.StartTime, &sc.EndTime, &sc.Status, &sc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrScheduleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get schedule %d: %w", id, err)
	}
	return &sc, nil
}

func (s *ScheduleStore) ListByIncubator(ctx context.Context, incubatorID int64) ([]*model.Schedule, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status, created_at FROM schedules WHERE incubator_id = ? ORDER BY start_time`, incubatorID)
	if err != nil {
		return nil, fmt.Errorf("list schedules by incubator: %w", err)
	}
	defer rows.Close()
	var out []*model.Schedule
	for rows.Next() {
		var sc model.Schedule
		if err := rows.Scan(&sc.ID, &sc.IncubatorID, &sc.BatchID, &sc.Phase, &sc.TargetTemp, &sc.TargetHumidity, &sc.StartTime, &sc.EndTime, &sc.Status, &sc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		out = append(out, &sc)
	}
	return out, rows.Err()
}

func (s *ScheduleStore) ListByBatch(ctx context.Context, batchID int64) ([]*model.Schedule, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status, created_at FROM schedules WHERE batch_id = ? ORDER BY start_time`, batchID)
	if err != nil {
		return nil, fmt.Errorf("list schedules by batch: %w", err)
	}
	defer rows.Close()
	var out []*model.Schedule
	for rows.Next() {
		var sc model.Schedule
		if err := rows.Scan(&sc.ID, &sc.IncubatorID, &sc.BatchID, &sc.Phase, &sc.TargetTemp, &sc.TargetHumidity, &sc.StartTime, &sc.EndTime, &sc.Status, &sc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		out = append(out, &sc)
	}
	return out, rows.Err()
}

func (s *ScheduleStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	if !model.ValidScheduleStatus(status) {
		return &model.ValidationError{Field: "status", Message: "invalid schedule status"}
	}
	res, err := s.db.ExecContext(ctx, `UPDATE schedules SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update schedule status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrScheduleNotFound
	}
	return nil
}

func (s *ScheduleStore) FindConflict(ctx context.Context, incubatorID int64, start, end time.Time) (*model.Schedule, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status, created_at FROM schedules WHERE incubator_id = ? AND status = 'active' AND start_time < ? AND end_time > ? LIMIT 1`,
		incubatorID, end, start)
	var sc model.Schedule
	err := row.Scan(&sc.ID, &sc.IncubatorID, &sc.BatchID, &sc.Phase, &sc.TargetTemp, &sc.TargetHumidity, &sc.StartTime, &sc.EndTime, &sc.Status, &sc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrScheduleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find conflict: %w", err)
	}
	return &sc, nil
}

func (s *ScheduleStore) ListActive(ctx context.Context) ([]*model.Schedule, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, batch_id, phase, target_temp, target_humidity, start_time, end_time, status, created_at FROM schedules WHERE status = 'active' ORDER BY start_time`)
	if err != nil {
		return nil, fmt.Errorf("list active schedules: %w", err)
	}
	defer rows.Close()
	var out []*model.Schedule
	for rows.Next() {
		var sc model.Schedule
		if err := rows.Scan(&sc.ID, &sc.IncubatorID, &sc.BatchID, &sc.Phase, &sc.TargetTemp, &sc.TargetHumidity, &sc.StartTime, &sc.EndTime, &sc.Status, &sc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		out = append(out, &sc)
	}
	return out, rows.Err()
}
