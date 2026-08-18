package equipmentoverview

type Cache struct{ states map[string]DeviceState }

func (cache *Cache) Save(states map[string]DeviceState) {
	cache.states = states
}

func (cache *Cache) Load() map[string]DeviceState {
	return cache.states
}
