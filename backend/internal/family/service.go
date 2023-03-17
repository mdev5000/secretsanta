package family

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/types"
)

type Service struct {
	store *store
}

func NewService(familyColl *mongo.Collection) *Service {
	s := store{coll: familyColl}
	return &Service{
		store: &s,
	}
}

func (s *Service) Create(ctx context.Context, family *types.Family) error {
	// @todo validate the family
	if err := s.store.Create(ctx, family); err != nil {
		return fmt.Errorf("failed to create new family")
	}
	return nil
}
