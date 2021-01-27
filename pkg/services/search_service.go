package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
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
	return nil, errors.New("not implemented yet")
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *SearchService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
