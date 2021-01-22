package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestEventService(t *testing.T) {
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

	// Create a new event request (this should fail)
	createRequest := client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  0,
		Description: uuid.New().String(),
		FromDate:    time.Now().AddDate(1, 0, 0).Unix(),
		ToDate:      time.Now().Unix(),
	}

	// This should fail because the from-date is greater than the to-date
	_, err = eventService.Create(ctx, createRequest)
	is.Equal(err.Error(), api.ErrInvalidDates.Error())

	// Fix the from-date so it is less than the to-date
	createRequest.FromDate = time.Now().AddDate(-2, 0, 0).Unix()
	// Try creating the event again (should also fail because of the importance is less than 1)
	_, err = eventService.Create(ctx, createRequest)
	is.Equal(err.Error(), api.ErrInvalidImportance.Error())
	// Set importance to 6 and try again (should also fail because of the importance is greater than 5)
	createRequest.Importance = 6
	_, err = eventService.Create(ctx, createRequest)
	is.Equal(err.Error(), api.ErrInvalidImportance.Error())

	// Fix importance and try again (should succeed)
	createRequest.Importance = 3
	event, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)
	is.Equal(event.Created.Importance, createRequest.Importance)
	is.Equal(event.Created.Description, createRequest.Description)
	is.Equal(event.Created.FromDate, createRequest.FromDate)
	is.Equal(event.Created.ToDate, createRequest.ToDate)

	gotten, err := eventService.Get(ctx, client.EventGetRequest{ID: event.Created.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(gotten.Event.ID, event.Created.ID)
	is.Equal(gotten.Event.CreatedAt, event.Created.CreatedAt)
	is.Equal(gotten.Event.UpdatedAt, event.Created.UpdatedAt)
	is.Equal(gotten.Event.DeletedAt, event.Created.DeletedAt)
	is.Equal(gotten.Event.Importance, event.Created.Importance)
	is.Equal(gotten.Event.Description, event.Created.Description)
	is.Equal(gotten.Event.FromDate, event.Created.FromDate)
	is.Equal(gotten.Event.ToDate, event.Created.ToDate)

	// Create a update-request
	updateRequest := client.EventUpdateRequest{
		ID:          event.Created.ID,
		CaseID:      testCase.ID,
		Importance:  0,
		Description: uuid.New().String(),
		FromDate:    time.Now().AddDate(1, 0, 0).Unix(),
		ToDate:      time.Now().Unix(),
	}

	// update the created event (should fail because of invalid date)
	_, err = eventService.Update(ctx, updateRequest)
	is.Equal(err.Error(), api.ErrInvalidDates.Error())

	// fix date and try to update it again (should fail because of importance)
	updateRequest.FromDate = time.Now().AddDate(-2, 0, 0).Unix()
	_, err = eventService.Update(ctx, updateRequest)
	is.Equal(err.Error(), api.ErrInvalidImportance.Error())

	// change the importance to 6 and try again (should also fail because the importance is invalid)
	updateRequest.Importance = 6
	_, err = eventService.Update(ctx, updateRequest)
	is.Equal(err.Error(), api.ErrInvalidImportance.Error())

	// Set a valid importance and make sure the data is valid
	updateRequest.Importance = 1
	updatedEvent, err := eventService.Update(ctx, updateRequest)
	is.NoErr(err)
	is.Equal(updatedEvent.Updated.ID, updateRequest.ID)
	is.Equal(updatedEvent.Updated.Importance, updateRequest.Importance)
	is.Equal(updatedEvent.Updated.Description, updateRequest.Description)
	is.Equal(updatedEvent.Updated.FromDate, updateRequest.FromDate)
	is.Equal(updatedEvent.Updated.ToDate, updateRequest.ToDate)

	// try updating with all other valid values for importance (should succeed)
	for _, v := range []int{2, 3, 4, 5} {
		updateRequest.Importance = v
		updatedEvent, err = eventService.Update(ctx, updateRequest)
		is.NoErr(err)
	}

	// Get the event again and make sure that it has been updated
	gotten, err = eventService.Get(ctx, client.EventGetRequest{ID: event.Created.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(gotten.Event.ID, updatedEvent.Updated.ID)
	is.Equal(gotten.Event.CreatedAt, updatedEvent.Updated.CreatedAt)
	is.Equal(gotten.Event.UpdatedAt, updatedEvent.Updated.UpdatedAt)
	is.Equal(gotten.Event.DeletedAt, updatedEvent.Updated.DeletedAt)
	is.Equal(gotten.Event.Importance, updatedEvent.Updated.Importance)
	is.Equal(gotten.Event.Description, updatedEvent.Updated.Description)
	is.Equal(gotten.Event.FromDate, updatedEvent.Updated.FromDate)
	is.Equal(gotten.Event.ToDate, updatedEvent.Updated.ToDate)

	// Get all keywords for the case where the event belongs
	keywords, err := caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 0)

	// Add two keywords to the first event
	keywordsRequest := client.KeywordsAddRequest{ID: gotten.Event.ID, CaseID: testCase.ID, Keywords: []string{"healthy", "green"}}
	_, err = eventService.KeywordsAdd(ctx, keywordsRequest)
	is.NoErr(err)

	// Get the keywords for the case again and make sure they have been added for the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 2)
	is.Equal(keywordsRequest.Keywords[0], keywords.Keywords[0])
	is.Equal(keywordsRequest.Keywords[1], keywords.Keywords[1])

	// Get the entity again and make sure the keywords has been added there as well
	gotten, err = eventService.Get(ctx, client.EventGetRequest{CaseID: testCase.ID, ID: gotten.Event.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Event.Keywords), 2)
	is.Equal(gotten.Event.Keywords, keywordsRequest.Keywords)
	// make sure the rest of the data is valid as well
	is.Equal(gotten.Event.ID, updatedEvent.Updated.ID)
	is.Equal(gotten.Event.CreatedAt, updatedEvent.Updated.CreatedAt)
	is.Equal(gotten.Event.DeletedAt, updatedEvent.Updated.DeletedAt)
	is.Equal(gotten.Event.Importance, updatedEvent.Updated.Importance)
	is.Equal(gotten.Event.Description, updatedEvent.Updated.Description)
	is.Equal(gotten.Event.FromDate, updatedEvent.Updated.FromDate)
	is.Equal(gotten.Event.ToDate, updatedEvent.Updated.ToDate)

	// List all events for the case and make sure that the created event is there
	list, err := eventService.List(ctx, client.EventListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(list.Events), 1)
	is.Equal(list.Events[0].ID, gotten.Event.ID)
	is.Equal(list.Events[0].CreatedAt, gotten.Event.CreatedAt)
	is.Equal(list.Events[0].UpdatedAt, gotten.Event.UpdatedAt)
	is.Equal(list.Events[0].DeletedAt, gotten.Event.DeletedAt)
	is.Equal(list.Events[0].Importance, gotten.Event.Importance)
	is.Equal(list.Events[0].Description, gotten.Event.Description)
	is.Equal(list.Events[0].FromDate, gotten.Event.FromDate)
	is.Equal(list.Events[0].ToDate, gotten.Event.ToDate)

	// Create a new event and make sure it has been added to the list
	event2, err := eventService.Create(ctx, client.EventCreateRequest{CaseID: testCase.ID, Importance: 3})
	is.NoErr(err)
	list, err = eventService.List(ctx, client.EventListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(list.Events), 2)

	// Add the keywords + another one, to event2 as well
	keywordsRequest.ID = event2.Created.ID
	keywordsRequest.Keywords = append(keywordsRequest.Keywords, "another one")
	_, err = eventService.KeywordsAdd(ctx, keywordsRequest)
	is.NoErr(err)

	// Make sure the keywords has been added to event2
	gotten, err = eventService.Get(ctx, client.EventGetRequest{CaseID: testCase.ID, ID: event2.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Event.Keywords), 3)
	is.Equal(gotten.Event.Keywords, keywordsRequest.Keywords)

	// We should now have 3 (unqiue) keywords in the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)
	is.Equal(keywordsRequest.Keywords[0], keywords.Keywords[0])
	is.Equal(keywordsRequest.Keywords[1], keywords.Keywords[1])
	is.Equal(keywordsRequest.Keywords[2], keywords.Keywords[2])

	// Remove the keyword first keyword from event2
	removeRequest := client.KeywordsRemoveRequest{ID: event2.Created.ID, CaseID: testCase.ID, Keywords: []string{keywords.Keywords[0]}}
	_, err = eventService.KeywordsRemove(ctx, removeRequest)
	is.NoErr(err)

	// Make sure it was removed from event2
	gotten, err = eventService.Get(ctx, client.EventGetRequest{CaseID: testCase.ID, ID: event2.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Event.Keywords), 2)
	is.Equal(gotten.Event.Keywords[0], keywordsRequest.Keywords[1])
	is.Equal(gotten.Event.Keywords[1], keywordsRequest.Keywords[2])

	// Make sure it is still available in the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)

	// Make sure it is still available in the first event
	gotten, err = eventService.Get(ctx, client.EventGetRequest{CaseID: testCase.ID, ID: event.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Event.Keywords), 2)
	is.Equal(gotten.Event.Keywords[0], keywordsRequest.Keywords[0])
	is.Equal(gotten.Event.Keywords[1], keywordsRequest.Keywords[1])

	// Remove it from the first event as well
	removeRequest.ID = event.Created.ID
	_, err = eventService.KeywordsRemove(ctx, removeRequest)
	is.NoErr(err)

	// Make sure it was removed from the case (since no event has the keyword anymore)
	keywords5, err := caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords5.Keywords), 2)
	is.Equal(keywords5.Keywords[0], keywordsRequest.Keywords[1])
	is.Equal(keywords5.Keywords[1], keywordsRequest.Keywords[2])

	// Make sure it was removed from the first event as well
	gotten, err = eventService.Get(ctx, client.EventGetRequest{CaseID: testCase.ID, ID: event.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Event.Keywords), 1)
	is.Equal(gotten.Event.Keywords[0], keywordsRequest.Keywords[1])

	// Delete all events
	for _, event := range list.Events {
		_, err = eventService.Delete(ctx, client.EventDeleteRequest{ID: event.ID, CaseID: testCase.ID})
		is.NoErr(err)
	}

	// Make sure all events were deleted
	list, err = eventService.List(ctx, client.EventListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(list.Events), 0)
}
