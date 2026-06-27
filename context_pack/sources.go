package context_pack

// SourceName is the canonical identifier for a context source — the `source` half of the
// (scope, source) merge key (see ManifestEntry.Source and the kit context_pack_models merge).
// Always reference a const below rather than hand-writing the string, so the merge key stays
// typo-safe: a drifted string would orphan a source's entries instead of replacing them.
//
// It is an ALIAS for string (not a distinct type) so it stays directly assignable to
// ManifestEntry.Source and the field remains OPEN to future / third-party source names.
//
// A source names a KIND of context (a data type), NOT the connector or cloud that produced it — the
// cloud / account / region / cluster is carried entirely by the Scope (its ID). So the same source
// name is reused across clouds and runners: the live-infra map is "services-observed" whether the
// AWS or GCP connector produced it, told apart by scope. That keeps SourceFile a clean 1:1 mapping
// and lets the agent read a uniform filename (/work/context/services-observed.json) on any cloud.
type SourceName = string

const (
	// SourceRepoCatalog — the org's repositories enumerated from the GitHub App. Org scope, Discovered.
	SourceRepoCatalog SourceName = "repo-catalog"
	// SourceServiceRepoMap — service↔repo DECLARED by each repo's deploy-config. Org scope, Derived.
	SourceServiceRepoMap SourceName = "service-repo-map"
	// SourceServicesObserved — service↔repo OBSERVED from live infrastructure (the AWS ECS connector
	// today; GCP/Azure reuse this name at their own scopes). Cluster scope, Discovered.
	SourceServicesObserved SourceName = "services-observed"
)

// SourceFile is the canonical on-disk file each source owns, 1:1 — so artifact names are unique per
// source by construction. That uniqueness is what makes the merge's artifact↔entry link (by
// Path == Name) unambiguous: no two sources can collide on a filename. A writer sets the manifest
// entry's Source + Path and the artifact's Name all from here.
//
// Filenames are the data type (cloud-neutral; the scope carries the cloud), so the agent reads a
// uniform name regardless of which cloud's connector produced it. Future sources extend this map
// with their own file, e.g. "dns.json", "registry.json", "secrets.json", "databases.json",
// "cluster.json", "answered.json".
var SourceFile = map[SourceName]string{
	SourceRepoCatalog:      "repos.json",
	SourceServiceRepoMap:   "services-declared.json",
	SourceServicesObserved: "services-observed.json",
}
