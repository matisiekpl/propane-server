package service

import (
	"github.com/matisiekpl/propane-server/internal/client"
	"github.com/matisiekpl/propane-server/internal/dto"
	"github.com/matisiekpl/propane-server/internal/repository"
)

type Services interface {
	User() UserService
}

type services struct {
	userService UserService
}

func NewServices(repositories repository.Repositories, config dto.Config, clients client.Clients) Services {
	userService := newUserService(repositories.User(), config)
	return &services{
		userService: userService,
	}
}

func (s services) User() UserService {
	return s.userService
}
