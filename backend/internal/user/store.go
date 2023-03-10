package user

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/google/uuid"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionUsers = "users"

	FieldID           = "id"
	FieldUsername     = "username"
	FieldFirstname    = "firstname"
	FieldLastname     = "lastname"
	FieldPasswordHash = "passwordHash"
	FieldUpdatedAt    = "updatedAt"
)

type store struct {
	users *mongo.Collection
}

func (s *store) FindByUsername(ctx context.Context, username string) (*User, error) {
	return s.findOneByField(ctx, FieldUsername, username)
}

func (s *store) FindByID(ctx context.Context, id ID) (*User, error) {
	return s.findOneByField(ctx, FieldID, id)
}

func (s *store) findOneByField(ctx context.Context, field string, value interface{}) (*User, error) {
	r := s.users.FindOne(ctx, bson.D{
		{field, value},
	}, options.FindOne())
	if r.Err() != nil {
		return nil, r.Err()
	}
	var user User
	if err := r.Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	return &user, nil
}

func (s *store) Create(ctx context.Context, u *User, passwordHash []byte) error {
	id := uuid.New()
	updatedAt := time.Now().UTC()
	_, err := s.users.InsertOne(ctx, bson.D{
		{FieldID, id},
		{FieldUsername, u.Username},
		{FieldPasswordHash, passwordHash},
		{FieldFirstname, u.Firstname},
		{FieldLastname, u.Lastname},
		{FieldUpdatedAt, updatedAt},
	})
	if err != nil {
		return err
	}
	u.ID = id
	u.UpdatedAt = updatedAt
	return nil
}

func (s *store) Count(ctx context.Context) (int64, error) {
	return s.users.CountDocuments(ctx, bson.D{})
}
