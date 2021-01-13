package services

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
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
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if !s.validType(r.Type) {
		return nil, api.ErrInvalidEntityType
	}

	entity := api.Entity{
		Title:    r.Title,
		PhotoURL: r.PhotoURL,
		Type:     r.Type,
		Custom:   r.Custom,
	}

	if err := s.db.CreateEntity(ctx, r.CaseID, &entity); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EntityCreateResponse{Created: entity}, nil
}

// Update updates an existing entity
func (s *EntityService) Update(ctx context.Context, r api.EntityUpdateRequest) (*api.EntityUpdateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if !s.validType(r.Type) {
		return nil, api.ErrInvalidEntityType
	}

	entity := api.Entity{
		Title:    r.Title,
		PhotoURL: r.PhotoURL,
		Type:     r.Type,
		Custom:   r.Custom,
	}
	entity.ID = r.ID

	if err := s.db.UpdateEntity(ctx, r.CaseID, &entity); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EntityUpdateResponse{Updated: entity}, nil
}

// Delete deletes an existing entity
func (s *EntityService) Delete(ctx context.Context, r api.EntityDeleteRequest) (*api.EntityDeleteResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if err := s.db.DeleteEntity(ctx, r.CaseID, r.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EntityDeleteResponse{}, nil
}

// Get the specified entity
func (s *EntityService) Get(ctx context.Context, r api.EntityGetRequest) (*api.EntityGetResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	entity, err := s.db.GetEntityByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.EntityGetResponse{Entity: *entity}, nil
}

// List all entities
func (s *EntityService) List(ctx context.Context, r api.EntityListRequest) (*api.EntityListResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	entities, err := s.db.GetEntities(ctx, r.CaseID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EntityListResponse{Entities: entities}, nil
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

// validType checks if the specified type (string) is valid
func (s *EntityService) validType(entityType string) bool {
	for _, t := range s.types {
		if t == entityType {
			return true
		}
	}
	return false
}
