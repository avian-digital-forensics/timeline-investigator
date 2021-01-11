package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
)

// EntityService holds the dependencies
// for handling the Entity API
type EntityService struct {
	db          datastore.Service
	caseService *CaseService
	// types hold all the entity-types
	// available for usage
	types []string
}

// NewEntityService creates a new entity service
func NewEntityService(db datastore.Service, caseService *CaseService) *EntityService {
	return &EntityService{
		db:          db,
		caseService: caseService,
		types:       []string{"organization", "location"},
	}
}

// Create creates a new entity
func (s *EntityService) Create(ctx context.Context, r api.EntityCreateRequest) (*api.EntityCreateResponse, error) {
	return nil, errors.New("not implemented")
}

// Update updates an existing entity
func (s *EntityService) Update(ctx context.Context, r api.EntityUpdateRequest) (*api.EntityUpdateResponse, error) {
	return nil, errors.New("not implemented")
}

// Delete deletes an existing entity
func (s *EntityService) Delete(ctx context.Context, r api.EntityDeleteRequest) (*api.EntityDeleteResponse, error) {
	return nil, errors.New("not implemented")
}

// Get the specified entity
func (s *EntityService) Get(ctx context.Context, r api.EntityGetRequest) (*api.EntityGetResponse, error) {
	return nil, errors.New("not implemented")
}

// List all entities
func (s *EntityService) List(ctx context.Context, r api.EntityListRequest) (*api.EntityListResponse, error) {
	return nil, errors.New("not implemented")
}

// Types returns the existing entity-types
func (s *EntityService) Types(ctx context.Context, r api.EntityTypesRequest) (*api.EntityTypesResponse, error) {
	return &api.EntityTypesResponse{EntityTypes: s.types}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *EntityService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
