package logs

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"
)

func TestLogTagConstants(t *testing.T) {
	cases := map[string]string{
		LogStageSetup:     "setup",
		LogStageImplement: "implement",
		LogStageReview:    "review",
		LogStageQA:        "qa",
		LogStagePR:        "pr",
		LogKindSection:    "section",
		LogKindText:       "text",
		LogKindTool:       "tool",
		LogKindResult:     "result",
		LogKindJSON:       "json",
		LogKindDone:       "done",
		LogKindError:      "error",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("constant = %q, want %q", got, want)
		}
	}
}

func gobRoundTrip(t *testing.T, in any, out any) {
	t.Helper()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := gob.NewDecoder(&buf).Decode(out); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestAddJobLogDtoV1GobRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		in   AddJobLogDtoV1
	}{
		{"untagged", AddJobLogDtoV1{ID: "abc", Message: "hello", Ts: 42}},
		{"tagged", AddJobLogDtoV1{ID: "abc", Message: "pulling", Ts: 42, Stage: LogStageSetup,
			Kind: LogKindText, Group: &JobLogGroupV1{Name: "docker progress", Count: 12}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := AddJobLogsArgsV1{JobLogs: []AddJobLogDtoV1{tt.in}}
			var out AddJobLogsArgsV1
			gobRoundTrip(t, in, &out)
			if !reflect.DeepEqual(in, out) {
				t.Fatalf("round trip = %+v, want %+v", out, in)
			}
			if tt.in.Group == nil && out.JobLogs[0].Group != nil {
				t.Fatalf("nil Group decoded as %+v", out.JobLogs[0].Group)
			}
		})
	}
}

func TestAddJobLogDtoV1DecodesOldPayload(t *testing.T) {
	// The struct as sent by runners before tags existed; gob matches by name.
	type AddJobLogDtoV1Old struct {
		ID           string
		Message      string
		ErrorMessage string
		Ts           int64
	}
	old := AddJobLogDtoV1Old{ID: "abc", Message: "hello", ErrorMessage: "boom", Ts: 7}
	var out AddJobLogDtoV1
	gobRoundTrip(t, old, &out)
	want := AddJobLogDtoV1{ID: "abc", Message: "hello", ErrorMessage: "boom", Ts: 7}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("decoded = %+v, want %+v", out, want)
	}
}
