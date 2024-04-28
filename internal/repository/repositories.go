package repository

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repositories interface {
	Measurement() MeasurementRepository
	Setting() SettingRepository
}

type repositories struct {
	measurementRepository MeasurementRepository
	settingRepository     SettingRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	err := db.AutoMigrate(&model.Measurement{})
	if err != nil {
		logrus.Panic(err)
	}
	err = db.AutoMigrate(&model.Setting{})
	if err != nil {
		logrus.Panic(err)
	}
	measurementRepository := newMeasurementRepository(db)
	settingRepository := newSettingRepository(db)
	return &repositories{
		measurementRepository: measurementRepository,
		settingRepository:     settingRepository,
	}
}

func (r repositories) Measurement() MeasurementRepository {
	return r.measurementRepository
}

func (r repositories) Setting() SettingRepository {
	return r.settingRepository
}
