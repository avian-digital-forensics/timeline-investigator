package tests

import (
	"context"
	b64 "encoding/base64"
	"log"
	"testing"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/matryer/is"
)

func TestSearchService(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(testURL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	testUser, err := newTestUser(ctx, client.NewTestService(httpClient, ""))
	is.NoErr(err)
	defer testUser.delete(ctx)

	caseService := client.NewCaseService(httpClient, testUser.Token)
	eventService := client.NewEventService(httpClient, testUser.Token)
	entityService := client.NewEntityService(httpClient, testUser.Token)
	fileService := client.NewFileService(httpClient, testUser.Token)
	personService := client.NewPersonService(httpClient, testUser.Token)
	searchService := client.NewSearchService(httpClient, testUser.Token)

	testCase, err := testUser.newTestCase(ctx, caseService)
	is.NoErr(err)

	file1, err := fileService.New(ctx, client.FileNewRequest{
		CaseID:      testCase.ID,
		Name:        "d1.txt",
		Description: "test1",
		Mime:        "text/plain",
		Data:        b64.URLEncoding.EncodeToString([]byte("data-1")),
	})
	is.NoErr(err)
	file2, err := fileService.New(ctx, client.FileNewRequest{
		CaseID:      testCase.ID,
		Name:        "d2.txt",
		Description: "test2",
		Mime:        "text/plain",
		Data:        b64.URLEncoding.EncodeToString([]byte("other-2")),
	})
	is.NoErr(err)

	processed1, err := fileService.Process(ctx, client.FileProcessRequest{ID: file1.New.ID, CaseID: testCase.ID})
	is.NoErr(err)
	processed2, err := fileService.Process(ctx, client.FileProcessRequest{ID: file2.New.ID, CaseID: testCase.ID})
	is.NoErr(err)

	event1, err := eventService.Create(ctx, client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  1,
		Description: "event-1",
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	})
	is.NoErr(err)
	event2, err := eventService.Create(ctx, client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  2,
		Description: "event-2",
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	})
	is.NoErr(err)

	entity1, err := entityService.Create(ctx, client.EntityCreateRequest{
		CaseID:   testCase.ID,
		Title:    "Avian APS",
		PhotoURL: "https://randomURL.com",
		Type:     "organization",
	})
	is.NoErr(err)
	entity2, err := entityService.Create(ctx, client.EntityCreateRequest{
		CaseID:   testCase.ID,
		Title:    "Timeline Investigator",
		PhotoURL: "https://anotherURL.com",
		Type:     "location",
	})
	is.NoErr(err)

	person1, err := personService.Create(ctx, client.PersonCreateRequest{
		CaseID:        testCase.ID,
		FirstName:     "Pär Gustaf",
		LastName:      "S",
		EmailAddress:  "pg@avian.dk",
		PostalAddress: "Vegagatan 6A, 113 29, Stockholm, SE",
		WorkAddress:   "Applebys Plads 7, 1411 Copenhagen, DK",
		TelephoneNo:   "+46765550125",
	})
	is.NoErr(err)
	person2, err := personService.Create(ctx, client.PersonCreateRequest{
		CaseID:        testCase.ID,
		FirstName:     "Simon",
		LastName:      "Jansson",
		EmailAddress:  "sja@avian.dk",
		PostalAddress: "Vegagatan 6A, 113 29, Stockholm, SE",
		WorkAddress:   "Applebys Plads 7, 1411 Copenhagen, DK",
		TelephoneNo:   "+46765550125",
	})
	is.NoErr(err)

	// Wait 10 seconds for the file to be indexed
	log.Println("WAIT : 10 seconds for processed file to be indexed")
	time.Sleep(1 * time.Second)
	for _, s := range []int{9, 8, 7, 6, 5, 4, 3, 2, 1} {
		log.Println("WAIT :", s)
		time.Sleep(1 * time.Second)
	}

	log.Println(event1, event2, entity1, entity2, processed1, processed2, person1, person2)

	resp, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "Simon"})
	is.NoErr(err)
	is.Equal(len(resp.Persons), 1)
	is.Equal(resp.Persons[0].FirstName, "Simon")

	resp2, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "pg@avian.dk"})
	is.NoErr(err)
	is.Equal(len(resp2.Persons), 1)
	is.Equal(resp2.Persons[0].EmailAddress, "pg@avian.dk")

	// Search for Avian (should return two persons with avian-email and the organization Avian APS)
	resp3, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "Avian"})
	is.NoErr(err)
	is.Equal(len(resp3.Entities), 1)
	is.Equal(len(resp3.Persons), 2)
	is.Equal(len(resp3.Persons), 2)
	is.Equal(resp3.Entities[0].Title, "Avian APS")
	is.Equal(resp3.Persons[0].EmailAddress, "pg@avian.dk")
	is.Equal(resp3.Persons[1].EmailAddress, "sja@avian.dk")

	// Search for content in files
	resp4, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "data"})
	is.NoErr(err)
	is.Equal(len(resp4.Processed.([]interface{})), 1) // should get one of the processed-files (content = data-1)
	resp5, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "other"})
	is.NoErr(err)
	is.Equal(len(resp5.Processed.([]interface{})), 1) // should get one of the processed-files (content = other-2)

	// Add keywords to the first entity and two (both) events
	_, err = entityService.KeywordsAdd(ctx, client.KeywordsAddRequest{ID: entity1.Created.ID, CaseID: testCase.ID, Keywords: []string{"first"}})
	is.NoErr(err)
	_, err = eventService.KeywordsAdd(ctx, client.KeywordsAddRequest{ID: event1.Created.ID, CaseID: testCase.ID, Keywords: []string{"first"}})
	is.NoErr(err)
	_, err = eventService.KeywordsAdd(ctx, client.KeywordsAddRequest{ID: event2.Created.ID, CaseID: testCase.ID, Keywords: []string{"first"}})
	is.NoErr(err)

	// Search for the keyword
	resp6, err := searchService.SearchWithText(ctx, client.SearchTextRequest{CaseID: testCase.ID, Text: "first"})
	is.NoErr(err)
	is.Equal(len(resp6.Entities), 1)
	is.Equal(len(resp6.Events), 2)
	is.Equal(resp6.Entities[0].ID, entity1.Created.ID)
	is.Equal(resp6.Entities[0].Title, entity1.Created.Title)
	is.Equal(resp6.Entities[0].PhotoURL, entity1.Created.PhotoURL)
	is.Equal(resp6.Entities[0].Type, entity1.Created.Type)
	is.Equal(resp6.Events[0].ID, event1.Created.ID)
	is.Equal(resp6.Events[0].Importance, event1.Created.Importance)
	is.Equal(resp6.Events[0].FromDate, event1.Created.FromDate)
	is.Equal(resp6.Events[0].ToDate, event1.Created.ToDate)
	is.Equal(resp6.Events[1].ID, event2.Created.ID)
	is.Equal(resp6.Events[1].Importance, event2.Created.Importance)
	is.Equal(resp6.Events[1].FromDate, event2.Created.FromDate)
	is.Equal(resp6.Events[1].ToDate, event2.Created.ToDate)
}
