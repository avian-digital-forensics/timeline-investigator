package services

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/filestore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

// FileService handles files
type FileService struct {
	db          datastore.Service
	store       filestore.Service
	caseService *CaseService
}

// NewFileService creates a new file-service
func NewFileService(
	db datastore.Service,
	store filestore.Service,
	caseService *CaseService,
) *FileService {
	return &FileService{
		db:          db,
		store:       store,
		caseService: caseService,
	}
}

// New uploads a file to the backend
func (s *FileService) New(ctx context.Context, r api.FileNewRequest) (*api.FileNewResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, err
		}
		return nil, api.ErrNotAllowed
	}

	f, err := s.store.Upload(r.CaseID, r.Name, []byte(r.Data))
	if err != nil {
		return nil, err
	}

	file := api.File{
		Name:        f.Name,
		Mime:        r.Mime,
		Description: r.Description,
		Path:        f.Path,
		Size:        f.Size,
		Processed:   false,
	}

	if err := s.db.CreateFile(ctx, r.CaseID, &file); err != nil {
		return nil, err
	}

	return &api.FileNewResponse{New: file}, nil
}

// Update updates the information for a file
func (s *FileService) Update(ctx context.Context, r api.FileUpdateRequest) (*api.FileUpdateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, err
		}
		return nil, api.ErrNotAllowed
	}

	file, err := s.db.UpdateFile(ctx, r.CaseID, r.ID, r.Description)
	if err != nil {
		return nil, err
	}

	return &api.FileUpdateResponse{Updated: *file}, nil
}

// Delete deletes the specified file
func (s *FileService) Delete(ctx context.Context, r api.FileDeleteRequest) (*api.FileDeleteResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, err
		}
		return nil, api.ErrNotAllowed
	}

	file, err := s.db.GetFile(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, err
	}

	if err := s.db.DeleteFile(ctx, r.CaseID, file.ID); err != nil {
		return nil, err
	}

	if err := s.store.Delete(r.CaseID, file.Name); err != nil {
		return nil, err
	}

	return &api.FileDeleteResponse{}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *FileService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}
