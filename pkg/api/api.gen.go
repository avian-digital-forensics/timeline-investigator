// Code generated by oto; DO NOT EDIT.

package api

import (
	"github.com/pacedotdev/oto/otohttp"

	context "context"

	http "net/http"
)

// CaseService is the API to handle cases
type CaseService interface {
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// Delete deletes the specified case
	Delete(context.Context, CaseDeleteRequest) (*CaseDeleteResponse, error)
	// Get returns the requested case
	Get(context.Context, CaseGetRequest) (*CaseGetResponse, error)
	// List the cases for a specified user
	List(context.Context, CaseListRequest) (*CaseListResponse, error)
	// New creates a new case
	New(context.Context, CaseNewRequest) (*CaseNewResponse, error)
	// Update updates the specified case
	Update(context.Context, CaseUpdateRequest) (*CaseUpdateResponse, error)
}

// FileService is the API for handling files
type FileService interface {
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// Delete deletes the specified file
	Delete(context.Context, FileDeleteRequest) (*FileDeleteResponse, error)
	// New uploads a file to the backend
	New(context.Context, FileNewRequest) (*FileNewResponse, error)
	// Update updates the information for a file
	Update(context.Context, FileUpdateRequest) (*FileUpdateResponse, error)
}

// ProcessService is the API - that handles evidence-processing
type ProcessService interface {
	// Abort aborts the specified processing-job
	Abort(context.Context, ProcessAbortRequest) (*ProcessAbortResponse, error)
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// Jobs returns the status of all processing-jobs in the specified case
	Jobs(context.Context, ProcessJobsRequest) (*ProcessJobsResponse, error)
	// Pause pauses the specified processing-job
	Pause(context.Context, ProcessPauseRequest) (*ProcessPauseResponse, error)
	// Start starts a processing with the specified files
	Start(context.Context, ProcessStartRequest) (*ProcessStartResponse, error)
}

// TestService is used for testing-purposes
type TestService interface {
	// CreateUser creates a test-user in Firebase
	CreateUser(context.Context, TestCreateUserRequest) (*TestCreateUserResponse, error)
	// DeleteUser deletes a test-user in Firebase
	DeleteUser(context.Context, TestDeleteUserRequest) (*TestDeleteUserResponse, error)
}

type caseServiceServer struct {
	server      *otohttp.Server
	caseService CaseService
	test        bool
}

// Register adds the CaseService to the otohttp.Server.
func RegisterCaseService(server *otohttp.Server, caseService CaseService) {
	handler := &caseServiceServer{
		server:      server,
		caseService: caseService,
	}

	server.Register("CaseService", "Delete", handler.handleDelete)
	server.Register("CaseService", "Get", handler.handleGet)
	server.Register("CaseService", "List", handler.handleList)
	server.Register("CaseService", "New", handler.handleNew)
	server.Register("CaseService", "Update", handler.handleUpdate)
}

