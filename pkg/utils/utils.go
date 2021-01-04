package utils

import (
	"context"
	"net/http"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
)

// GetToken gets token from request
func GetToken(req *http.Request) string {
	return req.Header.Get("Authorization")
}

var userType api.User

// SetUser sets the user to the context
func SetUser(ctx context.Context, user api.User) context.Context {
	return context.WithValue(ctx, userType, user)
}

// GetUser gets the user from the context
func GetUser(ctx context.Context) api.User {
	user := ctx.Value(userType)
	if user == nil {
		return api.User{}
	}
	return user.(api.User)
}
