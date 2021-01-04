package utils_test

import (
	"context"
	"testing"

	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/utils"
	"github.com/matryer/is"
)

func TestSetUser(t *testing.T) {
	is := is.New(t)

	user1 := api.User{UID: "1"}
	ctx1 := context.Background()
	ctx1 = utils.SetUser(ctx1, user1)

	user2 := api.User{UID: "2"}
	ctx2 := context.Background()
	ctx2 = utils.SetUser(ctx2, user2)

	user3 := api.User{}
	ctx3 := context.Background()
	ctx3 = utils.SetUser(ctx3, user3)

	is.Equal(user1.UID, utils.GetUser(ctx1).UID)
	is.Equal(user2.UID, utils.GetUser(ctx2).UID)
	is.Equal(user3.UID, utils.GetUser(ctx3).UID)
}
