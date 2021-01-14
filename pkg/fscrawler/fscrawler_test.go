package fscrawler_test

import (
	"context"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/fscrawler"
	"github.com/matryer/is"
)

func TestClient(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	// Create client
	fs := fscrawler.New("http://localhost:8080/fscrawler")

	// Test status
	ok, err := fs.Ping(ctx)
	is.NoErr(err)
	is.Equal(ok, true)

	// Test processing
	process := fs.NewProcess("test-1.txt").WithIndex("test").WithID("id")
	err = process.Start()
	is.NoErr(err)
}
