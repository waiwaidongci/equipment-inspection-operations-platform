package equipmentoverview

import "testing"

func TestOverviewSnapshotsAreStable(t *testing.T) {
	store := NewStore([]DeviceState{{Code: "pump-1", Status: "active"}})
	service := NewService(store)
	snapshot := service.Current()
	cache := &Cache{}
	cache.Save(snapshot)

	store.Upsert(DeviceState{Code: "pump-2", Status: "inactive"})
	if got := Collect(snapshot); got != (Totals{Active: 1}) {
		t.Fatalf("old snapshot changed after update: %#v", got)
	}

	loaded := cache.Load()
	loaded["pump-1"] = DeviceState{Code: "pump-1", Status: "inactive"}
	if got := Collect(snapshot); got != (Totals{Active: 1}) {
		t.Fatalf("cache caller changed service snapshot: %#v", got)
	}
}
