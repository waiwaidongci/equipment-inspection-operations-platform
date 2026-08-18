package inspectionstate

func AfterRetry(success bool) Status {
	if success {
		return Completed
	}
	return Failed
}
