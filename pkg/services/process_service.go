package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
)

// ProcessService is the service
// that handles processes
type ProcessService struct {
}

// NewProcessService creates a new process-service
func NewProcessService() *ProcessService {
	return &ProcessService{}
}

// Start starts a processing with the specified files
func (s *ProcessService) Start(ctx context.Context, r api.ProcessStartRequest) (*api.ProcessStartResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Jobs returns the status of all processing-jobs
// in the specified case
func (s *ProcessService) Jobs(ctx context.Context, r api.ProcessJobsRequest) (*api.ProcessJobsResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Abort aborts the specified processing-job
func (s *ProcessService) Abort(ctx context.Context, r api.ProcessAbortRequest) (*api.ProcessAbortResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Pause pauses the specified processing-job
func (s *ProcessService) Pause(ctx context.Context, r api.ProcessPauseRequest) (*api.ProcessPauseResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *ProcessService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}
