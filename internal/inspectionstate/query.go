package inspectionstate

func InProgress(status Status) bool { return status == Pending || status == Retrying }
