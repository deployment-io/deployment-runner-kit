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
			Branch:           "main",
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
			VariableNames:    []string{"API_KEY", "DATABASE_URL"},
			PreviewDatabase: &PreviewDatabaseV1{
				Mode: PreviewDatabaseModeUse, Engine: PreviewDatabaseEnginePostgres, Version: "16",
				Variables: []string{"DATABASE_URL"},
			},
			PreviewMigrateCommand: "./migrate-preview",
			BuildSecrets:          map[string]string{"NPM_TOKEN": "npm-secret"},
			ProductionEnvironment: true,
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

// A runner built before BuildSecrets decodes a reply carrying them: gob skips the field.
func TestWebServicePreviewSettingsBuildSecretsIgnoredByOlderRunner(t *testing.T) {
	type olderWebServicePreviewSettingsV1 struct {
		DeploymentName        string
		PreviewMigrateCommand string
	}
	type olderReply struct {
		Services      []olderWebServicePreviewSettingsV1
		NamespaceName string
	}
	in := WebServicePreviewSettingsReplyV1{
		Services: []WebServicePreviewSettingsV1{{
			DeploymentName:        "api",
			PreviewMigrateCommand: "./migrate-preview",
			BuildSecrets:          map[string]string{"NPM_TOKEN": "npm-secret"},
		}},
		NamespaceName: "previews-a1b2",
	}
	var out olderReply
	gobRoundTrip(t, in, &out)
	want := olderReply{
		Services:      []olderWebServicePreviewSettingsV1{{DeploymentName: "api", PreviewMigrateCommand: "./migrate-preview"}},
		NamespaceName: "previews-a1b2",
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("got %+v\nwant %+v", out, want)
	}
}

func TestStaticSiteBuildSettingsReplyGobRoundTrip(t *testing.T) {
	in := StaticSiteBuildSettingsReplyV1{
		Sites: []StaticSiteBuildSettingsV1{{
			DeploymentName:        "web",
			EnvironmentName:       "production",
			Branch:                "main",
			ProductionEnvironment: true,
			RootDirectory:         "apps/web",
			BuildCommand:          "npm run build",
			PublishDirectory:      "dist",
			IsSpa:                 true,
		}},
		ConfigurationSet: true,
		Variables:        map[string]string{"API_KEY": "secret"},
		Files:            []PreviewFileV1{{Name: ".env", Contents: "A=b"}},
	}
	var out StaticSiteBuildSettingsReplyV1
	gobRoundTrip(t, in, &out)
	if !reflect.DeepEqual(in, out) {
		t.Errorf("got %+v\nwant %+v", out, in)
	}
}

// A runner built before Branch and ProductionEnvironment decodes replies carrying
// them: gob skips the fields.
func TestBranchAndProductionEnvironmentIgnoredByOlderRunner(t *testing.T) {
	type olderService struct {
		DeploymentName  string
		EnvironmentName string
		RootDirectory   string
	}
	type olderWebReply struct{ Services []olderService }
	type olderStaticReply struct{ Sites []olderService }

	var web olderWebReply
	gobRoundTrip(t, WebServicePreviewSettingsReplyV1{Services: []WebServicePreviewSettingsV1{{
		DeploymentName: "api", EnvironmentName: "production", Branch: "main", ProductionEnvironment: true, RootDirectory: "api",
	}}}, &web)
	want := []olderService{{DeploymentName: "api", EnvironmentName: "production", RootDirectory: "api"}}
	if !reflect.DeepEqual(web.Services, want) {
		t.Errorf("web services: got %+v\nwant %+v", web.Services, want)
	}

	var static olderStaticReply
	gobRoundTrip(t, StaticSiteBuildSettingsReplyV1{Sites: []StaticSiteBuildSettingsV1{{
		DeploymentName: "api", EnvironmentName: "production", Branch: "main", ProductionEnvironment: true, RootDirectory: "api",
	}}}, &static)
	if !reflect.DeepEqual(static.Sites, want) {
		t.Errorf("static sites: got %+v\nwant %+v", static.Sites, want)
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
