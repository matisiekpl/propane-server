package dto

import "time"

type Payload struct {
	AmmoniaLevel int64     `json:"ammonia_level"`
	PropaneLevel int64     `json:"propane_level"`
	MeasuredAt   time.Time `json:"measured_at"`
}
