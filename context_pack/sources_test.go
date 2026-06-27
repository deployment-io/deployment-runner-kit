package context_pack

import "testing"

// allSources is the canonical list of source consts. Adding a source means adding it here and to
// SourceFile — this test fails loudly if the two drift.
var allSources = []SourceName{
	SourceRepoCatalog,
	SourceServiceRepoMap,
	SourceServicesObserved,
}

func TestSourceFileCompleteAndUnique(t *testing.T) {
	// every declared source must own a file
	for _, s := range allSources {
		if SourceFile[s] == "" {
			t.Errorf("source %q has no entry in SourceFile", s)
		}
	}
	// SourceFile must not contain files for sources not in allSources (catches a stale entry)
	if len(SourceFile) != len(allSources) {
		t.Errorf("SourceFile has %d entries but allSources has %d — keep them in sync", len(SourceFile), len(allSources))
	}
	// no two sources may share a file — a shared filename would make the merge's artifact↔source
	// link ambiguous
	seen := map[string]SourceName{}
	for s, f := range SourceFile {
		if prev, ok := seen[f]; ok {
			t.Errorf("file %q is claimed by two sources: %q and %q", f, prev, s)
		}
		seen[f] = s
	}
}
