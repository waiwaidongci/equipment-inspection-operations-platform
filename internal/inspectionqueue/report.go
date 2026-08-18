package inspectionqueue

type Report struct {
	Scheduled []Candidate
}

func BuildReport(queue Queue, deferred []Candidate) Report {
	scheduled := queue.Items
	scheduled = append(scheduled, deferred...)
	return Report{Scheduled: scheduled}
}
