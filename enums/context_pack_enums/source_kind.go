package context_pack_enums

// SourceKind is the provenance of a context item. This single field lets the system absorb
// new context types and acquisition modes as cases rather than bolt-ons. Stored in the pack
// (and thus in context_packs), so values are fixed and new kinds are appended before
// MaxSourceKind — never renumbered.
type SourceKind uint

const (
	Discovered    SourceKind = iota + 1 // runner scan, deterministic
	Fetched                             // live API snapshot (short TTL)
	Answered                            // human-confirmed (durable)
	Derived                             // structural extract / LLM
	MaxSourceKind                       // always add new source kinds before MaxSourceKind
)

var sourceKindToString = map[SourceKind]string{
	Discovered: "discovered",
	Fetched:    "fetched",
	Answered:   "answered",
	Derived:    "derived",
}

func (s SourceKind) String() string {
	return sourceKindToString[s]
}
