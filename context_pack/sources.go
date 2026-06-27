package context_pack

// SourceName is the canonical identifier for a context source — the `source` half of the
// (scope, source) merge key (see ManifestEntry.Source and the kit context_pack_models merge), and
// also the stem of the source's on-disk filename (see File). Always reference a const below rather
// than hand-writing the string, so the merge key + filename stay typo-safe.
//
// It is an ALIAS for string (not a distinct type) so it stays directly assignable to
// ManifestEntry.Source and the field remains OPEN to future / third-party source names.
//
// A source names a KIND of context (a data type), NOT the connector or cloud that produced it — the
// cloud / account / region / cluster is carried entirely by the Scope (its ID). So the same source
// name is reused across clouds and runners: the live-infra map is "services-observed" whether the
// AWS or GCP connector produced it, told apart by scope.
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

// File is the on-disk filename a source's artifact uses: the source name + ".json". Deriving it
// (rather than a separate map) keeps filename and source 1:1 by construction — unique per source
// automatically, and impossible to drift. That 1:1-ness is what makes the merge's artifact↔entry
// link (Path == Name) unambiguous. Writers set the manifest entry's Source + Path and the artifact's
// Name all from File(source).
//
// All context artifacts are JSON today; if a source ever needs another format, give it its own
// case here rather than reintroducing a hand-maintained name map.
func File(source SourceName) string {
	return source + ".json"
}
