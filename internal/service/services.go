package service

import (
	"github.com/matisiekpl/propane-server/internal/client"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/repository"
)

type Services interface {
	Measurement() MeasurementService
}

type services struct {
	measurementService MeasurementService
}

func NewServices(repositories repository.Repositories, config dto.Config, clients client.Clients) Services {
	measurementService := newMeasurementService(repositories.Measurement())
	return &services{
		measurementService: measurementService,
	}
}

func (s services) Measurement() MeasurementService {
	return s.measurementService
}
