package services

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

const (
	minImportance = 1
	maxImportance = 5
)

// EventService holds dependencies
// for the event-api
type EventService struct {
	db          datastore.Service
	caseService *CaseService
}

// NewEventService creates a new event-service
func NewEventService(db datastore.Service, caseService *CaseService) *EventService {
	return &EventService{db: db, caseService: caseService}
}

// Create creates a new event
func (s *EventService) Create(ctx context.Context, r api.EventCreateRequest) (*api.EventCreateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Check that the timespan is valid
	if r.FromDate > r.ToDate {
		return nil, api.ErrInvalidDates
	}

	// Check that the importance is valid
	if r.Importance < minImportance || r.Importance > maxImportance {
		return nil, api.ErrInvalidImportance
	}

	event := api.Event{
		Importance:  r.Importance,
		Description: r.Description,
		FromDate:    r.FromDate,
		ToDate:      r.ToDate,
	}

	if err := s.db.CreateEvent(ctx, r.CaseID, &event); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EventCreateResponse{Created: event}, nil
}

// Update updates an existing event
func (s *EventService) Update(ctx context.Context, r api.EventUpdateRequest) (*api.EventUpdateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Check that the timespan is valid
	if r.FromDate > r.ToDate {
		return nil, api.ErrInvalidDates
	}

	// Check that the importance is valid
	if r.Importance < minImportance || r.Importance > maxImportance {
		return nil, api.ErrInvalidImportance
	}

	event := api.Event{
		Importance:  r.Importance,
		Description: r.Description,
		FromDate:    r.FromDate,
		ToDate:      r.ToDate,
	}
	event.ID = r.ID

	if err := s.db.UpdateEvent(ctx, r.CaseID, &event); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EventUpdateResponse{Updated: event}, nil
}

// Delete deletes an existing event
func (s *EventService) Delete(ctx context.Context, r api.EventDeleteRequest) (*api.EventDeleteResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if err := s.db.DeleteEvent(ctx, r.CaseID, r.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EventDeleteResponse{}, nil
}

// Get the specified event
func (s *EventService) Get(ctx context.Context, r api.EventGetRequest) (*api.EventGetResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	event, err := s.db.GetEventByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.EventGetResponse{Event: *event}, nil
}

// List all events
func (s *EventService) List(ctx context.Context, r api.EventListRequest) (*api.EventListResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	events, err := s.db.GetEvents(ctx, r.CaseID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EventListResponse{Events: events}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *EventService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
