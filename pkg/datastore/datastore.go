package datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore/internal"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	indexCase    = "cases"
	indexEntity  = "entities"
	indexEvent   = "events"
	indexLink    = "links"
	indexPerson  = "persons"
	indexProcess = "processes"
	indexKeyword = "keywords"
)

// Service is the interface for the datastore
type Service interface {
	// Case-methods
	CreateCase(ctx context.Context, caze *api.Case) error
	UpdateCase(ctx context.Context, caze *api.Case) error
	GetCasesByEmail(ctx context.Context, email string) ([]api.Case, error)
	GetCase(ctx context.Context, id string) (*api.Case, error)
	DeleteCase(ctx context.Context, id string) error

	// Event-methods
	CreateEvent(ctx context.Context, caseID string, event *api.Event) error
	UpdateEvent(ctx context.Context, caseID string, event *api.Event) error
	DeleteEvent(ctx context.Context, caseID, eventID string) error
	GetEventByID(ctx context.Context, caseID, eventID string) (*api.Event, error)
	GetEventsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Event, error)
	GetEvents(ctx context.Context, caseID string) ([]api.Event, error)
	SearchEvents(ctx context.Context, caseID, prefix string) ([]api.Event, error)

	// Entity-methods
	CreateEntity(ctx context.Context, caseID string, entity *api.Entity) error
	UpdateEntity(ctx context.Context, caseID string, entity *api.Entity) error
	DeleteEntity(ctx context.Context, caseID, entityID string) error
	GetEntityByID(ctx context.Context, caseID, entityID string) (*api.Entity, error)
	GetEntitiesByIDs(ctx context.Context, caseID string, ids []string) ([]api.Entity, error)
	GetEntities(ctx context.Context, caseID string) ([]api.Entity, error)
	SearchEntities(ctx context.Context, caseID, prefix string) ([]api.Entity, error)

	// File-methods
	CreateFile(ctx context.Context, caseID string, file *api.File) error
	UpdateFile(ctx context.Context, caseID string, file *api.File) error
	DeleteFile(ctx context.Context, caseID, fileID string) error
	GetFileByID(ctx context.Context, caseID, fileID string) (*api.File, error)
	GetFilesByIDs(ctx context.Context, caseID string, ids []string) ([]api.File, error)
	SearchFiles(ctx context.Context, caseID, prefix string) ([]api.File, error)

	// Link-methods
	CreateLink(ctx context.Context, caseID string, link *api.Link) error
	UpdateLink(ctx context.Context, caseID string, link *api.Link) error
	GetLinkByID(ctx context.Context, caseID, id string) (*api.Link, error)
	GetLinksByIDs(ctx context.Context, caseID string, ids []string) ([]api.Link, error)
	DeleteLink(ctx context.Context, caseID, id string) error

	// Person-methods
	CreatePerson(ctx context.Context, caseID string, person *api.Person) error
	UpdatePerson(ctx context.Context, caseID string, person *api.Person) error
	DeletePerson(ctx context.Context, caseID, personID string) error
	GetPersonByID(ctx context.Context, caseID, personID string) (*api.Person, error)
	GetPersonsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Person, error)
	GetPersons(ctx context.Context, caseID string) ([]api.Person, error)
	SearchPersons(ctx context.Context, caseID, prefix string) ([]api.Person, error)

	// Keyword-methods
	SaveKeyword(ctx context.Context, caseID string, keyword *api.Keyword) error
	DeleteKeyword(ctx context.Context, caseID, keywordID string) error
	GetKeywordByID(ctx context.Context, caseID string, id string) (*api.Keyword, error)
	GetKeywordsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Keyword, error)
	GetKeywords(ctx context.Context, caseID string) ([]string, error)
	SearchKeywords(ctx context.Context, caseID, name string) ([]api.Keyword, error)

	// Process-methods
	ProcessIndex(caseID string) string
	GetProcessedFiles(ctx context.Context, caseID string) (interface{}, error)
	GetProcessedFile(ctx context.Context, caseID, id string) (interface{}, error)
	GetProcessedFilesByIDs(ctx context.Context, caseID string, ids []string) (interface{}, error)
	SearchProcessedFiles(ctx context.Context, caseID, wildcard string) (interface{}, error)
	// CreateProcess(ctx context.Context, process *api.Process) error
	// UpdateProcess(ctx context.Context, process *api.Process) error
	// GetProcess(ctx context.Context, id string) (*api.Process, error)
}

type svc struct {
	es   *elasticsearch.Client
	urls []string
}

