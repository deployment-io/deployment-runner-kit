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
	// Explicit values, like parameters_enums.Key and Provider. These are not
	// persisted — Task.AgentType stores the STRING — but a value that depends
	// on the lines above it is a trap either way, and the enum is swept by
	// AllAgentTypes below rather than counted through.
	ClaudeCode AgentType = 1
	Codex      AgentType = 2
	Opencode   AgentType = 3
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
	// Membership, not a range. A sentinel has to be hand-bumped, so adding an
	// agent and forgetting reads as INVALID, while reserving a value makes the
	// gap read as valid. Neither fails loudly. Same reason Provider dropped
	// MaxProvider.
	_, ok := agentTypeToString[h]
	return ok
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

// agentTypeToDisplayName is what a human calls the agent. Separate from
// String(), which is the AGENT_TYPE wire value agentbox switches on — letting a
// copy edit change what gets sent to a container is the class of bug this
// package exists to end.
var agentTypeToDisplayName = map[AgentType]string{
	ClaudeCode: "Claude Code",
	Codex:      "Codex",
	Opencode:   "opencode",
}

// DisplayName returns the human-facing name, falling back to the wire value so
// an unnamed agent still renders as something.
func (h AgentType) DisplayName() string {
	if name, ok := agentTypeToDisplayName[h]; ok {
		return name
	}
	return h.String()
}

// batchOnlyAgentTypes cannot run an interactive Assistant session.
//
// A capability of the AGENT, so it belongs here rather than in whichever client
// renders a session picker. opencode has no interactive mode in agentbox: it
// runs a batch invocation and exits, so a session would have nothing to talk
// to.
var batchOnlyAgentTypes = map[AgentType]bool{
	Opencode: true,
}

// SupportsInteractiveSession reports whether this agent can back an Assistant
// session, as opposed to batch Task Steps only.
func (h AgentType) SupportsInteractiveSession() bool {
	return !batchOnlyAgentTypes[h]
}

// agentTypePriority decides which agent runs a model that SEVERAL can run.
//
// Explicit, because this used to be the enum's declaration order: AgentTypesFor
// counted up from ClaudeCode, so "claude-code wins a shared model" was encoded
// in the fact that it happened to be 1. Renumbering or inserting an agent would
// have silently changed which harness ran an existing Task.
//
// Lower wins. Claude Code leads as the most established harness; opencode is
// last because it is the provider-agnostic one and the least specific answer
// when something else can run the model too.
var agentTypePriority = map[AgentType]int{
	ClaudeCode: 0,
	Codex:      1,
	Opencode:   2,
}

// Priority returns the tie-break rank for a model several agents can run.
func (h AgentType) Priority() int {
	if p, ok := agentTypePriority[h]; ok {
		return p
	}
	return len(agentTypePriority) // unranked agents sort last, deterministically
}
