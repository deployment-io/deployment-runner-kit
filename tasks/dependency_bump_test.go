package tasks

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/deployment-io/deployment-runner-kit/types"
)

func TestValidDependencyEcosystem(t *testing.T) {
	for _, name := range []string{"go", "npm", "cargo", "python-uv", "a" + strings.Repeat("b", 31)} {
		if !ValidDependencyEcosystem(name) {
			t.Errorf("ValidDependencyEcosystem(%q) = false; want true", name)
		}
	}
	for _, name := range []string{"", "Go", "1go", "go!", "a" + strings.Repeat("b", 32), "-go", "go\n"} {
		if ValidDependencyEcosystem(name) {
			t.Errorf("ValidDependencyEcosystem(%q) = true; want false", name)
		}
	}
}

func TestDependencyWireNames(t *testing.T) {
	cases := map[string]string{
		DependencyEcosystemGo:             "go",
		DependencyReleaseCommit:           "commit",
		DependencyReleaseRegistry:         "registry",
		DependencyBumpStateBumping:        "bumping",
		DependencyBumpStateBumped:         "bumped",
		DependencyBumpStateAlreadyCurrent: "already_current",
		DependencyBumpStateNotRequired:    "not_required",
		DependencyBumpStateFailed:         "failed",
		DependencyBumpStateWaitingRelease: "waiting_release",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("wire name %q; want %q", got, want)
		}
	}
	want := []string{"bumped", "already_current", "not_required", "failed"}
	if !reflect.DeepEqual(DependencyBumpResultStates, want) {
		t.Errorf("DependencyBumpResultStates = %v; want %v", DependencyBumpResultStates, want)
	}
}

// The JSON field names are what kit stamps in the DependencyBumpRequest
// parameter and the runner reads; changing one is a wire break.
func TestDependencyBumpRequestJSON(t *testing.T) {
	in := DependencyBumpRequestV1{
		Repository: 1,
		Dependencies: []DependencyBumpTargetV1{{
			Ecosystem:   DependencyEcosystemGo,
			Name:        "github.com/deployment-io/kit",
			Description: "the Go module github.com/deployment-io/kit",
			Source:      0,
			Commit:      strings.Repeat("a", 40),
			Release:     DependencyReleaseCommit,
		}},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 || fields["repository"] == nil || fields["dependencies"] == nil {
		t.Errorf("request fields = %v; want repository, dependencies", fields)
	}
	targets, _ := fields["dependencies"].([]interface{})
	if len(targets) != 1 {
		t.Fatalf("dependencies = %v", fields["dependencies"])
	}
	target, _ := targets[0].(map[string]interface{})
	want := []string{"ecosystem", "name", "description", "source", "commit", "release"}
	if len(target) != len(want) {
		t.Errorf("target fields = %v; want %v", target, want)
	}
	for _, k := range want {
		if _, ok := target[k]; !ok {
			t.Errorf("missing JSON field %q in %s", k, data)
		}
	}
	var out DependencyBumpRequestV1
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip: got %+v, want %+v", out, in)
	}
}

func TestRecordTaskDependenciesArgsGob(t *testing.T) {
	auth := types.AuthArgsV1{OrganizationID: "o", Token: "t"}
	for _, in := range []RecordTaskDependenciesArgsV1{
		{AuthArgsV1: auth, Dependencies: TaskDependenciesDtoV1{TaskID: "task", JobID: "job"}},
		{AuthArgsV1: auth, Dependencies: TaskDependenciesDtoV1{
			TaskID:   "task",
			JobID:    "job",
			Provided: []DependencyProvidedV1{{Repository: 0, Ecosystem: "cargo", Name: "kit"}},
			Required: []DependencyRequiredV1{{Repository: 1, Ecosystem: "cargo", Name: "kit",
				Description: "the crate kit", ProvidedBy: 0, Version: "v1.2.3"}},
		}},
	} {
		var out RecordTaskDependenciesArgsV1
		gobRoundTrip(t, in, &out)
		if !reflect.DeepEqual(in, out) {
			t.Errorf("round trip: got %+v, want %+v", out, in)
		}
	}
	var reply RecordTaskDependenciesReplyV1
	gobRoundTrip(t, RecordTaskDependenciesReplyV1{Done: true}, &reply)
	if !reply.Done {
		t.Error("reply Done lost")
	}
}

func TestUpdateTaskDependencyBumpArgsGob(t *testing.T) {
	auth := types.AuthArgsV1{OrganizationID: "o", Token: "t"}
	for _, in := range []UpdateTaskDependencyBumpArgsV1{
		{AuthArgsV1: auth, Update: UpdateTaskDependencyBumpDtoV1{TaskID: "task", JobID: "job", Repository: 1}},
		{AuthArgsV1: auth, Update: UpdateTaskDependencyBumpDtoV1{
			TaskID:     "task",
			JobID:      "job",
			Repository: 1,
			CommitSHA:  strings.Repeat("b", 40),
			PRURL:      "https://github.com/acme/web/pull/7",
			PRNumber:   7,
			Results: []DependencyBumpResultDtoV1{
				{Ecosystem: "go", Name: "m", State: DependencyBumpStateBumped, FromVersion: "v1", ToVersion: "v2"},
				{Ecosystem: "cargo", Name: "c", State: DependencyBumpStateFailed, Error: "e"},
			},
		}},
	} {
		var out UpdateTaskDependencyBumpArgsV1
		gobRoundTrip(t, in, &out)
		if !reflect.DeepEqual(in, out) {
			t.Errorf("round trip: got %+v, want %+v", out, in)
		}
	}
	var reply UpdateTaskDependencyBumpReplyV1
	gobRoundTrip(t, UpdateTaskDependencyBumpReplyV1{Done: true}, &reply)
	if !reply.Done {
		t.Error("reply Done lost")
	}
}
