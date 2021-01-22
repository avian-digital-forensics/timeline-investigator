package services

import (
	"context"
	"log"
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

	entity, err := s.db.GetEntityByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	entity.Title = r.Title
	entity.PhotoURL = r.PhotoURL
	entity.Type = r.Type
	entity.Custom = r.Custom

	if err := s.db.UpdateEntity(ctx, r.CaseID, entity); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.EntityUpdateResponse{Updated: *entity}, nil
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

// KeywordsAdd adds keywords to an entity
func (s *EntityService) KeywordsAdd(ctx context.Context, r api.KeywordsAddRequest) (*api.KeywordsAddResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the entity to add the keyword to
	entity, err := s.db.GetEntityByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be added to the entity
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Create a map of the keywords that were found from the db
	// and add the entity ID to each keyword
	var keywordFound = make(map[string]bool)
	for _, keyword := range keywords {
		keywordFound[keyword.Name] = true
		keyword.EntityIDs = append(keyword.EntityIDs, entity.ID)
	}

	// Add the keywords from the request to the entity
	// and append the keywords that didn't already exist
	// to the keyword-slice
	for _, keyword := range r.Keywords {
		if !keywordFound[keyword] {
			keywords = append(keywords, api.Keyword{Name: keyword, EntityIDs: []string{entity.ID}})
		}
		entity.Keywords = append(entity.Keywords, keyword)
	}

	// Save the keywords with the entity ID
	// TODO / FIXME: Use bulk-indexer instead
	for _, keyword := range keywords {
		if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
			return nil, api.Error(err, api.ErrCannotPerformOperation)
		}
	}

	// Update the entity with the added keywords
	if err := s.db.UpdateEntity(ctx, r.CaseID, entity); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsAddResponse{OK: true}, nil
}

// KeywordsRemove removes keywords from an entity
func (s *EntityService) KeywordsRemove(ctx context.Context, r api.KeywordsRemoveRequest) (*api.KeywordsRemoveResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Create a map of the keywords to remove
	var keywordToRemove = make(map[string]bool)
	for _, keyword := range r.Keywords {
		keywordToRemove[keyword] = true
	}

	// Get the entity to remove the keywords from
	entity, err := s.db.GetEntityByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be removed from the entity
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Remove the keywords from the entity
	for i, keyword := range entity.Keywords {
		if keywordToRemove[keyword] {
			entity.Keywords = append(entity.Keywords[:i], entity.Keywords[i+1:]...)
		}
	}

	// Remove the entityID from the keywords
	for ki, keyword := range keywords {
		if keywordToRemove[keyword.Name] {
			for ei, id := range keyword.EntityIDs {
				if id == entity.ID {
					keywords[ki].EntityIDs = append(
						keyword.EntityIDs[:ei],
						keyword.EntityIDs[ei+1:]...,
					)
				}
			}
		}
	}

	// Save the keywords (or delete if empty)
	// TODO/FIXME: use bulk indexer
	for _, keyword := range keywords {
		log.Println(keyword)
		toDelete := len(keyword.EntityIDs) == 0 && len(keyword.PersonIDs) == 0 && len(keyword.EventIDs) == 0 && len(keyword.FileIDs) == 0
		log.Println(toDelete)
		if toDelete {
			if err := s.db.DeleteKeyword(ctx, r.CaseID, keyword.Name); err != nil {
				return nil, api.Error(err, api.ErrCannotPerformOperation)
			}
		} else if !toDelete {
			if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
				return nil, api.Error(err, api.ErrCannotPerformOperation)
			}
		}
	}

	// Update the entity without the removed keyword
	if err := s.db.UpdateEntity(ctx, r.CaseID, entity); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsRemoveResponse{}, nil
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
