package equipmentoverview

type Totals struct {
	Active   int
	Inactive int
}

func Collect(states map[string]DeviceState) Totals {
	var totals Totals
	for _, state := range states {
		if state.Status == "active" {
			totals.Active++
		} else {
			totals.Inactive++
		}
	}
	return totals
}