// NewService creates a new database-service
func NewService(elasticURLs ...string) (Service, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: elasticURLs})
	if err != nil {
		return nil, err
	}
	if _, err := es.Ping(); err != nil {
		return nil, err
	}
	return svc{es: es, urls: elasticURLs}, nil
}

func (s svc) CreateCase(ctx context.Context, caze *api.Case) error {
	caze.ID = internal.NewID()
	caze.CreatedAt = time.Now().Unix()
	if err := s.save(ctx, indexCase, caze.ID, caze); err != nil {
		return fmt.Errorf("failed to save Case : %v", err)
	}
	if err := s.newIndex(ctx, indexKeyword+"-"+caze.ID); err != nil {
		return fmt.Errorf("failed to create keywords-index for Case : %v", err)
	}
	return nil
}

func (s svc) UpdateCase(ctx context.Context, caze *api.Case) error {
	caze.UpdatedAt = time.Now().Unix()
	if err := s.save(ctx, indexCase, caze.ID, caze); err != nil {
		return fmt.Errorf("failed to save Case : %v", err)
	}
	return nil
}

// Returns all cases with a investigator with the given email.
func (s svc) GetCasesByEmail(ctx context.Context, email string) ([]api.Case, error) {
	search, err := s.search(ctx, indexCase)
	if err != nil {
		return nil, err
	}

	var cases []api.Case
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var caze api.Case
		if err := json.Unmarshal(source, &caze); err != nil {
			return nil, fmt.Errorf("Case json.Unmarshal: %v", err)
		}

		for _, investigator := range caze.Investigators {
			if investigator == email {
				cases = append(cases, caze)
			}
		}
	}

	return cases, nil
}

func (s svc) GetCase(ctx context.Context, id string) (*api.Case, error) {
	resp, err := s.searchByID(ctx, indexCase, id)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Case: %v", err)
	}

	var caze api.Case
	if err := json.Unmarshal(resp, &caze); err != nil {
		return nil, fmt.Errorf("- json.Unmarshal: %v", err)
	}

	return &caze, nil
}

// Deletes the case with the given ID.
// Returns an error if no such case exists.
func (s svc) DeleteCase(ctx context.Context, id string) error {
	caze, err := s.GetCase(ctx, id)
	if err != nil {
		return fmt.Errorf("Cannot find case: %w", err)
	}

	events, err := s.GetEvents(ctx, id)
	if err == nil {
		for _, event := range events {
			s.DeleteEvent(ctx, id, event.ID)
		}
	}

	entities, err := s.GetEntities(ctx, id)
	if err == nil {
		for _, entity := range entities {
			s.DeleteEntity(ctx, id, entity.ID)
		}
	}

	for _, file := range caze.Files {
		s.DeleteEvent(ctx, id, file.ID)
	}

	links, err := s.GetLinks(ctx, id)
	if err == nil {
		for _, link := range links {
			s.DeleteLink(ctx, id, link.ID)
		}
	}

	persons, err := s.GetPersons(ctx, id)
	if err == nil {
		for _, person := range persons {
			s.DeletePerson(ctx, id, person.ID)
		}
	}

	keywords, err := s.GetKeywords(ctx, id)
	if err == nil {
		for _, keyword := range keywords {
			s.DeleteKeyword(ctx, id, keyword)
		}
	}

	err = s.delete(ctx, indexCase, id)
	if err != nil {
		return fmt.Errorf("Error deleting case: %w", err)
	}
	return nil
}

func (s svc) CreateEvent(ctx context.Context, caseID string, event *api.Event) error {
	event.ID = internal.NewID()
	event.CreatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	if err := s.save(ctx, index, event.ID, event); err != nil {
		return fmt.Errorf("failed to save ); : %v", err)
	}
	return nil
}

func (s svc) UpdateEvent(ctx context.Context, caseID string, event *api.Event) error {
	event.UpdatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	if err := s.save(ctx, index, event.ID, event); err != nil {
		return fmt.Errorf("failed to save ); : %v", err)
	}
	return nil
}

func (s svc) GetEventByID(ctx context.Context, caseID, eventID string) (*api.Event, error) {
	resp, err := s.searchByID(ctx, indexEvent+"-"+caseID, eventID)
	if err != nil {
		return nil, fmt.Errorf("cannot find Event in Case: %v", err)
	}

	var event api.Event
	if err := json.Unmarshal(resp, &event); err != nil {
		return nil, fmt.Errorf("Event json.Unmarshal: %v", err)
	}

	return &event, nil
}

