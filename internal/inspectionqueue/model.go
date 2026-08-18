package inspectionqueue

type Candidate struct {
	DeviceID int64
	Priority int
}

type Queue struct {
	Items []Candidate
}

func newQueue(items []Candidate) Queue {
	return Queue{Items: items}
}
