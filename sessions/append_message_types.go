package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// AppendMessageDtoV1 is one incremental assistant-message update streamed from
// the runner to deployment-server during an interactive Assistant session.
//
// A session Job produces many of these (unlike the one-shot devops agent).
// Deltas of a single assistant message share MessageID and arrive with
// IsDone=false (Content is the new text fragment emitted by agentbox); the
// terminating update sets IsDone=true. deployment-server maps JobID→Thread and
// appends a message_stream record per update (persisting the Message when
// IsDone), reusing the same SSE path the devops agent already uses.
type AppendMessageDtoV1 struct {
	JobID     string
	MessageID string // stable per assistant message; the UI coalesces deltas by it
	Content   string // text fragment (delta) for this update
	IsDone    bool   // marks the assistant message complete
	// TurnEnd marks the agent's turn boundary: it finished its reply and is
	// waiting for the next user message. Sent as its own update (no Content,
	// IsDone=true, fresh MessageID); deployment-server records it as a
	// message_enums.TurnEnd message so the UI can gate its composer on the
	// boundary (ping-pong) and group a turn's messages.
	TurnEnd bool
}

type AppendMessageArgsV1 struct {
	types.AuthArgsV1
	Messages []AppendMessageDtoV1
}

type AppendMessageReplyV1 struct {
	Done bool
}
