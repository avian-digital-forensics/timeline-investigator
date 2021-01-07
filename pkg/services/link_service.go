package services

import (
	"context"
	"errors"
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

// CreateEvent creates a link for an event
// with multiple objects
func (s *LinkService) CreateEvent(ctx context.Context, r api.LinkEventCreateRequest) (*api.LinkEventCreateResponse, error) {
	return nil, errors.New("not implemented")
}

// GetEvent gets an event with its links
func (s *LinkService) GetEvent(ctx context.Context, r api.LinkEventCreateRequest) (*api.LinkEventCreateResponse, error) {
	return nil, errors.New("not implemented")
}

// DeleteEvent deletes all links to the specified event
func (s *LinkService) DeleteEvent(ctx context.Context, r api.LinkEventDeleteRequest) (*api.LinkEventDeleteResponse, error) {
	return nil, errors.New("not implemented")
}

// UpdateEvent updates links for the specified event
func (s *LinkService) UpdateEvent(ctx context.Context, r api.LinkEventUpdateRequest) (*api.LinkEventUpdateResponse, error) {
	return nil, errors.New("not implemented")
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *LinkService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
