package tests

import (
	"context"
	"os"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/google/uuid"
)

var (
	testURL    = os.Getenv("TEST_URL")
	testSecret = os.Getenv("TEST_SECRET")
)

// testUser is a user created
// with the test-service during
// development to use for testing
// the system of the Timeline-Investigator
type testUser struct {
	ID          string
	Token       string
	caseIDs     []string
	testService *client.TestService
	caseService *client.CaseService
}

// newTestUser creates a test-user
func newTestUser(ctx context.Context, testService *client.TestService) (*testUser, error) {
	var uid = uuid.New().String()
	request := client.TestCreateUserRequest{
		Name:     uid,
		ID:       uid,
		Email:    uid + "@test.com",
		Password: uuid.New().String(),
		Secret:   testSecret,
	}
	created, err := testService.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return &testUser{ID: uid, Token: created.Token, testService: testService}, nil
}

func (u *testUser) delete(ctx context.Context) error {
	/*
		for _, caseID := range u.caseIDs {
			deleteRequest := client.CaseDeleteRequest{ID: caseID}
			if _, err := u.caseService.Delete(ctx, deleteRequest); err != nil {
				return err
			}
		}
	*/

	request := client.TestDeleteUserRequest{ID: u.ID, Secret: testSecret}
	if _, err := u.testService.DeleteUser(ctx, request); err != nil {
		return err
	}
	return nil
}

func (u *testUser) newTestCase(ctx context.Context, caseService *client.CaseService) (*client.Case, error) {
	resp, err := caseService.New(ctx, client.CaseNewRequest{Name: uuid.New().String()})
	if err != nil {
		return nil, err
	}
	u.caseIDs = append(u.caseIDs, resp.New.ID)
	u.caseService = caseService
	return &resp.New, nil
}