func (s *caseServiceServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	var request CaseDeleteRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.caseService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.caseService.Delete(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *caseServiceServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var request CaseGetRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.caseService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.caseService.Get(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *caseServiceServer) handleList(w http.ResponseWriter, r *http.Request) {
	var request CaseListRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.caseService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.caseService.List(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *caseServiceServer) handleNew(w http.ResponseWriter, r *http.Request) {
	var request CaseNewRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.caseService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.caseService.New(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *caseServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var request CaseUpdateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.caseService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.caseService.Update(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

type fileServiceServer struct {
	server      *otohttp.Server
	fileService FileService
	test        bool
}

// Register adds the FileService to the otohttp.Server.
func RegisterFileService(server *otohttp.Server, fileService FileService) {
	handler := &fileServiceServer{
		server:      server,
		fileService: fileService,
	}

	server.Register("FileService", "Delete", handler.handleDelete)
	server.Register("FileService", "New", handler.handleNew)
	server.Register("FileService", "Update", handler.handleUpdate)
}

func (s *fileServiceServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	var request FileDeleteRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.fileService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.fileService.Delete(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *fileServiceServer) handleNew(w http.ResponseWriter, r *http.Request) {
	var request FileNewRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.fileService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.fileService.New(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *fileServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var request FileUpdateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.fileService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.fileService.Update(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

type processServiceServer struct {
	server         *otohttp.Server
	processService ProcessService
	test           bool
}

// Register adds the ProcessService to the otohttp.Server.
func RegisterProcessService(server *otohttp.Server, processService ProcessService) {
	handler := &processServiceServer{
		server:         server,
		processService: processService,
	}
	server.Register("ProcessService", "Abort", handler.handleAbort)

	server.Register("ProcessService", "Jobs", handler.handleJobs)
	server.Register("ProcessService", "Pause", handler.handlePause)
	server.Register("ProcessService", "Start", handler.handleStart)
}

func (s *processServiceServer) handleAbort(w http.ResponseWriter, r *http.Request) {
	var request ProcessAbortRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.processService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.processService.Abort(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *processServiceServer) handleJobs(w http.ResponseWriter, r *http.Request) {
	var request ProcessJobsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.processService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.processService.Jobs(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *processServiceServer) handlePause(w http.ResponseWriter, r *http.Request) {
	var request ProcessPauseRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.processService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.processService.Pause(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *processServiceServer) handleStart(w http.ResponseWriter, r *http.Request) {
	var request ProcessStartRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.processService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.processService.Start(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

type testServiceServer struct {
	server      *otohttp.Server
	testService TestService
	test        bool
}

// Register adds the TestService to the otohttp.Server.
func RegisterTestService(server *otohttp.Server, testService TestService) {
	handler := &testServiceServer{
		server:      server,
		testService: testService,
	}
	server.Register("TestService", "CreateUser", handler.handleCreateUser)
	server.Register("TestService", "DeleteUser", handler.handleDeleteUser)
}

func (s *testServiceServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request TestCreateUserRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx := r.Context()
	response, err := s.testService.CreateUser(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *testServiceServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var request TestDeleteUserRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx := r.Context()
	response, err := s.testService.DeleteUser(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

// Base model for the database
type Base struct {
	// ID is the identifier for the object
	ID string `json:"id"`
	// CreatedAt - when the object was created
	CreatedAt int64 `json:"createdAt"`
	// UpdatedAt - when the object was updated
	UpdatedAt int64 `json:"updatedAt"`
	// DeletedAt - when the object was deleted
	DeletedAt int64 `json:"deletedAt"`
}

// File holds information about an uploaded file
type File struct {
	Base
	// Name of the file
	Name string `json:"name"`
	// Mime is the mime-type of the file
	Mime string `json:"mime"`
	// Description of the file
	Description string `json:"description"`
	// Path to where the file is stored
	Path string `json:"path"`
	// Size of the file in bytes
	Size int `json:"size"`
	// Processed is if the file has been processed or not
	Processed bool `json:"processed"`
}

// Process holds information about a job that processes data to app
type Process struct {
	Base
}

// Case is an object to hold data for a specific investigation
type Case struct {
	Base
	// CreatorID is the user-id of the user who created the case (super admin)
	CreatorID string `json:"creatorID"`
	// Name of the case
	Name string `json:"name"`
	// Description of the case
	Description string `json:"description"`
	// FromDate is the unix-date for the start of the primary timespan for the case
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-date for the end of the primary timespan for the case
	ToDate int64 `json:"toDate"`
	// Investigators of the case (users who has access to the case)
	Investigators []string `json:"investigators"`
	// Files that exists in the case
	Files []File `json:"files"`
	// Processes that exists in the case
	Processes []Process `json:"processes"`
}

// CaseDeleteRequest is the input-object for deleting an existing case
type CaseDeleteRequest struct {
	// ID of the case to delete
	ID string `json:"id"`
}

// CaseDeleteResponse is the output-object for deleting an existing case
type CaseDeleteResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// CaseGetRequest is the input-object for getting a specified case
type CaseGetRequest struct {
	// ID of the case to get
	ID string `json:"id"`
}

// CaseGetResponse is the output-object for getting a specified case
type CaseGetResponse struct {
	Case Case `json:"case"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// CaseListRequest is the input-object for listing cases for a specified user
type CaseListRequest struct {
	// UserID of the user to list cases for
	UserID string `json:"userID"`
}

// CaseListResponse is the output-object for listing cases for a specified user
type CaseListResponse struct {
	Cases []Case `json:"cases"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// CaseNewRequest is the input-object for creating a new case
type CaseNewRequest struct {
	// Name of the case
	Name string `json:"name"`
	// description of the case to create
	Description string `json:"description"`
	// FromDate is the unix-date for the start of the primary timespan for the case
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-date for the end of the primary timespan for the case
	ToDate int64 `json:"toDate"`
}

// CaseNewResponse is the output-object for creating a new case
type CaseNewResponse struct {
	New Case `json:"new"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// CaseUpdateRequest is the input-object for updating an existing case
type CaseUpdateRequest struct {
	// ID of the case to update
	ID string `json:"id"`
	// Name of the case
	Name string `json:"name"`
	// description of the case to create
	Description string `json:"description"`
	// FromDate is the unix-date for the start of the primary timespan for the case
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-date for the end of the primary timespan for the case
	ToDate int64 `json:"toDate"`
}

// CaseUpdateResponse is the output-object for updating an existing case
type CaseUpdateResponse struct {
	Updated Case `json:"updated"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// CaseUploadRequest is the input-object for uploading an evidence to the case
type CaseUploadRequest struct {
	// ID of the case to upload
	ID string `json:"id"`
	// Name of the item to upload
	Name string `json:"name"`
}

// FileDeleteRequest is the input-object for deleting a file
type FileDeleteRequest struct {
	// ID of the file to delete
	ID string `json:"id"`
	// CaseID of the case where the file to delete belongs
	CaseID string `json:"caseID"`
}

// FileDeleteResponse is the output-object for deleting a file
type FileDeleteResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// FileNewRequest is the input-object for creating a new file
type FileNewRequest struct {
	// CaseID of the case to upload the file
	CaseID string `json:"caseID"`
	// Name of the file
	Name string `json:"name"`
	// Description of the file
	Description string `json:"description"`
	// Mime is the mime-type of the file
	Mime string `json:"mime"`
	// Data of the file (base64 encoded)
	Data string `json:"data"`
}

// FileNewResponse is the output-object for creating a new file
type FileNewResponse struct {
	New File `json:"new"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// FileUpdateRequest is the input-object for updating a files information
type FileUpdateRequest struct {
	// ID of the file to update
	ID string `json:"id"`
	// CaseID of the case where the file to update belongs
	CaseID string `json:"caseID"`
	// Description of the file
	Description string `json:"description"`
}

// FileUpdateResponse is the output-object for updating a files information
type FileUpdateResponse struct {
	Updated File `json:"updated"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// ProcessAbortRequest is the input-object for aborting a processing-job
type ProcessAbortRequest struct {
	// ID of the processing-job to abort
	ID string `json:"id"`
	// CaseID of the case the processing-job belongs to
	CaseID string `json:"caseID"`
}

// ProcessAbortResponse is the output-object for aborting a processing-job
type ProcessAbortResponse struct {
	Aborted Process `json:"aborted"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// ProcessJobsRequest is the input-object for getting all processing-jobs for a
// case
type ProcessJobsRequest struct {
	// CaseID of the case to get the processing-jobs for
	CaseID string `json:"caseID"`
}

// ProcessJobsResponse is the output-object for getting all processing-jobs for a
// case
type ProcessJobsResponse struct {
	Processes []Process `json:"processes"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// ProcessPauseRequest is the input-object for pausing a processing-job
type ProcessPauseRequest struct {
	// ID of the processing-job to pause
	ID string `json:"id"`
	// CaseID of the case the processing-job belongs to
	CaseID string `json:"caseID"`
}

// ProcessPauseResponse is the output-object for pausing a processing-job
type ProcessPauseResponse struct {
	Paused Process `json:"paused"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// ProcessStartRequest is the input-object for starting a processing-job
type ProcessStartRequest struct {
	// CaseID of the case to start the processing for
	CaseID string `json:"caseID"`
	// FileIDs of the files to process
	FileIDs []string `json:"fileIDs"`
}

// ProcessStartResponse is the output-object for starting a processing-job
type ProcessStartResponse struct {
	Started Process `json:"started"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// TestCreateUserRequest is the input-object for creating a test-user
type TestCreateUserRequest struct {
	// Name of the user to create
	Name string `json:"name"`
	// ID of the user to create
	ID string `json:"id"`
	// Email of the user to create
	Email string `json:"email"`
	// Password for the new user
	Password string `json:"password"`
	// Secret for using the test-service
	Secret string `json:"secret,omitempty"`
}

// TestCreateUserResponse is the output-object for creating a test-user
type TestCreateUserResponse struct {
	// Token for the created user
	Token string `json:"token"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// TestDeleteUserRequest is the input-object for deleting a test-user
type TestDeleteUserRequest struct {
	// ID of the user to delete
	ID string `json:"id"`
	// Secret for using the test-service
	Secret string `json:"secret,omitempty"`
}

// TestDeleteUserResponse is the output-object for deleting a test-user
type TestDeleteUserResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}
