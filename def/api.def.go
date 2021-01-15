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

// EntityService is the API to handle entities
type EntityService interface {
	// Create creates a new entity
	Create(EntityCreateRequest) EntityCreateResponse

	// Update updates an existing entity
	Update(EntityUpdateRequest) EntityUpdateResponse

	// Delete deletes an existing entity
	Delete(EntityDeleteRequest) EntityDeleteResponse

	// Get the specified entity
	Get(EntityGetRequest) EntityGetResponse

	// List all entities
	List(EntityListRequest) EntityListResponse

	// Types returns the existing entity-types
	Types(EntityTypesRequest) EntityTypesResponse

	// Authenticate is a middleware
	// in the http-handler
	//
	// NOTE : Only for Go-servers
	Authenticate(*http.Request) context.Context
}

// EventService is the API to handle events
type EventService interface {
	// Create creates a new event
	Create(EventCreateRequest) EventCreateResponse

	// Update updates an existing event
	Update(EventUpdateRequest) EventUpdateResponse

	// Delete deletes an existing event
	Delete(EventDeleteRequest) EventDeleteResponse

	// Get the specified event
	Get(EventGetRequest) EventGetResponse

	// List all events
	List(EventListRequest) EventListResponse

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

	// Open opens a file
	Open(FileOpenRequest) FileOpenResponse

	// Process processes a file
	Process(FileProcessRequest) FileProcessResponse

	// Processed gets information for a processed file
	Processed(FileProcessedRequest) FileProcessedResponse

	// Processes gets information for all proccesed
	// files in the specified case
	Processes(FileProcessesRequest) FileProcessesResponse

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

// LinkService is a API for creating links
// between objects
type LinkService interface {
	// CreateEvent creates a link for an event
	// with multiple objects
	CreateEvent(LinkEventCreateRequest) LinkEventCreateResponse

	// GetEvent gets an event with its links
	GetEvent(LinkEventGetRequest) LinkEventGetResponse

	// DeleteEvent deletes all links to the specified event
	DeleteEvent(LinkEventDeleteRequest) LinkEventDeleteResponse

	// UpdateEvent updates links for the specified event
	UpdateEvent(LinkEventUpdateRequest) LinkEventUpdateResponse

	// Authenticate is a middleware
	// in the http-handler
	//
	// NOTE : Only for Go-servers
	Authenticate(*http.Request) context.Context
}

// PersonService is the API to handle entities
type PersonService interface {
	// Create creates a new person
	Create(PersonCreateRequest) PersonCreateResponse

	// Update updates an existing person
	Update(PersonUpdateRequest) PersonUpdateResponse

	// Delete deletes an existing person
	Delete(PersonDeleteRequest) PersonDeleteResponse

	// Get the specified person
	Get(PersonGetRequest) PersonGetResponse

	// List all entities for a case
	List(PersonListRequest) PersonListResponse

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

// TestService is used for testing-purposes
type TestService interface {
	// CreateUser creates a test-user in Firebase
	CreateUser(TestCreateUserRequest) TestCreateUserResponse

	// DeleteUser deletes a test-user in Firebase
	DeleteUser(TestDeleteUserRequest) TestDeleteUserResponse
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

// Entity is an object that can be
// of different types. For example,
// organization or location
type Entity struct {
	Base

	// Title of the entity
	//
	// example: "Avian APS"
	Title string

	// PhotoURL of the entity.
	// NOTE: can currently be any string,
	// but in the future have it be uploaded
	// and served by the file-service with some security
	//
	// example: "api.google.com/logo.png"
	PhotoURL string

	// Type of the entity
	//
	// example: "organization"
	Type string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// EntityCreateRequest is the input-object
// for creating an entity
type EntityCreateRequest struct {
	// CaseID of the case to create
	// the new entity to
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Title of the entity
	//
	// example: "Avian APS"
	Title string

	// PhotoURL of the entity.
	// NOTE: can currently be any string,
	// but in the future have it be uploaded
	// and served by the file-service with some security
	//
	// example: "api.google.com/logo.png"
	PhotoURL string

	// Type of the entity
	//
	// example: "organization"
	Type string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// EntityCreateResponse is the output-object
// for creating an entity
type EntityCreateResponse struct {
	Created Entity
}

// EntityUpdateRequest is the input-object
// for updating an existing entity
type EntityUpdateRequest struct {
	// ID of the entity to update
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to update
	// the existing entity to
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Title of the entity
	//
	// example: "Avian APS"
	Title string

	// PhotoURL of the entity.
	//
	// NOTE: can currently be any string,
	// but in the future have it be uploaded
	// and served by the file-service with some security
	//
	// example: "api.google.com/logo.png"
	PhotoURL string

	// Type of the entity
	//
	// example: "organization"
	Type string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// EntityUpdateResponse is the output-object
// for updating an existing entity
type EntityUpdateResponse struct {
	Updated Entity
}

// EntityGetRequest is the input-object
// for getting an existing entity
type EntityGetRequest struct {
	// ID of the entity to get
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to get the entity for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EntityGetResponse is the output-object
// for getting an existing entity
type EntityGetResponse struct {
	Entity Entity
}

// EntityDeleteRequest is the input-object
// for deleting an existing entity
type EntityDeleteRequest struct {
	// ID of the entity to delete
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to delete
	// the new entity to
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EntityDeleteResponse is the output-object
// for updating an existing entity
type EntityDeleteResponse struct{}

// EntityListRequest is the input-object
// for deleting an existing entity
type EntityListRequest struct {
	// CaseID of the case to list the entities for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EntityListResponse is the output-object
// for updating an existing entity
type EntityListResponse struct {
	Entities []Entity
}

// EntityTypesRequest is the input-object
// for getting all entity-types
type EntityTypesRequest struct{}

// EntityTypesResponse is the output-object
// for getting all entity-types
type EntityTypesResponse struct {
	// EntityTypes are the existing
	// entity-types in the system
	//
	// example: ["organization", "location"]
	EntityTypes []string
}

// Event is an important happening
// that needs investigation.
type Event struct {
	Base

	// Set the importance of the event,
	// defined by a number between 1 - 5.
	//
	// example: 3
	Importance int

	// Desription of the event.
	//
	// example: "This needs investigation."
	Description string

	// FromDate is the unix-timestamp of when
	// the event started
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-timestamp of when
	// the event finished
	//
	// example: 1257894000
	ToDate int64
}

// EventCreateRequest is the input-object
// for creating an event
type EventCreateRequest struct {
	// CaseID of the case to create
	// the event for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Set the importance of the event,
	// defined by a number between 1 - 5.
	//
	// example: 3
	Importance int

	// Desription of the event.
	//
	// example: "This needs investigation."
	Description string

	// FromDate is the unix-timestamp of when
	// the event started
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-timestamp of when
	// the event finished
	//
	// example: 1257894000
	ToDate int64
}

// EventCreateResponse is the output-object
// for creating an event
type EventCreateResponse struct {
	Created Event
}

// EventUpdateRequest is the input-object
// for updating an existing event
type EventUpdateRequest struct {
	// ID of the event to update
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the event
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// Set the importance of the event,
	// defined by a number between 1 - 5.
	//
	// example: 3
	Importance int

	// Desription of the event.
	//
	// example: "This needs investigation."
	Description string

	// FromDate is the unix-timestamp of when
	// the event started
	//
	// example: 1100127600
	FromDate int64

	// ToDate is the unix-timestamp of when
	// the event finished
	//
	// example: 1257894000
	ToDate int64
}

// EventUpdateResponse is the output-object
// for updating an existing event
type EventUpdateResponse struct {
	Updated Event
}

// EventDeleteRequest is the input-object
// for deleting an existing event
type EventDeleteRequest struct {
	// ID of the event to Delete
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the event
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EventDeleteResponse is the output-object
// for deleting an existing event
type EventDeleteResponse struct{}

// EventGetRequest is the input-object
// for getting an existing event
type EventGetRequest struct {
	// ID of the event to get
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the event
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EventGetResponse is the output-object
// for deleting an existing event
type EventGetResponse struct {
	Event Event
}

// EventListRequest is the input-object
// for listing all existing events for a case
type EventListRequest struct {
	// CaseID to list the events for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// EventListResponse is the output-object
// for listing all existing events for a case
type EventListResponse struct {
	Events []Event
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

	// ProcessedAt is the unix-timestamp
	// for when (if) the item was processed
	//
	// example: 1257894000
	ProcessedAt int64
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
	// (decided by frontend)
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

// FileProcessRequest is the input-object
// for processing a file in a case
type FileProcessRequest struct {
	// ID of the file to process
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to process the file in
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// FileProcessResponse is the output-object
// for processing a file in a case
type FileProcessResponse struct {
	Processed File
}

// FileProcessedRequest is the input-object
// for getting a processed file in a case
type FileProcessedRequest struct {
	// ID of the processed file
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to the processed file
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// FileProcessedResponse is the output-object
// for get a processed file in a case
type FileProcessedResponse struct {
	ID        string
	Processed interface{}
}

// FileProcessesRequest is the input-object
// for getting a Processes file in a case
type FileProcessesRequest struct {
	// CaseID of the case to the get all the processes
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// FileProcessesResponse is the output-object
// for get a Processes file in a case
type FileProcessesResponse struct {
	Processes interface{}
}

// FileOpenRequest is the input-object
// for opening a file in a case
type FileOpenRequest struct {
	// ID of the file to open
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case to open the file
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// FileOpenResponse is the output-object
// for opening a file in a case
type FileOpenResponse struct {
	// Data contains the b64-encoded
	// data for the file
	//
	// example: "c2FtcGxlCmRhdGEKMQ=="
	Data string
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

// LinkEvent is a link for an event between different objects
type LinkEvent struct {
	Base

	// From which event has been linked
	From Event

	// Events that has been linked
	Events []Event
}

// LinkEventCreateRequest is the input-object
// for linking objects with an event
type LinkEventCreateRequest struct {
	// CaseID for the event
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// FromID is the ID of the event to hold the link
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	FromID string

	// EventIDs of the events to be linked
	//
	// example: ["7a1713b0249d477d92f5e10124a59861"]
	EventIDs []string

	// Bidirectional means that he link also should be
	// created for the "ToID"
	//
	// example: true
	Bidirectional bool
}

// LinkEventCreateResponse is the output-object
// for linking objects with an event
type LinkEventCreateResponse struct {
	Linked LinkEvent
}

// LinkEventUpdateRequest is the input-object
// for updating linked objects with an event
type LinkEventUpdateRequest struct {
	// EventID is the ID of the event to hold the link
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	EventID string

	// CaseID for the event
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// EventAddIDs of the events to be linked
	//
	// example: ["7a1713b0249d477d92f5e10124a59861"]
	EventAddIDs []string

	// EventRemoveIDs of the events to be removed
	//
	// example: ["7a1713b0249d477d92f5e10124a59861"]
	EventRemoveIDs []string
}

// LinkEventUpdateResponse is the output-object
// for linking objects with an event
type LinkEventUpdateResponse struct {
	Updated LinkEvent
}

// LinkEventGetRequest is the input-object
// for getting a linked Event
type LinkEventGetRequest struct {
	// CaseID of the case where the event
	// belongs
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// EventID of the Event to get
	// all links for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	EventID string
}

// LinkEventGetResponse is the output-object
// for getting a linked Event
type LinkEventGetResponse struct {
	Link LinkEvent
}

// LinkEventDeleteRequest is the input-object
// for removing a linked event
type LinkEventDeleteRequest struct {
	// CaseID of the case where the linked event
	// belongs
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// EventID of the Event to delete the link for
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	EventID string
}

// LinkEventDeleteResponse is the output-object
// for removing a linked event
type LinkEventDeleteResponse struct{}

// Person is a human related to a case
type Person struct {
	Base

	// FirstName(s) of the person
	//
	// example: "Simon"
	FirstName string

	// LastName(s) of the person
	//
	// example: "Jansson"
	LastName string

	// EmailAddress of the person
	//
	// example: "sja@avian.dk"
	EmailAddress string

	// PostalAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	PostalAddress string

	// WorkAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	WorkAddress string

	// TelephoneNo of the person
	//
	// example: "+46765550125"
	TelephoneNo string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// PersonCreateRequest is the input-object
// for creating a person
type PersonCreateRequest struct {
	// CaseID of the case where
	// the person should be created
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// FirstName(s) of the person
	//
	// example: "Simon"
	FirstName string

	// LastName(s) of the person
	//
	// example: "Jansson"
	LastName string

	// EmailAddress of the person
	//
	// example: "sja@avian.dk"
	EmailAddress string

	// PostalAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	PostalAddress string

	// WorkAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	WorkAddress string

	// TelephoneNo of the person
	//
	// example: "+46765550125"
	TelephoneNo string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// PersonCreateResponse is the output-object
// for creating a person
type PersonCreateResponse struct {
	Created Person
}

// PersonUpdateRequest is the input-object
// for updating an existing person
type PersonUpdateRequest struct {
	// ID of the person to update
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case where
	// the person should be updated
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string

	// FirstName(s) of the person
	//
	// example: "Simon"
	FirstName string

	// LastName(s) of the person
	//
	// example: "Jansson"
	LastName string

	// EmailAddress of the person
	//
	// example: "sja@avian.dk"
	EmailAddress string

	// PostalAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	PostalAddress string

	// WorkAddress of the person
	//
	// example: "Applebys Plads 7, 1411 Copenhagen, Denmark"
	WorkAddress string

	// TelephoneNo of the person
	//
	// example: "+46765550125"
	TelephoneNo string

	// Custom is a free form with key-value pairs
	// specified by the user.
	Custom map[string]interface{}
}

// PersonUpdateResponse is the output-object
// for updating an existing person
type PersonUpdateResponse struct {
	Updated Person
}

// PersonGetRequest is the input-object
// for getting an existing person
type PersonGetRequest struct {
	// ID of the person to get
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case where
	// the person should be gotten from
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// PersonGetResponse is the output-object
// for getting an existing person
type PersonGetResponse struct {
	Person Person
}

// PersonDeleteRequest is the input-object
// for deleting an existing person
type PersonDeleteRequest struct {
	// ID of the person to delete
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	ID string

	// CaseID of the case where
	// the person should be deleted
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// PersonDeleteResponse is the output-object
// for deleting an existing person
type PersonDeleteResponse struct{}

// PersonListRequest is the input-object
// for listing all persons for a case
type PersonListRequest struct {
	// CaseID of the case to listen all persons
	//
	// example: "7a1713b0249d477d92f5e10124a59861"
	CaseID string
}

// PersonListResponse is the output-object
// for listing all persons for a case
type PersonListResponse struct {
	Persons []Person
}

// Process holds information about
// a job that processes data to app
type Process struct {
	Base

	// Files for the process
	Files []string
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

// TestCreateUserRequest is the input-object
// for creating a test-user
type TestCreateUserRequest struct {
	// Name of the user to create
	//
	// example: "Simon"
	Name string

	// ID of the user to create
	//
	// example: "aaef42k4t2"
	ID string

	// Email of the user to create
	//
	// example: "sja@avian.dk"
	Email string

	// Password for the new user
	//
	// example: "supersecret"
	Password string

	// Secret for using the test-service
	//
	// example: "supersecret"
	Secret string
}

// TestCreateUserResponse is the output-object
// for creating a test-user
type TestCreateUserResponse struct {
	// Token for the created user
	//
	// example: "er324235tt...."
	Token string
}

// TestDeleteUserRequest is the input-object
// for deleting a test-user
type TestDeleteUserRequest struct {
	// ID of the user to delete
	//
	// example: "aaef42k4t2"
	ID string

	// Secret for using the test-service
	//
	// example: "supersecret"
	Secret string
}

// TestDeleteUserResponse is the output-object
// for deleting a test-user
type TestDeleteUserResponse struct{}

// User holds information for a user
// in the timeline-investigator
type User struct {
	DisplayName string
	Email       string
	PhoneNumber string
	PhotoURL    string
	ProviderID  string
	UID         string
}
