package controller

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/service"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type MeasurementController interface {
	Insert(c echo.Context) error
	List(c echo.Context) error
	InsertSamples()
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

func (m *measurementController) List(c echo.Context) error {
	from, err := strconv.Atoi(c.QueryParam("from"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "from parameter is not a valid unix timestamp")
	}
	to, err := strconv.Atoi(c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "to parameter is not a valid unix timestamp")
	}

	start := time.Unix(int64(from), 0)
	end := time.Unix(int64(to), 0)
	measurements, err := m.measurementService.GetByDate(start, end)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, measurements)
}

func (m *measurementController) InsertSamples() {
	for {
		ammoniaLevel := randomNumberBetween(0, 100)
		propaneLevel := randomNumberBetween(0, 100)
		measuredAt := time.Now()
		measurement, err := m.measurementService.Insert(ammoniaLevel, propaneLevel, measuredAt)
		if err != nil {
			logrus.Panic(err)
		}
		message, err := json.Marshal(measurement)
		if err != nil {
			logrus.Panic(err)
		}
		m.broadcaster(message)
		time.Sleep(10 * time.Second)
	}
}

func randomNumberBetween(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}
