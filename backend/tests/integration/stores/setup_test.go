package stores

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/family"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/transactions"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_canSetupANewApp(t *testing.T) {
	ctx := context.Background()
	db := acquireDb(t)
	familiesCol := db.Collection(family.CollectionFamilies)
	usersCol := db.Collection(user.CollectionUsers)
	deleteAll(t, ctx, familiesCol, usersCol)

	usersSvc := user.NewService(usersCol)
	svc := setup.NewService(
		transactions.NoTransactions(),
		usersSvc,
		family.NewService(familiesCol),
	)

	adminPassword := "adminPassword01"

	called := time.Now().UTC()
	s := &setup.Data{
		DefaultAdmin: &types.User{
			Username:  "admin",
			Firstname: "Admina",
			Lastname:  "Strator",
		},
		DefaultAdminPassword: []byte("adminPassword01"), // this will get purged from memory
		DefaultFamily: &types.Family{
			Name:        "Default",
			Description: "Default family",
		},
	}
	err := svc.Setup(ctx, s)
	require.NoError(t, err)

	require.NotNil(t, s.DefaultAdmin.ID)
	require.True(t, called.Before(s.DefaultAdmin.UpdatedAt))

	require.NotNil(t, s.DefaultFamily.ID)
	require.True(t, called.Before(s.DefaultFamily.UpdatedAt))

	requireCountEq(t, ctx, usersCol, 1)
	requireCountEq(t, ctx, familiesCol, 1)

	newAdmin, err := usersSvc.Login(ctx, s.DefaultAdmin.Username, []byte(adminPassword))
	require.NoError(t, err)
	rqUsersMatches(t, newAdmin, s.DefaultAdmin)

	// @todo test later with find
	//newFamily, err := svc.FindByID(ctx, f.ID)
	//require.NoError(t, err)
	//
	//compare.Equal(t, newFamily, &f, compare.IgnoreFields(types.Family{}, "UpdatedAt"))
}
