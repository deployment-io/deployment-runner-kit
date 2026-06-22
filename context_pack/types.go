// Package context_pack defines the shared, serializable shape of an infra/repo "context
// pack" — the artifact a BuildInfraContext runner command produces and deployment-server
// persists into the per-org context_packs collection.
//
// The pack is structured-canonical: typed Artifacts (the source of truth, queryable once in
// BSON) indexed by a Manifest. The agent-facing markdown is not stored — it's regenerated from
// the Artifacts at materialization. The runner marshals the per-scope packs to JSON into
// JobOutput; the server unmarshals them and stores one record per (org, scope).
package context_pack

import "github.com/deployment-io/deployment-runner-kit/enums/context_pack_enums"

// ManifestEntry indexes one file in the pack. It does quadruple duty: routing (Summary),
// freshness (SyncedTs + TtlS), correctness honesty (Confidence), and drift (Hash).
type ManifestEntry struct {
	Path       string                        `json:"path,omitempty" bson:"path,omitempty"`
	Source     string                        `json:"source,omitempty" bson:"source,omitempty"` // connector that produced it
	SourceKind context_pack_enums.SourceKind `json:"sourceKind,omitempty" bson:"sourceKind,omitempty"`
	Hash       string                        `json:"hash,omitempty" bson:"hash,omitempty"` // content hash, for drift detection
	SyncedTs   int64                         `json:"syncedTs,omitempty" bson:"syncedTs,omitempty"`
	TtlS       int64                         `json:"ttlS,omitempty" bson:"ttlS,omitempty"` // per-connector TTL, seconds
	Confidence context_pack_enums.Confidence `json:"confidence,omitempty" bson:"confidence,omitempty"`
	Summary    string                        `json:"summary,omitempty" bson:"summary,omitempty"` // one-line, redaction-clean
}

// Manifest is the pack index — the linchpin the agent reads first to route, and that the
// dashboard and drift-diffing read for freshness and coverage.
type Manifest struct {
	PackVersion int             `json:"packVersion,omitempty" bson:"packVersion,omitempty"`
	BuiltTs     int64           `json:"builtTs,omitempty" bson:"builtTs,omitempty"`
	Files       []ManifestEntry `json:"files,omitempty" bson:"files,omitempty"`
	// Gaps records what we could not see (auth can-i honesty) rather than implying
	// completeness — e.g. "RBAC in ns 'payments': forbidden".
	Gaps []string `json:"gaps,omitempty" bson:"gaps,omitempty"`
}

// Artifact is one typed structured artifact (e.g. repos.json, infra-fingerprint.json) — the
// canonical, queryable form. Data is left open (object or array) so each source owns its own
// schema; it round-trips JSON -> server -> queryable BSON.
type Artifact struct {
	Name string      `json:"name,omitempty" bson:"name,omitempty"`
	Data interface{} `json:"data,omitempty" bson:"data,omitempty"`
}

// Pack is one scope's slice of context — the canonical, queryable data we persist: structured
// Artifacts indexed by a Manifest. A BuildInfraContext run emits one Pack per scope (see
// ScopedPack); the agent's /work/context/ is composed from the packs whose scope covers its
// target.
//
// The agent-facing markdown is deliberately NOT stored here. It's a deterministic projection
// of the Artifacts, regenerated cheaply at materialization by the renderer — so derived data
// isn't stored twice, and renderer improvements apply to old packs automatically. LLM-derived
// narrative prose, which can't be cheaply regenerated, is stored as an Artifact instead.
type Pack struct {
	Manifest  Manifest   `json:"manifest" bson:"manifest"`
	Artifacts []Artifact `json:"artifacts,omitempty" bson:"artifacts,omitempty"`
}

// Scope is the breadth a stored context record applies to: a level plus the id of the owning
// entity — the org id for Org, an environment/cluster id otherwise. Records are keyed by Scope
// so org-wide content is stored once and per-cluster content once — never duplicated per
// environment.
//
// The id is always a real, non-empty entity id (the Org scope's id is the org id, set on the
// runner), so omitempty is safe and matches the codebase. And even when a key field is zero,
// the upsert filter carries scope.level/id, so they're written into the doc on insert anyway —
// the same reason _id can use omitempty (Mongo fills it in on insert).
type Scope struct {
	Level context_pack_enums.ScopeLevel `json:"level,omitempty" bson:"level,omitempty"`
	ID    string                        `json:"id,omitempty" bson:"id,omitempty"`
}

// ScopedPack is the unit a BuildInfraContext run emits and deployment-server persists: one
// scope's content, stored as a single context_packs record per (org, scope). A run produces a
// []ScopedPack — the sources' outputs grouped by their scope.
type ScopedPack struct {
	Scope Scope `json:"scope" bson:"scope"`
	Pack  Pack  `json:"pack" bson:"pack"`
}
