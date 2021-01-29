package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

// SearchService holds the dependencies
// for the search-service
type SearchService struct {
	db          datastore.Service
	caseService *CaseService
}

// NewSearchService creates a new search-service
func NewSearchService(db datastore.Service, caseService *CaseService) *SearchService {
	return &SearchService{
		db:          db,
		caseService: caseService,
	}
}

// SearchWithTimespan returns events from the selected timespan
func (s *SearchService) SearchWithTimespan(ctx context.Context, r api.SearchTimespanRequest) (*api.SearchTimespanResponse, error) {
	return nil, errors.New("not implemented yet")
}

// SearchWithText returns data in the case that is related to the text
func (s *SearchService) SearchWithText(ctx context.Context, r api.SearchTextRequest) (*api.SearchTextResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if len(r.Text) < 3 {
		return nil, api.Error(
			errors.New("specify at least 3 characters for text-search"),
			api.ErrCannotPerformOperation,
		)
	}

	// Get keywords based on the search
	keywords, err := s.db.SearchKeywords(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// store all object-ids in slices
	// that were gotten from the keywords
	var eventIDs []string
	var entityIDs []string
	var fileIDs []string
	var personIDs []string
	for _, keyword := range keywords {
		eventIDs = append(eventIDs, keyword.EventIDs...)
		entityIDs = append(entityIDs, keyword.EntityIDs...)
		fileIDs = append(fileIDs, keyword.FileIDs...)
		personIDs = append(personIDs, keyword.PersonIDs...)
	}

	// Get events based on the search
	events, err := s.db.SearchEvents(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, err
	}

	// Get entities based on the search
	entities, err := s.db.SearchEntities(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, err
	}

	// Get persons based on the search
	persons, err := s.db.SearchPersons(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, err
	}

	// Get files based on the search
	files, err := s.db.SearchFiles(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, err
	}

	// Get processed-files based on the search
	processed, err := s.db.SearchProcessedFiles(ctx, r.CaseID, r.Text)
	if err != nil {
		return nil, err
	}

	// Get all unique eventIDs
	var eventFound = make(map[string]bool)
	for i := range events {
		eventFound[events[i].ID] = true
	}
	for i := range eventIDs {
		if eventFound[eventIDs[i]] {
			eventIDs = append(eventIDs[:i], eventIDs[i+1:]...)
		}
	}
	// Get the unqiue events from the keywords
	keywordEvents, err := s.db.GetEventsByIDs(ctx, r.CaseID, eventIDs)
	if err != nil {
		return nil, err
	}
	// append the events to the first slice
	events = append(events, keywordEvents...)

	// Get all unique entityIDs
	var entityFound = make(map[string]bool)
	for i := range entities {
		entityFound[entities[i].ID] = true
	}
	for i := range entityIDs {
		if entityFound[entityIDs[i]] {
			entityIDs = append(entityIDs[:i], entityIDs[i+1:]...)
		}
	}
	// Get the unqiue entities from the keywords
	keywordEntities, err := s.db.GetEntitiesByIDs(ctx, r.CaseID, entityIDs)
	if err != nil {
		return nil, err
	}
	// append the entities to the first slice
	entities = append(entities, keywordEntities...)

	// Get all unique personIDs
	var personFound = make(map[string]bool)
	for i := range persons {
		personFound[persons[i].ID] = true
	}
	for i := range personIDs {
		if personFound[personIDs[i]] {
			personIDs = append(personIDs[:i], personIDs[i+1:]...)
		}
	}
	// Get the unqiue persons from the keywords
	keywordPersons, err := s.db.GetPersonsByIDs(ctx, r.CaseID, personIDs)
	if err != nil {
		return nil, err
	}
	// append the persons to the first slice
	persons = append(persons, keywordPersons...)

	// Get all unique fileIDs
	var fileFound = make(map[string]bool)
	for i := range files {
		fileFound[files[i].ID] = true
	}
	for i := range fileIDs {
		if fileFound[fileIDs[i]] {
			fileIDs = append(fileIDs[:i], fileIDs[i+1:]...)
		}
		//fileFound[fileIDs[i]] = true
	}
	// Get the unqiue files from the keywords
	keywordFiles, err := s.db.GetFilesByIDs(ctx, r.CaseID, fileIDs)
	if err != nil {
		return nil, err
	}
	// append the files to the first slice
	files = append(files, keywordFiles...)

	return &api.SearchTextResponse{
		Events:    events,
		Entities:  entities,
		Files:     files,
		Persons:   persons,
		Processed: processed,
	}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *SearchService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
