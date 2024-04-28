package repository

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"gorm.io/gorm"
)

type SettingRepository interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type settingRepository struct {
	db *gorm.DB
}

func newSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db}
}

func (s *settingRepository) Get(key string) (string, error) {
	var setting model.Setting
	err := s.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *settingRepository) Set(key, value string) error {
	setting := model.Setting{
		Key:   key,
		Value: value,
	}
	return s.db.Create(&setting).Error
}
