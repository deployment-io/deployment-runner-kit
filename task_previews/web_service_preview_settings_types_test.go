package task_previews

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"
)

// gobRoundTrip encodes v and decodes it into out, as net/rpc's default codec does
// (deployment-server serves with rpc.ServeConn).
func gobRoundTrip(t *testing.T, v, out interface{}) {
	t.Helper()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := gob.NewDecoder(&buf).Decode(out); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestWebServicePreviewSettingsReplyGobRoundTrip(t *testing.T) {
	in := WebServicePreviewSettingsReplyV1{
		Services: []WebServicePreviewSettingsV1{{
			DeploymentName:   "api",
			EnvironmentName:  "production",
			Region:           "us-east-1",
			RootDirectory:    "services/api",
			StartCommand:     "./server",
			PreDeployCommand: "./migrate",
			DockerFilePath:   "Dockerfile",
			DockerBuildArgs:  []string{"A=b"},
			Port:             8080,
			HealthCheckPath:  "/healthz",
			Cpu:              "256",
			Memory:           "512",
		}},
		RuntimeBuiltServices: []string{"worker"},
		ConfigurationSet:     true,
		Variables:            map[string]string{"API_KEY": "secret"},
		Files:                []PreviewFileV1{{Name: ".env", Contents: "A=b"}},
		Network: WebServicePreviewNetworkV1{
			VpcID:                   "vpc-1",
			VpcCidr:                 "10.0.0.0/16",
			PublicSubnets:           []string{"subnet-a"},
			PrivateSubnets:          []string{"subnet-b"},
			EcsClusterArn:           "arn:aws:ecs:us-east-1:1:cluster/c",
			EcsTaskExecutionRoleArn: "arn:aws:iam::1:role/r",
		},
		NetworkReady:  true,
		NamespaceName: "previews-a1b2",
	}
	var out WebServicePreviewSettingsReplyV1
	gobRoundTrip(t, in, &out)
	if !reflect.DeepEqual(in, out) {
		t.Errorf("got %+v\nwant %+v", out, in)
	}
}

func TestEnsureTaskPreviewReplyGobRoundTrip(t *testing.T) {
	in := EnsureTaskPreviewReplyV1{
		PreviewID:          "p",
		ExistingDistID:     "E1",
		ExistingDomain:     "d1.cloudfront.net",
		DeploymentName:     "task-preview-t-api",
		EcrRepositoryUri:   "1.dkr.ecr.us-east-1.amazonaws.com/r",
		AlbSecurityGroupId: "sg-1",
		TargetGroupArn:     "tg",
		ListenerArnPort80:  "l80",
		LoadBalancerArn:    "lb",
		LoadBalancerDns:    "lb.elb.amazonaws.com",
		EcsServiceArn:      "svc",
	}
	var out EnsureTaskPreviewReplyV1
	gobRoundTrip(t, in, &out)
	if in != out {
		t.Errorf("got %+v\nwant %+v", out, in)
	}
}

// A zero reply (nothing matched, no network) is valid on the wire.
func TestZeroRepliesGobRoundTrip(t *testing.T) {
	var settings WebServicePreviewSettingsReplyV1
	gobRoundTrip(t, WebServicePreviewSettingsReplyV1{}, &settings)
	if settings.NetworkReady || len(settings.Services) != 0 || settings.NamespaceName != "" {
		t.Errorf("got %+v, want the zero reply", settings)
	}
	var ensure EnsureTaskPreviewReplyV1
	gobRoundTrip(t, EnsureTaskPreviewReplyV1{}, &ensure)
	if ensure != (EnsureTaskPreviewReplyV1{}) {
		t.Errorf("got %+v, want the zero reply", ensure)
	}
}
