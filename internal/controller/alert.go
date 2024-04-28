package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/service"
	"net/http"
)

type AlertController interface {
	SetFirebaseToken(c echo.Context) error
	SetThresholds(c echo.Context) error
	GetThresholds(c echo.Context) error
}

type alertController struct {
	alertService        service.AlertService
	notificationService service.NotificationService
}

func newAlertController(alertService service.AlertService, notificationService service.NotificationService) AlertController {
	return &alertController{
		alertService:        alertService,
		notificationService: notificationService,
	}
}

func (a alertController) SetFirebaseToken(c echo.Context) error {
	var payload dto.SetFirebaseTokenRequest

	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := a.notificationService.SetFirebaseToken(payload.Token); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, payload)
}

func (a alertController) SetThresholds(c echo.Context) error {
	var payload dto.Thresholds

	if err := c.Bind(&payload); err != nil {
		return err
	}

	if err := a.alertService.SetThresholds(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, payload)
}

func (a alertController) GetThresholds(c echo.Context) error {
	thresholds, err := a.alertService.GetThresholds()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, thresholds)
}
