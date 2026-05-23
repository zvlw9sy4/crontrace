package visualizer_test

import (
	"testing"

	"github.com/user/crontrace/internal/visualizer"
)

func TestDescribe_KnownExpressions(t *testing.T) {
	cases := []struct {
		expr string
		want string
	}{
		{"* * * * *", "Every minute"},
		{"0 * * * *", "Every hour"},
		{"0 0 * * *", "Daily at midnight"},
		{"0 12 * * *", "Daily at noon"},
		{"0 0 * * 0", "Weekly on Sunday"},
		{"0 0 1 * *", "Monthly on 1st"},
		{"0 0 1 1 *", "Yearly on Jan 1st"},
		{"*/5 * * * *", "Every 5 minutes"},
		{"*/15 * * * *", "Every 15 minutes"},
		{"*/30 * * * *", "Every 30 minutes"},
		{"0 */2 * * *", "Every 2 hours"},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			got := visualizer.Describe(tc.expr)
			if got != tc.want {
				t.Errorf("Describe(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestDescribe_FallbackCustom(t *testing.T) {
	cases := []string{
		"3 14 1 * *",
		"45 23 * * 5",
		"0 8 15 6 *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			got := visualizer.Describe(expr)
			if got == "" {
				t.Errorf("Describe(%q) returned empty string", expr)
			}
		})
	}
}

func TestDescribe_InvalidFieldCount(t *testing.T) {
	got := visualizer.Describe("* * *")
	if got != "Custom schedule" {
		t.Errorf("expected 'Custom schedule' for malformed expr, got %q", got)
	}
}
