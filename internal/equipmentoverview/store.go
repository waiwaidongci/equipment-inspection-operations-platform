package equipmentoverview

import "sync"

type DeviceState struct {
	Code   string
	Status string
}

type Store struct {
	mu      sync.RWMutex
	devices map[string]DeviceState
}

func NewStore(states []DeviceState) *Store {
	store := &Store{devices: make(map[string]DeviceState, len(states))}
	for _, state := range states {
		store.devices[state.Code] = state
	}
	return store
}

func (store *Store) Snapshot() map[string]DeviceState {
	return store.devices
}

func (store *Store) Upsert(state DeviceState) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.devices[state.Code] = state
}
