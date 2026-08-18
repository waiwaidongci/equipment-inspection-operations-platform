package equipmentoverview

type Cache struct{ states map[string]DeviceState }

func cloneStates(states map[string]DeviceState) map[string]DeviceState {
	copyOfStates := make(map[string]DeviceState, len(states))
	for code, state := range states {
		copyOfStates[code] = state
	}
	return copyOfStates
}

func (cache *Cache) Save(states map[string]DeviceState) {
	cache.states = cloneStates(states)
}

func (cache *Cache) Load() map[string]DeviceState {
	return cloneStates(cache.states)
}
