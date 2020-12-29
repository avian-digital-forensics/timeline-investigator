package authentication_test

import (
	"context"
	"os"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/authentication"

	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	credentialsFile = os.Getenv("AUTH_CREDENTIALS_FILE")
	apiKey          = os.Getenv("AUTH_API_KEY")
)

func TestAuthentication(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	authService, err := authentication.New(ctx, credentialsFile, apiKey)
	is.NoErr(err)

	var uid = uuid.New().String()
	user, err := authService.Create(ctx, uid, uid+"@email.com", uid, uid)
	is.NoErr(err)
	defer is.NoErr(authService.Delete(ctx, user.UID))

	customToken, err := authService.GetCustomToken(ctx, user.UID)
	is.NoErr(err)

	idToken, err := authService.VerifyCustomToken(ctx, customToken)
	is.NoErr(err)

	is.NoErr(authService.Verify(ctx, idToken))
}
