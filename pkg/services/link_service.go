package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

// Create creates a link for an object
// with multiple objects
func (s *LinkService) Create(ctx context.Context, r api.LinkCreateRequest) (*api.LinkCreateResponse, error) {
	if _, err := s.db.GetLinkByID(ctx, r.CaseID, r.FromID); err == nil {
		return nil, api.Error(errors.New("link already exists"), api.ErrCannotPerformOperation)
	}

	// Get the object to create links for
	getObject := func(dest interface{}) error {
		decode := func(source, dest interface{}) error {
			data, err := json.Marshal(source)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, dest)
		}
		if from, _ := s.db.GetEventByID(ctx, r.CaseID, r.FromID); from != nil {
			return decode(from, dest)
		}
		if from, _ := s.db.GetPersonByID(ctx, r.CaseID, r.FromID); from != nil {
			return decode(from, dest)
		}
		if from, _ := s.db.GetEntityByID(ctx, r.CaseID, r.FromID); from != nil {
			return decode(from, dest)
		}
		if from, _ := s.db.GetFileByID(ctx, r.CaseID, r.FromID); from != nil {
			return decode(from, dest)
		}
		return errors.New("the object cannot be found")
	}

	var from interface{}
	if err := getObject(&from); err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}
	log.Println(from)

	events, err := s.db.GetEventsByIDs(ctx, r.CaseID, r.EventIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	persons, err := s.db.GetPersonsByIDs(ctx, r.CaseID, r.PersonIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	entities, err := s.db.GetEntitiesByIDs(ctx, r.CaseID, r.EntityIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	files, err := s.db.GetFilesByIDs(ctx, r.CaseID, r.FileIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Create the link
	link := api.Link{
		From:     from,
		Events:   events,
		Persons:  persons,
		Entities: entities,
		Files:    files,
	}

	if err := s.db.CreateLink(ctx, r.CaseID, &link); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.LinkCreateResponse{Linked: link}, nil
}

// Get gets an  with its links
func (s *LinkService) Get(ctx context.Context, r api.LinkGetRequest) (*api.LinkGetResponse, error) {
	link, err := s.db.GetLinkByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.LinkGetResponse{Link: *link}, nil
}

// Delete deletes all links to the specified object
func (s *LinkService) Delete(ctx context.Context, r api.LinkDeleteRequest) (*api.LinkDeleteResponse, error) {
	if err := s.db.DeleteLink(ctx, r.CaseID, r.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}
	return &api.LinkDeleteResponse{}, nil
}

// Add adds links for the specified object
func (s *LinkService) Add(ctx context.Context, r api.LinkAddRequest) (*api.LinkAddResponse, error) {
	link, err := s.db.GetLinkByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// create maps for the existing linked objects
	var eventMap = make(map[string]bool)
	for _, event := range link.Events {
		eventMap[event.ID] = true
	}
	var personMap = make(map[string]bool)
	for _, person := range link.Persons {
		personMap[person.ID] = true
	}
	var entityMap = make(map[string]bool)
	for _, entity := range link.Entities {
		entityMap[entity.ID] = true
	}
	var fileMap = make(map[string]bool)
	for _, file := range link.Files {
		fileMap[file.ID] = true
	}

	// remove the ids that already exist
	for i, id := range r.EventIDs {
		if eventMap[id] {
			r.EventIDs = append(r.EventIDs[:i], r.EventIDs[i+1:]...)
		}
	}
	for i, id := range r.PersonIDs {
		if personMap[id] {
			r.PersonIDs = append(r.PersonIDs[:i], r.PersonIDs[i+1:]...)
		}
	}
	for i, id := range r.EntityIDs {
		if entityMap[id] {
			r.EntityIDs = append(r.EntityIDs[:i], r.EntityIDs[i+1:]...)
		}
	}
	for i, id := range r.FileIDs {
		if fileMap[id] {
			r.FileIDs = append(r.FileIDs[:i], r.FileIDs[i+1:]...)
		}
	}

	// Get the objects to link
	events, err := s.db.GetEventsByIDs(ctx, r.CaseID, r.EventIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	persons, err := s.db.GetPersonsByIDs(ctx, r.CaseID, r.PersonIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	entities, err := s.db.GetEntitiesByIDs(ctx, r.CaseID, r.EntityIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	files, err := s.db.GetFilesByIDs(ctx, r.CaseID, r.FileIDs)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// append the new links
	link.Events = append(link.Events, events...)
	link.Persons = append(link.Persons, persons...)
	link.Entities = append(link.Entities, entities...)
	link.Files = append(link.Files, files...)

	// Update the link
	if err := s.db.UpdateLink(ctx, r.CaseID, link); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.LinkAddResponse{AddedLinks: *link}, nil
}

// Remove removes links for the specified object
func (s *LinkService) Remove(ctx context.Context, r api.LinkRemoveRequest) (*api.LinkRemoveResponse, error) {
	link, err := s.db.GetLinkByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// create maps for the IDs that should be removed
	var eventMap = make(map[string]bool)
	for _, id := range r.EventIDs {
		eventMap[id] = true
	}
	var personMap = make(map[string]bool)
	for _, id := range r.PersonIDs {
		personMap[id] = true
	}
	var entityMap = make(map[string]bool)
	for _, id := range r.EntityIDs {
		entityMap[id] = true
	}
	var fileMap = make(map[string]bool)
	for _, id := range r.FileIDs {
		fileMap[id] = true
	}

	// remove the ids that already exist
	for i, event := range link.Events {
		if eventMap[event.ID] {
			link.Events = append(link.Events[:i], link.Events[i+1:]...)
		}
	}
	for i, person := range link.Persons {
		if personMap[person.ID] {
			link.Persons = append(link.Persons[:i], link.Persons[i+1:]...)
		}
	}
	for i, entity := range link.Entities {
		if entityMap[entity.ID] {
			link.Entities = append(link.Entities[:i], link.Entities[i+1:]...)
		}
	}
	for i, file := range link.Files {
		if fileMap[file.ID] {
			link.Files = append(link.Files[:i], link.Files[i+1:]...)
		}
	}
	// Update the link
	if err := s.db.UpdateLink(ctx, r.CaseID, link); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.LinkRemoveResponse{RemovedLinks: *link}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *LinkService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
