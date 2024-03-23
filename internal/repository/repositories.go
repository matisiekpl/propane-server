package repository

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repositories interface {
	Measurement() MeasurementRepository
}

type repositories struct {
	measurementRepository MeasurementRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	err := db.AutoMigrate(&model.Measurement{})
	if err != nil {
		logrus.Panic(err)
	}
	measurementRepository := newMeasurementRepository(db)
	return &repositories{
		measurementRepository: measurementRepository,
	}
}

func (r repositories) Measurement() MeasurementRepository {
	return r.measurementRepository
}
