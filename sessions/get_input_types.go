package sessions

import "github.com/deployment-io/deployment-runner-kit/types"

// UserMessageDtoV1 is one user turn the runner pulls from deployment-server to
// feed the live agent (written into agentbox's .agentbox-input/messages/). The
// user types in the browser, app-server persists the message, and the runner
// polls for undelivered turns and forwards them to the running session.
type UserMessageDtoV1 struct {
	ID      string
	Content string
	Ts      int64
}

// GetInputArgsV1 asks for user turns newer than AfterTs for the session Job's
// thread. The runner keeps AfterTs as an in-memory high-water mark for the
// session's lifetime (the command runs for the whole session), so the server
// stays stateless — it returns User messages with Ts > AfterTs, ordered.
type GetInputArgsV1 struct {
	types.AuthArgsV1
	JobID   string
	AfterTs int64
}

type GetInputReplyV1 struct {
	Messages []UserMessageDtoV1
}
