package services

import (
	"context"
	"errors"
	"log"
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
	log.Println("getting user from context")
	currentUser := utils.GetUser(ctx)

	log.Println("creating new case")
	caze := api.Case{
		CreatorID:     currentUser.UID,
		Name:          r.Name,
		Description:   r.Description,
		FromDate:      r.FromDate,
		ToDate:        r.ToDate,
		Investigators: []string{currentUser.Email},
	}

	if err := s.db.CreateCase(ctx, &caze); err != nil {
		return nil, err
	}

	log.Println("returning new case")
	return &api.CaseNewResponse{New: caze}, nil
}

// Get returns the requested case
func (s *CaseService) Get(ctx context.Context, r api.CaseGetRequest) (*api.CaseGetResponse, error) {
	caze, err := s.db.GetCase(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	return &api.CaseGetResponse{Case: *caze}, nil
}

// Update updates the specified case
func (s *CaseService) Update(ctx context.Context, r api.CaseUpdateRequest) (*api.CaseUpdateResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// Delete deletes the specified case
func (s *CaseService) Delete(ctx context.Context, r api.CaseDeleteRequest) (*api.CaseDeleteResponse, error) {
	return nil, errors.New("Not implemented yet")
}

// List the cases for a specified user
func (s *CaseService) List(ctx context.Context, r api.CaseListRequest) (*api.CaseListResponse, error) {
	user, err := s.auth.GetUserByID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	cases, err := s.db.GetCasesByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	return &api.CaseListResponse{Cases: cases}, nil
}

// Authenticate is a middleware
// in the http-handler
//
// NOTE : Only for Go-servers
func (s *CaseService) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	log.Println("getting user by token")
	usr, err := s.auth.GetUserByToken(ctx, utils.GetToken(r))
	if err != nil {
		return nil, err
	}

	log.Println("setting user to context")
	return utils.SetUser(ctx, api.User{
		DisplayName: usr.DisplayName,
		Email:       usr.Email,
		PhoneNumber: usr.PhoneNumber,
		PhotoURL:    usr.PhotoURL,
		ProviderID:  usr.ProviderID,
		UID:         usr.UID,
	}), nil
}
