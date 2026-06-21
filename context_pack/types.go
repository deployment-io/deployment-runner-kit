// Package context_pack defines the shared, serializable shape of an infra/repo "context
// pack" — the artifact a BuildInfraContext runner command produces and deployment-server
// persists into the per-org context_packs collection.
//
// The pack is structured-canonical: typed Artifacts (the source of truth, queryable once in
// BSON) indexed by a Manifest, plus a derived Markdown projection (the view the agent reads,
// gzipped at the storage layer). The runner marshals a Pack to JSON into JobOutput; the
// server unmarshals it and stores it.
package context_pack

// SourceKind is the provenance of a context item. This single field lets the system absorb
// new context types and acquisition modes as cases rather than bolt-ons.
type SourceKind string

const (
	Discovered SourceKind = "discovered" // runner scan, deterministic
	Fetched    SourceKind = "fetched"    // live API snapshot (short TTL)
	Answered   SourceKind = "answered"   // human-confirmed (durable)
	Derived    SourceKind = "derived"    // structural extract / LLM
)

// Confidence is the per-item correctness signal. Detection heuristics can be wrong, so
// consumers treat low-confidence facts as "confirm live" before acting on them.
type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
	ConfidenceNone   Confidence = "none"
)

// ManifestEntry indexes one file in the pack. It does quadruple duty: routing (Summary),
// freshness (SyncedTs + TtlS), correctness honesty (Confidence), and drift (Hash).
type ManifestEntry struct {
	Path       string     `json:"path" bson:"path"`
	Source     string     `json:"source" bson:"source"` // connector that produced it
	SourceKind SourceKind `json:"sourceKind" bson:"sourceKind"`
	Hash       string     `json:"hash" bson:"hash"` // content hash, for drift detection
	SyncedTs   int64      `json:"syncedTs" bson:"syncedTs"`
	TtlS       int64      `json:"ttlS" bson:"ttlS"` // per-connector TTL, seconds
	Confidence Confidence `json:"confidence" bson:"confidence"`
	Summary    string     `json:"summary" bson:"summary"` // one-line, redaction-clean
}

// Manifest is the pack index — the linchpin the agent reads first to route, and that the
// dashboard and drift-diffing read for freshness and coverage.
type Manifest struct {
	PackVersion int             `json:"packVersion" bson:"packVersion"`
	Environment string          `json:"environment" bson:"environment"`
	BuiltTs     int64           `json:"builtTs" bson:"builtTs"`
	Files       []ManifestEntry `json:"files" bson:"files"`
	// Gaps records what we could not see (auth can-i honesty) rather than implying
	// completeness — e.g. "RBAC in ns 'payments': forbidden".
	Gaps []string `json:"gaps,omitempty" bson:"gaps,omitempty"`
}

// Artifact is one typed structured artifact (e.g. repos.json, infra-fingerprint.json) — the
// canonical, queryable form. Data is left open (object or array) so each source owns its own
// schema; it round-trips JSON -> server -> queryable BSON.
type Artifact struct {
	Name string      `json:"name" bson:"name"`
	Data interface{} `json:"data" bson:"data"`
}

// MarkdownFile is a derived markdown projection of the structured data — the view the agent
// reads. Regenerable from the Artifacts; gzipped at the storage layer.
type MarkdownFile struct {
	Path    string `json:"path" bson:"path"`
	Content string `json:"content" bson:"content"`
}

// Pack is the full output of a BuildInfraContext run: structured-canonical Artifacts indexed
// by a Manifest, plus the derived Markdown projection.
type Pack struct {
	Manifest  Manifest       `json:"manifest" bson:"manifest"`
	Artifacts []Artifact     `json:"artifacts" bson:"artifacts"`
	Markdown  []MarkdownFile `json:"markdown" bson:"markdown"`
}
