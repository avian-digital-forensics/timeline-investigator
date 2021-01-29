package datastore_test

import (
	"context"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/matryer/is"
)

func TestSearch(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	caseID := "39cfa9d974fb4f759c79a051e60e5b9e"
	prefix := "g"

	db, err := datastore.NewService("http://localhost:9200")
	is.NoErr(err)
	keywords, err := db.SearchKeywords(ctx, caseID, prefix)
	is.NoErr(err)
	log.Println(keywords)
}
