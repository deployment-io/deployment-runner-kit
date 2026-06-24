package context_pack

import "github.com/deployment-io/deployment-runner-kit/types"

// MaterializeContextArgsV1 asks deployment-server for the context to drop into an agent run's
// /work/context, for the given scopes (Org always; environment/cluster added later).
type MaterializeContextArgsV1 struct {
	types.AuthArgsV1
	Scopes []Scope
}

// ContextFileV1 is one file to write under /work/context (Path is relative, e.g. "repos.json").
// The server renders packs to files so the runner only writes bytes — and so the RPC payload
// stays gob-simple: no Artifact.Data interface{} crosses the wire.
type ContextFileV1 struct {
	Path    string
	Content string
}

// MaterializeContextReplyV1 is the rendered context for the requested scopes.
type MaterializeContextReplyV1 struct {
	Files []ContextFileV1
}
