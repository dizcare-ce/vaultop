package lease

import (
	"strings"
	"testing"
	"time"
)

func TestReport_ActiveOnly(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	leases := []Lease{
		{Key: "api/key", IssuedAt: now, Duration: time.Hour, Renewals: 0},
		{Key: "db/pass", IssuedAt: now.Add(-2 * time.Hour), Duration: time.Hour, Renewals: 1},
	}
	var buf strings.Builder
	err := Report(&buf, leases, ReportOptions{Now: now})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api/key") {
		t.Error("expected api/key in output")
	}
	if strings.Contains(out, "db/pass") {
		t.Error("expired lease db/pass should be excluded")
	}
}

func TestReport_ShowExpired_IncludesAll(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	leases := []Lease{
		{Key: "api/key", IssuedAt: now, Duration: time.Hour},
		{Key: "db/pass", IssuedAt: now.Add(-2 * time.Hour), Duration: time.Hour},
	}
	var buf strings.Builder
	err := Report(&buf, leases, ReportOptions{Now: now, ShowExpired: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api/key") {
		t.Error("expected api/key in output")
	}
	if !strings.Contains(out, "db/pass") {
		t.Error("expected db/pass in output when ShowExpired=true")
	}
	if !strings.Contains(out, "expired") {
		t.Error("expected 'expired' status in output")
	}
}

func TestReport_SortedByKey(t *testing.T) {
	now := time.Now()
	leases := []Lease{
		{Key: "z/secret", IssuedAt: now, Duration: time.Hour},
		{Key: "a/token", IssuedAt: now, Duration: time.Hour},
		{Key: "m/cert", IssuedAt: now, Duration: time.Hour},
	}
	var buf strings.Builder
	_ = Report(&buf, leases, ReportOptions{Now: now})
	out := buf.String()
	posA := strings.Index(out, "a/token")
	posM := strings.Index(out, "m/cert")
	posZ := strings.Index(out, "z/secret")
	if posA > posM || posM > posZ {
		t.Errorf("leases not sorted alphabetically: positions a=%d m=%d z=%d", posA, posM, posZ)
	}
}

func TestReport_RenewalsVisible(t *testing.T) {
	now := time.Now()
	leases := []Lease{
		{Key: "k", IssuedAt: now, Duration: time.Hour, Renewals: 3},
	}
	var buf strings.Builder
	_ = Report(&buf, leases, ReportOptions{Now: now})
	out := buf.String()
	if !strings.Contains(out, "3") {
		t.Error("expected renewal count 3 in output")
	}
}
