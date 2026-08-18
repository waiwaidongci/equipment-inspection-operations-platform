package inspectionstate

type Status string

const (
	Pending   Status = "pending"
	Retrying  Status = "retrying"
	Completed Status = "completed"
	Failed    Status = "failed"
)

func Transition(from, to Status) bool {
	switch from {
	case Pending:
		return to == Retrying || to == Failed
	case Retrying:
		return to == Completed || to == Failed
	case Failed:
		return to == Retrying
	default:
		return false
	}
}
