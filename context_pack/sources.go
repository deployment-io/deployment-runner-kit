package context_pack

// SourceName is the canonical identifier for a context source — the PRODUCER of an artifact. It is
// the `source` half of the (scope, source) merge key and the stem of the on-disk filename (File).
// Always reference a const below rather than hand-writing the string, so the merge key + filename
// stay typo-safe.
//
// A source is named for its producer, not the data type:
//   - server-side derivations are named for what they derive — "repo-catalog", "service-repo-map"
//     (their "producer" is the derivation itself; there's no external connector);
//   - runner-side CLOUD CONNECTORS are named cloud+service — "aws-ecs", and later "gcp-cloudrun",
//     "aws-route53", "gcp-clouddns", etc.
//
// Naming by connector keeps it honest (different clouds genuinely produce different data) and matches
// the connector being the natural unit of TTL, gaps, and auth. The "observed vs declared" distinction
// the agent reconciles on is carried by SourceKind (Discovered vs Derived), NOT the name — so a
// cross-cloud "what's running" view groups the Discovered connectors regardless of their names, and
// connectors producing the same kind of record share a record schema for that (the cloud lives in a
// record field + the Scope). The cloud / account / region / cluster is carried by the Scope (its ID),
// never the source name.
//
// Alias for string (not a distinct type) so it stays assignable to ManifestEntry.Source and the field
// remains open to future / third-party source names.
type SourceName = string

const (
	// SourceRepoCatalog — the org's repositories enumerated from the GitHub App. Org scope, Discovered.
	SourceRepoCatalog SourceName = "repo-catalog"
	// SourceServiceRepoMap — service↔repo DECLARED by each repo's deploy-config. Org scope, Derived.
	SourceServiceRepoMap SourceName = "service-repo-map"
	// SourceAwsEcs — service↔repo OBSERVED from live AWS ECS. Cluster scope, Discovered. Sibling cloud
	// connectors follow the cloud+service convention (gcp-cloudrun, azure-containerapps, …).
	SourceAwsEcs SourceName = "aws-ecs"
	// SourceAnsweredServiceMap — service↔repo CONFIRMED by a human in the dashboard: the durable
	// "answered" layer. Org scope, Answered (SourceKind). Wins reconciliation over declared/observed
	// and survives re-scans; app-server writes it when a user confirms or corrects an observed
	// service's repo.
	SourceAnsweredServiceMap SourceName = "answered-service-map"
	// SourceManagedDeployments — deployment.io's OWN managed deployments for the org, fetched from our
	// system of record (the deployments collection). Org scope, Fetched (SourceKind). Grounds the agent
	// on what deployment.io already runs, and is the authoritative layer reconciliation uses to
	// RECOGNIZE an observed ECS service as already-managed — joining Deployment.EcsServiceArn ↔ the
	// observed service's arn, then enriching to the real name/repo/environment instead of a
	// rediscovered image-name row. Written server-side by app-server (it holds the Deployment records).
	SourceManagedDeployments SourceName = "managed-deployments"
)

// File is the on-disk filename a source's artifact uses: the source name + ".json". Deriving it
// (rather than a separate map) keeps filename and source 1:1 by construction — unique per source
// automatically, and impossible to drift. That 1:1-ness is what makes the merge's artifact↔entry
// link (Path == Name) unambiguous. Writers set the manifest entry's Source + Path and the artifact's
// Name all from File(source).
//
// All context artifacts are JSON today; if a source ever needs another format, give it its own case
// here rather than reintroducing a hand-maintained name map.
func File(source SourceName) string {
	return source + ".json"
}
