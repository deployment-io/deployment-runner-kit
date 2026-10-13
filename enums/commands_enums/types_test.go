package commands_enums

import "testing"

func TestRunVerifyStage(t *testing.T) {
	if RunVerifyStage != 34 {
		t.Errorf("RunVerifyStage = %d; want 34", RunVerifyStage)
	}
	if got := RunVerifyStage.String(); got != "Run Verify Stage" {
		t.Errorf("RunVerifyStage.String() = %q; want %q", got, "Run Verify Stage")
	}
	if RunReviewStage != 33 {
		t.Errorf("RunReviewStage = %d; want 33", RunReviewStage)
	}
}
