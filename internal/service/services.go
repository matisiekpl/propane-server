package service

import (
	"firebase.google.com/go/messaging"
	"github.com/matisiekpl/propane-server/internal/client"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/repository"
)

type Services interface {
	Measurement() MeasurementService
	Notification() NotificationService
	Alert() AlertService
}

type services struct {
	measurementService  MeasurementService
	notificationService NotificationService
	alertService        AlertService
}

func NewServices(repositories repository.Repositories, config dto.Config, clients client.Clients, messaging *messaging.Client) Services {
	notificationService := newNotificationService(repositories.Setting(), messaging)
	alertService := newAlertService(repositories.Setting(), notificationService)
	measurementService := newMeasurementService(repositories.Measurement(), alertService)
	return &services{
		measurementService:  measurementService,
		notificationService: notificationService,
		alertService:        alertService,
	}
}

func (s services) Measurement() MeasurementService {
	return s.measurementService
}

func (s services) Notification() NotificationService {
	return s.notificationService
}

func (s services) Alert() AlertService {
	return s.alertService
}
