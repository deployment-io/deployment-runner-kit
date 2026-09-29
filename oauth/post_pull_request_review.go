package oauth

import "github.com/deployment-io/deployment-runner-kit/types"

// PostPullRequestReviewArgsV1 is the runner→deployment-server RPC payload for
// posting the Review stage's open findings as inline comments on a pull
// request, as ONE non-blocking review. Same identifier conventions as
// OpenPullRequestArgsV1: RepoName is the provider-side repo identifier
// ("owner/repo" on GitHub), interpreted via installation.OauthProvider.
//
// The server (the provider implementation) drops every comment that does not
// point at a line of the pull request's diff, because a provider typically
// rejects the whole review over one such comment.
type PostPullRequestReviewArgsV1 struct {
	types.AuthArgsV1
	InstallationID string
	RepoName       string
	PRNumber       int
	// Body is the review's summary text. Every
	// PullRequestReviewPostedPlaceholder in it is replaced with the number of
	// comments actually posted, after the server has dropped those not on the
	// diff.
	Body     string
	Comments []PullRequestReviewCommentV1
}

// PullRequestReviewPostedPlaceholder in PostPullRequestReviewArgsV1.Body is
// replaced server-side with the number of comments posted. It matches the
// provider implementation's placeholder.
const PullRequestReviewPostedPlaceholder = "{{posted}}"

// PullRequestReviewCommentV1 is one inline comment. Path is repo-relative;
// Line is a line of the pull request's head. StartLine is the first line of a
// multi-line comment, 0 for a single-line one.
type PullRequestReviewCommentV1 struct {
	Path      string
	Line      int
	StartLine int
	Body      string
}

// PostPullRequestReviewDtoV1 is the normalized response. Posted counts the
// comments that went out; Dropped the ones that were not on the diff.
// Unsupported reports a provider that cannot post inline reviews — NOT an
// error: nothing was posted and the caller carries on.
type PostPullRequestReviewDtoV1 struct {
	Posted      int
	Dropped     int
	Unsupported bool
}
