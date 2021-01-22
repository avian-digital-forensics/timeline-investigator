package tests

import (
	"context"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/matryer/is"
)

func TestPersonService(t *testing.T) {
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

	personService := client.NewPersonService(httpClient, testUser.Token)

	// Make a create person request
	request := client.PersonCreateRequest{
		CaseID:        testCase.ID,
		FirstName:     "Simon",
		LastName:      "Jansson",
		EmailAddress:  "sja@avian.dk",
		PostalAddress: "Vegagatan 6A, 113 29, Stockholm, SE",
		WorkAddress:   "Applebys Plads 7, 1411 Copenhagen, DK",
		TelephoneNo:   "+46765550125",
		Custom: map[string]interface{}{
			"Age":            "24",
			"Favourite song": "Purple Haze",
		},
	}

	// Create the person and validate the data
	person, err := personService.Create(ctx, request)
	is.NoErr(err)
	is.Equal(person.Created.FirstName, request.FirstName)
	is.Equal(person.Created.LastName, request.LastName)
	is.Equal(person.Created.EmailAddress, request.EmailAddress)
	is.Equal(person.Created.PostalAddress, request.PostalAddress)
	is.Equal(person.Created.WorkAddress, request.WorkAddress)
	is.Equal(person.Created.TelephoneNo, request.TelephoneNo)
	is.Equal(person.Created.Custom, request.Custom)

	// Get the person and validate the data
	gotten, err := personService.Get(ctx, client.PersonGetRequest{ID: person.Created.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(gotten.Person.ID, person.Created.ID)
	is.Equal(gotten.Person.CreatedAt, person.Created.CreatedAt)
	is.Equal(gotten.Person.UpdatedAt, person.Created.UpdatedAt)
	is.Equal(gotten.Person.DeletedAt, person.Created.DeletedAt)
	is.Equal(gotten.Person.FirstName, person.Created.FirstName)
	is.Equal(gotten.Person.LastName, person.Created.LastName)
	is.Equal(gotten.Person.EmailAddress, person.Created.EmailAddress)
	is.Equal(gotten.Person.PostalAddress, person.Created.PostalAddress)
	is.Equal(gotten.Person.WorkAddress, person.Created.WorkAddress)
	is.Equal(gotten.Person.TelephoneNo, person.Created.TelephoneNo)
	is.Equal(gotten.Person.Custom, person.Created.Custom)

	// Make an update request
	updateRequest := client.PersonUpdateRequest{
		ID:            person.Created.ID,
		CaseID:        testCase.ID,
		FirstName:     "Smino",
		LastName:      "Jazz",
		EmailAddress:  "contact@simonjansson.dev",
		PostalAddress: "Vegagatan 6B, 113 29, Stockholm, SE",
		WorkAddress:   "Applebys Plads 9, 1411 Copenhagen, DK",
		TelephoneNo:   "+46765350125",
		Custom: map[string]interface{}{
			"Age":              "25",
			"Favourite artist": "Jimi Hendrix",
		},
	}

	// Update the person and validate the data
	updated, err := personService.Update(ctx, updateRequest)
	is.NoErr(err)
	is.Equal(updated.Updated.ID, updateRequest.ID)
	is.Equal(updated.Updated.FirstName, updateRequest.FirstName)
	is.Equal(updated.Updated.LastName, updateRequest.LastName)
	is.Equal(updated.Updated.EmailAddress, updateRequest.EmailAddress)
	is.Equal(updated.Updated.PostalAddress, updateRequest.PostalAddress)
	is.Equal(updated.Updated.WorkAddress, updateRequest.WorkAddress)
	is.Equal(updated.Updated.TelephoneNo, updateRequest.TelephoneNo)
	is.Equal(updated.Updated.Custom, updateRequest.Custom)

	// Get all keywords for the case where the person belongs
	keywords, err := caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 0)

	// Add two keywords to the first person
	keywordsRequest := client.KeywordsAddRequest{ID: person.Created.ID, CaseID: testCase.ID, Keywords: []string{"healthy", "green"}}
	_, err = personService.KeywordsAdd(ctx, keywordsRequest)
	is.NoErr(err)

	// Get the keywords for the case again and make sure they have been added for the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 2)
	is.Equal(keywordsRequest.Keywords[0], keywords.Keywords[0])
	is.Equal(keywordsRequest.Keywords[1], keywords.Keywords[1])

	// Get the person again and make sure the keywords has been added there as well
	gotten, err = personService.Get(ctx, client.PersonGetRequest{CaseID: testCase.ID, ID: person.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Person.Keywords), 2)
	is.Equal(gotten.Person.Keywords, keywordsRequest.Keywords)
	// make sure the rest of the data is valid as well
	is.Equal(gotten.Person.ID, updated.Updated.ID)
	is.Equal(gotten.Person.DeletedAt, updated.Updated.DeletedAt)
	is.Equal(gotten.Person.FirstName, updated.Updated.FirstName)
	is.Equal(gotten.Person.LastName, updated.Updated.LastName)
	is.Equal(gotten.Person.EmailAddress, updated.Updated.EmailAddress)
	is.Equal(gotten.Person.PostalAddress, updated.Updated.PostalAddress)
	is.Equal(gotten.Person.WorkAddress, updated.Updated.WorkAddress)
	is.Equal(gotten.Person.TelephoneNo, updated.Updated.TelephoneNo)
	is.Equal(gotten.Person.Custom, updated.Updated.Custom)

	// Create two new persons
	person2, err := personService.Create(ctx, client.PersonCreateRequest{CaseID: testCase.ID, FirstName: "Simon 2"})
	is.NoErr(err)
	person3, err := personService.Create(ctx, client.PersonCreateRequest{CaseID: testCase.ID, FirstName: "Simon 3"})
	is.NoErr(err)

	// List all persons in the case
	all, err := personService.List(ctx, client.PersonListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(all.Persons), 3)

	// Add the keywords + another one, to person2 as well
	keywordsRequest.ID = person2.Created.ID
	keywordsRequest.Keywords = append(keywordsRequest.Keywords, "another one")
	_, err = personService.KeywordsAdd(ctx, keywordsRequest)
	is.NoErr(err)

	// Make sure the keywords has been added to person2
	gotten, err = personService.Get(ctx, client.PersonGetRequest{CaseID: testCase.ID, ID: person2.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Person.Keywords), 3)
	is.Equal(gotten.Person.Keywords, keywordsRequest.Keywords)

	// We should now have 3 (unqiue) keywords in the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)
	is.Equal(keywordsRequest.Keywords[0], keywords.Keywords[0])
	is.Equal(keywordsRequest.Keywords[1], keywords.Keywords[1])
	is.Equal(keywordsRequest.Keywords[2], keywords.Keywords[2])

	// Remove the keyword first keyword from person2
	removeRequest := client.KeywordsRemoveRequest{ID: person2.Created.ID, CaseID: testCase.ID, Keywords: []string{keywords.Keywords[0]}}
	_, err = personService.KeywordsRemove(ctx, removeRequest)
	is.NoErr(err)

	// Make sure it was removed from person2
	gotten, err = personService.Get(ctx, client.PersonGetRequest{CaseID: testCase.ID, ID: person2.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Person.Keywords), 2)
	is.Equal(gotten.Person.Keywords[0], keywordsRequest.Keywords[1])
	is.Equal(gotten.Person.Keywords[1], keywordsRequest.Keywords[2])

	// Make sure it is still available in the case
	keywords, err = caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords.Keywords), 3)

	// Make sure it is still available in the first person
	gotten, err = personService.Get(ctx, client.PersonGetRequest{CaseID: testCase.ID, ID: person.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Person.Keywords), 2)
	is.Equal(gotten.Person.Keywords[0], keywordsRequest.Keywords[0])
	is.Equal(gotten.Person.Keywords[1], keywordsRequest.Keywords[1])

	// Remove it from the first person as well
	removeRequest.ID = person.Created.ID
	_, err = personService.KeywordsRemove(ctx, removeRequest)
	is.NoErr(err)

	// Make sure it was removed from the case (since no person has the keyword anymore)
	keywords5, err := caseService.Keywords(ctx, client.CaseKeywordsRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(keywords5.Keywords), 2)
	is.Equal(keywords5.Keywords[0], keywordsRequest.Keywords[1])
	is.Equal(keywords5.Keywords[1], keywordsRequest.Keywords[2])

	// Make sure it was removed from the first person as well
	gotten, err = personService.Get(ctx, client.PersonGetRequest{CaseID: testCase.ID, ID: person.Created.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Person.Keywords), 1)
	is.Equal(gotten.Person.Keywords[0], keywordsRequest.Keywords[1])

	// Validate the list in a for loop
	for _, person := range all.Persons {
		var check client.Person
		if person.ID == gotten.Person.ID {
			check = gotten.Person
		} else if person.ID == person2.Created.ID {
			check = person2.Created
		} else {
			check = person3.Created
		}

		is.Equal(person.ID, check.ID)
		is.Equal(person.CreatedAt, check.CreatedAt)
		is.Equal(person.DeletedAt, check.DeletedAt)
		is.Equal(person.FirstName, check.FirstName)
		is.Equal(person.LastName, check.LastName)
		is.Equal(person.EmailAddress, check.EmailAddress)
		is.Equal(person.PostalAddress, check.PostalAddress)
		is.Equal(person.WorkAddress, check.WorkAddress)
		is.Equal(person.TelephoneNo, check.TelephoneNo)
		is.Equal(person.Custom, check.Custom)
	}

	// Delete the first created person
	_, err = personService.Delete(ctx, client.PersonDeleteRequest{ID: person.Created.ID, CaseID: testCase.ID})
	is.NoErr(err)

	// List all persons and make sure the correct person was deleted
	all, err = personService.List(ctx, client.PersonListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(all.Persons), 2)
	for _, p := range all.Persons {
		is.Equal(p.ID == person.Created.ID, false)
	}
}
