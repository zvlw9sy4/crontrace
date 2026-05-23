package parser

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		expr string
	}{
		{"* * * * *"},
		{"0 12 * * *"},
		{"*/15 * * * *"},
		{"0 9-17 * * 1-5"},
		{"30 6 1,15 * *"},
		{"0 0 * * 0"},
		{"5 4 * * sun"},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			c, err := Parse(tc.expr)
			if err != nil {
				t.Errorf("expected no error for %q, got: %v", tc.expr, err)
			}
			if c != nil && c.Raw != tc.expr {
				t.Errorf("Raw mismatch: want %q got %q", tc.expr, c.Raw)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []struct {
		expr    string
		reason string
	}{
		{"* * * *", "only 4 fields"},
		{"60 * * * *", "minute out of range"},
		{"* 24 * * *", "hour out of range"},
		{"* * 0 * *", "day out of range (0)"},
		{"* * * 13 *", "month out of range"},
		{"* * * * 8", "weekday out of range"},
		{"*/0 * * * *", "step value zero"},
		{"1-2-3 * * * *", "malformed range"},
	}

	for _, tc := range cases {
		t.Run(tc.reason, func(t *testing.T) {
			_, err := Parse(tc.expr)
			if err == nil {
				t.Errorf("expected error for %q (%s), got nil", tc.expr, tc.reason)
			}
		})
	}
}

func TestParse_Fields(t *testing.T) {
	c, err := Parse("30 6 15 3 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Minute != "30" || c.Hour != "6" || c.Day != "15" || c.Month != "3" || c.Weekday != "1" {
		t.Errorf("field mismatch: %+v", c)
	}
}
