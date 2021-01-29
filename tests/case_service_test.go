package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"

	"github.com/matryer/is"
)

// TestCaseService tests the CaseService
func TestCaseService(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(testURL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}
	testUser, err := newTestUser(ctx, client.NewTestService(httpClient, ""))
	is.NoErr(err)
	defer testUser.delete(ctx)

	// Create a new CaseService
	service := client.NewCaseService(httpClient, testUser.Token)

	// Create a new case request
	newRequest := client.CaseNewRequest{
		Name:        "Pighvar",
		Description: "New case",
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	}

	// Create the new case (send the request)
	caze, err := service.New(ctx, newRequest)
	is.NoErr(err)
	is.Equal(caze.New.Name, newRequest.Name)
	is.Equal(caze.New.Description, newRequest.Description)
	is.Equal(caze.New.FromDate, newRequest.FromDate)
	is.Equal(caze.New.ToDate, newRequest.ToDate)

	// Get the case
	getResponse, err := service.Get(ctx, client.CaseGetRequest{ID: caze.New.ID})
	is.NoErr(err)
	is.Equal(getResponse.Case.ID, caze.New.ID)
	is.Equal(getResponse.Case.CreatedAt, caze.New.CreatedAt)
	is.Equal(getResponse.Case.UpdatedAt, caze.New.UpdatedAt)
	is.Equal(getResponse.Case.DeletedAt, caze.New.DeletedAt)
	is.Equal(getResponse.Case.Name, caze.New.Name)
	is.Equal(getResponse.Case.Description, caze.New.Description)
	is.Equal(getResponse.Case.FromDate, caze.New.FromDate)
	is.Equal(getResponse.Case.ToDate, caze.New.ToDate)

	// Get all cases for the test-user
	all, err := service.List(ctx, client.CaseListRequest{UserID: testUser.ID})
	is.NoErr(err)
	is.Equal(len(all.Cases), 1)
	is.Equal(all.Cases[0].ID, getResponse.Case.ID)
	is.Equal(all.Cases[0].CreatedAt, getResponse.Case.CreatedAt)
	is.Equal(all.Cases[0].UpdatedAt, getResponse.Case.UpdatedAt)
	is.Equal(all.Cases[0].DeletedAt, getResponse.Case.DeletedAt)
	is.Equal(all.Cases[0].Name, getResponse.Case.Name)
	is.Equal(all.Cases[0].Description, getResponse.Case.Description)
	is.Equal(all.Cases[0].FromDate, getResponse.Case.FromDate)
	is.Equal(all.Cases[0].ToDate, getResponse.Case.ToDate)

	// Delete the case.
	_, err = service.Delete(ctx, client.CaseDeleteRequest{ID: caze.New.ID})
	is.NoErr(err)
	_, getError := service.Get(ctx, client.CaseGetRequest{ID: caze.New.ID})
	is.True(getError != nil)
}
