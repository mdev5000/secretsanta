package user

import (
	"context"
	"github.com/mdev5000/secretsanta/testutil/compare"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"

	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/stretchr/testify/require"
)

func Test_canCreateNewUsers(t *testing.T) {
	db := acquireDb(t)
	col := db.Collection(user.CollectionUsers)
	svc := user.NewService(col)
	u := user.User{
		Username:  "bob",
		Firstname: "Bob",
		Lastname:  "Test",
	}
	called := time.Now().UTC()
	ctx := context.Background()
	err := svc.Create(ctx, &u, []byte("mypassword"))
	require.NoError(t, err)
	require.NotNil(t, u.ID)
	require.True(t, called.Before(u.UpdatedAt))

	numDocs, err := col.CountDocuments(ctx, bson.D{})
	require.NoError(t, err)
	require.Equal(t, int64(1), numDocs)

	newUser, err := svc.FindByID(ctx, u.ID)
	require.NoError(t, err)

	compare.Equal(t, newUser, &u, compare.IgnoreFields(user.User{}, "UpdatedAt", "PasswordHash"))
}

func Test_canLoginNewUsers(t *testing.T) {
	db := acquireDb(t)
	col := db.Collection(user.CollectionUsers)
	svc := user.NewService(col)
	u := user.User{
		Username:  "bob",
		Firstname: "Bob",
		Lastname:  "Test",
	}
	ctx := context.Background()
	err := svc.Create(ctx, &u, []byte("mypassword"))
	require.NoError(t, err)

	t.Run("no error when password ok", func(t *testing.T) {
		user2, err := svc.Login(ctx, u.Username, []byte("mypassword"))
		require.NoError(t, err)
		require.Equal(t, user2.ID, u.ID)
	})

	t.Run("error when password is invalid", func(t *testing.T) {
		user3, err := svc.Login(ctx, u.Username, []byte("badpassword"))
		require.Error(t, err)
		require.Nil(t, user3)
	})
}
