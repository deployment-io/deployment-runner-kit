package commands_enums

import "testing"

func TestRunQAStage(t *testing.T) {
	if RunQAStage != 34 {
		t.Errorf("RunQAStage = %d; want 34", RunQAStage)
	}
	if got := RunQAStage.String(); got != "Run QA Stage" {
		t.Errorf("RunQAStage.String() = %q; want %q", got, "Run QA Stage")
	}
	if RunReviewStage != 33 {
		t.Errorf("RunReviewStage = %d; want 33", RunReviewStage)
	}
}
