package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type Store struct {
	db           *sql.DB
	Incubators   *IncubatorStore
	Batches      *BatchStore
	EggTrays     *EggTrayStore
	Readings     *ReadingStore
	Schedules    *ScheduleStore
	HatchRecords *HatchRecordStore
	Alerts       *AlertStore
	Maintenance  *MaintenanceStore
}

func NewStore(dsn string) (*Store, error) {
	if !strings.Contains(dsn, "_pragma") {
		if strings.Contains(dsn, "?") {
			dsn += "&_pragma=foreign_keys(1)"
		} else {
			dsn += "?_pragma=foreign_keys(1)"
		}
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set journal_mode: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set foreign_keys: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	s.Incubators = NewIncubatorStore(db)
	s.Batches = NewBatchStore(db)
	s.EggTrays = NewEggTrayStore(db)
	s.Readings = NewReadingStore(db)
	s.Schedules = NewScheduleStore(db)
	s.HatchRecords = NewHatchRecordStore(db)
	s.Alerts = NewAlertStore(db)
	s.Maintenance = NewMaintenanceStore(db)
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS incubators (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			capacity INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			location TEXT NOT NULL DEFAULT '',
			installed_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS batches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			incubator_id INTEGER NOT NULL,
			egg_count INTEGER NOT NULL DEFAULT 0,
			species TEXT NOT NULL DEFAULT '',
			start_date DATETIME NOT NULL,
			expected_hatch_date DATETIME NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(incubator_id) REFERENCES incubators(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_batches_incubator ON batches(incubator_id);`,
		`CREATE TABLE IF NOT EXISTS egg_trays (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			batch_id INTEGER NOT NULL,
			tray_number INTEGER NOT NULL,
			egg_count INTEGER NOT NULL DEFAULT 0,
			weight REAL NOT NULL DEFAULT 0,
			source_farm TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(batch_id) REFERENCES batches(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_egg_trays_batch ON egg_trays(batch_id);`,
		`CREATE TABLE IF NOT EXISTS readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			incubator_id INTEGER NOT NULL,
			sensor_type TEXT NOT NULL,
			value REAL NOT NULL,
			recorded_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(incubator_id) REFERENCES incubators(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_readings_incubator ON readings(incubator_id, recorded_at DESC);`,
		`CREATE TABLE IF NOT EXISTS schedules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			incubator_id INTEGER NOT NULL,
			batch_id INTEGER NOT NULL,
			phase TEXT NOT NULL,
			target_temp REAL NOT NULL DEFAULT 37.5,
			target_humidity REAL NOT NULL DEFAULT 60.0,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(incubator_id) REFERENCES incubators(id),
			FOREIGN KEY(batch_id) REFERENCES batches(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_incubator ON schedules(incubator_id, start_time);`,
		`CREATE TABLE IF NOT EXISTS hatch_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			batch_id INTEGER NOT NULL,
			hatched_count INTEGER NOT NULL DEFAULT 0,
			healthy_count INTEGER NOT NULL DEFAULT 0,
			weak_count INTEGER NOT NULL DEFAULT 0,
			dead_count INTEGER NOT NULL DEFAULT 0,
			hatch_date DATETIME NOT NULL,
			graded_by TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(batch_id) REFERENCES batches(id)
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			incubator_id INTEGER NOT NULL,
			batch_id INTEGER NOT NULL DEFAULT 0,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_incubator ON alerts(incubator_id, status);`,
		`CREATE TABLE IF NOT EXISTS maintenance_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			incubator_id INTEGER NOT NULL CHECK(incubator_id > 0),
			description TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			scheduled_for DATETIME NOT NULL,
			completed_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("migrate: %w (stmt=%s)", err, stmt)
		}
	}
	return nil
}

func (s *Store) SeedIncubator(name, location string) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO incubators(name, capacity, status, location, installed_at) VALUES(?,?,?,?,?)`,
		name, 1000, "active", location, "2024-01-01T00:00:00Z",
	)
	if err != nil {
		return 0, fmt.Errorf("seed incubator: %w", err)
	}
	return res.LastInsertId()
}
