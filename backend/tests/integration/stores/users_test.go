package stores

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/mdev5000/secretsanta/testutil/compare"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"testing"
	"time"

	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/stretchr/testify/require"
)

func Test_canCreateNewUsers(t *testing.T) {
	ctx := context.Background()
	db := acquireDb(t)
	col := db.Collection(user.CollectionUsers)
	deleteAll(t, ctx, col)

	svc := user.NewService(col)
	u := types.User{
		Username:  "bob",
		Firstname: "Bob",
		Lastname:  "Test",
		FamilyIDs: []string{"family1", "family2"},
	}
	called := time.Now().UTC()
	err := svc.Create(ctx, &u, []byte("mypassword"))
	require.NoError(t, err)
	require.NotNil(t, u.ID)
	require.True(t, called.Before(u.UpdatedAt))

	numDocs, err := col.CountDocuments(ctx, bson.D{})
	require.NoError(t, err)
	require.Equal(t, int64(1), numDocs)

	newUser, err := svc.FindByID(ctx, u.ID)
	require.NoError(t, err)

	rqUsersMatches(t, newUser, &u)
}

func rqUsersMatches(t *testing.T, existing, expected *types.User) {
	compare.Equal(t, expected, existing, compare.IgnoreFields(types.User{}, "UpdatedAt", "PasswordHash"))
}

func Test_canLoginNewUsers(t *testing.T) {
	ctx := context.Background()
	db := acquireDb(t)
	col := db.Collection(user.CollectionUsers)
	deleteAll(t, ctx, col)

	svc := user.NewService(col)
	u := types.User{
		Username:  "bob",
		Firstname: "Bob",
		Lastname:  "Test",
	}
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

func Test_canFindUsersByFilter(t *testing.T) {
	ctx := context.Background()
	db := acquireDb(t)
	col := db.Collection(user.CollectionUsers)
	deleteAll(t, ctx, col)

	store := user.Store{Users: col}

	uFam1 := types.User{ID: "uFam1", FamilyIDs: []string{"family1"}}
	require.NoError(t, store.Create(ctx, &uFam1, nil))

	uFam2 := types.User{ID: "uFam2", FamilyIDs: []string{"family2"}}
	require.NoError(t, store.Create(ctx, &uFam2, nil))

	uFam1And2 := types.User{ID: "uFam1And2", FamilyIDs: []string{"family1", "family2"}}
	require.NoError(t, store.Create(ctx, &uFam1And2, nil))

	uFam3 := types.User{ID: "uFam3", FamilyIDs: []string{"family3"}}
	require.NoError(t, store.Create(ctx, &uFam3, nil))

	uFam4 := types.User{ID: "uFam4", FamilyIDs: []string{"family4"}}
	require.NoError(t, store.Create(ctx, &uFam4, nil))

	uFam2And4 := types.User{ID: "uFam2And4", FamilyIDs: []string{"family2", "family4"}}
	require.NoError(t, store.Create(ctx, &uFam2And4, nil))

	expectedIds := []string{
		uFam1.ID,
		uFam2.ID,
		uFam1And2.ID,
		uFam2And4.ID,
	}
	sort.Strings(expectedIds)

	filter := bson.M{
		"$or": bson.A{
			bson.M{user.FieldFamilyIds: bson.M{"$eq": "family1"}},
			bson.M{user.FieldFamilyIds: bson.M{"$eq": "family2"}},
		},
	}

	results, err := store.FindAll(ctx, filter)
	require.NoError(t, err)

	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	sort.Strings(ids)

	require.Equal(t, expectedIds, ids)
}
