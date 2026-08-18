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
