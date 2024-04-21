package service

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/matisiekpl/propane-server/internal/repository"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type MeasurementService interface {
	Insert(ammoniaLevel, propaneLevel int64, measuredAt time.Time) (model.Measurement, error)

	InsertSamples()
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

func (m *measurementService) InsertSamples() {
	for {
		ammoniaLevel := randomNumberBetween(0, 100)
		propaneLevel := randomNumberBetween(0, 100)
		measuredAt := time.Now()
		_, err := m.Insert(ammoniaLevel, propaneLevel, measuredAt)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)
	}
}

func randomNumberBetween(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}
