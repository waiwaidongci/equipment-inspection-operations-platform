package inspectionstate

func Visible(status Status) bool {
	return status == Pending || status == Retrying || status == Completed || status == Failed
}
