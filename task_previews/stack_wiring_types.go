package task_previews

import "github.com/deployment-io/deployment-runner-kit/types"

// StackWiringArgsV1 asks deployment-server how the Task's frontends reach its APIs
// in production, so a frontend preview can be built to call the API preview of the
// same change: which of a frontend's variables holds an API's address, and which
// of the API's variables lists the origins allowed to call it. deployment-server
// works it out by matching production values and answers with variable NAMES only.
type StackWiringArgsV1 struct {
	types.AuthArgsV1
	TaskID string
	// Repositories are every repository of the Task. Each CloneURL must be one of
	// the Task's own repositories (compared as StaticSiteBuildSettingsArgsV1.CloneURL
	// is); otherwise the call is refused.
	Repositories []StackRepositoryV1
}

// StackRepositoryV1 is one repository of the Task.
type StackRepositoryV1 struct {
	// Repository is the repository's work-dir name, e.g. "0-acme/web".
	Repository string
	CloneURL   string
}

// StackWiringReplyV1 holds the links deployment-server found. No link for a
// frontend and an API means it couldn't tell how they connect.
type StackWiringReplyV1 struct {
	// Links are sorted by FrontendRepository, then APIRepository, then Frontend,
	// then API. When a frontend and an API deployment of the same environment
	// matched for a pair of repositories, only same-environment links are returned
	// for that pair.
	Links []StackLinkV1
}

// StackLinkV1 is how one static-site deployment (the frontend) reaches one
// web-service deployment (the API) in production. A link has at least one
// APIURLVariable or AllowedOriginVariable.
type StackLinkV1 struct {
	FrontendRepository string
	// Frontend is the static-site deployment's name.
	Frontend      string
	APIRepository string
	// API is the web-service deployment's name.
	API string
	// APIURLVariables are the frontend's variables whose production value is an
	// http(s) URL on the API's host, sorted by Name.
	APIURLVariables []StackURLVariableV1
	// AllowedOriginVariables are the API's variables whose production value lists
	// the frontend's origin, sorted by Name.
	AllowedOriginVariables []StackOriginVariableV1
}

// StackURLVariableV1 is a frontend variable that holds the API's address.
type StackURLVariableV1 struct {
	Name string
	// Path is the path of the production URL ("" for none or "/"), never its
	// scheme, host, port, query or fragment. Only a plain route is returned (at
	// most 5 segments of letters, digits and ._~-, each 1–40 characters): a
	// variable whose path isn't one is left out.
	Path string
}

// StackOriginVariableV1 is an API variable that lists the origins allowed to call it.
type StackOriginVariableV1 struct {
	Name string
	// Separator is the separator the value uses between entries: "," or " ", or ""
	// when the value is the origin alone.
	Separator string
}
