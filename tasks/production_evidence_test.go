package tasks

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
)

// The JSON field names are what kit stamps in the ProductionEvidence
// parameter and the runner reads; changing one is a wire break.
func TestProductionEvidenceJSON(t *testing.T) {
	in := ProductionEvidenceV1{
		SourceTaskID: "t",
		Services: []ProductionEvidenceServiceV1{{
			Name: "api", DeploymentID: "d", BuildID: "b", Region: region_enums.EuWest1, From: 10, To: 20,
		}},
		Skipped: []ProductionEvidenceSkippedV1{{Name: "web", Reason: "r"}},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	checkFields := func(fields map[string]interface{}, want []string) {
		t.Helper()
		if len(fields) != len(want) {
			t.Errorf("fields = %v; want %v", fields, want)
		}
		for _, k := range want {
			if _, ok := fields[k]; !ok {
				t.Errorf("missing JSON field %q in %v", k, fields)
			}
		}
	}
	checkFields(fields, []string{"sourceTaskId", "services", "skipped"})
	checkFields(fields["services"].([]interface{})[0].(map[string]interface{}),
		[]string{"name", "deploymentId", "buildId", "region", "from", "to"})
	checkFields(fields["skipped"].([]interface{})[0].(map[string]interface{}), []string{"name", "reason"})

	var out ProductionEvidenceV1
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip: got %+v, want %+v", out, in)
	}
}
