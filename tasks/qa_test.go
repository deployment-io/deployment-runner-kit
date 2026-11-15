package tasks

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/deployment-io/deployment-runner-kit/types"
)

// The QA wire names are read by kit's qa_enums wire-mirror test and by
// deployment-server; changing one is a wire break.
func TestQAWireNames(t *testing.T) {
	cases := map[string]string{
		StageQA:               "qa",
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

func gobRoundTrip(t *testing.T, in, out interface{}) {
	t.Helper()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatal(err)
	}
	if err := gob.NewDecoder(&buf).Decode(out); err != nil {
		t.Fatal(err)
	}
}

func TestQADtosGobRoundTrip(t *testing.T) {
	for _, ids := range [][]string{nil, {"0123456789abcdef01234567", "76543210fedcba9876543210"}} {
		in := UpdateTaskStepQAArgsV1{Update: UpdateTaskStepQADtoV1{
			TaskID: "t", JobID: "j",
			Result: QAResultDtoV1{Status: QAStatusCompleted, Criteria: []QACriterionDtoV1{
				{Criterion: "c", State: QACriterionPassed, Reason: "r", EvidenceIDs: ids},
			}},
		}}
		var out UpdateTaskStepQAArgsV1
		gobRoundTrip(t, in, &out)
		if !reflect.DeepEqual(in, out) {
			t.Errorf("round trip: got %+v, want %+v", out, in)
		}
	}
}

// legacyQACriterionDtoV1 is QACriterionDtoV1 before EvidenceIDs: an older
// peer on either side of the wire.
type legacyQACriterionDtoV1 struct {
	Criterion string
	State     string
	Reason    string
}

func TestQACriterionEvidenceIDsAreCompatibleWithOlderPeers(t *testing.T) {
	var old legacyQACriterionDtoV1
	gobRoundTrip(t, QACriterionDtoV1{Criterion: "c", State: QACriterionFailed, EvidenceIDs: []string{"x"}}, &old)
	if old.Criterion != "c" || old.State != QACriterionFailed {
		t.Errorf("older server decode: %+v", old)
	}
	var cur QACriterionDtoV1
	gobRoundTrip(t, legacyQACriterionDtoV1{Criterion: "c", Reason: "r"}, &cur)
	if cur.Criterion != "c" || cur.Reason != "r" || cur.EvidenceIDs != nil {
		t.Errorf("older runner decode: %+v", cur)
	}
}

func TestUploadQAEvidenceDtosGobRoundTrip(t *testing.T) {
	in := UploadQAEvidenceArgsV1{
		AuthArgsV1: types.AuthArgsV1{OrganizationID: "o", Token: "k"},
		Upload: UploadQAEvidenceDtoV1{
			TaskID: "t", StepIndex: 2, JobID: "j", Round: 1, CheckID: "c1",
			Criteria: []int{0, 3}, Image: []byte{0xFF, 0xD8, 0xFF, 0xE0}, Width: 800, Height: 600,
		},
	}
	var out UploadQAEvidenceArgsV1
	gobRoundTrip(t, in, &out)
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip: got %+v, want %+v", out, in)
	}
	var reply UploadQAEvidenceReplyV1
	gobRoundTrip(t, UploadQAEvidenceReplyV1{EvidenceID: "abc"}, &reply)
	if reply.EvidenceID != "abc" {
		t.Errorf("reply = %+v", reply)
	}
}
