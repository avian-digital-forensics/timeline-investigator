// Code generated by oto; DO NOT EDIT.

package api

import (
	"github.com/pacedotdev/oto/otohttp"

	http "net/http"

	context "context"
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

// EventService is the API to handle events
type EventService interface {
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// Create creates a new event
	Create(context.Context, EventCreateRequest) (*EventCreateResponse, error)
	// Delete deletes an existing event
	Delete(context.Context, EventDeleteRequest) (*EventDeleteResponse, error)
	// Get the specified event
	Get(context.Context, EventGetRequest) (*EventGetResponse, error)
	// List all events
	List(context.Context, EventListRequest) (*EventListResponse, error)
	// Update updates an existing event
	Update(context.Context, EventUpdateRequest) (*EventUpdateResponse, error)
}

// FileService is the API for handling files
type FileService interface {
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// Delete deletes the specified file
	Delete(context.Context, FileDeleteRequest) (*FileDeleteResponse, error)
	// New uploads a file to the backend
	New(context.Context, FileNewRequest) (*FileNewResponse, error)
	// Open opens a file
	Open(context.Context, FileOpenRequest) (*FileOpenResponse, error)
	// Update updates the information for a file
	Update(context.Context, FileUpdateRequest) (*FileUpdateResponse, error)
}

// LinkService is a API for creating links between objects
type LinkService interface {
	// Authenticate is a middleware in the http-handler
	Authenticate(context.Context, *http.Request) (context.Context, error)
	// CreateEvent creates a link for an event with multiple objects
	CreateEvent(context.Context, LinkEventCreateRequest) (*LinkEventCreateResponse, error)
	// DeleteEvent deletes all links to the specified event
	DeleteEvent(context.Context, LinkEventDeleteRequest) (*LinkEventDeleteResponse, error)
	// GetEvent gets an event with its links
	GetEvent(context.Context, LinkEventGetRequest) (*LinkEventGetResponse, error)
	// UpdateEvent updates links for the specified event
	UpdateEvent(context.Context, LinkEventUpdateRequest) (*LinkEventUpdateResponse, error)
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

type eventServiceServer struct {
	server       *otohttp.Server
	eventService EventService
	test         bool
}

// Register adds the EventService to the otohttp.Server.
func RegisterEventService(server *otohttp.Server, eventService EventService) {
	handler := &eventServiceServer{
		server:       server,
		eventService: eventService,
	}

	server.Register("EventService", "Create", handler.handleCreate)
	server.Register("EventService", "Delete", handler.handleDelete)
	server.Register("EventService", "Get", handler.handleGet)
	server.Register("EventService", "List", handler.handleList)
	server.Register("EventService", "Update", handler.handleUpdate)
}

func (s *eventServiceServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	var request EventCreateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.eventService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.eventService.Create(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *eventServiceServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	var request EventDeleteRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.eventService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.eventService.Delete(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *eventServiceServer) handleGet(w http.ResponseWriter, r *http.Request) {
	var request EventGetRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.eventService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.eventService.Get(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *eventServiceServer) handleList(w http.ResponseWriter, r *http.Request) {
	var request EventListRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.eventService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.eventService.List(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *eventServiceServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var request EventUpdateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.eventService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.eventService.Update(ctx, request)
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
	server.Register("FileService", "Open", handler.handleOpen)
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

func (s *fileServiceServer) handleOpen(w http.ResponseWriter, r *http.Request) {
	var request FileOpenRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.fileService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.fileService.Open(ctx, request)
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

type linkServiceServer struct {
	server      *otohttp.Server
	linkService LinkService
	test        bool
}

// Register adds the LinkService to the otohttp.Server.
func RegisterLinkService(server *otohttp.Server, linkService LinkService) {
	handler := &linkServiceServer{
		server:      server,
		linkService: linkService,
	}

	server.Register("LinkService", "CreateEvent", handler.handleCreateEvent)
	server.Register("LinkService", "DeleteEvent", handler.handleDeleteEvent)
	server.Register("LinkService", "GetEvent", handler.handleGetEvent)
	server.Register("LinkService", "UpdateEvent", handler.handleUpdateEvent)
}

func (s *linkServiceServer) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var request LinkEventCreateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.linkService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.linkService.CreateEvent(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *linkServiceServer) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	var request LinkEventDeleteRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.linkService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.linkService.DeleteEvent(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *linkServiceServer) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	var request LinkEventGetRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.linkService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.linkService.GetEvent(ctx, request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *linkServiceServer) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	var request LinkEventUpdateRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	ctx, err := s.linkService.Authenticate(r.Context(), r)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.linkService.UpdateEvent(ctx, request)
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

// Event is an important happening that needs investigation.
type Event struct {
	Base
	// Set the importance of the event, defined by a number between 1 - 5.
	Importance int `json:"importance"`
	// Desription of the event.
	Description string `json:"description"`
	// FromDate is the unix-timestamp of when the event started
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-timestamp of when the event finished
	ToDate int64 `json:"toDate"`
}

// EventCreateRequest is the input-object for creating an event
type EventCreateRequest struct {
	// CaseID of the case to create the event for
	CaseID string `json:"caseID"`
	// Set the importance of the event, defined by a number between 1 - 5.
	Importance int `json:"importance"`
	// Desription of the event.
	Description string `json:"description"`
	// FromDate is the unix-timestamp of when the event started
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-timestamp of when the event finished
	ToDate int64 `json:"toDate"`
}

// EventCreateResponse is the output-object for creating an event
type EventCreateResponse struct {
	Created Event `json:"created"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// EventDeleteRequest is the input-object for deleting an existing event
type EventDeleteRequest struct {
	// ID of the event to Delete
	ID string `json:"id"`
	// CaseID of the event
	CaseID string `json:"caseID"`
}

// EventDeleteResponse is the output-object for deleting an existing event
type EventDeleteResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// EventGetRequest is the input-object for getting an existing event
type EventGetRequest struct {
	// ID of the event to get
	ID string `json:"id"`
	// CaseID of the event
	CaseID string `json:"caseID"`
}

// EventGetResponse is the output-object for deleting an existing event
type EventGetResponse struct {
	Event Event `json:"event"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// EventListRequest is the input-object for listing all existing events for a case
type EventListRequest struct {
	// CaseID to list the events for
	CaseID string `json:"caseID"`
}

// EventListResponse is the output-object for listing all existing events for a
// case
type EventListResponse struct {
	Events []Event `json:"events"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// EventUpdateRequest is the input-object for updating an existing event
type EventUpdateRequest struct {
	// ID of the event to update
	ID string `json:"id"`
	// CaseID of the event
	CaseID string `json:"caseID"`
	// Set the importance of the event, defined by a number between 1 - 5.
	Importance int `json:"importance"`
	// Desription of the event.
	Description string `json:"description"`
	// FromDate is the unix-timestamp of when the event started
	FromDate int64 `json:"fromDate"`
	// ToDate is the unix-timestamp of when the event finished
	ToDate int64 `json:"toDate"`
}

// EventUpdateResponse is the output-object for updating an existing event
type EventUpdateResponse struct {
	Updated Event `json:"updated"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
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
	// Mime is the mime-type of the file (decided by frontend)
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

// FileOpenRequest is the input-object for opening a file in a case
type FileOpenRequest struct {
	// ID of the file to open
	ID string `json:"id"`
	// CaseID of the case to open the file
	CaseID string `json:"caseID"`
}

// FileOpenResponse is the output-object for opening a file in a case
type FileOpenResponse struct {
	// Data contains the b64-encoded data for the file
	Data string `json:"data"`
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

// LinkEvent is a link for an event between different objects
type LinkEvent struct {
	Base
	// From which event has been linked
	From Event `json:"from"`
	// Events that has been linked
	Events []Event `json:"events"`
}

// LinkEventCreateRequest is the input-object for linking objects with an event
type LinkEventCreateRequest struct {
	// CaseID for the event
	CaseID string `json:"caseID"`
	// FromID is the ID of the event to hold the link
	FromID string `json:"fromID"`
	// EventIDs of the events to be linked
	EventIDs []string `json:"eventIDs"`
	// Bidirectional means that he link also should be created for the "ToID"
	Bidirectional bool `json:"bidirectional"`
}

// LinkEventCreateResponse is the output-object for linking objects with an event
type LinkEventCreateResponse struct {
	Linked LinkEvent `json:"linked"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// LinkEventDeleteRequest is the input-object for removing a linked event
type LinkEventDeleteRequest struct {
	// CaseID of the case where the linked event belongs
	CaseID string `json:"caseID"`
	// EventID of the Event to delete the link for
	EventID string `json:"eventID"`
}

// LinkEventDeleteResponse is the output-object for removing a linked event
type LinkEventDeleteResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// LinkEventGetRequest is the input-object for getting a linked Event
type LinkEventGetRequest struct {
	// CaseID of the case where the event belongs
	CaseID string `json:"caseID"`
	// EventID of the Event to get all links for
	EventID string `json:"eventID"`
}

// LinkEventGetResponse is the output-object for getting a linked Event
type LinkEventGetResponse struct {
	Link LinkEvent `json:"link"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// LinkEventUpdateRequest is the input-object for updating linked objects with an
// event
type LinkEventUpdateRequest struct {
	// EventID is the ID of the event to hold the link
	EventID string `json:"eventID"`
	// CaseID for the event
	CaseID string `json:"caseID"`
	// EventAddIDs of the events to be linked
	EventAddIDs []string `json:"eventAddIDs"`
	// EventRemoveIDs of the events to be removed
	EventRemoveIDs []string `json:"eventRemoveIDs"`
}

// LinkEventUpdateResponse is the output-object for linking objects with an event
type LinkEventUpdateResponse struct {
	Updated LinkEvent `json:"updated"`
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

// User holds information for a user in the timeline-investigator
type User struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	PhotoURL    string `json:"photoURL"`
	ProviderID  string `json:"providerID"`
	UID         string `json:"uID"`
}
