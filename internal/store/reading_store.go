package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/hatchery/internal/model"
)

type ReadingStore struct {
	db *sql.DB
}

func NewReadingStore(db *sql.DB) *ReadingStore {
	return &ReadingStore{db: db}
}

func (s *ReadingStore) Create(ctx context.Context, incubatorID int64, sensorType string, value float64, recordedAt time.Time) (int64, error) {
	if incubatorID <= 0 {
		return 0, &model.ValidationError{Field: "incubator_id", Message: "incubator_id is required"}
	}
	if !model.ValidSensorType(sensorType) {
		return 0, &model.ValidationError{Field: "sensor_type", Message: "invalid sensor type"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO readings(incubator_id, sensor_type, value, recorded_at) VALUES(?,?,?,?)`,
		incubatorID, sensorType, value, recordedAt)
	if err != nil {
		return 0, fmt.Errorf("create reading: %w", err)
	}
	return res.LastInsertId()
}

func (s *ReadingStore) GetByID(ctx context.Context, id int64) (*model.SensorReading, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, sensor_type, value, recorded_at, created_at FROM readings WHERE id = ?`, id)
	var r model.SensorReading
	err := row.Scan(&r.ID, &r.IncubatorID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrReadingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get reading %d: %w", id, err)
	}
	return &r, nil
}

func (s *ReadingStore) ListByIncubator(ctx context.Context, incubatorID int64, limit int) ([]*model.SensorReading, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, sensor_type, value, recorded_at, created_at FROM readings WHERE incubator_id = ? ORDER BY recorded_at DESC LIMIT ?`,
		incubatorID, limit)
	if err != nil {
		return nil, fmt.Errorf("list readings by incubator: %w", err)
	}
	defer rows.Close()
	var out []*model.SensorReading
	for rows.Next() {
		var r model.SensorReading
		if err := rows.Scan(&r.ID, &r.IncubatorID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan reading: %w", err)
		}
		out = append(out, &r)
	}
	return out, rows.Err()
}

func (s *ReadingStore) ListByTimeRange(ctx context.Context, incubatorID int64, from, to time.Time) ([]*model.SensorReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, incubator_id, sensor_type, value, recorded_at, created_at FROM readings WHERE incubator_id = ? AND recorded_at >= ? AND recorded_at <= ? ORDER BY recorded_at`,
		incubatorID, from, to)
	if err != nil {
		return nil, fmt.Errorf("list readings by time range: %w", err)
	}
	defer rows.Close()
	var out []*model.SensorReading
	for rows.Next() {
		var r model.SensorReading
		if err := rows.Scan(&r.ID, &r.IncubatorID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan reading: %w", err)
		}
		out = append(out, &r)
	}
	return out, rows.Err()
}

func (s *ReadingStore) BatchCreate(ctx context.Context, readings []*model.SensorReading) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO readings(incubator_id, sensor_type, value, recorded_at) VALUES(?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()
	for _, r := range readings {
		if _, err := stmt.ExecContext(ctx, r.IncubatorID, r.SensorType, r.Value, r.RecordedAt); err != nil {
			return fmt.Errorf("batch insert reading: %w", err)
		}
	}
	return tx.Commit()
}

func (s *ReadingStore) CountByIncubator(ctx context.Context, incubatorID int64) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE incubator_id = ?`, incubatorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count readings: %w", err)
	}
	return count, nil
}

func (s *ReadingStore) LatestByIncubator(ctx context.Context, incubatorID int64, sensorType string) (*model.SensorReading, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, incubator_id, sensor_type, value, recorded_at, created_at FROM readings WHERE incubator_id = ? AND sensor_type = ? ORDER BY recorded_at DESC LIMIT 1`,
		incubatorID, sensorType)
	var r model.SensorReading
	err := row.Scan(&r.ID, &r.IncubatorID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrReadingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("latest reading: %w", err)
	}
	return &r, nil
}
