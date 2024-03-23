package controller

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/service"
	"net/http"
	"time"
)

type MeasurementController interface {
	Insert(c echo.Context) error
}

type measurementController struct {
	measurementService service.MeasurementService
	broadcaster        Broadcaster
}

func newMeasurementController(measurementService service.MeasurementService, broadcaster Broadcaster) MeasurementController {
	return &measurementController{measurementService: measurementService, broadcaster: broadcaster}
}

func (m *measurementController) Insert(c echo.Context) error {
	var payload dto.Payload
	err := c.Bind(&payload)
	if err != nil {
		return err
	}

	measurement, err := m.measurementService.Insert(payload.AmmoniaLevel, payload.PropaneLevel, time.Unix(payload.MeasuredAt, 0))
	if err != nil {
		return err
	}

	message, err := json.Marshal(measurement)
	if err != nil {
		return err
	}
	m.broadcaster(message)

	return c.JSON(http.StatusCreated, measurement)
}
