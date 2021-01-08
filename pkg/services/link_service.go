package services

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
)

// LinkService holds the dependencies
// for handling Links between objects
type LinkService struct {
	db          datastore.Service
	caseService *CaseService
}

// NewLinkService creates a new LinkService
func NewLinkService(db datastore.Service, caseService *CaseService) *LinkService {
	return &LinkService{db: db, caseService: caseService}
}

// CreateEvent creates a link for an event
// with multiple objects
func (s *LinkService) CreateEvent(ctx context.Context, r api.LinkEventCreateRequest) (*api.LinkEventCreateResponse, error) {
	from, err := s.db.GetEventByID(ctx, r.CaseID, r.FromID)
	if err != nil {
		return nil, err
	}

	var events []api.Event
	for _, id := range r.EventIDs {
		event, err := s.db.GetEventByID(ctx, r.CaseID, id)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	link := api.LinkEvent{From: *from, Events: events}
	if err := s.db.CreateLinkEvent(ctx, r.CaseID, &link); err != nil {
		return nil, err
	}

	if !r.Bidirectional {
		return &api.LinkEventCreateResponse{Linked: link}, nil
	}

	for _, event := range events {
		if err := s.db.CreateLinkEvent(ctx, r.CaseID, &api.LinkEvent{From: event, Events: []api.Event{*from}}); err != nil {
			return nil, err
		}
	}

	return &api.LinkEventCreateResponse{Linked: link}, nil
}

// GetEvent gets an event with its links
func (s *LinkService) GetEvent(ctx context.Context, r api.LinkEventGetRequest) (*api.LinkEventGetResponse, error) {
	link, err := s.db.GetLinkEvent(ctx, r.CaseID, r.EventID)
	if err != nil {
		return nil, err
	}

	return &api.LinkEventGetResponse{Link: *link}, nil
}

// DeleteEvent deletes all links to the specified event
func (s *LinkService) DeleteEvent(ctx context.Context, r api.LinkEventDeleteRequest) (*api.LinkEventDeleteResponse, error) {
	if err := s.db.DeleteLinkEvent(ctx, r.CaseID, r.EventID); err != nil {
		return nil, err
	}
	return &api.LinkEventDeleteResponse{}, nil
}

// UpdateEvent updates links for the specified event
func (s *LinkService) UpdateEvent(ctx context.Context, r api.LinkEventUpdateRequest) (*api.LinkEventUpdateResponse, error) {
	link, err := s.db.GetLinkEvent(ctx, r.CaseID, r.EventID)
	if err != nil {
		return nil, err
	}

	var removeEvent = make(map[string]bool)
	for _, id := range r.EventRemoveIDs {
		removeEvent[id] = true
	}

	var events []api.Event
	for _, id := range r.EventAddIDs {
		event, err := s.db.GetEventByID(ctx, r.CaseID, id)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	for _, event := range link.Events {
		if !removeEvent[event.ID] {
			events = append(events, event)
		}
	}

	link.Events = events
	if err := s.db.UpdateLinkEvent(ctx, r.CaseID, link); err != nil {
		return nil, err
	}

	return &api.LinkEventUpdateResponse{Updated: *link}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *LinkService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