func (s svc) GetEventsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Event, error) {
	if len(ids) == 0 {
		return []api.Event{}, nil
	}

	resp, err := s.searchByIDs(ctx, indexEvent+"-"+caseID, ids)
	if err != nil {
		return nil, fmt.Errorf("Cannot find  in Case: %v", err)
	}

	var events []api.Event
	if err := json.Unmarshal(resp, &events); err != nil {
		return nil, fmt.Errorf("Events json.Unmarshal: %v", err)
	}

	return events, nil
}

func (s svc) GetEvents(ctx context.Context, caseID string) ([]api.Event, error) {
	search, err := s.search(ctx, indexEvent+"-"+caseID)
	if err != nil {
		return nil, fmt.Errorf("cannot search in events-document: %v", err)
	}

	var events []api.Event
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var event api.Event
		if err := json.Unmarshal(source, &event); err != nil {
			return nil, fmt.Errorf("Event json.Unmarshal: %v", err)
		}

		events = append(events, event)
	}
	return events, nil
}

func (s svc) SearchEvents(ctx context.Context, caseID, prefix string) ([]api.Event, error) {
	search, err := s.searchWithPrefix(ctx, indexKeyword+"-"+caseID, "description", prefix)
	if err != nil {
		return nil, err
	}

	var events []api.Event
	if err := json.Unmarshal(search, &events); err != nil {
		return nil, fmt.Errorf("Events json.Unmarshal: %v", err)
	}
	return events, nil
}

func (s svc) DeleteEvent(ctx context.Context, caseID, eventID string) error {
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	if err := s.delete(ctx, index, eventID); err != nil {
		return fmt.Errorf("cannot delete Event: %v", err)
	}
	return nil
}

func (s svc) CreateEntity(ctx context.Context, caseID string, entity *api.Entity) error {
	entity.ID = internal.NewID()
	entity.CreatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEntity, caseID)
	if err := s.save(ctx, index, entity.ID, entity); err != nil {
		return fmt.Errorf("failed to save entity : %v", err)
	}
	return nil
}

func (s svc) UpdateEntity(ctx context.Context, caseID string, entity *api.Entity) error {
	entity.UpdatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEntity, caseID)
	if err := s.save(ctx, index, entity.ID, entity); err != nil {
		return fmt.Errorf("failed to save entity : %v", err)
	}
	return nil
}

func (s svc) GetEntityByID(ctx context.Context, caseID, entityID string) (*api.Entity, error) {
	resp, err := s.searchByID(ctx, indexEntity+"-"+caseID, entityID)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Event in Case: %v", err)
	}

	var entity api.Entity
	if err := json.Unmarshal(resp, &entity); err != nil {
		return nil, fmt.Errorf("Entity json.Unmarshal: %v", err)
	}

	return &entity, nil
}

func (s svc) GetEntitiesByIDs(ctx context.Context, caseID string, ids []string) ([]api.Entity, error) {
	if len(ids) == 0 {
		return []api.Entity{}, nil
	}

	resp, err := s.searchByIDs(ctx, indexEntity+"-"+caseID, ids)
	if err != nil {
		return nil, fmt.Errorf("Cannot find  in Case: %v", err)
	}

	var entities []api.Entity
	if err := json.Unmarshal(resp, &entities); err != nil {
		return nil, fmt.Errorf("Entity json.Unmarshal: %v", err)
	}

	return entities, nil
}

func (s svc) GetEntities(ctx context.Context, caseID string) ([]api.Entity, error) {
	search, err := s.search(ctx, indexEntity+"-"+caseID)
	if err != nil {
		return nil, err
	}

	var entities []api.Entity
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var entity api.Entity
		if err := json.Unmarshal(source, &entity); err != nil {
			return nil, fmt.Errorf("Entity json.Unmarshal: %v", err)
		}

		entities = append(entities, entity)
	}
	return entities, nil
}

func (s svc) DeleteEntity(ctx context.Context, caseID, entityID string) error {
	index := fmt.Sprintf("%s-%s", indexEntity, caseID)
	if err := s.delete(ctx, index, entityID); err != nil {
		return fmt.Errorf("cannot delete Entity: %v", err)
	}
	return nil
}

func (s svc) SearchEntities(ctx context.Context, caseID, prefix string) ([]api.Entity, error) {
	search, err := s.searchWithPrefix(ctx, indexEntity+"-"+caseID, "title", prefix)
	if err != nil {
		return nil, err
	}

	var entities []api.Entity
	if err := json.Unmarshal(search, &entities); err != nil {
		return nil, fmt.Errorf("Entities json.Unmarshal: %v", err)
	}
	return entities, nil
}

