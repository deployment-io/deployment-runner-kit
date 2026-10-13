package tasks

import "testing"

// The QA wire names are read by kit's qa_enums wire-mirror test and by
// deployment-server; changing one is a wire break.
func TestQAWireNames(t *testing.T) {
	cases := map[string]string{
		StageVerify:           "qa",
		QACriterionPassed:     "passed",
		QACriterionFailed:     "failed",
		QACriterionNotChecked: "not_checked",
		QAStatusCompleted:     "completed",
		QAStatusSkipped:       "skipped",
		QAStatusCouldNotRun:   "couldnt_run",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("wire name = %q; want %q", got, want)
		}
	}
}
