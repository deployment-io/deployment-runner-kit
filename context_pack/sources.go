package context_pack

// SourceName is the canonical identifier for a context source — the `source` half of the
// (scope, source) merge key (see ManifestEntry.Source and the kit context_pack_models merge).
// Always reference one of the consts below rather than hand-writing the string, so the merge key
// stays typo-safe: a drifted string would orphan a source's entries instead of replacing them.
//
// It is an ALIAS for string (not a distinct type) on purpose — so it stays directly assignable to
// ManifestEntry.Source and the field remains OPEN to future / third-party source names without a
// closed enum or a cross-repo type change. The consts are the canonical set; the type is just a
// readable name for the concept.
//
// One source ≈ one artifact: a source's ManifestEntry.Source and its Artifact.Name travel together,
// and the merge replaces a source's manifest entries + artifacts as a unit. Account/region/cloud are
// carried by the Scope (its ID), never encoded into the source name — so the same connector keeps a
// single stable name across accounts, regions, and runners.
type SourceName = string

const (
	// SourceRepoCatalog — the org's repositories enumerated from the GitHub App (server-side, Org
	// scope). Discovered.
	SourceRepoCatalog SourceName = "repo-catalog"
	// SourceServiceRepoMap — the declared service↔repo map parsed from each repo's deploy-config
	// (server-side, Org scope). Derived.
	SourceServiceRepoMap SourceName = "service-repo-map"
	// SourceAwsEcs — observed service↔repo recovered from live AWS ECS (runner-side, Cluster scope).
	// Discovered (observed). Multi-cloud siblings follow the cloud+service convention, e.g.
	// "gcp-cloudrun", "azure-containerapps".
	SourceAwsEcs SourceName = "aws-ecs"
)