func (s svc) CreatePerson(ctx context.Context, caseID string, person *api.Person) error {
	person.ID = internal.NewID()
	person.CreatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexPerson, caseID)
	if err := s.save(ctx, index, person.ID, person); err != nil {
		return fmt.Errorf("failed to save Person : %v", err)
	}
	return nil
}

func (s svc) UpdatePerson(ctx context.Context, caseID string, person *api.Person) error {
	person.UpdatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexPerson, caseID)
	if err := s.save(ctx, index, person.ID, person); err != nil {
		return fmt.Errorf("failed to save Person : %v", err)
	}
	return nil
}

func (s svc) GetPersonByID(ctx context.Context, caseID, personID string) (*api.Person, error) {
	resp, err := s.searchByID(ctx, indexPerson+"-"+caseID, personID)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Event in Case: %v", err)
	}

	var person api.Person
	if err := json.Unmarshal(resp, &person); err != nil {
		return nil, fmt.Errorf("Person json.Unmarshal: %v", err)
	}

	return &person, nil
}

func (s svc) GetPersonsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Person, error) {
	if len(ids) == 0 {
		return []api.Person{}, nil
	}

	resp, err := s.searchByIDs(ctx, indexPerson+"-"+caseID, ids)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Persons in Case: %v", err)
	}

	var persons []api.Person
	if err := json.Unmarshal(resp, &persons); err != nil {
		return nil, fmt.Errorf("PersonsByIDs json.Unmarshal: %v", err)
	}

	return persons, nil
}

func (s svc) GetPersons(ctx context.Context, caseID string) ([]api.Person, error) {
	search, err := s.search(ctx, indexPerson+"-"+caseID)
	if err != nil {
		return nil, err
	}

	var persons []api.Person
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var person api.Person
		if err := json.Unmarshal(source, &person); err != nil {
			return nil, fmt.Errorf("Person json.Unmarshal: %v", err)
		}

		persons = append(persons, person)
	}
	return persons, nil
}

func (s svc) DeletePerson(ctx context.Context, caseID, personID string) error {
	index := fmt.Sprintf("%s-%s", indexPerson, caseID)
	if err := s.delete(ctx, index, personID); err != nil {
		return fmt.Errorf("cannot delete Person: %v", err)
	}
	return nil
}

func (s svc) SearchPersons(ctx context.Context, caseID, prefix string) ([]api.Person, error) {
	// search with the prefix for firstname
	search, err := s.searchWithPrefix(ctx, indexPerson+"-"+caseID, "firstName", prefix)
	if err != nil {
		return nil, err
	}

	// search with the prefix for lastname
	search2, err := s.searchWithPrefix(ctx, indexPerson+"-"+caseID, "lastName", prefix)
	if err != nil {
		return nil, err
	}

	// search with the prefix for emailAddress
	search3, err := s.searchWithPrefix(ctx, indexPerson+"-"+caseID, "emailAddress", prefix)
	if err != nil {
		return nil, err
	}

	// append the byte-slices to the first
	if bytes.Equal(search, []byte("null")) {
		search = nil
	}
	if !bytes.Equal(search2, []byte("null")) {
		search = append(search, search2...)
	}
	if !bytes.Equal(search3, []byte("null")) {
		search = append(search, search3...)
	}

	if search == nil {
		return []api.Person{}, nil
	}

	// unmarshal the data to structs
	var persons []api.Person
	if err := json.Unmarshal(search, &persons); err != nil {
		return nil, fmt.Errorf("Persons json.Unmarshal: %v", err)
	}

	// Remove the duplicates from the response
	var found = make(map[string]bool)
	for i := range persons {
		if found[persons[i].ID] {
			persons = append(persons[:i], persons[i+1:]...)
		}
		found[persons[i].ID] = true
	}

	return persons, nil
}

func (s svc) CreateFile(ctx context.Context, caseID string, file *api.File) error {
	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return err
	}

	file.ID = internal.NewID()
	file.CreatedAt = time.Now().Unix()

	caze.Files = append(caze.Files, *file)
	if err := s.UpdateCase(ctx, caze); err != nil {
		return err
	}

	return nil
}

func (s svc) UpdateFile(ctx context.Context, caseID string, file *api.File) error {
	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return err
	}

	for i, f := range caze.Files {
		if f.ID == file.ID {
			caze.Files[i] = *file
			if err := s.UpdateCase(ctx, caze); err != nil {
				return fmt.Errorf("cannot update file in case: %v", err)
			}
			return nil
		}
	}

	return errors.New("file not found")
}

