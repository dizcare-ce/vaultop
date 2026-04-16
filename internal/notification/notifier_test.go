package notification_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vaultop/vaultop/internal/notification"
)

func baseEvent(success bool) notification.Event {
	return notification.Event{
		SecretKey: "db/password",
		Provider:  "aws",
		Success:   success,
		RotatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestWriterNotifier_SuccessEvent(t *testing.T) {
	var buf bytes.Buffer
	n := notification.NewWriterNotifier(&buf)

	if err := n.Notify(baseEvent(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"provider=aws", "key=db/password", "status=OK"} {
		if !strings.Contains(line, want) {
			t.Errorf("expected %q in output %q", want, line)
		}
	}
}

func TestWriterNotifier_FailureEvent(t *testing.T) {
	var buf bytes.Buffer
	n := notification.NewWriterNotifier(&buf)

	e := baseEvent(false)
	e.Error = errors.New("permission denied")

	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"status=FAILED", "permission denied"} {
		if !strings.Contains(line, want) {
			t.Errorf("expected %q in output %q", want, line)
		}
	}
}

func TestWriterNotifier_FailureNoError(t *testing.T) {
	var buf bytes.Buffer
	n := notification.NewWriterNotifier(&buf)

	if err := n.Notify(baseEvent(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "status=FAILED") {
		t.Error("expected FAILED status in output")
	}
}

func TestNoop_Notify(t *testing.T) {
	var n notification.Noop
	if err := n.Notify(baseEvent(true)); err != nil {
		t.Fatalf("noop notifier returned error: %v", err)
	}
}
