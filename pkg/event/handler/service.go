package handler

import (
	"errors"
	"fmt"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/grammeaway/users_poc/pkg/creating"
	"github.com/grammeaway/users_poc/pkg/deleting"
	"github.com/grammeaway/users_poc/pkg/updating"
)

type Service interface {
	HandleEvent(event etg.Event) error
}

type service struct {
	creatingService creating.Service
	updatingService updating.Service
	deletingService deleting.Service
}

func NewEventHandler(creatingService creating.Service, updatingService updating.Service, deletingService deleting.Service) Service {
	return &service{creatingService, updatingService, deletingService}
}

func (srv *service) HandleEvent(event etg.Event) error {
	go func() error {
		for {
			if event.EventName == "creation_request" {
				err := srv.creatingService.CreateUser(event)
				return err
			}

			if event.EventName == "updating_request" {
				err := srv.updatingService.UpdateUser(event)
				return err
			}

			if event.EventName == "deletion_request" {
				fmt.Println("Things are being deleted")
				err := srv.deletingService.DeleteUser(event)
				return err
			}

			return errors.New("Event not of type handled by this service")
		}
	}()
	return nil
}
