package equipmentoverview

type Cache struct{ states map[string]DeviceState }

func copyStates(states map[string]DeviceState) map[string]DeviceState {
	copied := make(map[string]DeviceState, len(states))
	for code, state := range states {
		copied[code] = state
	}
	return copied
}

func (cache *Cache) Save(states map[string]DeviceState) {
	cache.states = copyStates(states)
}

func (cache *Cache) Load() map[string]DeviceState {
	return copyStates(cache.states)
}