func (s svc) DeleteFile(ctx context.Context, caseID, fileID string) error {
	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return err
	}

	for i, file := range caze.Files {
		if file.ID == fileID {
			caze.Files = append(caze.Files[:i], caze.Files[i+1:]...)
			if err := s.UpdateCase(ctx, caze); err != nil {
				return fmt.Errorf("cannot delete file in case: %v", err)
			}
			return nil
		}
	}

	return errors.New("file not found")
}

func (s svc) GetFileByID(ctx context.Context, caseID, fileID string) (*api.File, error) {
	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("cannot get case for file: %v", err)
	}

	for _, file := range caze.Files {
		if file.ID == fileID {
			return &file, nil
		}
	}

	return nil, errors.New("file not found")
}

func (s svc) GetFilesByIDs(ctx context.Context, caseID string, ids []string) ([]api.File, error) {
	if len(ids) == 0 {
		return []api.File{}, nil
	}

	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("cannot get case for file: %v", err)
	}

	// Create a hash-map for the IDs
	var idMap = make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	var files []api.File
	for _, f := range caze.Files {
		if idMap[f.ID] {
			files = append(files, f)
		}
	}

	return files, nil
}

func (s svc) SearchFiles(ctx context.Context, caseID, prefix string) ([]api.File, error) {
	// search with the prefix for name
	search, err := s.searchWithPrefix(ctx, indexProcess+"-"+caseID, "name", prefix)
	if err != nil {
		return nil, err
	}

	// search with the wildcard for description
	search2, err := s.searchWithWildcard(ctx, indexProcess+"-"+caseID, "description", prefix)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(search, []byte("null")) {
		search = nil
	}
	if !bytes.Equal(search2, []byte("null")) {
		search = append(search, search2...)
	}

	if search == nil {
		return nil, nil
	}

	// unmarshal the data to an interface
	var files []api.File
	if err := json.Unmarshal(search, &files); err != nil {
		return nil, fmt.Errorf("Files json.Unmarshal: %v", err)
	}

	// exclude the duplicate files
	var exists = make(map[string]bool)
	for i := range files {
		if exists[files[i].ID] {
			files = append(files[:i], files[i+1:]...)
		}
		exists[files[i].ID] = true
	}

	return files, nil
}

func (s svc) GetProcessedFile(ctx context.Context, caseID, id string) (interface{}, error) {
	resp, err := s.searchByID(ctx, s.ProcessIndex(caseID), id)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Processed File in Case: %v", err)
	}

	var processed interface{}
	if err := json.Unmarshal(resp, &processed); err != nil {
		return nil, fmt.Errorf("- json.Unmarshal: %v", err)
	}

	return &processed, nil
}

func (s svc) GetProcessedFiles(ctx context.Context, caseID string) (interface{}, error) {
	search, err := s.search(ctx, s.ProcessIndex(caseID))
	if err != nil {
		return nil, err
	}

	source, err := json.Marshal(search.Hits.Hits)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	var processes interface{}
	if err := json.Unmarshal(source, &processes); err != nil {
		return nil, fmt.Errorf("Processes json.Unmarshal: %v", err)
	}

	return &processes, nil
}

func (s svc) GetProcessedFilesByIDs(ctx context.Context, caseID string, ids []string) (interface{}, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	search, err := s.searchByIDs(ctx, s.ProcessIndex(caseID), ids)
	if err != nil {
		return nil, err
	}

	var processes interface{}
	if err := json.Unmarshal(search, &processes); err != nil {
		return nil, fmt.Errorf("Processes json.Unmarshal: %v", err)
	}

	return &processes, nil
}

func (s svc) CreateProcess(ctx context.Context, process *api.Process) error {
	process.ID = internal.NewID()
	process.CreatedAt = time.Now().Unix()
	if err := s.save(ctx, "processes", process.ID, process); err != nil {
		return fmt.Errorf("failed to save Process : %v", err)
	}
	return nil
}

func (s svc) UpdateProcess(ctx context.Context, process *api.Process) error {
	process.UpdatedAt = time.Now().Unix()
	if err := s.save(ctx, "processes", process.ID, process); err != nil {
		return fmt.Errorf("failed to save Process : %v", err)
	}
	return nil
}

func (s svc) SearchProcessedFiles(ctx context.Context, caseID, wildcard string) (interface{}, error) {
	// search with the wildcard for content
	search, err := s.searchWithWildcard(ctx, indexProcess+"-"+caseID, "content", wildcard)
	if err != nil {
		return nil, err
	}

	// unmarshal the data to an interface
	var processes interface{}
	if err := json.Unmarshal(search, &processes); err != nil {
		return nil, fmt.Errorf("Processes json.Unmarshal: %v", err)
	}

	return processes, nil
}

