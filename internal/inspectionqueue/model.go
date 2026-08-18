package inspectionqueue

type Candidate struct {
	DeviceID int64
	Priority int
}

type Queue struct {
	Items []Candidate
}

func newQueue(items []Candidate) Queue {
	owned := append([]Candidate(nil), items...)
	return Queue{Items: owned}
}
