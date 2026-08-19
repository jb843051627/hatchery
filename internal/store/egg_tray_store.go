package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/hatchery/internal/model"
)

type EggTrayStore struct {
	db *sql.DB
}

func NewEggTrayStore(db *sql.DB) *EggTrayStore {
	return &EggTrayStore{db: db}
}

func (s *EggTrayStore) Create(ctx context.Context, batchID int64, trayNumber int, eggCount int, weight float64, sourceFarm string) (int64, error) {
	if batchID <= 0 {
		return 0, &model.ValidationError{Field: "batch_id", Message: "batch_id is required"}
	}
	if eggCount <= 0 {
		return 0, &model.ValidationError{Field: "egg_count", Message: "egg_count must be positive"}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO egg_trays(batch_id, tray_number, egg_count, weight, source_farm) VALUES(?,?,?,?,?)`,
		batchID, trayNumber, eggCount, weight, sourceFarm)
	if err != nil {
		return 0, fmt.Errorf("create egg tray: %w", err)
	}
	return res.LastInsertId()
}

func (s *EggTrayStore) GetByID(ctx context.Context, id int64) (*model.EggTray, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, batch_id, tray_number, egg_count, weight, source_farm, created_at FROM egg_trays WHERE id = ?`, id)
	var t model.EggTray
	err := row.Scan(&t.ID, &t.BatchID, &t.TrayNumber, &t.EggCount, &t.Weight, &t.SourceFarm, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrEggTrayNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get egg tray %d: %w", id, err)
	}
	return &t, nil
}

func (s *EggTrayStore) ListByBatch(ctx context.Context, batchID int64) ([]*model.EggTray, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, batch_id, tray_number, egg_count, weight, source_farm, created_at FROM egg_trays WHERE batch_id = ? ORDER BY tray_number`, batchID)
	if err != nil {
		return nil, fmt.Errorf("list egg trays by batch: %w", err)
	}
	defer rows.Close()
	var out []*model.EggTray
	for rows.Next() {
		var t model.EggTray
		if err := rows.Scan(&t.ID, &t.BatchID, &t.TrayNumber, &t.EggCount, &t.Weight, &t.SourceFarm, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan egg tray: %w", err)
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (s *EggTrayStore) UpdateEggCount(ctx context.Context, id int64, count int) error {
	if count < 0 {
		return &model.ValidationError{Field: "egg_count", Message: "egg_count must be non-negative"}
	}
	res, err := s.db.ExecContext(ctx, `UPDATE egg_trays SET egg_count = ? WHERE id = ?`, count, id)
	if err != nil {
		return fmt.Errorf("update egg count: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrEggTrayNotFound
	}
	return nil
}

func (s *EggTrayStore) DeleteByBatch(ctx context.Context, batchID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM egg_trays WHERE batch_id = ?`, batchID)
	if err != nil {
		return fmt.Errorf("delete egg trays by batch: %w", err)
	}
	return nil
}

func (s *EggTrayStore) CountByBatch(ctx context.Context, batchID int64) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM egg_trays WHERE batch_id = ?`, batchID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count egg trays: %w", err)
	}
	return count, nil
}

func (s *EggTrayStore) SumEggsByBatch(ctx context.Context, batchID int64) (int, error) {
	var total sql.NullInt64
	err := s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(egg_count), 0) FROM egg_trays WHERE batch_id = ?`, batchID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum eggs by batch: %w", err)
	}
	return int(total.Int64), nil
}
