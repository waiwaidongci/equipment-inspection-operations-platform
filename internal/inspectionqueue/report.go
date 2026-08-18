package inspectionqueue

type Report struct {
	Scheduled []Candidate
}

func BuildReport(queue Queue, deferred []Candidate) Report {
	scheduled := make([]Candidate, 0, len(queue.Items)+len(deferred))
	scheduled = append(scheduled, queue.Items...)
	scheduled = append(scheduled, deferred...)
	return Report{Scheduled: scheduled}
}
