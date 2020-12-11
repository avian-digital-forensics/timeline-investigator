package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
)

type Service struct {
}

func (s *Service) Greet(ctx context.Context, r api.GreetRequest) (*api.GreetResponse, error) {
	return &api.GreetResponse{Greeting: fmt.Sprintf("Hello %s", r.Name)}, nil
}

func (s *Service) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	return ctx, nil
}
