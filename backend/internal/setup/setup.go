package setup

import "sync"

type Service struct {
	setupMutex sync.RWMutex
	isSetup    bool
}

// IsSetup determines if the backend has been set up or if the setup page should be shown.
func (s *Service) IsSetup() bool {
	s.setupMutex.RLock()
	if s.isSetup {
		s.setupMutex.RUnlock()
		return true
	}
	s.setupMutex.RUnlock()

	// @todo check database
	var dbResult bool
	dbResult = true
	if dbResult {
		s.setupMutex.Lock()
		s.isSetup = true
		s.setupMutex.Unlock()
		return true
	}

	return false
}
