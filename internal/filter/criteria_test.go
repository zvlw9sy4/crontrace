package filter

import (
	"testing"
	"time"
)

func TestCriteria_IsZero_Empty(t *testing.T) {
	c := Criteria{}
	if !c.IsZero() {
		t.Error("expected empty Criteria to be zero")
	}
}

func TestCriteria_IsZero_LabelSet(t *testing.T) {
	c := Criteria{LabelContains: "backup"}
	if c.IsZero() {
		t.Error("expected Criteria with LabelContains to be non-zero")
	}
}

func TestCriteria_IsZero_ExpressionPrefixSet(t *testing.T) {
	c := Criteria{ExpressionPrefix: "*/5"}
	if c.IsZero() {
		t.Error("expected Criteria with ExpressionPrefix to be non-zero")
	}
}

func TestCriteria_IsZero_AfterSet(t *testing.T) {
	c := Criteria{After: time.Now()}
	if c.IsZero() {
		t.Error("expected Criteria with After to be non-zero")
	}
}

func TestCriteria_IsZero_BeforeSet(t *testing.T) {
	c := Criteria{Before: time.Now()}
	if c.IsZero() {
		t.Error("expected Criteria with Before to be non-zero")
	}
}

func TestCriteria_IsZero_OnlyConflicting(t *testing.T) {
	c := Criteria{OnlyConflicting: true}
	if c.IsZero() {
		t.Error("expected Criteria with OnlyConflicting=true to be non-zero")
	}
}

func TestCriteria_IsZero_MaxNextRuns(t *testing.T) {
	c := Criteria{MaxNextRuns: 5}
	if c.IsZero() {
		t.Error("expected Criteria with MaxNextRuns>0 to be non-zero")
	}
}

func TestCriteria_IsZero_MaxNextRunsZero(t *testing.T) {
	c := Criteria{MaxNextRuns: 0}
	if !c.IsZero() {
		t.Error("expected Criteria with MaxNextRuns=0 to be zero")
	}
}

func TestCriteria_IsZero_MaxNextRunsNegative(t *testing.T) {
	c := Criteria{MaxNextRuns: -1}
	if !c.IsZero() {
		t.Error("expected Criteria with MaxNextRuns=-1 to be zero")
	}
}
