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
	indexCase  = "cases"
	indexEvent = "events"
)

// Service is the interface for the datastore
type Service interface {
	CreateCase(ctx context.Context, caze *api.Case) error
	UpdateCase(ctx context.Context, caze *api.Case) error
	GetCasesByEmail(ctx context.Context, email string) ([]api.Case, error)
	GetCase(ctx context.Context, id string) (*api.Case, error)
	CreateEvent(ctx context.Context, caseID string, event *api.Event) error
	UpdateEvent(ctx context.Context, caseID string, event *api.Event) error
	DeleteEvent(ctx context.Context, caseID, eventID string) error
	GetEventByID(ctx context.Context, caseID, eventID string) (*api.Event, error)
	GetEvents(ctx context.Context, caseID string) ([]api.Event, error)
	// CreateFile(ctx context.Context, file *api.File) error
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
	return s.save(ctx, indexCase, caze.ID, caze)
}

func (s svc) UpdateCase(ctx context.Context, caze *api.Case) error {
	caze.UpdatedAt = time.Now().Unix()
	return s.save(ctx, indexCase, caze.ID, caze)
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
			return nil, err
		}

		var caze api.Case
		if err := json.Unmarshal(source, &caze); err != nil {
			return nil, err
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
		return nil, err
	}

	var caze api.Case
	if err := json.Unmarshal(resp, &caze); err != nil {
		return nil, err
	}

	return &caze, nil
}

func (s svc) CreateEvent(ctx context.Context, caseID string, event *api.Event) error {
	event.ID = internal.NewID()
	event.CreatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	return s.save(ctx, index, event.ID, event)
}

func (s svc) UpdateEvent(ctx context.Context, caseID string, event *api.Event) error {
	event.UpdatedAt = time.Now().Unix()
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	return s.save(ctx, index, event.ID, event)
}

func (s svc) GetEventByID(ctx context.Context, caseID, eventID string) (*api.Event, error) {
	resp, err := s.searchByID(ctx, indexEvent+"-"+caseID, eventID)
	if err != nil {
		return nil, err
	}

	var event api.Event
	if err := json.Unmarshal(resp, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (s svc) GetEvents(ctx context.Context, caseID string) ([]api.Event, error) {
	search, err := s.search(ctx, indexEvent+"-"+caseID)
	if err != nil {
		return nil, err
	}

	var events []api.Event
	for _, hit := range search.Hits.Hits {
		source, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, err
		}

		var event api.Event
		if err := json.Unmarshal(source, &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func (s svc) DeleteEvent(ctx context.Context, caseID, eventID string) error {
	index := fmt.Sprintf("%s-%s", indexEvent, caseID)
	return s.delete(ctx, index, eventID)
}

func (s svc) CreateFile(ctx context.Context, file *api.File) error {
	file.ID = internal.NewID()
	file.CreatedAt = time.Now().Unix()
	return s.save(ctx, "files", file.ID, file)
}

func (s svc) CreateProcess(ctx context.Context, process *api.Process) error {
	process.ID = internal.NewID()
	process.CreatedAt = time.Now().Unix()
	return s.save(ctx, "processes", process.ID, process)
}

func (s svc) UpdateProcess(ctx context.Context, process *api.Process) error {
	process.UpdatedAt = time.Now().Unix()
	return s.save(ctx, "processes", process.ID, process)
}

func (s svc) GetProcess(ctx context.Context, id string) (*api.Process, error) {
	resp, err := s.searchByID(ctx, "processes", id)
	if err != nil {
		return nil, err
	}

	var process api.Process
	if err := json.Unmarshal(resp, &process); err != nil {
		return nil, err
	}

	return &process, nil
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
				return nil, err
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
