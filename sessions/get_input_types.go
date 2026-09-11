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
	// Attachments carry the files the user attached to this turn (Assistant
	// file upload, plans/PLAN_ASSISTANT_FILE_UPLOAD.md §4). The runner writes
	// each to <work>/uploads/<Path> BEFORE delivering the turn — Content's
	// pointer block names those paths — so the agent never sees a turn that
	// references a file that isn't there yet. Additive: gob drops the field on
	// a runner that predates it, which then delivers the turn text alone.
	Attachments []SessionAttachmentDtoV1
}

// SessionAttachmentDtoV1 is one attachment's extracted content. Raw uploads
// never reach the runner: Content is already-extracted text; Bytes is reserved
// for the bounded, re-encoded image copy (Phase 3b) and is nil until then.
// Exactly one of Content / Bytes is set. Path is server-assigned, relative,
// deterministic for a given message (a re-picked-up session rewrites the same
// paths instead of accumulating duplicates), and must still be anchored under
// the uploads dir by the writer — it derives from a user-supplied filename.
type SessionAttachmentDtoV1 struct {
	Name    string // sanitized display name, e.g. "pentest-report.pdf"
	Path    string // e.g. "a1b2c3-1-pentest-report.pdf.txt"
	Content string
	Bytes   []byte
	Width   int // images only
	Height  int // images only
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
