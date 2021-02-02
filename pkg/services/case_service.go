package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/authentication"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
)

// CaseService handles cases
type CaseService struct {
	db   datastore.Service
	auth authentication.Service
}

// NewCaseService creates a new case-service
func NewCaseService(db datastore.Service, auth authentication.Service) *CaseService {
	return &CaseService{db: db, auth: auth}
}

// New creates a new case
func (s *CaseService) New(ctx context.Context, r api.CaseNewRequest) (*api.CaseNewResponse, error) {
	if r.FromDate > r.ToDate {
		return nil, api.ErrInvalidDates
	}

	currentUser := utils.GetUser(ctx)

	caze := api.Case{
		CreatorID:     currentUser.UID,
		Name:          r.Name,
		Description:   r.Description,
		FromDate:      r.FromDate,
		ToDate:        r.ToDate,
		Investigators: []string{currentUser.Email},
	}

	if err := s.db.CreateCase(ctx, &caze); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.CaseNewResponse{New: caze}, nil
}

// Get returns the requested case
func (s *CaseService) Get(ctx context.Context, r api.CaseGetRequest) (*api.CaseGetResponse, error) {
	caze, err := s.db.GetCase(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("case - %v", api.ErrNotFound)
	}

	if !isAllowed(caze, utils.GetUser(ctx).Email) {
		return nil, api.ErrNotAllowed
	}

	return &api.CaseGetResponse{Case: *caze}, nil
}

// Update updates the specified case
func (s *CaseService) Update(ctx context.Context, r api.CaseUpdateRequest) (*api.CaseUpdateResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Delete deletes the specified case
func (s *CaseService) Delete(ctx context.Context, r api.CaseDeleteRequest) (*api.CaseDeleteResponse, error) {
	return nil, s.db.DeleteCase(ctx, r.ID)
}

// List the cases for a specified user
func (s *CaseService) List(ctx context.Context, r api.CaseListRequest) (*api.CaseListResponse, error) {
	currentUser := utils.GetUser(ctx)
	if currentUser.UID != r.UserID {
		// TODO : Check if user is system-admin (?)
		/*
			user, err := s.auth.GetUserByID(ctx, r.UserID)
			if err != nil {
				return nil, err
			}
		*/
		return nil, api.ErrNotAllowed
	}

	cases, err := s.db.GetCasesByEmail(ctx, currentUser.Email)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.CaseListResponse{Cases: cases}, nil
}

// Keywords lists all the keywords for the case
func (s *CaseService) Keywords(ctx context.Context, r api.CaseKeywordsRequest) (*api.CaseKeywordsResponse, error) {
	currentUser := utils.GetUser(ctx)
	if ok, err := s.isAllowed(ctx, r.ID, currentUser.Email); !ok {
		if err != nil {
			return nil, api.Error(err, api.ErrNotAllowed)
		}
		return nil, api.ErrNotAllowed
	}

	keywords, err := s.db.GetKeywords(ctx, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.CaseKeywordsResponse{Keywords: keywords}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *CaseService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	usr, err := s.auth.GetUserByToken(ctx, utils.GetToken(r))
	if err != nil {
		return nil, api.Error(err, api.ErrNotAllowed)
	}

	return utils.SetUser(ctx, api.User{
		DisplayName: usr.DisplayName,
		Email:       usr.Email,
		PhoneNumber: usr.PhoneNumber,
		PhotoURL:    usr.PhotoURL,
		ProviderID:  usr.ProviderID,
		UID:         usr.UID,
	}), nil
}

func isAllowed(caze *api.Case, email string) bool {
	for _, investigator := range caze.Investigators {
		if investigator == email {
			return true
		}
	}
	return false
}

func (s *CaseService) isAllowed(ctx context.Context, caseID, email string) (bool, error) {
	caze, err := s.db.GetCase(ctx, caseID)
	if err != nil {
		return false, err
	}
	return isAllowed(caze, email), nil
}
