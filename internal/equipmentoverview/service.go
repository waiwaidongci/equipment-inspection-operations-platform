package equipmentoverview

type Service struct{ store *Store }

func NewService(store *Store) *Service { return &Service{store: store} }

func (service *Service) Current() map[string]DeviceState {
	states := service.store.Snapshot()
	owned := make(map[string]DeviceState, len(states))
	for code, state := range states {
		owned[code] = state
	}
	return owned
}
