package stores

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/family"
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func Test_canCreateNewFamilies(t *testing.T) {
	ctx := context.Background()
	db := acquireDb(t)
	col := db.Collection(family.CollectionFamilies)
	deleteAll(t, ctx, col)

	svc := family.NewService(col)
	f := types.Family{
		Name:        "Family 1",
		Description: "some description",
	}
	called := time.Now().UTC()
	err := svc.Create(ctx, &f)
	require.NoError(t, err)

	require.NotNil(t, f.ID)
	require.True(t, called.Before(f.UpdatedAt))

	numDocs, err := col.CountDocuments(ctx, bson.D{})
	require.NoError(t, err)
	require.Equal(t, int64(1), numDocs)

	// @todo test later with find
	//newFamily, err := svc.FindByID(ctx, f.ID)
	//require.NoError(t, err)
	//
	//compare.Equal(t, newFamily, &f, compare.IgnoreFields(types.Family{}, "UpdatedAt"))
}
