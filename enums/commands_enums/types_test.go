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

func TestDeleteAwsNamespaces(t *testing.T) {
	if DeleteAwsNamespaces != 36 {
		t.Errorf("DeleteAwsNamespaces = %d; want 36", DeleteAwsNamespaces)
	}
	if got := DeleteAwsNamespaces.String(); got != "Delete AWS Namespaces" {
		t.Errorf("DeleteAwsNamespaces.String() = %q; want %q", got, "Delete AWS Namespaces")
	}
}

// RunNewAutomation and RunNewAgent are deprecated (the automations feature is
// gone) but old job documents still store them, so their values must not move.
func TestDeprecatedAutomationCommandValues(t *testing.T) {
	if RunNewAutomation != 24 {
		t.Errorf("RunNewAutomation = %d; want 24", RunNewAutomation)
	}
	if RunNewAgent != 25 {
		t.Errorf("RunNewAgent = %d; want 25", RunNewAgent)
	}
}
