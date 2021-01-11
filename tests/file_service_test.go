package tests

import (
	"context"
	b64 "encoding/base64"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/matryer/is"
)

func TestFileService(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(testURL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	// Create a test-user for testing the file-service
	testUser, err := newTestUser(ctx, client.NewTestService(httpClient, ""))
	is.NoErr(err)
	defer testUser.delete(ctx)

	// Create a test-case to use
	caseService := client.NewCaseService(httpClient, testUser.Token)
	testCase, err := testUser.newTestCase(ctx, caseService)
	is.NoErr(err)

	// Initialize the file-service
	fileService := client.NewFileService(httpClient, testUser.Token)

	// Create test-data
	var d1 = []byte("sample\ndata\n1")
	var d2 = []byte("sample\ndata\n2")
	var d3 = []byte("sample\ndata\n3")

	// Create a new file request
	newRequest := client.FileNewRequest{
		CaseID:      testCase.ID,
		Name:        "d1.txt",
		Description: "test1",
		Mime:        "text/plain",
		Data:        b64.URLEncoding.EncodeToString(d1),
	}

	// Upload the new file and make sure data is valid
	file, err := fileService.New(ctx, newRequest)
	is.NoErr(err)
	is.Equal(file.New.Name, newRequest.Name)
	is.Equal(file.New.Description, newRequest.Description)
	is.Equal(file.New.Mime, newRequest.Mime)
	is.Equal(file.New.Processed, false)

	// Open the file and make sure the data is valid
	open, err := fileService.Open(ctx, client.FileOpenRequest{ID: file.New.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(open.Data, newRequest.Data)
	// Decode the data we got from the open file
	// and make sure it is the same data as we uploaded
	data, err := b64.StdEncoding.DecodeString(open.Data)
	is.NoErr(err)
	is.Equal(data, d1)

	// Create a new update request
	updateRequest := client.FileUpdateRequest{ID: file.New.ID, CaseID: testCase.ID, Description: "hej"}

	// Update the created file
	updated, err := fileService.Update(ctx, updateRequest)
	is.NoErr(err)
	is.Equal(updated.Updated.Description, updateRequest.Description)
	is.Equal(updated.Updated.ID, file.New.ID)
	is.Equal(updated.Updated.CreatedAt, file.New.CreatedAt)
	is.Equal(updated.Updated.Mime, file.New.Mime)
	is.Equal(updated.Updated.Processed, file.New.Processed)

	// Open the case and make sure the file is listed
	gotten, err := caseService.Get(ctx, client.CaseGetRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Case.Files), 1)
	is.Equal(gotten.Case.Files[0], updated.Updated)

	// Upload more files to the case
	file2, err := fileService.New(ctx, client.FileNewRequest{CaseID: testCase.ID, Name: "d2.txt", Data: b64.URLEncoding.EncodeToString(d2)})
	is.NoErr(err)
	file3, err := fileService.New(ctx, client.FileNewRequest{CaseID: testCase.ID, Name: "d3.txt", Data: b64.URLEncoding.EncodeToString(d3)})
	is.NoErr(err)

	// Open the case againg and make sure the files has been added
	gotten, err = caseService.Get(ctx, client.CaseGetRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Case.Files), 3)
	is.Equal(gotten.Case.Files[0], updated.Updated)
	is.Equal(gotten.Case.Files[1], file2.New)
	is.Equal(gotten.Case.Files[2], file3.New)

	// Delete a file and make sure the file has been removed from the case
	_, err = fileService.Delete(ctx, client.FileDeleteRequest{CaseID: testCase.ID, ID: file.New.ID})
	is.NoErr(err)
	gotten, err = caseService.Get(ctx, client.CaseGetRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Case.Files), 2)
	is.Equal(gotten.Case.Files[0], file2.New)
	is.Equal(gotten.Case.Files[1], file3.New)
}
