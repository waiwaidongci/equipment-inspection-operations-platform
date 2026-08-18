package inspectionstate

func AfterRetry(success bool) Status {
	if success {
		return Retrying
	}
	return Failed
}
