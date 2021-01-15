package services

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/filestore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/fscrawler"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

// FileService handles files
type FileService struct {
	db          datastore.Service
	store       filestore.Service
	caseService *CaseService
	fs          *fscrawler.Client
}

// NewFileService creates a new file-service
func NewFileService(
	db datastore.Service,
	store filestore.Service,
	caseService *CaseService,
	fs *fscrawler.Client,
) *FileService {
	return &FileService{
		db:          db,
		store:       store,
		caseService: caseService,
		fs:          fs,
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
		ProcessedAt: 0,
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

// Process Processs a file from the backend
func (s *FileService) Process(ctx context.Context, r api.FileProcessRequest) (*api.FileProcessResponse, error) {
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

	// Process the file
	process := s.fs.NewProcess(file.Path).WithID(file.ID).WithIndex(s.db.ProcessIndex(r.CaseID))
	if err := process.Start(ctx); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	file.ProcessedAt = time.Now().Unix()

	if err := s.db.UpdateFile(ctx, r.CaseID, file); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileProcessResponse{Processed: *file}, nil
}

// Processed gets information for a processed file
func (s *FileService) Processed(ctx context.Context, r api.FileProcessedRequest) (*api.FileProcessedResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	processed, err := s.db.GetProcessedFile(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.FileProcessedResponse{ID: r.ID, Processed: processed}, nil
}

// Processes gets information for all proccesed files in the specified case
func (s *FileService) Processes(ctx context.Context, r api.FileProcessesRequest) (*api.FileProcessesResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	processes, err := s.db.GetProcessedFiles(ctx, r.CaseID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.FileProcessesResponse{Processes: processes}, nil
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

	// Get the file to update
	file, err := s.db.GetFile(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Set the description and updatedAt
	file.UpdatedAt = time.Now().Unix()
	file.Description = r.Description

	// update the file
	if err := s.db.UpdateFile(ctx, r.CaseID, file); err != nil {
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
