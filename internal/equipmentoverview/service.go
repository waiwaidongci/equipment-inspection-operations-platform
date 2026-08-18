package equipmentoverview

type Service struct{ store *Store }

func NewService(store *Store) *Service { return &Service{store: store} }

func (service *Service) Current() map[string]DeviceState {
	return service.store.Snapshot()
}
