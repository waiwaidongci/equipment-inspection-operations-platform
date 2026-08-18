package inspectionstate

func Visible(status Status) bool { return status != Retrying }
