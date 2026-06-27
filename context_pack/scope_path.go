package context_pack

import (
	"strings"

	"github.com/deployment-io/deployment-runner-kit/enums/context_pack_enums"
)

// ScopePath is the /work/context sub-path a scope's files live under. It exists because filenames
// are per-source (File(SourceAwsEcs) == "aws-ecs.json"), NOT per-(source, scope) — so two ECS
// clusters' packs both carry an "aws-ecs.json" and would clobber each other if rendered into the
// same directory. Namespacing by scope keeps the full path unique per (scope, source).
//
// It is the single rule materialize applies to EVERY pack — not per-source plumbing:
//
//	Org         -> ""                       (root; materialize is per-org, so an org segment is redundant)
//	Environment -> "environments/<id>"
//	Cluster     -> "clusters/<id>"
//	(other)     -> "scopes/<id>"            (safety net: a new level still namespaces, never collides)
//
// <id> is sanitizeID(scope.ID) — a cloud-neutral form of the raw scope ID (an AWS ARN, a GCP
// project/cluster path, …), so it's unique + deterministic on any cloud without per-cloud parsing.
// The path can be ugly-but-unique because readability comes from index.md (which lists the full path
// plus a human summary). Used by BOTH the materialize file-render and the index.md entries, so the
// file an agent opens and the path the index advertises can never drift. A new source needs zero
// path work (it emits a scope-agnostic File(source) and lands under its pack's ScopePath); a new
// scope level is one case here.
func ScopePath(scope Scope) string {
	switch scope.Level {
	case context_pack_enums.Org:
		return ""
	case context_pack_enums.Environment:
		return "environments/" + sanitizeID(scope.ID)
	case context_pack_enums.Cluster:
		return "clusters/" + sanitizeID(scope.ID)
	default:
		return "scopes/" + sanitizeID(scope.ID)
	}
}

// sanitizeID turns a raw scope ID into a safe, stable single path segment: lowercased, every run of
// characters outside [a-z0-9] collapsed to a single "-", and trimmed. Injective enough for real
// scope IDs — an ARN's separators (":", "/") sit at fixed structural positions, so distinct ARNs
// (differing in region / account / name) stay distinct after collapsing. Readability is index.md's
// job, so this optimizes for uniqueness + determinism, not prettiness.
func sanitizeID(id string) string {
	var b strings.Builder
	b.Grow(len(id))
	lastDash := false
	for _, r := range strings.ToLower(id) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
