package handler

import (
	"errors"

	e "github.com/Bachelor-project-f20/eventToGo"
	"github.com/grammeaway/users_poc/users/pkg/creating"
	"github.com/grammeaway/users_poc/users/pkg/deleting"
	"github.com/grammeaway/users_poc/users/pkg/updating"
)

type Service interface {
	HandleEvent(event e.Event) error
}

type service struct {
	creatingService creating.Service
	updatingService updating.Service
	deletingService deleting.Service
}

func NewEventHandler(creatingService creating.Service, updatingService updating.Service, deletingService deleting.Service) Service {
	return &service{creatingService, updatingService, deletingService}
}

func (srv *service) HandleEvent(event e.Event) error {
	if event.GetEventName() == "creation_request" {
		err := srv.creatingService.CreateUser(event)
		return err
	}

	if event.GetEventName() == "updating_request" {
		err := srv.updatingService.UpdateUser(event)
		return err
	}

	if event.GetEventName() == "deletion_request" {
		err := srv.deletingService.DeleteUser(event)
		return err
	}

	return errors.New("Event not of type handled by this service")

}
