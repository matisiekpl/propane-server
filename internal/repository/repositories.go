package repository

import (
	"github.com/matisiekpl/propane-server/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repositories interface {
	User() UserRepository
}

type repositories struct {
	userRepository UserRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		logrus.Panic(err)
	}
	userRepository := newUserRepository(db)
	return &repositories{
		userRepository: userRepository,
	}
}

func (r repositories) User() UserRepository {
	return r.userRepository
}
