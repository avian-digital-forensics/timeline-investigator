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

	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
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

	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
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
	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
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

	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Delete the keywords for the file
	if err := s.removeKeywords(ctx, r.CaseID, file, file.Keywords); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	if err := s.db.DeleteFile(ctx, r.CaseID, file.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	if err := s.store.Delete(r.CaseID, file.Name); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.FileDeleteResponse{}, nil
}

// KeywordsAdd adds keywords to a file
func (s *FileService) KeywordsAdd(ctx context.Context, r api.KeywordsAddRequest) (*api.KeywordsAddResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the file to add the keyword to
	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be added to the file
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Create a map of the keywords that were found from the db
	// and add the File ID to each keyword
	var keywordFound = make(map[string]bool)
	for i := range keywords {
		keywordFound[keywords[i].Name] = true
		keywords[i].FileIDs = append(keywords[i].FileIDs, file.ID)
	}

	// Add the keywords from the request to the file
	// and append the keywords that didn't already exist
	// to the keyword-slice
	for _, keyword := range r.Keywords {
		if !keywordFound[keyword] {
			keywords = append(keywords, api.Keyword{Name: keyword, FileIDs: []string{file.ID}})
		}
		file.Keywords = append(file.Keywords, keyword)
	}

	// Save the keywords with the File ID
	// TODO / FIXME: Use bulk-indexer instead
	for _, keyword := range keywords {
		if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
			return nil, api.Error(err, api.ErrCannotPerformOperation)
		}
	}

	// Update the file with the added keywords
	if err := s.db.UpdateFile(ctx, r.CaseID, file); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsAddResponse{OK: true}, nil
}

// KeywordsRemove removes keywords from a file
func (s *FileService) KeywordsRemove(ctx context.Context, r api.KeywordsRemoveRequest) (*api.KeywordsRemoveResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the file to remove the keywords from
	file, err := s.db.GetFileByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Delete the requested keywords for the file
	if err := s.removeKeywords(ctx, r.CaseID, file, r.Keywords); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Update the file without the removed keyword
	if err := s.db.UpdateFile(ctx, r.CaseID, file); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsRemoveResponse{}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *FileService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}

func (s *FileService) removeKeywords(ctx context.Context, caseID string, file *api.File, removeKeywords []string) error {
	// Create a map of the keywords to remove
	var keywordToRemove = make(map[string]bool)
	for _, keyword := range removeKeywords {
		keywordToRemove[keyword] = true
	}

	// Get the keywords that should be removed from the file
	keywords, err := s.db.GetKeywordsByIDs(ctx, caseID, removeKeywords)
	if err != nil {
		return err
	}

	// Remove the keywords from the file
	for i, keyword := range file.Keywords {
		if keywordToRemove[keyword] {
			file.Keywords = append(file.Keywords[:i], file.Keywords[i+1:]...)
		}
	}

	// Remove the FileID from the keywords
	for ki, keyword := range keywords {
		if keywordToRemove[keyword.Name] {
			for ei, id := range keyword.FileIDs {
				if id == file.ID {
					keywords[ki].FileIDs = append(
						keyword.FileIDs[:ei],
						keyword.FileIDs[ei+1:]...,
					)
				}
			}
		}
	}

	// Save the keywords (or delete if empty)
	// TODO/FIXME: use bulk indexer
	for _, keyword := range keywords {
		toDelete := len(keyword.EntityIDs) == 0 && len(keyword.PersonIDs) == 0 && len(keyword.EventIDs) == 0 && len(keyword.FileIDs) == 0
		if toDelete {
			if err := s.db.DeleteKeyword(ctx, caseID, keyword.Name); err != nil {
				return err
			}
		} else if !toDelete {
			if err := s.db.SaveKeyword(ctx, caseID, &keyword); err != nil {
				return err
			}
		}
	}
	return nil
}
