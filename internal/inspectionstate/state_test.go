package inspectionstate

import "testing"

func TestRetrySuccessCompletesAndIsVisible(t *testing.T) {
	if !Transition(Retrying, Completed) {
		t.Fatal("retrying task cannot complete")
	}
	if got := AfterRetry(true); got != Completed {
		t.Fatalf("status=%s", got)
	}
	if !InProgress(Retrying) {
		t.Fatal("retrying task disappeared from in-progress view")
	}
}
