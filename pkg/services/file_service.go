package services

import (
	"context"
	"encoding/base64"
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
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	data, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	f, err := s.store.Upload(r.CaseID, r.Name, data)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
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
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileNewResponse{New: file}, nil
}

// Open opens a file from the backend
func (s *FileService) Open(ctx context.Context, r api.FileOpenRequest) (*api.FileOpenResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	file, err := s.db.GetFile(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	content, err := s.store.GetContent(file.Path)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileOpenResponse{Data: base64.URLEncoding.EncodeToString(content)}, nil
}

// Update updates the information for a file
func (s *FileService) Update(ctx context.Context, r api.FileUpdateRequest) (*api.FileUpdateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	file, err := s.db.UpdateFile(ctx, r.CaseID, r.ID, r.Description)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileUpdateResponse{Updated: *file}, nil
}

// Delete deletes the specified file
func (s *FileService) Delete(ctx context.Context, r api.FileDeleteRequest) (*api.FileDeleteResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	file, err := s.db.GetFile(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	if err := s.db.DeleteFile(ctx, r.CaseID, file.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	if err := s.store.Delete(r.CaseID, file.Name); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileDeleteResponse{}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *FileService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
