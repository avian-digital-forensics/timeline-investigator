package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
)

// CaseService handles cases
type CaseService struct {
	db datastore.Service
}

// NewCaseService creates a new case-service
func NewCaseService(db datastore.Service) *CaseService {
	return &CaseService{db: db}
}

// New creates a new case
func (s *CaseService) New(ctx context.Context, r api.CaseNewRequest) (*api.CaseNewResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Get returns the requested case
func (s *CaseService) Get(ctx context.Context, r api.CaseGetRequest) (*api.CaseGetResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Update updates the specified case
func (s *CaseService) Update(ctx context.Context, r api.CaseUpdateRequest) (*api.CaseUpdateResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Delete deletes the specified case
func (s *CaseService) Delete(ctx context.Context, r api.CaseDeleteRequest) (*api.CaseDeleteResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// List the cases for a specified user
func (s *CaseService) List(ctx context.Context, r api.CaseListRequest) (*api.CaseListResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *CaseService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}
