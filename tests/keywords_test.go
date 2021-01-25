package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

// TestKeywords will test the keywords for events more specifically
func TestKeywords(t *testing.T) {
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
	testCase, err := testUser.newTestCase(ctx, caseService)
	is.NoErr(err)

	eventService := client.NewEventService(httpClient, testUser.Token)

	// Create a new event request
	createRequest := client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  1,
		Description: uuid.New().String(),
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	}

	// Create three events
	create1, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)
	create2, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)
	create3, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)

	event1 := create1.Created
	event2 := create2.Created
	event3 := create3.Created

	keywords, err := caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 0)

	// Add two keywords to the first event
	keywordsRequest1 := client.KeywordsAddRequest{ID: event1.ID, CaseID: testCase.ID, Keywords: []string{"healthy", "green"}}
	_, err = eventService.KeywordsAdd(ctx, keywordsRequest1)
	is.NoErr(err)

	// Add three keywords to the second event
	keywordsRequest2 := client.KeywordsAddRequest{ID: event2.ID, CaseID: testCase.ID, Keywords: append(keywordsRequest1.Keywords, "boom")}
	_, err = eventService.KeywordsAdd(ctx, keywordsRequest2)
	is.NoErr(err)

	// Add four keywords to the third event
	keywordsRequest3 := client.KeywordsAddRequest{ID: event3.ID, CaseID: testCase.ID, Keywords: append(keywordsRequest2.Keywords, "sauce")}
	_, err = eventService.KeywordsAdd(ctx, keywordsRequest3)
	is.NoErr(err)

	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 4)
	is.Equal(keywords.Keywords[0], keywordsRequest3.Keywords[0])
	is.Equal(keywords.Keywords[1], keywordsRequest3.Keywords[1])
	is.Equal(keywords.Keywords[2], keywordsRequest3.Keywords[2])
	is.Equal(keywords.Keywords[3], keywordsRequest3.Keywords[3])

	// Remove the first keyword from the third event
	_, err = eventService.KeywordsRemove(ctx, client.KeywordsRemoveRequest{ID: event3.ID, CaseID: testCase.ID, Keywords: []string{keywords.Keywords[0]}})
	is.NoErr(err)

	// Make the first keyword it was removed from the third event
	got3, err := eventService.Get(ctx, client.EventGetRequest{ID: event3.ID, CaseID: testCase.ID})
	is.NoErr(err)
	event3 = got3.Event
	is.Equal(event3.Keywords, keywords.Keywords[1:])

	// Make sure the first keyword still exists in the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 4)

	// Delete the third event
	_, err = eventService.Delete(ctx, client.EventDeleteRequest{ID: event3.ID, CaseID: testCase.ID})
	is.NoErr(err)

	// Make sure the last keyword was deleted (since it was only used by event3)
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)
	is.Equal(keywords.Keywords[0], keywordsRequest3.Keywords[0])
	is.Equal(keywords.Keywords[1], keywordsRequest3.Keywords[1])
	is.Equal(keywords.Keywords[2], keywordsRequest3.Keywords[2])

	// Get the other events again and make sure they still have the correct keywords
	got1, err := eventService.Get(ctx, client.EventGetRequest{ID: event1.ID, CaseID: testCase.ID})
	is.NoErr(err)
	got2, err := eventService.Get(ctx, client.EventGetRequest{ID: event2.ID, CaseID: testCase.ID})
	is.NoErr(err)

	event1 = got1.Event
	event2 = got2.Event

	is.Equal(event1.Keywords, keywordsRequest1.Keywords)
	is.Equal(event2.Keywords, keywordsRequest2.Keywords)

	// Delete the first event
	_, err = eventService.Delete(ctx, client.EventDeleteRequest{ID: event1.ID, CaseID: testCase.ID})
	is.NoErr(err)

	// Make sure that now keywords were deleted since all of them are in use by the second event
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)

	// Delete the second event and make sure that all keywords are deleted
	_, err = eventService.Delete(ctx, client.EventDeleteRequest{ID: event2.ID, CaseID: testCase.ID})
	is.NoErr(err)
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 0)
}
