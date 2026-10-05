package task_previews

import "github.com/deployment-io/deployment-runner-kit/types"

// ListTaskPreviewsArgsV1 asks deployment-server for the previews of the task's
// active ephemeral preview Environment — one entry per undeleted Deployment that
// has a URL. The runner uses it to validate the URLs its preview tools fetch and to
// list the previews in the task's pull request body.
type ListTaskPreviewsArgsV1 struct {
	types.AuthArgsV1
	TaskID string
}

// ListTaskPreviewsReplyV1 holds the task's previews, sorted by ServiceName. Empty
// when the task has no active preview Environment.
type ListTaskPreviewsReplyV1 struct {
	Previews []TaskPreviewV1
}

// TaskPreviewV1 is one service of a task's preview Environment.
type TaskPreviewV1 struct {
	// ServiceName is the ServiceName the runner passed to EnsureV1.
	ServiceName string
	// ServiceType is one of the ServiceType* tokens.
	ServiceType string
	// URL is the preview's default URL, e.g. "https://<distribution>.cloudfront.net".
	URL string
}
