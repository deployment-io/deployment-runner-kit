package context_pack

import (
	"testing"

	"github.com/deployment-io/deployment-runner-kit/enums/context_pack_enums"
)

func TestScopePath_OrgIsRoot(t *testing.T) {
	if got := ScopePath(Scope{Level: context_pack_enums.Org, ID: "anything"}); got != "" {
		t.Errorf("Org should materialize at root, got %q", got)
	}
}

func TestScopePath_ClusterNamespacedByArn(t *testing.T) {
	got := ScopePath(Scope{Level: context_pack_enums.Cluster, ID: "arn:aws:ecs:us-east-1:123456789:cluster/prod"})
	want := "clusters/arn-aws-ecs-us-east-1-123456789-cluster-prod"
	if got != want {
		t.Errorf("ScopePath = %q, want %q", got, want)
	}
}

// The whole point: distinct cluster scopes must not collide on disk (they share aws-ecs.json).
func TestScopePath_DistinctClustersDistinctPaths(t *testing.T) {
	a := ScopePath(Scope{Level: context_pack_enums.Cluster, ID: "arn:aws:ecs:us-east-1:111:cluster/web"})
	b := ScopePath(Scope{Level: context_pack_enums.Cluster, ID: "arn:aws:ecs:us-west-2:222:cluster/web"})
	if a == b {
		t.Errorf("different account/region clusters collided on path %q", a)
	}
}

func TestScopePath_UnknownLevelStillNamespaces(t *testing.T) {
	// A level this switch doesn't know yet must still get a unique, collision-free home.
	got := ScopePath(Scope{Level: context_pack_enums.MaxScopeLevel, ID: "acct/123"})
	if got != "scopes/acct-123" {
		t.Errorf("unknown level path = %q, want scopes/acct-123", got)
	}
}

func TestSanitizeID(t *testing.T) {
	cases := map[string]string{
		"arn:aws:ecs:us-east-1:123:cluster/prod": "arn-aws-ecs-us-east-1-123-cluster-prod",
		"UPPER/Case":                             "upper-case",
		"::leading-and-trailing::":               "leading-and-trailing",
		"a:b/c":                                  "a-b-c",
	}
	for in, want := range cases {
		if got := sanitizeID(in); got != want {
			t.Errorf("sanitizeID(%q) = %q, want %q", in, got, want)
		}
	}
}
