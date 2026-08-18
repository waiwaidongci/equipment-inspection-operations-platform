package inspectionqueue

import "sort"

func Plan(candidates []Candidate, active map[int64]bool) Queue {
	filtered := activeCandidates(candidates, active)
	sort.SliceStable(filtered, func(left, right int) bool {
		return filtered[left].Priority > filtered[right].Priority
	})
	return newQueue(filtered)
}
