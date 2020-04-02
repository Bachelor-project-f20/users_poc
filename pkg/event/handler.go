package handler

import (
	"errors"
	"fmt"
	"log"

	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/Bachelor-project-f20/users_poc/pkg/creating"
	"github.com/Bachelor-project-f20/users_poc/pkg/deleting"
	"github.com/Bachelor-project-f20/users_poc/pkg/updating"
)

type TestObject struct {
	Ok  bool
	Err error
}

type handler struct {
	errorChan       chan error
	testingChan     chan TestObject
	creatingService creating.Service
	updatingService updating.Service
	deletingService deleting.Service
}

func StartEventHandler(
	eventChan <-chan models.Event,
	creatingService creating.Service,
	updatingService updating.Service,
	deletingService deleting.Service) {

	errChan := make(chan error, 5)
	defer close(errChan)

	handler := handler{
		errChan,
		nil,
		creatingService,
		updatingService,
		deletingService,
	}

	handler.handleEvents(eventChan)
}

func TestingStartEventHandler(
	testingChan chan TestObject,
	eventChan <-chan models.Event,
	creatingService creating.Service,
	updatingService updating.Service,
	deletingService deleting.Service) {

	errChan := make(chan error, 5)
	defer close(errChan)

	handler := handler{
		errChan,
		testingChan,
		creatingService,
		updatingService,
		deletingService,
	}

	handler.handleEvents(eventChan)
}

func (h *handler) handleEvents(eventChan <-chan models.Event) {
	for {
		select {
		case event, open := <-eventChan:
			if !open {
				h.testErrors(false, "EventHandler, event channel closed. STOPPING")
				return
			}
			h.handleEvent(event)
		case err, open := <-h.errorChan:
			if !open {
				h.testErrors(false, "EventHandler, error channel closed. STOPPING")
				return
			}
			h.testErrors(false, fmt.Sprintf("ERROR: ", err))
		case <-h.testingChan:
			log.Println("Stopping eventHandler")
			return
		}
	}
}

func (h *handler) handleEvent(event models.Event) {
	go func() {
		eventType := models.UserEvents(int32(models.UserEvents_value[event.EventName]))
		var err error
		switch eventType {
		case models.UserEvents_CREATE_USER:
			err = h.creatingService.CreateUser(event)
		case models.UserEvents_UPDATE_USER:
			err = h.updatingService.UpdateUser(event)
		case models.UserEvents_DELETE_USER:
			err = h.deletingService.DeleteUser(event)
		default:
			err = errors.New("Event not of type handled by this service")
		}
		if err != nil {
			h.errorChan <- err
		}
		h.testErrors(true, "Event handled")
	}()
}

func (h *handler) testErrors(ok bool, msg string) {
	log.Println(msg)
	if h.testingChan != nil && !ok {
		h.testingChan <- TestObject{ok, errors.New(msg)}
	} else if h.testingChan != nil && ok {
		h.testingChan <- TestObject{ok, nil}
	}
}
