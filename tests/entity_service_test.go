package tests

import (
	"context"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/matryer/is"
)

func TestEntityService(t *testing.T) {
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

	entityService := client.NewEntityService(httpClient, testUser.Token)

	// Create a new entity request
	request := client.EntityCreateRequest{
		CaseID:   testCase.ID,
		Title:    "Avian APS",
		PhotoURL: "https://randomURL.com",
		Type:     "something", // This is invalid (should fail)
		Custom:   map[string]interface{}{"CEO": "Jacob Isaksen", "Employees": "150"},
	}

	// Try to create the new entity (should fail because of the type)
	_, err = entityService.Create(ctx, request)
	is.Equal(err.Error(), api.ErrInvalidEntityType.Error())

	// Get all entity-types
	types, err := entityService.Types(ctx, client.EntityTypesRequest{})
	is.NoErr(err)
	is.Equal(len(types.EntityTypes), 2)
	is.Equal(types.EntityTypes[0], "organization")
	is.Equal(types.EntityTypes[1], "location")

	// Use a valid entity-type and try to create the entity again (should succeed)
	request.Type = types.EntityTypes[0]
	entity, err := entityService.Create(ctx, request)
	is.NoErr(err)
	is.Equal(entity.Created.Title, request.Title)
	is.Equal(entity.Created.PhotoURL, request.PhotoURL)
	is.Equal(entity.Created.Type, request.Type)
	is.Equal(entity.Created.Custom, request.Custom)

	// Get the entity by ID
	gotten, err := entityService.Get(ctx, client.EntityGetRequest{CaseID: testCase.ID, ID: entity.Created.ID})
	is.NoErr(err)
	is.Equal(gotten.Entity.ID, entity.Created.ID)
	is.Equal(gotten.Entity.CreatedAt, entity.Created.CreatedAt)
	is.Equal(gotten.Entity.UpdatedAt, entity.Created.UpdatedAt)
	is.Equal(gotten.Entity.DeletedAt, entity.Created.DeletedAt)
	is.Equal(gotten.Entity.Title, entity.Created.Title)
	is.Equal(gotten.Entity.PhotoURL, entity.Created.PhotoURL)
	is.Equal(gotten.Entity.Type, entity.Created.Type)
	is.Equal(gotten.Entity.Custom, entity.Created.Custom)

	// Make an update-request
	updateRequest := client.EntityUpdateRequest{
		ID:       entity.Created.ID,
		CaseID:   testCase.ID,
		Title:    "New Title",
		PhotoURL: "https://NewRandomURL.com",
		Type:     types.EntityTypes[1],
		Custom:   map[string]interface{}{"CEO": "Dennis", "Employees": "200"},
	}

	// Update the entity
	updated, err := entityService.Update(ctx, updateRequest)
	is.NoErr(err)
	is.Equal(updated.Updated.ID, updateRequest.ID)
	is.Equal(updated.Updated.Title, updateRequest.Title)
	is.Equal(updated.Updated.PhotoURL, updateRequest.PhotoURL)
	is.Equal(updated.Updated.Type, updateRequest.Type)
	is.Equal(updated.Updated.Custom, updateRequest.Custom)

	// Get the entity again make sure it has been updated
	gotten, err = entityService.Get(ctx, client.EntityGetRequest{CaseID: testCase.ID, ID: entity.Created.ID})
	is.NoErr(err)
	is.Equal(gotten.Entity.ID, updated.Updated.ID)
	is.Equal(gotten.Entity.CreatedAt, updated.Updated.CreatedAt)
	is.Equal(gotten.Entity.UpdatedAt, updated.Updated.UpdatedAt)
	is.Equal(gotten.Entity.DeletedAt, updated.Updated.DeletedAt)
	is.Equal(gotten.Entity.Title, updated.Updated.Title)
	is.Equal(gotten.Entity.PhotoURL, updated.Updated.PhotoURL)
	is.Equal(gotten.Entity.Type, updated.Updated.Type)
	is.Equal(gotten.Entity.Custom, updated.Updated.Custom)

	// Create new entities
	entity2, err := entityService.Create(ctx, client.EntityCreateRequest{CaseID: testCase.ID, Title: "Entity 2", Type: "organization"})
	is.NoErr(err)
	entity3, err := entityService.Create(ctx, client.EntityCreateRequest{CaseID: testCase.ID, Title: "Entity 3", Type: "organization"})
	is.NoErr(err)

	// List all entities for the test-case and make sure the created entities is there and valid
	all, err := entityService.List(ctx, client.EntityListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(all.Entities), 3)
	// Validate the data in a for loop
	for i, entity := range all.Entities {
		var check client.Entity
		if i == 0 {
			check = gotten.Entity
		} else if i == 1 {
			check = entity2.Created
		} else {
			check = entity3.Created
		}

		is.Equal(entity.ID, check.ID)
		is.Equal(entity.CreatedAt, check.CreatedAt)
		is.Equal(entity.UpdatedAt, check.UpdatedAt)
		is.Equal(entity.DeletedAt, check.DeletedAt)
		is.Equal(entity.Title, check.Title)
		is.Equal(entity.PhotoURL, check.PhotoURL)
		is.Equal(entity.Type, check.Type)
		is.Equal(entity.Custom, check.Custom)
	}

	// Delete the first created entity
	_, err = entityService.Delete(ctx, client.EntityDeleteRequest{ID: entity.Created.ID, CaseID: testCase.ID})
	is.NoErr(err)

	// Make sure the deleted entity won't be listed
	all, err = entityService.List(ctx, client.EntityListRequest{CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(all.Entities), 2)
}
