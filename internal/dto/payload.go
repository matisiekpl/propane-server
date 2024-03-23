package dto

type Payload struct {
	AmmoniaLevel int64 `json:"ammonia_level"`
	PropaneLevel int64 `json:"propane_level"`
	MeasuredAt   int64 `json:"measured_at"`
}
