package def

import (
	"context"
	"net/http"
)

// CaseService is the API to handle cases
type CaseService interface {
	// New creates a new case
	New(CaseNewRequest) CaseNewResponse

	// Get returns the requested case
	Get(CaseGetRequest) CaseGetResponse

	// Update updates the specified case
	Update(CaseUpdateRequest) CaseUpdateResponse

	// Delete deletes the specified case
	Delete(CaseDeleteRequest) CaseDeleteResponse

	// List the cases for a specified user
	List(CaseListRequest) CaseListResponse

	// Authenticate is a middleware
	// in the http-handler
	//
	// NOTE : Only for Go-servers
	Authenticate(*http.Request) context.Context
}

// FileService is the API for handling files
type FileService interface {
	// New uploads a file to the backend
	New(FileNewRequest) FileNewResponse

	// Update updates the information for a file
	Update(FileUpdateRequest) FileUpdateResponse

	// Delete deletes the specified file
	Delete(FileDeleteRequest) FileDeleteResponse

	// Authenticate is a middleware
	// in the http-handler
	//
	// NOTE : Only for Go-servers
	Authenticate(*http.Request) context.Context
}

// ProcessService is the API -
// that handles evidence-processing
type ProcessService interface {
	// Start starts a processing with the specified files
	Start(ProcessStartRequest) ProcessStartResponse

	// Jobs returns the status of all processing-jobs
	// in the specified case
	Jobs(ProcessJobsRequest) ProcessJobsResponse

	// Abort aborts the specified processing-job
	Abort(ProcessAbortRequest) ProcessAbortResponse

	// Pause pauses the specified processing-job
	Pause(ProcessPauseRequest) ProcessPauseResponse

	// Authenticate is a middleware
	// in the http-handler
	//
	// NOTE : Only for Go-servers
	Authenticate(*http.Request) context.Context
}

// Case is an object to hold
// data for a specific investigation
type Case struct {
	Base

	// CreatorID is the user-id of the user
	// who created the case (super admin)
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CreatorID string

	// Name of the case
	//
	// example: "Case 1"
	Name string

	// Description of the case
	//
	// example: "This is a case"
	Description string

	// FromDate is the unix-date for the start
	// of the primary timespan for the case
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-date for the end
	// of the primary timespan for the case
	//
	// example: 1257894000
	ToDate int64

	// Investigators of the case
	// (users who has access to the case)
	// NOTE: defined by email
	//
	// example: ["sja@avian.dk", "jis@avian.dk"]
	Investigators []string

	// Files that exists in the case
	Files []File

	// Processes that exists in the case
	Processes []Process
}

// CaseNewRequest is the input-object
// for creating a new case
type CaseNewRequest struct {
	// Name of the case
	//
	// example: "Case 1"
	Name string

	// description of the case
	// to create
	//
	// example: "This is a case"
	Description string

	// FromDate is the unix-date for the start
	// of the primary timespan for the case
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-date for the end
	// of the primary timespan for the case
	//
	// example: 1257894000
	ToDate int64
}

// CaseNewResponse is the output-object
// for creating a new case
type CaseNewResponse struct {
	New Case
}

// CaseGetRequest is the input-object
// for getting a specified case
type CaseGetRequest struct {
	// ID of the case to get
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string
}

// CaseGetResponse is the output-object
// for getting a specified case
type CaseGetResponse struct {
	Case Case
}

// CaseUpdateRequest is the input-object
// for updating an existing case
type CaseUpdateRequest struct {
	// ID of the case to update
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// Name of the case
	//
	// example: "Case 1"
	Name string

	// description of the case
	// to create
	//
	// example: "This is a case"
	Description string

	// FromDate is the unix-date for the start
	// of the primary timespan for the case
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-date for the end
	// of the primary timespan for the case
	//
	// example: 1257894000
	ToDate int64
}

// CaseUpdateResponse is the output-object
// for updating an existing case
type CaseUpdateResponse struct {
	Updated Case
}

// CaseDeleteRequest is the input-object
// for deleting an existing case
type CaseDeleteRequest struct {
	// ID of the case to delete
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string
}

// CaseDeleteResponse is the output-object
// for deleting an existing case
type CaseDeleteResponse struct{}

