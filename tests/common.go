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
	ID    string
	Token string
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

	return &testUser{ID: uid, Token: created.Token}, nil
}

func (u *testUser) delete(ctx context.Context, testService *client.TestService) error {
	request := client.TestDeleteUserRequest{ID: u.ID, Secret: testSecret}
	if _, err := testService.DeleteUser(ctx, request); err != nil {
		return err
	}
	return nil
}
