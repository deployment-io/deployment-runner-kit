package context_pack_enums

// Confidence is the per-item correctness signal. Detection heuristics can be wrong, so
// consumers treat low-confidence facts as "confirm live" before acting on them. Ordered
// ascending (higher = more certain); appended before MaxConfidence, never renumbered.
type Confidence uint

const (
	ConfidenceNone   Confidence = iota + 1 // could not determine
	ConfidenceLow                          // weak heuristic
	ConfidenceMedium                       // probable
	ConfidenceHigh                         // human-confirmed or unambiguous
	MaxConfidence                          // always add new levels before MaxConfidence
)

var confidenceToString = map[Confidence]string{
	ConfidenceNone:   "none",
	ConfidenceLow:    "low",
	ConfidenceMedium: "medium",
	ConfidenceHigh:   "high",
}

func (c Confidence) String() string {
	return confidenceToString[c]
}
