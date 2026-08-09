package telegram

import "sync"

type PurchaseState struct {
	PlanID    string
	Locations map[string]string
}

type StateStorage struct {
	mu    sync.RWMutex
	store map[int64]*PurchaseState
}

func NewStorageState() *StateStorage {
	return &StateStorage{store: make(map[int64]*PurchaseState)}
}

func (s *StateStorage) Set(userID int64, state *PurchaseState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = state
}

func (s *StateStorage) Get(userID int64) *PurchaseState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[userID]
}

func (s *StateStorage) Clear(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
}
