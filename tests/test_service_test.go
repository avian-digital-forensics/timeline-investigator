package tests

import (
	"context"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestTestService(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(testURL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	// Create the test-service (doesn't need a token)
	testService := client.NewTestService(httpClient, "")

	// Create a new user without the test-secret
	var uid = uuid.New().String()
	createRequest := client.TestCreateUserRequest{
		Name:     uid,
		ID:       uid,
		Email:    uid + "@test.com",
		Password: uuid.New().String(),
	}
	_, err := testService.CreateUser(ctx, createRequest)
	is.Equal(err.Error(), api.ErrNotAllowed.Error())

	// Create the user with the test-secret
	createRequest.Secret = testSecret
	_, err = testService.CreateUser(ctx, createRequest)
	is.NoErr(err)

	// Delete the new user without the test-secret
	deleteRequest := client.TestDeleteUserRequest{ID: uid}
	_, err = testService.DeleteUser(ctx, deleteRequest)
	is.Equal(err.Error(), api.ErrNotAllowed.Error())

	// Delete the new user with the test-secret
	deleteRequest.Secret = testSecret
	_, err = testService.DeleteUser(ctx, deleteRequest)
	is.NoErr(err)
}
