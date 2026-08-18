package taskfanout

import (
	"errors"
	"testing"
)

func TestRunClosesAfterError(t *testing.T) {
	_, err := Run([]string{"ok", "bad"}, func(value string) error {
		if value == "bad" {
			return errors.New("bad")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected processing error")
	}
}

func TestRunMixedResultsReturnsErrorAndSuccess(t *testing.T) {
	out, err := Run([]string{"a", "bad", "c"}, func(value string) error {
		if value == "bad" {
			return errors.New("bad data")
		}
		return nil
	})
	if err == nil || err.Error() != "bad data" {
		t.Fatalf("expected 'bad data' error, got %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 successful results, got %d: %v", len(out), out)
	}
}

func TestRunEmptyInputs(t *testing.T) {
	out, err := Run([]string{}, func(value string) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %v", out)
	}
}

func TestRunAllSucceed(t *testing.T) {
	out, err := Run([]string{"a", "b", "c"}, func(value string) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d: %v", len(out), out)
	}
}
