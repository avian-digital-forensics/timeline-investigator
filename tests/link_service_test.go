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

func TestLinkService(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(testURL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	testUser, err := newTestUser(ctx, client.NewTestService(httpClient, ""))
	is.NoErr(err)
	defer testUser.delete(ctx)

	testCase, err := testUser.newTestCase(ctx, client.NewCaseService(httpClient, testUser.Token))
	is.NoErr(err)

	eventService := client.NewEventService(httpClient, testUser.Token)

	// Create a new event request
	createRequest := client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  3,
		Description: uuid.New().String(),
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	}

	// Create three new events
	event1, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)
	event2, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)
	event3, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)

	// init the link-service
	linkService := client.NewLinkService(httpClient, testUser.Token)

	// create a new bidirectional link request
	linkRequest := client.LinkEventCreateRequest{
		CaseID:        testCase.ID,
		FromID:        event1.Created.ID,
		EventIDs:      []string{event2.Created.ID, event3.Created.ID},
		Bidirectional: true,
	}

	// Create the link for the events and make sure the data is valid
	link, err := linkService.CreateEvent(ctx, linkRequest)
	is.NoErr(err)
	is.Equal(link.Linked.From, event1.Created)
	is.Equal(link.Linked.Events[0], event2.Created)
	is.Equal(link.Linked.Events[1], event3.Created)

	// Get the link for the first event (the "from-event" we created the link with) and make sure the data is valid
	gotten1, err := linkService.GetEvent(ctx, client.LinkEventGetRequest{CaseID: testCase.ID, EventID: event1.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten1.Link.Events), 2)
	is.Equal(gotten1.Link.From, link.Linked.From)
	is.Equal(gotten1.Link.Events[0], link.Linked.Events[0])
	is.Equal(gotten1.Link.Events[1], link.Linked.Events[1])

	// Since the link was bidrectional, get the links for the other events as well and make sure the data is valid
	gotten2, err := linkService.GetEvent(ctx, client.LinkEventGetRequest{CaseID: testCase.ID, EventID: event2.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten2.Link.Events), 1)
	is.Equal(gotten2.Link.From, event2.Created)
	is.Equal(gotten2.Link.Events[0], event1.Created)

	gotten3, err := linkService.GetEvent(ctx, client.LinkEventGetRequest{CaseID: testCase.ID, EventID: event3.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten3.Link.Events), 1)
	is.Equal(gotten3.Link.From, event3.Created)
	is.Equal(gotten3.Link.Events[0], event1.Created)

	// Create a new event
	event4, err := eventService.Create(ctx, createRequest)
	is.NoErr(err)

	// Link the new event (event4) and remove event2 from the link
	updateRequest := client.LinkEventUpdateRequest{
		EventID:        event1.Created.ID,
		CaseID:         testCase.ID,
		EventAddIDs:    []string{event4.Created.ID},
		EventRemoveIDs: []string{event2.Created.ID},
	}

	update, err := linkService.UpdateEvent(ctx, updateRequest)
	is.NoErr(err)
	is.Equal(len(update.Updated.Events), 2)
	is.Equal(update.Updated.From, event1.Created)
	is.Equal(update.Updated.Events[1], event3.Created)
	is.Equal(update.Updated.Events[0], event4.Created)

	// Delete the link for the event
	_, err = linkService.DeleteEvent(ctx, client.LinkEventDeleteRequest{CaseID: testCase.ID, EventID: event1.Created.ID})
	is.NoErr(err)
}
