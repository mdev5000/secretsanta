package user

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/types"
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
	FieldFamilyIds    = "familyIDs"
)

type Store struct {
	Users *mongo.Collection
}

func (s *Store) FindByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.findOneByField(ctx, FieldUsername, username)
}

func (s *Store) FindByID(ctx context.Context, id types.ID) (*types.User, error) {
	return s.findOneByField(ctx, FieldID, id)
}

func (s *Store) findOneByField(ctx context.Context, field string, value interface{}) (*types.User, error) {
	r := s.Users.FindOne(ctx, bson.D{
		{field, value},
	}, options.FindOne())
	if r.Err() != nil {
		return nil, r.Err()
	}
	var user types.User
	if err := r.Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	return &user, nil
}

func (s *Store) FindAll(ctx context.Context, filter bson.M) ([]*types.User, error) {
	r, err := s.Users.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	var users []*types.User
	err = r.All(ctx, &users)
	return users, err
}

func (s *Store) CreateWithNewId(ctx context.Context, u *types.User, passwordHash []byte) error {
	id := uuid.New().String()
	u.ID = id
	return s.Create(ctx, u, passwordHash)
}

func (s *Store) Create(ctx context.Context, u *types.User, passwordHash []byte) error {
	updatedAt := time.Now().UTC()
	_, err := s.Users.InsertOne(ctx, bson.D{
		{FieldID, u.ID},
		{FieldUsername, u.Username},
		{FieldPasswordHash, passwordHash},
		{FieldFirstname, u.Firstname},
		{FieldLastname, u.Lastname},
		{FieldFamilyIds, u.FamilyIDs},
		{FieldUpdatedAt, updatedAt},
	})
	if err != nil {
		return err
	}
	u.UpdatedAt = updatedAt
	return nil
}

func (s *Store) Count(ctx context.Context) (int64, error) {
	return s.Users.CountDocuments(ctx, bson.D{})
}
