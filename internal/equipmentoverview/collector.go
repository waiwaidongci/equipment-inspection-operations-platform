package equipmentoverview

type Totals struct {
	Active   int
	Inactive int
}

func Collect(states map[string]DeviceState) Totals {
	owned := make([]DeviceState, 0, len(states))
	for _, state := range states {
		owned = append(owned, state)
	}
	var totals Totals
	for _, state := range owned {
		if state.Status == "active" {
			totals.Active++
		} else {
			totals.Inactive++
		}
	}
	return totals
}
