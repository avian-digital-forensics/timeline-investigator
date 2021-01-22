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

// KeywordsAdd adds keywords to a person
func (s *PersonService) KeywordsAdd(ctx context.Context, r api.KeywordsAddRequest) (*api.KeywordsAddResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Get the person to add the keyword to
	person, err := s.db.GetPersonByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be added to the Person
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Create a map of the keywords that were found from the db
	// and add the Person ID to each keyword
	var keywordFound = make(map[string]bool)
	for _, keyword := range keywords {
		keywordFound[keyword.Name] = true
		keyword.PersonIDs = append(keyword.PersonIDs, person.ID)
	}

	// Add the keywords from the request to the Person
	// and append the keywords that didn't already exist
	// to the keyword-slice
	for _, keyword := range r.Keywords {
		if !keywordFound[keyword] {
			keywords = append(keywords, api.Keyword{Name: keyword, PersonIDs: []string{person.ID}})
		}
		person.Keywords = append(person.Keywords, keyword)
	}

	// Save the keywords with the person ID
	// TODO / FIXME: Use bulk-indexer instead
	for _, keyword := range keywords {
		if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
			return nil, api.Error(err, api.ErrCannotPerformOperation)
		}
	}

	// Update the person with the added keywords
	if err := s.db.UpdatePerson(ctx, r.CaseID, person); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsAddResponse{OK: true}, nil
}

// KeywordsRemove removes keywords from a Person
func (s *PersonService) KeywordsRemove(ctx context.Context, r api.KeywordsRemoveRequest) (*api.KeywordsRemoveResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.caseService.isAllowed(ctx, r.CaseID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	// Create a map of the keywords to remove
	var keywordToRemove = make(map[string]bool)
	for _, keyword := range r.Keywords {
		keywordToRemove[keyword] = true
	}

	// Get the Person to remove the keywords from
	person, err := s.db.GetPersonByID(ctx, r.CaseID, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrNotFound)
	}

	// Get the keywords that should be removed from the person
	keywords, err := s.db.GetKeywordsByIDs(ctx, r.CaseID, r.Keywords)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// Remove the keywords from the person
	for i, keyword := range person.Keywords {
		if keywordToRemove[keyword] {
			person.Keywords = append(person.Keywords[:i], person.Keywords[i+1:]...)
		}
	}

	// Remove the PersonID from the keywords
	for ki, keyword := range keywords {
		if keywordToRemove[keyword.Name] {
			for ei, id := range keyword.PersonIDs {
				if id == person.ID {
					keywords[ki].PersonIDs = append(
						keyword.PersonIDs[:ei],
						keyword.PersonIDs[ei+1:]...,
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
			if err := s.db.DeleteKeyword(ctx, r.CaseID, keyword.Name); err != nil {
				return nil, api.Error(err, api.ErrCannotPerformOperation)
			}
		} else if !toDelete {
			if err := s.db.SaveKeyword(ctx, r.CaseID, &keyword); err != nil {
				return nil, api.Error(err, api.ErrCannotPerformOperation)
			}
		}
	}

	// Update the person without the removed keyword
	if err := s.db.UpdatePerson(ctx, r.CaseID, person); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.KeywordsRemoveResponse{}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *PersonService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return s.caseService.Authenticate(ctx, r)
}
