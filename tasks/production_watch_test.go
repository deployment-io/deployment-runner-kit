package tasks

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/deployment-io/deployment-runner-kit/enums/region_enums"
)

// The JSON field names are what app-server stamps in the
// ProductionWatchTargets parameter and the runner reads; changing one is a
// wire break.
func TestProductionWatchTargetJSON(t *testing.T) {
	in := ProductionWatchTargetV1{
		Repository:               "0-acme/api",
		DeploymentID:             "d",
		Kind:                     ProductionWatchKindWebService,
		Region:                   region_enums.EuWest1,
		LoadBalancerArn:          "lb",
		TargetGroupArn:           "tg",
		EcsServiceArn:            "svc",
		ClusterName:              "c",
		CloudfrontDistributionID: "cf",
		URLs:                     []string{"https://a", "https://a/health"},
		DeployStartedAt:          10,
		DeployEndedAt:            20,
		BuildID:                  "b",
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	want := []string{"repository", "deploymentId", "kind", "region", "loadBalancerArn", "targetGroupArn",
		"ecsServiceArn", "clusterName", "cloudfrontDistributionId", "urls", "deployStartedAt", "deployEndedAt", "buildId"}
	if len(fields) != len(want) {
		t.Errorf("fields = %v; want %v", fields, want)
	}
	for _, k := range want {
		if _, ok := fields[k]; !ok {
			t.Errorf("missing JSON field %q in %s", k, data)
		}
	}
	var out ProductionWatchTargetV1
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip: got %+v, want %+v", out, in)
	}

	// A static site has no load balancer, service or cluster: omitted.
	data, err = json.Marshal(ProductionWatchTargetV1{Repository: "1-acme/web", Kind: ProductionWatchKindStaticSite})
	if err != nil {
		t.Fatal(err)
	}
	fields = map[string]interface{}{}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"loadBalancerArn", "targetGroupArn", "ecsServiceArn", "clusterName", "cloudfrontDistributionId", "urls"} {
		if _, ok := fields[k]; ok {
			t.Errorf("empty %q not omitted in %s", k, data)
		}
	}
}

func TestProductionWatchWireNames(t *testing.T) {
	cases := map[string]string{
		ProductionWatchKindWebService:          "web_service",
		ProductionWatchKindPrivateService:      "private_service",
		ProductionWatchKindStaticSite:          "static_site",
		ProductionWatchStateRunning:            "running",
		ProductionWatchStateDone:               "done",
		ProductionWatchVerdictNoChange:         "no_change",
		ProductionWatchVerdictErrorsUp:         "errors_up",
		ProductionWatchVerdictNotEnoughTraffic: "not_enough_traffic",
		ProductionWatchSignalUnhealthyTargets:  "unhealthy_targets",
		ProductionWatchSignalTaskCrashed:       "task_crashed",
		ProductionWatchSignalURLFailing:        "url_failing",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("wire name = %q; want %q", got, want)
		}
	}
	if ProductionWatchSignalErrorsUp != "errors_up" {
		t.Errorf("ProductionWatchSignalErrorsUp = %q", ProductionWatchSignalErrorsUp)
	}
	if ProductionWatchWindowSeconds != 1800 {
		t.Errorf("ProductionWatchWindowSeconds = %d; want 1800", ProductionWatchWindowSeconds)
	}
}

func TestProductionWatchDtosGobRoundTrip(t *testing.T) {
	for _, report := range []ProductionWatchReportDtoV1{
		{State: ProductionWatchStateRunning, PolledAt: 5, WindowEndsAt: 1805},
		{
			State: ProductionWatchStateDone, Verdict: ProductionWatchVerdictErrorsUp, PolledAt: 1805, WindowEndsAt: 1805,
			Targets: []ProductionWatchTargetReportDtoV1{{Repository: "0-a/b", RequestsAfter: 3, Errors5xxAfter: 1, FailingURLs: []string{"https://a"}}},
			Signals: []ProductionWatchSignalDtoV1{{Repository: "0-a/b", Kind: ProductionWatchSignalErrorsUp, Detail: "d", SeenAt: 7}},
		},
	} {
		in := UpdateTaskProductionWatchArgsV1{Update: UpdateTaskProductionWatchDtoV1{TaskID: "t", JobID: "j", Report: report}}
		var out UpdateTaskProductionWatchArgsV1
		gobRoundTrip(t, in, &out)
		if !reflect.DeepEqual(in, out) {
			t.Errorf("round trip: got %+v, want %+v", out, in)
		}
	}
	var reply UpdateTaskProductionWatchReplyV1
	gobRoundTrip(t, UpdateTaskProductionWatchReplyV1{Done: true, Stop: true}, &reply)
	if !reply.Done || !reply.Stop {
		t.Errorf("reply round trip: %+v", reply)
	}
}
