package tasks

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The production watch wire values are read by kit's wire-value test and by
// deployment-server; changing one is a wire break.
func TestProductionWatchWireValues(t *testing.T) {
	cases := []struct{ got, want string }{
		{ProductionWatchStateRunning, "running"},
		{ProductionWatchStateDone, "done"},
		{ProductionWatchVerdictNoChange, "no_change"},
		{ProductionWatchVerdictErrorsUp, "errors_up"},
		{ProductionWatchVerdictNotEnoughTraffic, "not_enough_traffic"},
		{ProductionWatchSignalUnhealthyTargets, "unhealthy_targets"},
		{ProductionWatchSignalTaskCrashed, "task_crashed"},
		{ProductionWatchSignalURLFailing, "url_failing"},
		{ProductionWatchSignalErrorsUp, "errors_up"},
		{ProductionWatchKindWebService, "web_service"},
		{ProductionWatchKindPrivateService, "private_service"},
		{ProductionWatchKindStaticSite, "static_site"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("wire value = %q; want %q", c.got, c.want)
		}
	}
	if ProductionWatchWindowSeconds != 1800 {
		t.Errorf("ProductionWatchWindowSeconds = %d; want 1800", ProductionWatchWindowSeconds)
	}
}

func TestProductionWatchTargetV1JSONRoundTrip(t *testing.T) {
	in := []ProductionWatchTargetV1{{
		Repository:               "0-acme/api",
		DeploymentID:             "65f0c0ffee0000000000000a",
		Kind:                     ProductionWatchKindWebService,
		Region:                   3,
		LoadBalancerArn:          "arn:aws:elasticloadbalancing:us-east-1:1:loadbalancer/app/x/y",
		TargetGroupArn:           "arn:aws:elasticloadbalancing:us-east-1:1:targetgroup/x/y",
		EcsServiceArn:            "arn:aws:ecs:us-east-1:1:service/c/s",
		ClusterName:              "c",
		CloudfrontDistributionID: "",
		URLs:                     []string{"https://api.acme.com", "https://api.acme.com/health"},
		DeployStartedAt:          1700000000,
		DeployEndedAt:            1700000300,
		BuildID:                  "65f0c0ffee0000000000000b",
	}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var m []map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"repository", "deploymentId", "kind", "region", "loadBalancerArn", "targetGroupArn", "ecsServiceArn", "clusterName", "urls", "deployStartedAt", "deployEndedAt", "buildId"} {
		if _, ok := m[0][k]; !ok {
			t.Errorf("json key %q missing in %s", k, b)
		}
	}
	var out []ProductionWatchTargetV1
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip = %+v; want %+v", out, in)
	}
}

func TestProductionWatchReportDtoV1JSONRoundTrip(t *testing.T) {
	in := ProductionWatchReportDtoV1{
		State:   ProductionWatchStateDone,
		Verdict: ProductionWatchVerdictErrorsUp,
		Targets: []ProductionWatchTargetReportDtoV1{{
			Repository:       "0-acme/api",
			RequestsBefore:   1000,
			RequestsAfter:    900,
			Errors5xxBefore:  1,
			Errors5xxAfter:   40,
			P95MsBefore:      120,
			P95MsAfter:       300,
			UnhealthyTargets: 1,
			CrashedTasks:     2,
			FailingURLs:      []string{"https://api.acme.com/health"},
		}},
		Signals: []ProductionWatchSignalDtoV1{{
			Repository: "0-acme/api",
			Kind:       ProductionWatchSignalErrorsUp,
			Detail:     "5xx 1 → 40",
			SeenAt:     1700000600,
		}},
		PolledAt:     1700002100,
		WindowEndsAt: 1700002100,
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out ProductionWatchReportDtoV1
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip = %+v; want %+v", out, in)
	}
}
