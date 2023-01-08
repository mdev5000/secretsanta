package user

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/mongo"
)

type Service struct {
	store store
}

func NewService(userCollection *mongo.Collection) *Service {
	userStore := store{
		users: userCollection,
	}

	return &Service{
		store: userStore,
	}
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.store.Count(ctx)
}

func (s *Service) Login(ctx context.Context, username string, password []byte) (*User, error) {
	u, err := s.store.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := comparePassword(u.PasswordHash, password); err != nil {
		return nil, err
	}
	// Purge password from memory.
	for i := range password {
		password[i] = 0
	}
	return u, nil
}

func (s *Service) Create(ctx context.Context, u *User, password []byte) error {
	passwordHash, err := hashPassword(password)
	// Purge password from memory.
	for i := range password {
		password[i] = 0
	}
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return s.store.Create(ctx, u, passwordHash)
}

func (s *Service) FindByID(ctx context.Context, id ID) (*User, error) {
	return s.store.FindByID(ctx, id)
}