func (s svc) GetProcess(ctx context.Context, id string) (*api.Process, error) {
	resp, err := s.searchByID(ctx, "processes", id)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Event in Case: %v", err)
	}

	var process api.Process
	if err := json.Unmarshal(resp, &process); err != nil {
		return nil, fmt.Errorf("Process json.Unmarshal: %v", err)
	}

	return &process, nil
}

func (s svc) CreateLink(ctx context.Context, caseID string, link *api.Link) error {
	link.ID = internal.NewID()
	link.CreatedAt = time.Now().Unix()
	if err := s.save(ctx, indexLink+"-"+caseID, link.ID, link); err != nil {
		return fmt.Errorf("failed to save Link : %v", err)
	}
	return nil
}

func (s svc) UpdateLink(ctx context.Context, caseID string, link *api.Link) error {
	link.UpdatedAt = time.Now().Unix()
	if err := s.save(ctx, indexLink+"-"+caseID, link.ID, link); err != nil {
		return fmt.Errorf("failed to save Link : %v", err)
	}
	return nil
}

func (s svc) GetLinkByID(ctx context.Context, caseID, id string) (*api.Link, error) {
	resp, err := s.searchByID(ctx, indexLink+"-"+caseID, id)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Link in Case: %v", err)
	}

	var link api.Link
	if err := json.Unmarshal(resp, &link); err != nil {
		return nil, fmt.Errorf("Link json.Unmarshal: %v", err)
	}

	return &link, nil
}

func (s svc) GetLinksByIDs(ctx context.Context, caseID string, ids []string) ([]api.Link, error) {
	if len(ids) == 0 {
		return []api.Link{}, nil
	}

	resp, err := s.searchByIDs(ctx, indexLink+"-"+caseID, ids)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Links in Case: %v", err)
	}

	var links []api.Link
	if err := json.Unmarshal(resp, &links); err != nil {
		return nil, fmt.Errorf("Link json.Unmarshal: %v", err)
	}

	return links, nil
}

// Returns all links in the specified case.
func (s svc) GetLinks(ctx context.Context, caseID string) ([]api.Link, error) {
	// Search in the service for links in the case.
	searchResponse, err := s.search(ctx, indexLink+"-"+caseID)
	if err != nil {
		return nil, err
	}

	// Convert the hits to links.
	var links []api.Link
	for _, hit := range searchResponse.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var link api.Link
		if err := json.Unmarshal(source, &link); err != nil {
			return nil, fmt.Errorf("Keyword json.Unmarshal: %v", err)
		}

		links = append(links, link)
	}

	return links, nil
}

func (s svc) DeleteLink(ctx context.Context, caseID, id string) error {
	index := fmt.Sprintf("%s-%s", indexLink, caseID)
	if err := s.delete(ctx, index, id); err != nil {
		return fmt.Errorf("cannot delete Link for : %v", err)
	}
	return nil
}

func (s svc) SaveKeyword(ctx context.Context, caseID string, keyword *api.Keyword) error {
	if err := s.save(ctx, indexKeyword+"-"+caseID, keyword.Name, keyword); err != nil {
		return fmt.Errorf("failed to save Keyword in case : %v", err)
	}
	return nil
}

func (s svc) DeleteKeyword(ctx context.Context, caseID, id string) error {
	index := fmt.Sprintf("%s-%s", indexKeyword, caseID)
	if err := s.delete(ctx, index, id); err != nil {
		return fmt.Errorf("cannot delete Keyword in case : %v", err)
	}
	return nil
}

func (s svc) GetKeywordByID(ctx context.Context, caseID, id string) (*api.Keyword, error) {
	resp, err := s.searchByID(ctx, indexKeyword+"-"+caseID, id)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Keyword in Case: %v", err)
	}

	var keyword api.Keyword
	if err := json.Unmarshal(resp, &keyword); err != nil {
		return nil, fmt.Errorf("Keyword json.Unmarshal: %v", err)
	}

	return &keyword, nil
}

func (s svc) GetKeywordsByIDs(ctx context.Context, caseID string, ids []string) ([]api.Keyword, error) {
	if len(ids) == 0 {
		return []api.Keyword{}, nil
	}

	resp, err := s.searchByIDs(ctx, indexKeyword+"-"+caseID, ids)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Keywords in Case: %v", err)
	}

	var keywords []api.Keyword
	if err := json.Unmarshal(resp, &keywords); err != nil {
		return nil, fmt.Errorf("Keyword json.Unmarshal: %v", err)
	}

	return keywords, nil
}

