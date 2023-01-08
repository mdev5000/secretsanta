package setup

import (
	"context"
	"github.com/mdev5000/secretsanta/internal/user"
	"sync"
)

type Service struct {
	users      *user.Service
	setupMutex sync.RWMutex
	isSetup    bool
}

type Data struct {
	DefaultAdmin         *user.User
	DefaultAdminPassword []byte
	DefaultFamily        string
}

func NewSetupService(users *user.Service) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) Setup(ctx context.Context, data Data) error {
	return s.users.Create(ctx, data.DefaultAdmin, data.DefaultAdminPassword)
	// @todo setup default family
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
