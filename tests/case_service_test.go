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
		Name:        "Simon",
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
	gotten, err := service.Get(ctx, client.CaseGetRequest{ID: caze.New.ID})
	is.NoErr(err)
	is.Equal(gotten.Case.ID, caze.New.ID)
	is.Equal(gotten.Case.CreatedAt, caze.New.CreatedAt)
	is.Equal(gotten.Case.UpdatedAt, caze.New.UpdatedAt)
	is.Equal(gotten.Case.DeletedAt, caze.New.DeletedAt)
	is.Equal(gotten.Case.Name, caze.New.Name)
	is.Equal(gotten.Case.Description, caze.New.Description)
	is.Equal(gotten.Case.FromDate, caze.New.FromDate)
	is.Equal(gotten.Case.ToDate, caze.New.ToDate)

	// Get all cases for the test-user
	all, err := service.List(ctx, client.CaseListRequest{UserID: testUser.ID})
	is.NoErr(err)
	is.Equal(len(all.Cases), 1)
	is.Equal(all.Cases[0].ID, gotten.Case.ID)
	is.Equal(all.Cases[0].CreatedAt, gotten.Case.CreatedAt)
	is.Equal(all.Cases[0].UpdatedAt, gotten.Case.UpdatedAt)
	is.Equal(all.Cases[0].DeletedAt, gotten.Case.DeletedAt)
	is.Equal(all.Cases[0].Name, gotten.Case.Name)
	is.Equal(all.Cases[0].Description, gotten.Case.Description)
	is.Equal(all.Cases[0].FromDate, gotten.Case.FromDate)
	is.Equal(all.Cases[0].ToDate, gotten.Case.ToDate)
}