func (s svc) GetKeywords(ctx context.Context, caseID string) ([]string, error) {
	search, err := s.search(ctx, indexKeyword+"-"+caseID)
	if err != nil {
		return nil, err
	}

	var keywords []string
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var keyword api.Keyword
		if err := json.Unmarshal(source, &keyword); err != nil {
			return nil, fmt.Errorf("Keyword json.Unmarshal: %v", err)
		}

		keywords = append(keywords, keyword.Name)
	}

	return keywords, nil
}

func (s svc) SearchKeywords(ctx context.Context, caseID, prefix string) ([]api.Keyword, error) {
	search, err := s.searchWithPrefix(ctx, indexKeyword+"-"+caseID, "name", prefix)
	if err != nil {
		return nil, err
	}

	var keywords []api.Keyword
	if err := json.Unmarshal(search, &keywords); err != nil {
		return nil, fmt.Errorf("Keywords json.Unmarshal: %v", err)
	}
	return keywords, nil
}

// ProcessIndex returns the elastic-index for the processes in the specified case
func (svc) ProcessIndex(caseID string) string { return fmt.Sprintf("%s-%s", indexProcess, caseID) }

// newIndex creates a new index
func (s svc) newIndex(ctx context.Context, index string) error {
	// Create a request
	req, err := http.NewRequest(http.MethodPut, s.urls[0]+"/"+index, nil)
	if err != nil {
		return fmt.Errorf("failed to create http-request: %v", err)
	}
	req.WithContext(ctx)

	// Perform the request
	res, err := s.es.Perform(req)
	if err != nil {
		return fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return decodeErrorHTTP(res)
	}

	// Deserialize the response into a map.
	var r internal.Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Cannot parse the response body: %v", err)
	}

	return nil
}

func (s svc) save(ctx context.Context, index, id string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(dataJSON),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, s.es)
	if err != nil {
		return fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return decodeError(res)
	}

	// Deserialize the response into a map.
	var r internal.Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Cannot parse the response body: %v", err)
	}

	return nil
}

func (s svc) delete(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, s.es)
	if err != nil {
		return fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return decodeError(res)
	}

	// Deserialize the response into a map.
	var r internal.Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("Cannot parse the response body: %v", err)
	}

	return nil
}

