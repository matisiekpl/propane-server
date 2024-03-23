package repository

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"gorm.io/gorm"
	"time"
)

type MeasurementRepository interface {
	Insert(measurement *model.Measurement) error
	Find(start, end time.Time) ([]model.Measurement, error)
}

type measurementRepository struct {
	db *gorm.DB
}

func newMeasurementRepository(db *gorm.DB) MeasurementRepository {
	return &measurementRepository{db}
}

func (m *measurementRepository) Insert(measurement *model.Measurement) error {
	return m.db.Create(measurement).Error
}

func (m *measurementRepository) Find(start, end time.Time) ([]model.Measurement, error) {
	var measurements []model.Measurement
	err := m.db.Model(&model.Measurement{}).Where("measured_at > ? and measured_at < ?", start, end).Find(&measurements).Error
	if err != nil {
		return nil, err
	}
	return measurements, nil
}
