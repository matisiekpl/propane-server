package service

import (
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/matisiekpl/propane-server/internal/repository"
	"strconv"
)

type AlertService interface {
	SetThresholds(thresholds dto.Thresholds) error
	GetThresholds() (dto.Thresholds, error)

	Check(measurement model.Measurement) error
}

type alertService struct {
	settingRepository   repository.SettingRepository
	notificationService NotificationService
}

func newAlertService(settingRepository repository.SettingRepository, notificationService NotificationService) AlertService {
	return &alertService{
		settingRepository:   settingRepository,
		notificationService: notificationService,
	}
}

const (
	AmmoniaThresholdKey = "ammoniaThreshold"
	PropaneThresholdKey = "propaneThreshold"

	AmmoniaThresholdNotifiedKey = "ammoniaThresholdNotified"
	PropaneThresholdNotifiedKey = "propaneThresholdNotified"
)

func (a alertService) SetThresholds(thresholds dto.Thresholds) error {
	err := a.settingRepository.Set(AmmoniaThresholdKey, strconv.Itoa(thresholds.AmmoniaThreshold))
	if err != nil {
		return err
	}

	err = a.settingRepository.Set(PropaneThresholdKey, strconv.Itoa(thresholds.PropaneThreshold))
	if err != nil {
		return err
	}

	return nil
}

func (a alertService) GetThresholds() (dto.Thresholds, error) {
	ammoniaThreshold, err := a.settingRepository.Get(AmmoniaThresholdKey)
	if err != nil {
		return dto.Thresholds{}, err
	}

	propaneThreshold, err := a.settingRepository.Get(PropaneThresholdKey)
	if err != nil {
		return dto.Thresholds{}, err
	}

	parsedAmmoniaThreshold, err := strconv.ParseInt(ammoniaThreshold, 10, 64)
	if err != nil {
		return dto.Thresholds{}, err
	}

	parsedPropaneThreshold, err := strconv.ParseInt(propaneThreshold, 10, 64)
	if err != nil {
		return dto.Thresholds{}, err
	}

	return dto.Thresholds{
		AmmoniaThreshold: int(parsedAmmoniaThreshold),
		PropaneThreshold: int(parsedPropaneThreshold),
	}, nil
}

func (a alertService) CheckThresholds(ammonia int, propane int) (bool, bool, error) {
	thresholds, err := a.GetThresholds()
	if err != nil {
		return false, false, err
	}

	return ammonia > thresholds.AmmoniaThreshold, propane > thresholds.PropaneThreshold, nil
}

func (a alertService) Check(measurement model.Measurement) error {
	amm, pro, err := a.CheckThresholds(int(measurement.AmmoniaLevel), int(measurement.PropaneLevel))
	if err != nil {
		return err
	}

	if amm {
		notified, err := a.settingRepository.Get(AmmoniaThresholdNotifiedKey)
		if err != nil {
			return err
		}
		if notified != "true" {
			err = a.notificationService.SendNotification("Ammonia level is too high!")
			if err != nil {
				return err
			}
			err = a.settingRepository.Set(AmmoniaThresholdNotifiedKey, "true")
			if err != nil {
				return err
			}
		}
	} else {
		err = a.settingRepository.Set(AmmoniaThresholdNotifiedKey, "false")
		if err != nil {
			return err
		}
	}

	if pro {
		notified, err := a.settingRepository.Get(PropaneThresholdNotifiedKey)
		if err != nil {
			return err
		}
		if notified != "true" {
			err = a.notificationService.SendNotification("Propane level is too high!")
			if err != nil {
				return err
			}
			err = a.settingRepository.Set(PropaneThresholdNotifiedKey, "true")
			if err != nil {
				return err
			}
		}
	} else {
		err = a.settingRepository.Set(PropaneThresholdNotifiedKey, "false")
		if err != nil {
			return err
		}
	}

	return nil
}
