package model

import "time"

type EggTray struct {
	ID         int64     `json:"id"`
	BatchID    int64     `json:"batch_id"`
	TrayNumber int       `json:"tray_number"`
	EggCount   int       `json:"egg_count"`
	Weight     float64   `json:"weight"`
	SourceFarm string    `json:"source_farm"`
	CreatedAt  time.Time `json:"created_at"`
}
