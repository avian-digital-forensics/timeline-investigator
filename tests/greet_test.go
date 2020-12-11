package tests

import (
	"context"
	"log"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/tests/client"

	"github.com/matryer/is"
)

// TestGreet test the greet-method
func TestGreet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	httpClient := client.New("http://localhost:8080/api/")
	httpClient.Debug = func(s string) {
		log.Println(s)
	}

	service := client.NewService(httpClient, "")
	resp, err := service.Greet(ctx, client.GreetRequest{Name: "Simon"})
	is.NoErr(err)
	is.Equal(resp.Greeting, "Hello Simon")
}
