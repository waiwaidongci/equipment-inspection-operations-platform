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
	store.mu.RLock()
	defer store.mu.RUnlock()
	copyOfDevices := make(map[string]DeviceState, len(store.devices))
	for code, state := range store.devices {
		copyOfDevices[code] = state
	}
	return copyOfDevices
}

func (store *Store) Upsert(state DeviceState) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.devices[state.Code] = state
}
