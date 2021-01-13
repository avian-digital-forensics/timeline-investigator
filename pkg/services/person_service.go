package services

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

// PersonService holds dependencies
// for handling the Person API
type PersonService struct {
	db          datastore.Service
	caseService *CaseService
}

// NewPersonService creates a new person-service
func NewPersonService(db datastore.Service, caseService *CaseService) *PersonService {
	return &PersonService{db: db, caseService: caseService}
}

// Create creates a new Person
func (s *PersonService) Create(ctx context.Context, r api.PersonCreateRequest) (*api.PersonCreateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	person := api.Person{
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		EmailAddress:  r.EmailAddress,
		PostalAddress: r.PostalAddress,
		WorkAddress:   r.WorkAddress,
		TelephoneNo:   r.TelephoneNo,
		Custom:        r.Custom,
	}

	if err := s.db.CreatePerson(ctx, r.CaseID, &person); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.PersonCreateResponse{Created: person}, nil
}

// Update updates an existing Person
func (s *PersonService) Update(ctx context.Context, r api.PersonUpdateRequest) (*api.PersonUpdateResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	person := api.Person{
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		EmailAddress:  r.EmailAddress,
		PostalAddress: r.PostalAddress,
		WorkAddress:   r.WorkAddress,
		TelephoneNo:   r.TelephoneNo,
		Custom:        r.Custom,
	}
	person.ID = r.ID

	if err := s.db.UpdatePerson(ctx, r.CaseID, &person); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.PersonUpdateResponse{Updated: person}, nil
}

// Delete deletes an existing Person
func (s *PersonService) Delete(ctx context.Context, r api.PersonDeleteRequest) (*api.PersonDeleteResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	if err := s.db.DeletePerson(ctx, r.CaseID, r.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.PersonDeleteResponse{}, nil
}

// Get the specified Person
func (s *PersonService) Get(ctx context.Context, r api.PersonGetRequest) (*api.PersonGetResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	person, err := s.db.GetPersonByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	return &api.PersonGetResponse{Person: *person}, nil
}

// List all entities
func (s *PersonService) List(ctx context.Context, r api.PersonListRequest) (*api.PersonListResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	persons, err := s.db.GetPersons(ctx, r.CaseID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.PersonListResponse{Persons: persons}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *PersonService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
