package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"

	"github.com/matryer/is"
)

var (
	URL = os.Getenv("TEST_URL")
)

// TestCaseNew test the New-method
func TestCaseNew(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New(URL)
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	service := client.NewCaseService(httpClient, "")
	_, err := service.New(ctx, client.CaseNewRequest{Name: "Simon"})
	is.Equal(err.Error(), "Not implemented yet")
}
