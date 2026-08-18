package inspectionqueue

import "sort"

func Plan(candidates []Candidate, active map[int64]bool) Queue {
	filtered := activeCandidates(append([]Candidate(nil), candidates...), active)
	sort.Slice(filtered, func(left, right int) bool {
		return filtered[left].Priority > filtered[right].Priority
	})
	return newQueue(filtered)
}
