package tests

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

	// create dependend services
	eventService := client.NewEventService(httpClient, testUser.Token)
	personService := client.NewPersonService(httpClient, testUser.Token)
	entityService := client.NewEntityService(httpClient, testUser.Token)
	fileService := client.NewFileService(httpClient, testUser.Token)

	// Make the requests for creating dependencies
	eventRequest := client.EventCreateRequest{
		CaseID:      testCase.ID,
		Importance:  3,
		Description: uuid.New().String(),
		FromDate:    time.Now().Unix(),
		ToDate:      time.Now().AddDate(1, 0, 0).Unix(),
	}

	personRequest := client.PersonCreateRequest{CaseID: testCase.ID, FirstName: uuid.New().String()}
	entityRequest := client.EntityCreateRequest{CaseID: testCase.ID, Title: uuid.New().String(), Type: "organization"}
	fileRequest := client.FileNewRequest{CaseID: testCase.ID, Name: uuid.New().String(), Data: b64.URLEncoding.EncodeToString([]byte("link-test"))}

	// Create three new events
	event1, err := eventService.Create(ctx, eventRequest)
	is.NoErr(err)
	event2, err := eventService.Create(ctx, eventRequest)
	is.NoErr(err)
	event3, err := eventService.Create(ctx, eventRequest)
	is.NoErr(err)

	// Create three new persons
	person1, err := personService.Create(ctx, personRequest)
	is.NoErr(err)
	person2, err := personService.Create(ctx, personRequest)
	is.NoErr(err)
	person3, err := personService.Create(ctx, personRequest)
	is.NoErr(err)

	// Create three new entities
	entity1, err := entityService.Create(ctx, entityRequest)
	is.NoErr(err)
	entity2, err := entityService.Create(ctx, entityRequest)
	is.NoErr(err)
	entity3, err := entityService.Create(ctx, entityRequest)
	is.NoErr(err)

	// Create three new files
	file1, err := fileService.New(ctx, fileRequest)
	is.NoErr(err)
	file2, err := fileService.New(ctx, fileRequest)
	is.NoErr(err)
	file3, err := fileService.New(ctx, fileRequest)
	is.NoErr(err)

	// init the link-service
	linkService := client.NewLinkService(httpClient, testUser.Token)

	// Decode function to decode "from"-body in link
	decodeJSON := func(source, dest interface{}) error {
		data, err := json.Marshal(source)
		if err != nil {
			return fmt.Errorf("json.Marshal: %v", err)
		}
		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}
		return nil
	}

	// Create links for the events
	for _, event := range []client.Event{event1.Created, event2.Created, event3.Created} {
		// create a new bidirectional link request
		createRequest := client.LinkCreateRequest{
			CaseID:    testCase.ID,
			FromID:    event.ID,
			PersonIDs: []string{person1.Created.ID, person2.Created.ID},
			EntityIDs: []string{entity1.Created.ID, entity2.Created.ID},
			FileIDs:   []string{file1.New.ID, file2.New.ID},
		}

		// add one of the eventIDs as well
		var linkedEvent client.Event
		for _, add := range []client.Event{event1.Created, event2.Created, event3.Created} {
			if event.ID != add.ID {
				createRequest.EventIDs = append(createRequest.EventIDs, add.ID)
				linkedEvent = add
				break
			}
		}

		// Create the link and make sure the data is valid
		link, err := linkService.Create(ctx, createRequest)
		is.NoErr(err)
		is.Equal(link.Linked.Events, []client.Event{linkedEvent})
		is.Equal(link.Linked.Persons, []client.Person{person1.Created, person2.Created})
		is.Equal(link.Linked.Files, []client.File{file1.New, file2.New})
		is.Equal(link.Linked.Entities, []client.Entity{entity1.Created, entity2.Created})

		// also make sure the from-object is valid
		var from client.Event
		is.NoErr(decodeJSON(link.Linked.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, event)                         // from-object should be equal to the current event

		// Get the link and make sure it is valid
		gotten, err := linkService.Get(ctx, client.LinkGetRequest{ID: link.Linked.ID, CaseID: testCase.ID})
		is.NoErr(err)
		is.Equal(gotten.Link, link.Linked)

		// Add the rest of the links
		addRequest := client.LinkAddRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			PersonIDs: []string{person3.Created.ID},
			EntityIDs: []string{entity3.Created.ID},
			FileIDs:   []string{file3.New.ID},
		}
		// add the next eventID as well
		var addedEvent client.Event
		for _, add := range []client.Event{event1.Created, event2.Created, event3.Created} {
			if event.ID != add.ID && add.ID != createRequest.EventIDs[0] {
				addRequest.EventIDs = append(addRequest.EventIDs, add.ID)
				addedEvent = add
				break
			}
		}
		added, err := linkService.Add(ctx, addRequest)
		is.NoErr(err)
		is.Equal(added.AddedLinks.Events, []client.Event{linkedEvent, addedEvent})
		is.Equal(added.AddedLinks.Persons, []client.Person{person1.Created, person2.Created, person3.Created})
		is.Equal(added.AddedLinks.Files, []client.File{file1.New, file2.New, file3.New})
		is.Equal(added.AddedLinks.Entities, []client.Entity{entity1.Created, entity2.Created, entity3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(added.AddedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, event)                              // from-object should be equal to the current event

		// Remove the first linked objects
		removeRequest := client.LinkRemoveRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			EventIDs:  []string{linkedEvent.ID},
			PersonIDs: []string{person1.Created.ID},
			EntityIDs: []string{entity1.Created.ID},
			FileIDs:   []string{file1.New.ID},
		}

		// Send the request and validate the data
		removed, err := linkService.Remove(ctx, removeRequest)
		is.NoErr(err)
		is.Equal(removed.RemovedLinks.Events, []client.Event{addedEvent})
		is.Equal(removed.RemovedLinks.Persons, []client.Person{person2.Created, person3.Created})
		is.Equal(removed.RemovedLinks.Files, []client.File{file2.New, file3.New})
		is.Equal(removed.RemovedLinks.Entities, []client.Entity{entity2.Created, entity3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(removed.RemovedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, event)                                  // from-object should be equal to the current event

		// Delete the link
		_, err = linkService.Delete(ctx, client.LinkDeleteRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.NoErr(err)

		// Get the link and make sure it fails since it is removed
		_, err = linkService.Get(ctx, client.LinkGetRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.Equal(true, strings.HasPrefix(err.Error(), "not found: "))
	}

	// Create links for the persons
	for _, person := range []client.Person{person1.Created, person2.Created, person3.Created} {
		// create a new bidirectional link request
		createRequest := client.LinkCreateRequest{
			CaseID:    testCase.ID,
			FromID:    person.ID,
			EventIDs:  []string{event1.Created.ID, event2.Created.ID},
			EntityIDs: []string{entity1.Created.ID, entity2.Created.ID},
			FileIDs:   []string{file1.New.ID, file2.New.ID},
		}

		// add one of the personIDs as well
		var linkedPerson client.Person
		for _, add := range []client.Person{person1.Created, person2.Created, person3.Created} {
			if person.ID != add.ID {
				createRequest.PersonIDs = append(createRequest.PersonIDs, add.ID)
				linkedPerson = add
				break
			}
		}

		// Create the link and make sure the data is valid
		link, err := linkService.Create(ctx, createRequest)
		is.NoErr(err)
		is.Equal(link.Linked.Persons, []client.Person{linkedPerson})
		is.Equal(link.Linked.Events, []client.Event{event1.Created, event2.Created})
		is.Equal(link.Linked.Files, []client.File{file1.New, file2.New})
		is.Equal(link.Linked.Entities, []client.Entity{entity1.Created, entity2.Created})

		// also make sure the from-object is valid
		var from client.Person
		is.NoErr(decodeJSON(link.Linked.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, person)                        // from-object should be equal to the current person

		// Get the link and make sure it is valid
		gotten, err := linkService.Get(ctx, client.LinkGetRequest{ID: link.Linked.ID, CaseID: testCase.ID})
		is.NoErr(err)
		is.Equal(gotten.Link, link.Linked)

		// Add the rest of the links
		addRequest := client.LinkAddRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			EventIDs:  []string{event3.Created.ID},
			EntityIDs: []string{entity3.Created.ID},
			FileIDs:   []string{file3.New.ID},
		}
		// add the next personID as well
		var addedPerson client.Person
		for _, add := range []client.Person{person1.Created, person2.Created, person3.Created} {
			if person.ID != add.ID && add.ID != createRequest.PersonIDs[0] {
				addRequest.PersonIDs = append(addRequest.PersonIDs, add.ID)
				addedPerson = add
				break
			}
		}
		added, err := linkService.Add(ctx, addRequest)
		is.NoErr(err)
		is.Equal(added.AddedLinks.Persons, []client.Person{linkedPerson, addedPerson})
		is.Equal(added.AddedLinks.Events, []client.Event{event1.Created, event2.Created, event3.Created})
		is.Equal(added.AddedLinks.Files, []client.File{file1.New, file2.New, file3.New})
		is.Equal(added.AddedLinks.Entities, []client.Entity{entity1.Created, entity2.Created, entity3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(added.AddedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, person)                             // from-object should be equal to the current person

		// Remove the first linked objects
		removeRequest := client.LinkRemoveRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			PersonIDs: []string{linkedPerson.ID},
			EventIDs:  []string{event1.Created.ID},
			EntityIDs: []string{entity1.Created.ID},
			FileIDs:   []string{file1.New.ID},
		}

		// Send the request and validate the data
		removed, err := linkService.Remove(ctx, removeRequest)
		is.NoErr(err)
		is.Equal(removed.RemovedLinks.Persons, []client.Person{addedPerson})
		is.Equal(removed.RemovedLinks.Events, []client.Event{event2.Created, event3.Created})
		is.Equal(removed.RemovedLinks.Files, []client.File{file2.New, file3.New})
		is.Equal(removed.RemovedLinks.Entities, []client.Entity{entity2.Created, entity3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(removed.RemovedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, person)                                 // from-object should be equal to the current person

		// Delete the link
		_, err = linkService.Delete(ctx, client.LinkDeleteRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.NoErr(err)

		// Get the link and make sure it fails since it is removed
		_, err = linkService.Get(ctx, client.LinkGetRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.Equal(true, strings.HasPrefix(err.Error(), "not found: "))
	}

	// Create links for the entitys
	for _, entity := range []client.Entity{entity1.Created, entity2.Created, entity3.Created} {
		// create a new bidirectional link request
		createRequest := client.LinkCreateRequest{
			CaseID:    testCase.ID,
			FromID:    entity.ID,
			EventIDs:  []string{event1.Created.ID, event2.Created.ID},
			PersonIDs: []string{person1.Created.ID, person2.Created.ID},
			FileIDs:   []string{file1.New.ID, file2.New.ID},
		}

		// add one of the entityIDs as well
		var linkedEntity client.Entity
		for _, add := range []client.Entity{entity1.Created, entity2.Created, entity3.Created} {
			if entity.ID != add.ID {
				createRequest.EntityIDs = append(createRequest.EntityIDs, add.ID)
				linkedEntity = add
				break
			}
		}

		// Create the link and make sure the data is valid
		link, err := linkService.Create(ctx, createRequest)
		is.NoErr(err)
		is.Equal(link.Linked.Entities, []client.Entity{linkedEntity})
		is.Equal(link.Linked.Events, []client.Event{event1.Created, event2.Created})
		is.Equal(link.Linked.Files, []client.File{file1.New, file2.New})
		is.Equal(link.Linked.Persons, []client.Person{person1.Created, person2.Created})

		// also make sure the from-object is valid
		var from client.Entity
		is.NoErr(decodeJSON(link.Linked.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, entity)                        // from-object should be equal to the current entity

		// Get the link and make sure it is valid
		gotten, err := linkService.Get(ctx, client.LinkGetRequest{ID: link.Linked.ID, CaseID: testCase.ID})
		is.NoErr(err)
		is.Equal(gotten.Link, link.Linked)

		// Add the rest of the links
		addRequest := client.LinkAddRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			EventIDs:  []string{event3.Created.ID},
			PersonIDs: []string{person3.Created.ID},
			FileIDs:   []string{file3.New.ID},
		}
		// add the next entityID as well
		var addedEntity client.Entity
		for _, add := range []client.Entity{entity1.Created, entity2.Created, entity3.Created} {
			if entity.ID != add.ID && add.ID != createRequest.EntityIDs[0] {
				addRequest.EntityIDs = append(addRequest.EntityIDs, add.ID)
				addedEntity = add
				break
			}
		}
		added, err := linkService.Add(ctx, addRequest)
		is.NoErr(err)
		is.Equal(added.AddedLinks.Entities, []client.Entity{linkedEntity, addedEntity})
		is.Equal(added.AddedLinks.Events, []client.Event{event1.Created, event2.Created, event3.Created})
		is.Equal(added.AddedLinks.Files, []client.File{file1.New, file2.New, file3.New})
		is.Equal(added.AddedLinks.Persons, []client.Person{person1.Created, person2.Created, person3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(added.AddedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, entity)                             // from-object should be equal to the current entity

		// Remove the first linked objects
		removeRequest := client.LinkRemoveRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			EntityIDs: []string{linkedEntity.ID},
			EventIDs:  []string{event1.Created.ID},
			PersonIDs: []string{person1.Created.ID},
			FileIDs:   []string{file1.New.ID},
		}

		// Send the request and validate the data
		removed, err := linkService.Remove(ctx, removeRequest)
		is.NoErr(err)
		is.Equal(removed.RemovedLinks.Entities, []client.Entity{addedEntity})
		is.Equal(removed.RemovedLinks.Events, []client.Event{event2.Created, event3.Created})
		is.Equal(removed.RemovedLinks.Files, []client.File{file2.New, file3.New})
		is.Equal(removed.RemovedLinks.Persons, []client.Person{person2.Created, person3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(removed.RemovedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, entity)                                 // from-object should be equal to the current entity

		// Delete the link
		_, err = linkService.Delete(ctx, client.LinkDeleteRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.NoErr(err)

		// Get the link and make sure it fails since it is removed
		_, err = linkService.Get(ctx, client.LinkGetRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.Equal(true, strings.HasPrefix(err.Error(), "not found: "))
	}

	// Create links for the files
	for _, file := range []client.File{file1.New, file2.New, file3.New} {
		// create a new bidirectional link request
		createRequest := client.LinkCreateRequest{
			CaseID:    testCase.ID,
			FromID:    file.ID,
			EventIDs:  []string{event1.Created.ID, event2.Created.ID},
			PersonIDs: []string{person1.Created.ID, person2.Created.ID},
			EntityIDs: []string{entity1.Created.ID, entity2.Created.ID},
		}

		// add one of the FileIDs as well
		var linkedFile client.File
		for _, add := range []client.File{file1.New, file2.New, file3.New} {
			if file.ID != add.ID {
				createRequest.FileIDs = append(createRequest.FileIDs, add.ID)
				linkedFile = add
				break
			}
		}

		// Create the link and make sure the data is valid
		link, err := linkService.Create(ctx, createRequest)
		is.NoErr(err)
		is.Equal(link.Linked.Files, []client.File{linkedFile})
		is.Equal(link.Linked.Events, []client.Event{event1.Created, event2.Created})
		is.Equal(link.Linked.Entities, []client.Entity{entity1.Created, entity2.Created})
		is.Equal(link.Linked.Persons, []client.Person{person1.Created, person2.Created})

		// also make sure the from-object is valid
		var from client.File
		is.NoErr(decodeJSON(link.Linked.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, file)                          // from-object should be equal to the current file

		// Get the link and make sure it is valid
		gotten, err := linkService.Get(ctx, client.LinkGetRequest{ID: link.Linked.ID, CaseID: testCase.ID})
		is.NoErr(err)
		is.Equal(gotten.Link, link.Linked)

		// Add the rest of the links
		addRequest := client.LinkAddRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			EventIDs:  []string{event3.Created.ID},
			PersonIDs: []string{person3.Created.ID},
			EntityIDs: []string{entity3.Created.ID},
		}
		// add the next fileID as well
		var addedFile client.File
		for _, add := range []client.File{file1.New, file2.New, file3.New} {
			if file.ID != add.ID && add.ID != createRequest.FileIDs[0] {
				addRequest.FileIDs = append(addRequest.FileIDs, add.ID)
				addedFile = add
				break
			}
		}
		added, err := linkService.Add(ctx, addRequest)
		is.NoErr(err)
		is.Equal(added.AddedLinks.Files, []client.File{linkedFile, addedFile})
		is.Equal(added.AddedLinks.Events, []client.Event{event1.Created, event2.Created, event3.Created})
		is.Equal(added.AddedLinks.Entities, []client.Entity{entity1.Created, entity2.Created, entity3.Created})
		is.Equal(added.AddedLinks.Persons, []client.Person{person1.Created, person2.Created, person3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(added.AddedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, file)                               // from-object should be equal to the current file

		// Remove the first linked objects
		removeRequest := client.LinkRemoveRequest{
			ID:        link.Linked.ID,
			CaseID:    testCase.ID,
			FileIDs:   []string{linkedFile.ID},
			EventIDs:  []string{event1.Created.ID},
			PersonIDs: []string{person1.Created.ID},
			EntityIDs: []string{entity1.Created.ID},
		}

		// Send the request and validate the data
		removed, err := linkService.Remove(ctx, removeRequest)
		is.NoErr(err)
		is.Equal(removed.RemovedLinks.Files, []client.File{addedFile})
		is.Equal(removed.RemovedLinks.Events, []client.Event{event2.Created, event3.Created})
		is.Equal(removed.RemovedLinks.Entities, []client.Entity{entity2.Created, entity3.Created})
		is.Equal(removed.RemovedLinks.Persons, []client.Person{person2.Created, person3.Created})

		// also make sure the from-object is valid
		is.NoErr(decodeJSON(removed.RemovedLinks.From, &from)) // should succeed to decode "from"-body
		is.Equal(from, file)                                   // from-object should be equal to the current file

		// Delete the link
		_, err = linkService.Delete(ctx, client.LinkDeleteRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.NoErr(err)

		// Get the link and make sure it fails since it is removed
		_, err = linkService.Get(ctx, client.LinkGetRequest{ID: removed.RemovedLinks.ID, CaseID: testCase.ID})
		is.Equal(true, strings.HasPrefix(err.Error(), "not found: "))
	}
}
