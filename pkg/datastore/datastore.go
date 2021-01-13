package datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore/internal"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	indexCase   = "cases"
	indexEntity = "entities"
	indexEvent  = "events"
	indexLink   = "links"
	indexPerson = "persons"
)

// Service is the interface for the datastore
type Service interface {
	// Case-methods
	CreateCase(ctx context.Context, caze *api.Case) error
	UpdateCase(ctx context.Context, caze *api.Case) error
	GetCasesByEmail(ctx context.Context, email string) ([]api.Case, error)
	GetCase(ctx context.Context, id string) (*api.Case, error)

	// Event-methods
	CreateEvent(ctx context.Context, caseID string, event *api.Event) error
	UpdateEvent(ctx context.Context, caseID string, event *api.Event) error
	DeleteEvent(ctx context.Context, caseID, eventID string) error
	GetEventByID(ctx context.Context, caseID, eventID string) (*api.Event, error)
	GetEvents(ctx context.Context, caseID string) ([]api.Event, error)

	// Entity-methods
	CreateEntity(ctx context.Context, caseID string, entity *api.Entity) error
	UpdateEntity(ctx context.Context, caseID string, entity *api.Entity) error
	DeleteEntity(ctx context.Context, caseID, entityID string) error
	GetEntityByID(ctx context.Context, caseID, entityID string) (*api.Entity, error)
	GetEntities(ctx context.Context, caseID string) ([]api.Entity, error)

	// File-methods
	CreateFile(ctx context.Context, caseID string, file *api.File) error
	UpdateFile(ctx context.Context, caseID, fileID, description string) (*api.File, error)
	DeleteFile(ctx context.Context, caseID, fileID string) error
	GetFile(ctx context.Context, caseID, fileID string) (*api.File, error)

	// Link-methods
	CreateLinkEvent(ctx context.Context, caseID string, link *api.LinkEvent) error
	UpdateLinkEvent(ctx context.Context, caseID string, link *api.LinkEvent) error
	GetLinkEvent(ctx context.Context, caseID, eventID string) (*api.LinkEvent, error)
	DeleteLinkEvent(ctx context.Context, caseID, eventID string) error

	// Person-methods
	CreatePerson(ctx context.Context, caseID string, Person *api.Person) error
	UpdatePerson(ctx context.Context, caseID string, Person *api.Person) error
	DeletePerson(ctx context.Context, caseID, personID string) error
	GetPersonByID(ctx context.Context, caseID, personID string) (*api.Person, error)
	GetPersons(ctx context.Context, caseID string) ([]api.Person, error)

	// CreateProcess(ctx context.Context, process *api.Process) error
	// UpdateProcess(ctx context.Context, process *api.Process) error
	// GetProcess(ctx context.Context, id string) (*api.Process, error)
}

type svc struct {
	es *elasticsearch.Client
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
	return svc{es: es}, nil
}

func (s svc) CreateCase(ctx context.Context, caze *api.Case) error {
	caze.ID = internal.NewID()
	caze.CreatedAt = time.Now().Unix()
	if err := s.save(ctx, indexCase, caze.ID, caze); err != nil {
		return fmt.Errorf("failed to save Case : %v", err)
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
		return nil, fmt.Errorf("Cannot find Event in Case: %v", err)
	}

	var caze api.Case
	if err := json.Unmarshal(resp, &caze); err != nil {
		return nil, fmt.Errorf("- json.Unmarshal: %v", err)
	}

	return &caze, nil
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
			return nil, fmt.Errorf("Case json.Unmarshal: %v", err)
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
			return nil, fmt.Errorf("Case json.Unmarshal: %v", err)
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

func (s svc) UpdateFile(ctx context.Context, caseID, fileID, description string) (*api.File, error) {
	caze, err := s.GetCase(ctx, caseID)
	if err != nil {
		return nil, err
	}

	for i, file := range caze.Files {
		if file.ID == fileID {
			file.UpdatedAt = time.Now().Unix()
			file.Description = description
			caze.Files[i] = file
			if err := s.UpdateCase(ctx, caze); err != nil {
				return nil, fmt.Errorf("cannot update file in case: %v", err)
			}
			return &file, nil
		}
	}

	return nil, errors.New("file not found")
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

func (s svc) GetFile(ctx context.Context, caseID, fileID string) (*api.File, error) {
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

func (s svc) CreateLinkEvent(ctx context.Context, caseID string, link *api.LinkEvent) error {
	link.ID = internal.NewID()
	link.CreatedAt = time.Now().Unix()
	if err := s.save(ctx, indexLink+"-"+caseID, link.From.ID, link); err != nil {
		return fmt.Errorf("failed to save Link : %v", err)
	}
	return nil
}

func (s svc) UpdateLinkEvent(ctx context.Context, caseID string, link *api.LinkEvent) error {
	link.UpdatedAt = time.Now().Unix()
	if err := s.save(ctx, indexLink+"-"+caseID, link.From.ID, link); err != nil {
		return fmt.Errorf("failed to save Link : %v", err)
	}
	return nil
}

func (s svc) GetLinkEvent(ctx context.Context, caseID, eventID string) (*api.LinkEvent, error) {
	resp, err := s.searchByID(ctx, indexLink+"-"+caseID, eventID)
	if err != nil {
		return nil, fmt.Errorf("Cannot find Event in Case: %v", err)
	}

	var link api.LinkEvent
	if err := json.Unmarshal(resp, &link); err != nil {
		return nil, fmt.Errorf("Link json.Unmarshal: %v", err)
	}

	return &link, nil
}

func (s svc) DeleteLinkEvent(ctx context.Context, caseID, eventID string) error {
	index := fmt.Sprintf("%s-%s", indexLink, caseID)
	if err := s.delete(ctx, index, eventID); err != nil {
		return fmt.Errorf("cannot delete Link for Event: %v", err)
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
	// Perform the search request.
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(index),
		//s.es.Search.WithBody(bytes.NewReader(queryJSON)),
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
