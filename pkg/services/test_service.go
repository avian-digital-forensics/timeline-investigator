package services

import (
	"context"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/authentication"
)

// TestService is used during development
// for performing tests against the system
type TestService struct {
	auth   authentication.Service
	secret string
}

// NewTestService creates a new TestService
func NewTestService(auth authentication.Service, secret string) *TestService {
	return &TestService{auth: auth, secret: secret}
}

// CreateUser creates a user for testing
func (s *TestService) CreateUser(ctx context.Context, r api.TestCreateUserRequest) (*api.TestCreateUserResponse, error) {
	if r.Secret != s.secret {
		return nil, api.ErrNotAllowed
	}

	user, err := s.auth.Create(ctx, r.ID, r.Email, r.Name, r.Password)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	customToken, err := s.auth.GetCustomToken(ctx, user.UID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	idToken, err := s.auth.VerifyCustomToken(ctx, customToken)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.TestCreateUserResponse{Token: idToken}, nil
}

// DeleteUser deletes a test-user
func (s *TestService) DeleteUser(ctx context.Context, r api.TestDeleteUserRequest) (*api.TestDeleteUserResponse, error) {
	if r.Secret != s.secret {
		return nil, api.ErrNotAllowed
	}

	user, err := s.auth.GetUserByID(ctx, r.ID)
	if err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	// make sure the user is a test-user before deleting
	if test, ok := user.CustomClaims["Test"]; !test.(bool) || !ok {
		return nil, api.Error(err, api.ErrNotAllowed)
	}

	if err := s.auth.Delete(ctx, r.ID); err != nil {
		return nil, api.Error(err, api.ErrCannotPerformOperation)
	}

	return &api.TestDeleteUserResponse{}, nil
}
