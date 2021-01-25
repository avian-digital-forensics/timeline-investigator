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

	// Get the event to delete
	event, err := s.db.GetEventByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Delete the keywords for the event
	if err := s.removeKeywords(ctx, r.CaseID, event, event.Keywords); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}
	if err := s.db.DeleteEvent(ctx, r.CaseID, event.ID); err != nil {
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

// KeywordsAdd adds keywords to an event
func (s *EventService) KeywordsAdd(ctx context.Context, r api.KeywordsAddRequest) (*api.KeywordsAddResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the event to add the keyword to
	event, err := s.db.GetEventByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be added to the Event
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Create a map of the keywords that were found from the db
	// and add the Event ID to each keyword
	var keywordFound = make(map[string]bool)
	for i := range keywords {
		keywordFound[keywords[i].Name] = true
		keywords[i].EventIDs = append(keywords[i].EventIDs, event.ID)
	}

	// Add the keywords from the request to the Event
	// and append the keywords that didn't already exist
	// to the keyword-slice
	for _, keyword := range r.Keywords {
		if !keywordFound[keyword] {
			keywords = append(keywords, api.Keyword{Name: keyword, EventIDs: []string{event.ID}})
		}
		event.Keywords = append(event.Keywords, keyword)
	}

	// Save the keywords with the Event ID
	// TODO / FIXME: Use bulk-indexer instead
	for _, keyword := range keywords {
		if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
			return nil, api.Error(err, api.ErrCannotPerformOperation)
		}
	}

	// Update the Event with the added keywords
	if err := s.db.UpdateEvent(ctx, r.CaseID, event); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsAddResponse{OK: true}, nil
}

// KeywordsRemove removes keywords from an event
func (s *EventService) KeywordsRemove(ctx context.Context, r api.KeywordsRemoveRequest) (*api.KeywordsRemoveResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the event to remove the keywords from
	event, err := s.db.GetEventByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	if err := s.removeKeywords(ctx, r.CaseID, event, r.Keywords); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Update the Event without the removed keyword
	if err := s.db.UpdateEvent(ctx, r.CaseID, event); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsRemoveResponse{}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *EventService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}

func (s *EventService) removeKeywords(ctx context.Context, caseID string, event *api.Event, removeKeywords []string) error {
	// Create a map of the keywords to remove
	var keywordToRemove = make(map[string]bool)
	for _, keyword := range removeKeywords {
		keywordToRemove[keyword] = true
	}

	// Get the keywords that should be removed from the event
	keywords, err := s.db.GetKeywordsByIDs(ctx, caseID, removeKeywords)
	if err != nil {
		return err
	}

	// check if all keywords in the event should be removed
	if len(removeKeywords) == len(event.Keywords) && len(removeKeywords) == len(keywords) {
		event.Keywords = nil
	}

	for i, keyword := range event.Keywords {
		if keywordToRemove[keyword] {
			event.Keywords = append(event.Keywords[:i], event.Keywords[i+1:]...)
		}
	}

	// Remove the EventID from the keywords
	for ki, keyword := range keywords {
		if keywordToRemove[keyword.Name] {
			for ei, id := range keyword.EventIDs {
				if id == event.ID {
					keywords[ki].EventIDs = append(
						keyword.EventIDs[:ei],
						keyword.EventIDs[ei+1:]...,
					)
				}
			}
		}
	}

	// Save the keywords (or delete if empty)
	// TODO/FIXME: use bulk indexer
	for _, keyword := range keywords {
		toDelete := len(keyword.EntityIDs) == 0 && len(keyword.PersonIDs) == 0 && len(keyword.EventIDs) == 0 && len(keyword.FileIDs) == 0
		if toDelete {
			if err := s.db.DeleteKeyword(ctx, caseID, keyword.Name); err != nil {
				return err
			}
		} else if !toDelete {
			if err := s.db.SaveKeyword(ctx, caseID, &keyword); err != nil {
				return err
			}
		}
	}
	return nil
}
