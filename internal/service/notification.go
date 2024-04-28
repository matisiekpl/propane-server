package service

import (
	"context"
	"firebase.google.com/go/messaging"
	"github.com/matisiekpl/propane-server/internal/repository"
	"github.com/sirupsen/logrus"
)

type NotificationService interface {
	SetFirebaseToken(token string) error
	SendNotification(text string) error
}

type notificationService struct {
	settingRepository repository.SettingRepository
	messaging         *messaging.Client
}

func newNotificationService(settingRepository repository.SettingRepository, messaging *messaging.Client) NotificationService {
	return &notificationService{settingRepository: settingRepository, messaging: messaging}
}

const FirebaseTokenKey = "firebaseToken"

func (n notificationService) SetFirebaseToken(token string) error {
	return n.settingRepository.Set(FirebaseTokenKey, token)
}

func (n notificationService) SendNotification(text string) error {
	token, err := n.settingRepository.Get(FirebaseTokenKey)
	if err != nil {
		return err
	}

	logrus.Infof("Sending notification to token: %s", token)

	if token == "" {
		return nil
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Air Monitoring System",
			Body:  text,
		},
		Token: token,
	}

	_, err = n.messaging.Send(context.Background(), message)
	return err
}
