package context_pack_enums

// ScopeLevel is the breadth a stored context record applies to. Records are keyed by scope so
// org-wide data (e.g. the repo catalog) is stored once and per-cluster infra once, rather than
// duplicated into every environment's pack. Stored, so values are fixed and new levels are
// appended before MaxScopeLevel — never renumbered.
type ScopeLevel uint

const (
	Org           ScopeLevel = iota + 1 // applies to the whole organization (e.g. repo catalog)
	Environment                         // a specific deployment environment
	Cluster                             // a specific K8s / compute cluster
	MaxScopeLevel                       // always add new scope levels before MaxScopeLevel
)

var scopeLevelToString = map[ScopeLevel]string{
	Org:         "org",
	Environment: "environment",
	Cluster:     "cluster",
}

func (s ScopeLevel) String() string {
	return scopeLevelToString[s]
}
