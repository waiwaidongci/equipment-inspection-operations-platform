package inspectionqueue

import (
	"reflect"
	"testing"
)

func TestQueuePlanningKeepsSourceAndQueueIndependent(t *testing.T) {
	source := []Candidate{{DeviceID: 1, Priority: 10}, {DeviceID: 2, Priority: 30}, {DeviceID: 3, Priority: 20}}
	wantSource := append([]Candidate(nil), source...)
	queue := Plan(source, map[int64]bool{1: true, 3: true})

	if !reflect.DeepEqual(source, wantSource) {
		t.Fatalf("planning changed source candidates: got %#v want %#v", source, wantSource)
	}
	if want := []Candidate{{DeviceID: 3, Priority: 20}, {DeviceID: 1, Priority: 10}}; !reflect.DeepEqual(queue.Items, want) {
		t.Fatalf("unexpected queue: got %#v want %#v", queue.Items, want)
	}

	report := BuildReport(queue, []Candidate{{DeviceID: 9, Priority: 1}})
	report.Scheduled[0].Priority = 99
	if queue.Items[0].Priority != 20 {
		t.Fatalf("report mutation leaked into queue: got %d", queue.Items[0].Priority)
	}
}