func (s svc) searchWithPrefix(ctx context.Context, index, field, prefix string) ([]byte, error) {
	query := internal.QueryRequest{
		Query: internal.Query{
			Bool: &internal.Bool{
				Must: []internal.Must{{
					MatchPhrasePrefix: map[string]interface{}{
						field: map[string]string{
							"query": prefix,
						},
					},
				}},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request.
	scrollDuration := time.Minute
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		s.es.Search.WithBody(bytes.NewReader(queryJSON)),
		s.es.Search.WithSort("_doc"),
		s.es.Search.WithSize(10),
		s.es.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return nil, fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	var hits []interface{}
	for _, hit := range search.Hits.Hits {
		hits = append(hits, hit.Source)
	}

	// Perform the scroll requests in sequence
	for len(search.Hits.Hits) > 0 {
		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := s.es.Scroll(s.es.Scroll.WithScrollID(search.ScrollID), s.es.Scroll.WithScroll(scrollDuration))
		if err != nil {
			return nil, fmt.Errorf("search scrolling failed: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, decodeError(res)
		}

		if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
			return nil, fmt.Errorf("Cannot parse the response body: %v", err)
		}

		for _, hit := range search.Hits.Hits {
			hits = append(hits, hit.Source)
		}
	}

	dataJSON, err := json.Marshal(hits)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}
	return dataJSON, nil
}

func (s svc) searchWithWildcard(ctx context.Context, index, field, wildcard string) ([]byte, error) {
	query := internal.QueryRequest{
		Query: internal.Query{
			Bool: &internal.Bool{
				Must: []internal.Must{{
					Wildcard: map[string]interface{}{
						field: map[string]string{
							"value": wildcard,
						},
					},
				}},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request.
	scrollDuration := time.Minute
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		s.es.Search.WithBody(bytes.NewReader(queryJSON)),
		s.es.Search.WithSort("_doc"),
		s.es.Search.WithSize(10),
		s.es.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return nil, fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	var hits []interface{}
	for _, hit := range search.Hits.Hits {
		hits = append(hits, hit.Source)
	}

	// Perform the scroll requests in sequence
	for len(search.Hits.Hits) > 0 {
		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := s.es.Scroll(s.es.Scroll.WithScrollID(search.ScrollID), s.es.Scroll.WithScroll(scrollDuration))
		if err != nil {
			return nil, fmt.Errorf("search scrolling failed: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, decodeError(res)
		}

		if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
			return nil, fmt.Errorf("Cannot parse the response body: %v", err)
		}

		for _, hit := range search.Hits.Hits {
			hits = append(hits, hit.Source)
		}
	}

	dataJSON, err := json.Marshal(hits)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}
	return dataJSON, nil
}

func (s svc) searchByIDs(ctx context.Context, index string, ids []string) ([]byte, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var query internal.QueryRequest
	query.Query.IDs = map[string][]string{"values": ids}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request.
	scrollDuration := time.Minute
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		s.es.Search.WithBody(bytes.NewReader(queryJSON)),
		s.es.Search.WithSort("_doc"),
		s.es.Search.WithSize(10),
		s.es.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return nil, fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	// make a hash-map of the ids
	var idMap = make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	var data []interface{}
	for _, hit := range search.Hits.Hits {
		if idMap[hit.ID] {
			data = append(data, hit.Source)
		}
	}

	// Perform the scroll requests in sequence
	for len(search.Hits.Hits) > 0 {
		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := s.es.Scroll(s.es.Scroll.WithScrollID(search.ScrollID), s.es.Scroll.WithScroll(scrollDuration))
		if err != nil {
			return nil, fmt.Errorf("search scrolling failed: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, decodeError(res)
		}

		if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
			return nil, fmt.Errorf("Cannot parse the response body: %v", err)
		}

		for _, hit := range search.Hits.Hits {
			if idMap[hit.ID] {
				data = append(data, hit.Source)
			}
		}
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	return dataJSON, nil
}

func (s svc) searchByID(ctx context.Context, index, id string) ([]byte, error) {
	var query internal.QueryRequest
	query.Query.Match = map[string]string{"_id": id}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request.
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		s.es.Search.WithBody(bytes.NewReader(queryJSON)),
		//es.Search.WithTrackTotalHits(true),
		//es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	for _, hit := range search.Hits.Hits {
		if hit.ID == id {
			dataJSON, err := json.Marshal(hit.Source)
			if err != nil {
				return nil, fmt.Errorf("json.Marshal: %v", err)
			}
			return dataJSON, nil
		}
	}

	return nil, errors.New("not found")
}

func (s svc) search(ctx context.Context, index string) (*internal.Response, error) {
	scrollDuration := time.Minute

	res, err := s.es.Search(
		s.es.Search.WithIndex(index),
		s.es.Search.WithSort("_doc"),
		s.es.Search.WithSize(10),
		s.es.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	// init a variable to store all hits
	hits := search.Hits.Hits

	// Perform the scroll requests in sequence
	for len(search.Hits.Hits) > 0 {
		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := s.es.Scroll(s.es.Scroll.WithScrollID(search.ScrollID), s.es.Scroll.WithScroll(scrollDuration))
		if err != nil {
			return nil, fmt.Errorf("search scrolling failed: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, decodeError(res)
		}

		if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
			return nil, fmt.Errorf("Cannot parse the response body: %v", err)
		}

		// append the hits from the latest scroll
		hits = append(hits, search.Hits.Hits...)
	}

	// Make sure the hits are unique
	var found = make(map[string]bool)
	search.Hits.Hits = nil
	for i := range hits {
		if !found[hits[i].ID] {
			search.Hits.Hits = append(search.Hits.Hits, hits[i])
		}
		found[hits[i].ID] = true
	}

	return &search, nil
}

func (s svc) searchByName(ctx context.Context, index, name string) (*internal.Response, error) {
	var query internal.QueryRequest
	query.Query.Match = map[string]string{"name": name}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform the search request.
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		s.es.Search.WithBody(bytes.NewReader(queryJSON)),
		//es.Search.WithTrackTotalHits(true),
		//es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Cannot get response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, decodeError(res)
	}

	var search internal.Response
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("Cannot parse the response body: %v", err)
	}

	return &search, nil
}

func decodeError(res *esapi.Response) error {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return fmt.Errorf("Cannot parse the response body: %v", err)
	}
	return fmt.Errorf("[%s] %s: %s",
		res.Status(),
		e["error"].(map[string]interface{})["type"],
		e["error"].(map[string]interface{})["reason"],
	)
}

func decodeErrorHTTP(res *http.Response) error {
	var e interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return fmt.Errorf("Cannot parse the response body: %v", err)
	}
	jsonData, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}
	return fmt.Errorf("[%d] %s",
		res.StatusCode,
		string(jsonData),
	)
}
