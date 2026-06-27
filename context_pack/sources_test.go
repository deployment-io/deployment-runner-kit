package context_pack

import "testing"

// allSources is the canonical list of source consts. Names must stay unique — they key the merge
// AND derive the on-disk filename (File), so a duplicate would both clobber the merge and collide
// on disk.
var allSources = []SourceName{
	SourceRepoCatalog,
	SourceServiceRepoMap,
	SourceServicesObserved,
}

func TestSourceNamesUnique(t *testing.T) {
	seen := map[SourceName]bool{}
	for _, s := range allSources {
		if seen[s] {
			t.Errorf("duplicate source name %q", s)
		}
		seen[s] = true
	}
}

func TestFileDerivesFromSource(t *testing.T) {
	if got := File(SourceServicesObserved); got != "services-observed.json" {
		t.Errorf("File(%q) = %q, want services-observed.json", SourceServicesObserved, got)
	}
}
