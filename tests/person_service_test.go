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

	testCase, err := testUser.newTestCase(ctx, client.NewCaseService(httpClient, testUser.Token))
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

	// Create two new persons
	person2, err := personService.Create(ctx, client.PersonCreateRequest{CaseID: testCase.ID, FirstName: "Simon 2"})
	is.NoErr(err)
	person3, err := personService.Create(ctx, client.PersonCreateRequest{CaseID: testCase.ID, FirstName: "Simon 3"})
	is.NoErr(err)

	// List all persons in the case
	all, err := personService.List(ctx, client.PersonListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(all.Persons), 3)

	// Validate the list in a for loop
	for i, person := range all.Persons {
		var check client.Person
		if i == 0 {
			check = updated.Updated
		} else if i == 1 {
			check = person2.Created
		} else {
			check = person3.Created
		}

		is.Equal(person.ID, check.ID)
		is.Equal(person.CreatedAt, check.CreatedAt)
		is.Equal(person.UpdatedAt, check.UpdatedAt)
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