// CaseListRequest is the input-object for
// listing cases for a specified user
type CaseListRequest struct {
	// UserID of the user to list
	// cases for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	UserID string
}

// CaseListResponse is the output-object for
// listing cases for a specified user
type CaseListResponse struct {
	Cases []Case
}

// CaseUploadRequest is the input-object for
// uploading an evidence to the case
type CaseUploadRequest struct {
	// ID of the case to upload
	ID string

	// Name of the item to upload
	Name string
}

// File holds information
// about an uploaded file
type File struct {
	Base

	// Name of the file
	//
	// example: "text-file.txt"
	Name string

	// Mime is the mime-type of the file
	//
	// example: "@file/plain"
	Mime string

	// Description of the file
	//
	// example: "This file contains evidence"
	Description string

	// Path to where the file is stored
	//
	// example: "/filestore/text-file.txt"
	Path string

	// Size of the file in bytes
	//
	// example: 450060
	Size int

	// Processed is if the file has been
	// processed or not
	//
	// example: false
	Processed bool
}

// FileNewRequest is the input-object
// for creating a new file
type FileNewRequest struct {
	// CaseID of the case to upload the file
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Name of the file
	//
	// example: "text-file.txt"
	Name string

	// Description of the file
	//
	// example: "This file contains evidence"
	Description string

	// Mime is the mime-type of the file
	//
	// example: "@file/plain"
	Mime string

	// Data of the file (base64 encoded)
	//
	// example: "iVBORw0KGgoAAAANSUhEUgAAA1IAAAEeCA......."
	Data string
}

// FileNewResponse is the output-object
// for creating a new file
type FileNewResponse struct {
	New File
}

// FileUpdateRequest is the input-object
// for updating a files information
type FileUpdateRequest struct {
	// ID of the file to update
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case where the file
	// to update belongs
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Description of the file
	//
	// example: "This file contains evidence"
	Description string
}

// FileUpdateResponse is the output-object
// for updating a files information
type FileUpdateResponse struct {
	Updated File
}

// FileDeleteRequest is the input-object
// for deleting a file
type FileDeleteRequest struct {
	// ID of the file to delete
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case where the file
	// to delete belongs
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// FileDeleteResponse is the output-object
// for deleting a file
type FileDeleteResponse struct{}

// Process holds information about
// a job that processes data to app
type Process struct {
	Base

	// TODO : Add more fields
}

// ProcessStartRequest is the input-object
// for starting a processing-job
type ProcessStartRequest struct {
	// CaseID of the case to start
	// the processing for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// FileIDs of the files to process
	//
	// example: ["7a1713b0249d477d92f5e10124a59861", "7a1713b0249d477d92f5e10124a59861"]
	FileIDs []string
}

// ProcessStartResponse is the output-object
// for starting a processing-job
type ProcessStartResponse struct {
	Started Process
}

// ProcessJobsRequest is the input-object
// for getting all processing-jobs for a case
type ProcessJobsRequest struct {
	// CaseID of the case to get
	// the processing-jobs for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// ProcessJobsResponse is the output-object
// for getting all processing-jobs for a case
type ProcessJobsResponse struct {
	Processes []Process
}

// ProcessAbortRequest is the input-object
// for aborting a processing-job
type ProcessAbortRequest struct {
	// ID of the processing-job to abort
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case the processing-job belongs to
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// ProcessAbortResponse is the output-object
// for aborting a processing-job
type ProcessAbortResponse struct {
	Aborted Process
}

// ProcessPauseRequest is the input-object
// for pausing a processing-job
type ProcessPauseRequest struct {
	// ID of the processing-job to pause
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case the processing-job belongs to
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// ProcessPauseResponse is the output-object
// for pausing a processing-job
type ProcessPauseResponse struct {
	Paused Process
}

// Base model for the database
type Base struct {
	// ID is the identifier for the object
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CreatedAt - when
	// the object was created
	//
	// example: 1257894000
	CreatedAt int64

	// UpdatedAt - when
	// the object was updated
	//
	// example: 0
	UpdatedAt int64

	// DeletedAt - when
	// the object was deleted
	//
	// example: 0
	DeletedAt int64
}
