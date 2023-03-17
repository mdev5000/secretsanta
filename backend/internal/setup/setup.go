package setup

import (
	"context"
	"fmt"
	"github.com/mdev5000/secretsanta/internal/types"
	"sync"
)

type TransactionMgr interface {
	WithTransaction(context.Context, func(ctx context.Context) error) error
}

type UserService interface {
	Create(ctx context.Context, admin *types.User, password []byte) error
	Count(context.Context) (int64, error)
}

type FamilyService interface {
	Create(ctx context.Context, family *types.Family) error
}

type Service struct {
	users          UserService
	family         FamilyService
	transactionMgr TransactionMgr
	setupMutex     sync.RWMutex
	isSetup        bool
}

type Data struct {
	DefaultAdmin         *types.User
	DefaultAdminPassword []byte
	DefaultFamily        *types.Family
}

func NewService(transactionMgr TransactionMgr, users UserService, family FamilyService) *Service {
	return &Service{
		transactionMgr: transactionMgr,
		users:          users,
		family:         family,
	}
}

func (s *Service) Setup(ctx context.Context, data *Data) error {
	err := s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.users.Create(ctx, data.DefaultAdmin, data.DefaultAdminPassword); err != nil {
			return err
		}
		if err := s.family.Create(ctx, data.DefaultFamily); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to run setup: %w", err)
	}
	return nil
}

// IsSetup determines if the backend has been set up or if the setup page should be shown.
func (s *Service) IsSetup(ctx context.Context) (bool, error) {
	s.setupMutex.RLock()
	if s.isSetup {
		s.setupMutex.RUnlock()
		return true, nil
	}
	s.setupMutex.RUnlock()

	userCount, err := s.users.Count(ctx)
	if err != nil {
		return false, err
	}

	if userCount > 0 {
		s.setupMutex.Lock()
		s.isSetup = true
		s.setupMutex.Unlock()
		return true, nil
	}

	return false, nil
}
