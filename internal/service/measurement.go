package service

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/matisiekpl/propane-server/internal/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type MeasurementService interface {
	Insert(ammoniaLevel, propaneLevel int64, measuredAt time.Time) (model.Measurement, error)
}

type measurementService struct {
	measurementRepository repository.MeasurementRepository
}

func newMeasurementService(measurementRepository repository.MeasurementRepository) MeasurementService {
	return &measurementService{measurementRepository: measurementRepository}
}

func (m *measurementService) Insert(ammoniaLevel, propaneLevel int64, measuredAt time.Time) (model.Measurement, error) {
	measurement := model.Measurement{
		AmmoniaLevel: ammoniaLevel,
		PropaneLevel: propaneLevel,
		MeasuredAt:   measuredAt,
	}
	err := m.measurementRepository.Insert(&measurement)
	if err != nil {
		return model.Measurement{}, err
	}
	logrus.Infof("Inserted measurement: %v", measurement)
	return measurement, nil
}
