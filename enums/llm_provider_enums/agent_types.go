package llm_provider_enums

import "fmt"

// AgentType is the coding-agent CLI that runs a model — the third axis,
// independent of provider and model.
//
// NAMING: identifiers say "agent", prose says "harness". That is the existing
// house convention, not an accident — "agent" is what the persisted field
// (Task.AgentType), the agentbox wire contract (AGENT_TYPE) and the dashboard
// all call it, while "harness" is the clarifying word used in comments because
// "agent" is heavily overloaded here (kit has five agent_*_models packages for
// the unrelated Automations feature). An earlier revision of this package used
// Harness as an identifier, which turned a useful clarifying word into a second
// name for one concept; renamed to match the rest of the codebase rather than
// oblige five repos to follow.
//
// It is what agentbox receives as AGENT_TYPE, and the strings below are that
// contract. An empty AGENT_TYPE defaults to claude-code inside agentbox, which
// is why ResolveAgentType maps "" rather than rejecting it: Tasks created before
// Codex support have no stored agent type.
//
// Not persisted as an integer — Task documents store the string (String()
// below), so the strings are the wire format and the numbering is free.
type AgentType uint

const (
	ClaudeCode AgentType = iota + 1
	Codex
	Opencode

	MaxAgentType // always add harnesses before MaxAgentType
)

var agentTypeToString = map[AgentType]string{
	ClaudeCode: "claude-code",
	Codex:      "codex",
	Opencode:   "opencode",
}

var stringToAgentType = func() map[string]AgentType {
	m := make(map[string]AgentType, len(agentTypeToString))
	for k, v := range agentTypeToString {
		m[v] = k
	}
	return m
}()

func (h AgentType) String() string {
	return agentTypeToString[h]
}

func (h AgentType) IsValid() bool {
	return h > 0 && h < MaxAgentType
}

// ResolveAgentType parses an AGENT_TYPE string. The empty string resolves to
// ClaudeCode for backward compatibility with Tasks created before the agent
// type existed — the same defaulting agentbox applies.
func ResolveAgentType(s string) (AgentType, error) {
	if s == "" {
		return ClaudeCode, nil
	}
	if h, ok := stringToAgentType[s]; ok {
		return h, nil
	}
	return 0, fmt.Errorf("unknown agent harness %q", s)
}
