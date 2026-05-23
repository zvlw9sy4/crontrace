package history

import (
	"strings"
	"testing"
	"time"
)

func TestRenderHistoryTable_HeaderPresent(t *testing.T) {
	var sb strings.Builder
	RenderHistoryTable(&sb, nil)
	out := sb.String()
	for _, col := range []string{"LABEL", "EXPRESSION", "RAN AT", "DURATION", "STATUS", "MESSAGE"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestRenderHistoryTable_ShowsEntries(t *testing.T) {
	var sb strings.Builder
	entries := []Entry{
		{Label: "myjob", Expression: "0 * * * *", RanAt: baseTime, Duration: 500 * time.Millisecond, Success: true},
		{Label: "failjob", Expression: "*/5 * * * *", RanAt: baseTime.Add(time.Minute), Duration: time.Second, Success: false, Message: "timeout"},
	}
	RenderHistoryTable(&sb, entries)
	out := sb.String()
	if !strings.Contains(out, "myjob") {
		t.Error("expected myjob in output")
	}
	if !strings.Contains(out, "FAIL") {
		t.Error("expected FAIL status in output")
	}
	if !strings.Contains(out, "timeout") {
		t.Error("expected message 'timeout' in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK status in output")
	}
}

func TestRenderHistoryTable_DashForEmptyMessage(t *testing.T) {
	var sb strings.Builder
	entries := []Entry{
		{Label: "job", Expression: "* * * * *", RanAt: baseTime, Duration: time.Second, Success: true},
	}
	RenderHistoryTable(&sb, entries)
	if !strings.Contains(sb.String(), "-") {
		t.Error("expected dash for empty message")
	}
}

func TestRenderSuccessSummary_Output(t *testing.T) {
	s := NewStore()
	_ = s.Record(makeEntry("alpha", "0 * * * *", true, 0))
	_ = s.Record(makeEntry("alpha", "0 * * * *", false, time.Minute))
	_ = s.Record(makeEntry("beta", "*/5 * * * *", true, 0))

	var sb strings.Builder
	RenderSuccessSummary(&sb, s)
	out := sb.String()

	if !strings.Contains(out, "alpha") {
		t.Error("expected alpha in summary")
	}
	if !strings.Contains(out, "50.0%") {
		t.Error("expected 50.0% success rate for alpha")
	}
	if !strings.Contains(out, "beta") {
		t.Error("expected beta in summary")
	}
	if !strings.Contains(out, "100.0%") {
		t.Error("expected 100.0% for beta")
	}
}
