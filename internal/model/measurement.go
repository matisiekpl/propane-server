package model

import (
	"gorm.io/gorm"
	"time"
)

type Measurement struct {
	gorm.Model
	AmmoniaLevel int64
	PropaneLevel int64
	MeasuredAt   time.Time
}
