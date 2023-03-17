package family

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	CollectionFamilies = "families"

	FieldID          = "id"
	FieldName        = "name"
	FieldDescription = "description"
	FieldUpdatedAt   = "updatedAt"
)

type store struct {
	coll *mongo.Collection
}

func (s *store) Create(ctx context.Context, f *types.Family) error {
	id := uuid.New().String()
	updatedAt := time.Now().UTC()

	_, err := s.coll.InsertOne(ctx, bson.D{
		{FieldID, id},
		{FieldName, f.Name},
		{FieldDescription, f.Description},
		{FieldUpdatedAt, updatedAt},
	})
	if err != nil {
		return fmt.Errorf("failed to create family: %w", err)
	}
	f.ID = id
	f.UpdatedAt = updatedAt
	return nil

}
