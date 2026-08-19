package model

import "time"

type HatchRecord struct {
	ID            int64     `json:"id"`
	BatchID       int64     `json:"batch_id"`
	HatchedCount  int       `json:"hatched_count"`
	HealthyCount  int       `json:"healthy_count"`
	WeakCount     int       `json:"weak_count"`
	DeadCount     int       `json:"dead_count"`
	HatchDate     time.Time `json:"hatch_date"`
	GradedBy      string    `json:"graded_by"`
	CreatedAt     time.Time `json:"created_at"`
}

func (r *HatchRecord) HatchRate() float64 {
	if r.HatchedCount == 0 {
		return 0
	}
	return float64(r.HealthyCount) / float64(r.HatchedCount) * 100
}
