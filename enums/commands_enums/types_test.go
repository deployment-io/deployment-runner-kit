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

func TestRunProductionQA(t *testing.T) {
	if RunProductionQA != 35 {
		t.Errorf("RunProductionQA = %d; want 35", RunProductionQA)
	}
	if got := RunProductionQA.String(); got != "Run Production QA" {
		t.Errorf("RunProductionQA.String() = %q; want %q", got, "Run Production QA")
	}
}
