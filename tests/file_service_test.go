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
	is.Equal(file.New.ProcessedAt, int64(0))

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
	is.Equal(updated.Updated.ProcessedAt, file.New.ProcessedAt)

	if testing.Short() {
		t.Skip("skipping for process in short mode.")
	}

	// Process the file
	processed, err := fileService.Process(ctx, client.FileProcessRequest{ID: file.New.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(processed.Processed.Description, updated.Updated.Description)
	is.Equal(processed.Processed.ID, updated.Updated.ID)
	is.Equal(processed.Processed.CreatedAt, updated.Updated.CreatedAt)
	is.Equal(processed.Processed.Mime, updated.Updated.Mime)
	is.Equal(false, processed.Processed.ProcessedAt == updated.Updated.ProcessedAt)

	// Open the case and make sure the file is listed
	gotten, err := caseService.Get(ctx, client.CaseGetRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Case.Files), 1)
	is.Equal(gotten.Case.Files[0], processed.Processed) // processed latest info for the first file

	// Wait 10 seconds for the file to be indexed
	log.Println("WAIT : 10 seconds for processed file to be indexed")
	time.Sleep(1 * time.Second)
	for _, s := range []int{9, 8, 7, 6, 5, 4, 3, 2, 1} {
		log.Println("WAIT :", s)
		time.Sleep(1 * time.Second)
	}

	// get the processed information for the first file
	processInfo, err := fileService.Processed(ctx, client.FileProcessedRequest{CaseID: testCase.ID, ID: file.New.ID})
	is.NoErr(err)
	is.Equal(processInfo.ID, file.New.ID)

	// Upload more files to the case
	file2, err := fileService.New(ctx, client.FileNewRequest{CaseID: testCase.ID, Name: "d2.txt", Data: b64.URLEncoding.EncodeToString(d2)})
	is.NoErr(err)
	file3, err := fileService.New(ctx, client.FileNewRequest{CaseID: testCase.ID, Name: "d3.txt", Data: b64.URLEncoding.EncodeToString(d3)})
	is.NoErr(err)

	// Open the case againg and make sure the files has been added
	gotten, err = caseService.Get(ctx, client.CaseGetRequest{ID: testCase.ID})
	is.NoErr(err)
	is.Equal(len(gotten.Case.Files), 3)
	is.Equal(gotten.Case.Files[0], processed.Processed)
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

	// Process the second file
	processed2, err := fileService.Process(ctx, client.FileProcessRequest{ID: file2.New.ID, CaseID: testCase.ID})
	is.NoErr(err)
	is.Equal(processed2.Processed.ID, file2.New.ID)

	// Get all processes for the case (should be two)
	_, err = fileService.Processes(ctx, client.FileProcessesRequest{CaseID: testCase.ID})
	is.NoErr(err)
}
