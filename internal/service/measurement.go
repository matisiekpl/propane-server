package service

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/matisiekpl/propane-server/internal/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type MeasurementService interface {
	Insert(ammoniaLevel, propaneLevel int64, measuredAt time.Time) (model.Measurement, error)
	GetByDate(start, end time.Time) ([]model.Measurement, error)
}

type measurementService struct {
	measurementRepository repository.MeasurementRepository
	alertService          AlertService
}

func newMeasurementService(measurementRepository repository.MeasurementRepository, alertService AlertService) MeasurementService {
	return &measurementService{measurementRepository: measurementRepository, alertService: alertService}
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

	err = m.alertService.Check(measurement)
	if err != nil {
		logrus.Errorf("Error checking measurement: %v", err)
	}

	return measurement, nil
}

func (m *measurementService) GetByDate(start, end time.Time) ([]model.Measurement, error) {
	measurements, err := m.measurementRepository.Find(start, end)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}
