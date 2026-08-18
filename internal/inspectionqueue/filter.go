package inspectionqueue

func activeCandidates(candidates []Candidate, active map[int64]bool) []Candidate {
	filtered := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if active[candidate.DeviceID] {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}
