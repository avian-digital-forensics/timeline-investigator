package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
)

// FileService handles files
type FileService struct {
}

// NewFileService creates a new file-service
func NewFileService() *FileService {
	return &FileService{}
}

// New uploads a file to the backend
func (s *FileService) New(ctx context.Context, r api.FileNewRequest) (*api.FileNewResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Update updates the information for a file
func (s *FileService) Update(ctx context.Context, r api.FileUpdateRequest) (*api.FileUpdateResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Delete deletes the specified file
func (s *FileService) Delete(ctx context.Context, r api.FileDeleteRequest) (*api.FileDeleteResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *FileService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}
