package llm_provider_enums

import "fmt"

// Harness is the coding-agent CLI that runs a model — the third axis,
// independent of provider and model.
//
// It is what agentbox receives as AGENT_TYPE, and the strings below are that
// contract. An empty AGENT_TYPE defaults to claude-code inside agentbox, which
// is why ResolveHarness maps "" rather than rejecting it: Tasks created before
// Codex support have no stored agent type.
//
// Not persisted as an integer — Task documents store the string (String()
// below), so the strings are the wire format and the numbering is free.
type Harness uint

const (
	ClaudeCode Harness = iota + 1
	Codex
	Opencode

	MaxHarness // always add harnesses before MaxHarness
)

var harnessToString = map[Harness]string{
	ClaudeCode: "claude-code",
	Codex:      "codex",
	Opencode:   "opencode",
}

var stringToHarness = func() map[string]Harness {
	m := make(map[string]Harness, len(harnessToString))
	for k, v := range harnessToString {
		m[v] = k
	}
	return m
}()

func (h Harness) String() string {
	return harnessToString[h]
}

func (h Harness) IsValid() bool {
	return h > 0 && h < MaxHarness
}

// ResolveHarness parses an AGENT_TYPE string. The empty string resolves to
// ClaudeCode for backward compatibility with Tasks created before the agent
// type existed — the same defaulting agentbox applies.
func ResolveHarness(s string) (Harness, error) {
	if s == "" {
		return ClaudeCode, nil
	}
	if h, ok := stringToHarness[s]; ok {
		return h, nil
	}
	return 0, fmt.Errorf("unknown agent harness %q", s)
}
