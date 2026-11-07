package tasks

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// oneOfEach returns a valid event of every type, keyed by type.
func oneOfEach() map[string]TaskEventV1 {
	return map[string]TaskEventV1{
		TaskEventStageStarted:          {Type: TaskEventStageStarted, Stage: TaskEventStageSetup, StageDetail: &TaskEventStageV1{}},
		TaskEventStageFinished:         {Type: TaskEventStageFinished, Stage: TaskEventStageReview, StageDetail: &TaskEventStageV1{Outcome: "No findings"}},
		TaskEventCheckoutDone:          {Type: TaskEventCheckoutDone, Stage: TaskEventStageSetup, Checkout: &TaskEventCheckoutV1{Repos: []string{"a/b"}}},
		TaskEventSelfCheck:             {Type: TaskEventSelfCheck, Stage: TaskEventStageImplement, Repo: "a/b", SelfCheck: &TaskEventSelfCheckV1{Commands: []string{"go test ./..."}, Passed: true}},
		TaskEventImplementFinished:     {Type: TaskEventImplementFinished, Stage: TaskEventStageImplement, Implement: &TaskEventImplementV1{Turns: 3}},
		TaskEventReviewRoundFinished:   {Type: TaskEventReviewRoundFinished, Stage: TaskEventStageReview, Round: 1, ReviewRound: &TaskEventReviewRoundV1{Findings: []TaskEventFindingV1{{Key: "k"}}}},
		TaskEventFixRoundFinished:      {Type: TaskEventFixRoundFinished, Stage: TaskEventStageReview, Round: 1, FixRound: &TaskEventFixRoundV1{Addressed: []string{"k"}}},
		TaskEventPreviewDeployFinished: {Type: TaskEventPreviewDeployFinished, Stage: TaskEventStageQA, PreviewDeploy: &TaskEventPreviewDeployV1{Kind: TaskEventPreviewKindSite, URL: "https://p.example.com", OK: true}},
		TaskEventQACheckFinished:       {Type: TaskEventQACheckFinished, Stage: TaskEventStageQA, QACheck: &TaskEventQACheckV1{Method: "GET", Path: "/me"}},
		TaskEventQACriterionFinished:   {Type: TaskEventQACriterionFinished, Stage: TaskEventStageQA, QACriterion: &TaskEventQACriterionV1{State: QACriterionPassed}},
		TaskEventQARoundFinished:       {Type: TaskEventQARoundFinished, Stage: TaskEventStageQA, QARound: &TaskEventQARoundV1{Final: true, Passed: 2}},
		TaskEventPROpened:              {Type: TaskEventPROpened, Stage: TaskEventStagePR, PullRequest: &TaskEventPullRequestV1{Number: 7, URL: "https://github.com/a/b/pull/7"}},
	}
}

func TestTaskEventWireNames(t *testing.T) {
	cases := map[string]string{
		TaskEventStageStarted:          "stage_started",
		TaskEventStageFinished:         "stage_finished",
		TaskEventCheckoutDone:          "checkout_done",
		TaskEventSelfCheck:             "self_check",
		TaskEventImplementFinished:     "implement_finished",
		TaskEventReviewRoundFinished:   "review_round_finished",
		TaskEventFixRoundFinished:      "fix_round_finished",
		TaskEventPreviewDeployFinished: "preview_deploy_finished",
		TaskEventQACheckFinished:       "qa_check_finished",
		TaskEventQACriterionFinished:   "qa_criterion_finished",
		TaskEventQARoundFinished:       "qa_round_finished",
		TaskEventPROpened:              "pr_opened",
		TaskEventStageSetup:            "setup",
		TaskEventStageImplement:        "implement",
		TaskEventStageReview:           "review",
		TaskEventStageQA:               "qa",
		TaskEventStagePR:               "pr",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("wire name = %q; want %q", got, want)
		}
	}
}

func TestTaskEventBoundKeepsEachPayload(t *testing.T) {
	for typ, e := range oneOfEach() {
		e := e
		if !e.Bound() {
			t.Errorf("%s: Bound() = false; want true", typ)
			continue
		}
		set, matching, _ := e.payloads()
		if set != 1 || !matching {
			t.Errorf("%s: payload lost after Bound (set=%d matching=%v)", typ, set, matching)
		}
	}
	e := oneOfEach()[TaskEventPreviewDeployFinished]
	e.Bound()
	if e.PreviewDeploy.URL != "https://p.example.com" || !e.PreviewDeploy.OK {
		t.Errorf("preview payload changed: %+v", e.PreviewDeploy)
	}
}

func TestTaskEventBoundRejects(t *testing.T) {
	unknown := oneOfEach()[TaskEventPROpened]
	unknown.Type = "pr_merged"
	wrongStage := oneOfEach()[TaskEventPROpened]
	wrongStage.Stage = "verify"
	noStage := oneOfEach()[TaskEventPROpened]
	noStage.Stage = ""
	nilPayload := oneOfEach()[TaskEventPROpened]
	nilPayload.PullRequest = nil
	mismatched := TaskEventV1{Type: TaskEventPROpened, Stage: TaskEventStagePR, QARound: &TaskEventQARoundV1{}}
	extra := oneOfEach()[TaskEventPROpened]
	extra.QARound = &TaskEventQARoundV1{}
	cases := map[string]TaskEventV1{
		"unknown type": unknown, "wrong stage": wrongStage, "empty stage": noStage,
		"nil payload": nilPayload, "mismatched payload": mismatched, "extra payload": extra,
	}
	for name, e := range cases {
		if e.Bound() {
			t.Errorf("%s: Bound() = true; want false", name)
		}
	}
}

func TestTaskEventIsVerdict(t *testing.T) {
	verdict := map[string]bool{
		TaskEventStageStarted:          false,
		TaskEventStageFinished:         true,
		TaskEventCheckoutDone:          false,
		TaskEventSelfCheck:             false,
		TaskEventImplementFinished:     true,
		TaskEventReviewRoundFinished:   true,
		TaskEventFixRoundFinished:      true,
		TaskEventPreviewDeployFinished: true,
		TaskEventQACheckFinished:       false,
		TaskEventQACriterionFinished:   true,
		TaskEventQARoundFinished:       true,
		TaskEventPROpened:              true,
	}
	if len(verdict) != 12 {
		t.Fatalf("want all twelve types, have %d", len(verdict))
	}
	for typ, want := range verdict {
		e := TaskEventV1{Type: typ}
		if got := e.IsVerdict(); got != want {
			t.Errorf("IsVerdict(%s) = %v; want %v", typ, got, want)
		}
	}
}

func checkRunes(t *testing.T, name, got string, n int) {
	t.Helper()
	if !utf8.ValidString(got) {
		t.Errorf("%s: invalid UTF-8 after cut", name)
	}
	if c := utf8.RuneCountInString(got); c != n {
		t.Errorf("%s: %d runes; want %d", name, c, n)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("%s: cut without …", name)
	}
}

func TestTaskEventBoundCaps(t *testing.T) {
	long := strings.Repeat("é", 5000) // 2 bytes per rune: a byte cut would split one
	many := func(n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = long
		}
		return out
	}

	e := TaskEventV1{Type: TaskEventReviewRoundFinished, Stage: TaskEventStageReview, Repo: long, Now: long, Round: -2, Seq: -1, Ts: -5,
		ReviewRound: &TaskEventReviewRoundV1{
			Findings: make([]TaskEventFindingV1, 150),
			Previous: make([]TaskEventFindingStatusV1, 150),
			Coverage: make([]TaskEventCoverageV1, 30),
		}}
	e.ReviewRound.Findings[0] = TaskEventFindingV1{Key: long, Parameter: long, Severity: long, Location: long, What: long, Why: long}
	e.ReviewRound.Previous[0] = TaskEventFindingStatusV1{Key: long, Status: long}
	e.ReviewRound.Coverage[0] = TaskEventCoverageV1{Parameter: long, State: long, Reason: long}
	if !e.Bound() {
		t.Fatal("Bound() = false")
	}
	checkRunes(t, "Repo", e.Repo, 200)
	checkRunes(t, "Now", e.Now, 200)
	if e.Round != 0 || e.Seq != 0 || e.Ts != 0 {
		t.Errorf("negative ints kept: round=%d seq=%d ts=%d", e.Round, e.Seq, e.Ts)
	}
	rr := e.ReviewRound
	if len(rr.Findings) != 100 || len(rr.Previous) != 100 || len(rr.Coverage) != 20 {
		t.Errorf("list caps: findings=%d previous=%d coverage=%d", len(rr.Findings), len(rr.Previous), len(rr.Coverage))
	}
	f := rr.Findings[0]
	checkRunes(t, "Key", f.Key, 200)
	checkRunes(t, "Location", f.Location, 500)
	checkRunes(t, "What", f.What, 1000)
	checkRunes(t, "Why", f.Why, 1000)
	checkRunes(t, "Status", rr.Previous[0].Status, 200)
	checkRunes(t, "Coverage.Reason", rr.Coverage[0].Reason, 200)
	checkRunes(t, "Coverage.Parameter", rr.Coverage[0].Parameter, 200)

	impl := TaskEventV1{Type: TaskEventImplementFinished, Stage: TaskEventStageImplement,
		Implement: &TaskEventImplementV1{Turns: -1, ChangesSummary: "a" + strings.Repeat("é", 8000), Diffstat: make([]TaskEventDiffstatV1, 25)}}
	impl.Implement.Diffstat[0] = TaskEventDiffstatV1{Repo: long, Files: -1, Insertions: -1, Deletions: -1}
	if !impl.Bound() {
		t.Fatal("implement: Bound() = false")
	}
	cs := impl.Implement.ChangesSummary
	if len(cs) > 8*1024 || !utf8.ValidString(cs) || !strings.HasSuffix(cs, "…") {
		t.Errorf("ChangesSummary: %d bytes, valid=%v", len(cs), utf8.ValidString(cs))
	}
	if impl.Implement.Turns != 0 || len(impl.Implement.Diffstat) != 20 {
		t.Errorf("implement caps: turns=%d diffstat=%d", impl.Implement.Turns, len(impl.Implement.Diffstat))
	}
	if d := impl.Implement.Diffstat[0]; d.Files != 0 || d.Insertions != 0 || d.Deletions != 0 {
		t.Errorf("negative diffstat kept: %+v", d)
	}
	checkRunes(t, "Diffstat.Repo", impl.Implement.Diffstat[0].Repo, 200)

	sc := TaskEventV1{Type: TaskEventSelfCheck, Stage: TaskEventStageImplement, SelfCheck: &TaskEventSelfCheckV1{Commands: many(15), Summary: long}}
	sc.Bound()
	if len(sc.SelfCheck.Commands) != 10 {
		t.Errorf("Commands: %d entries", len(sc.SelfCheck.Commands))
	}
	checkRunes(t, "Commands[0]", sc.SelfCheck.Commands[0], 200)
	checkRunes(t, "SelfCheck.Summary", sc.SelfCheck.Summary, 200)

	co := TaskEventV1{Type: TaskEventCheckoutDone, Stage: TaskEventStageSetup, Checkout: &TaskEventCheckoutV1{Repos: many(25)}}
	co.Bound()
	if len(co.Checkout.Repos) != 20 {
		t.Errorf("Repos: %d entries", len(co.Checkout.Repos))
	}

	fix := TaskEventV1{Type: TaskEventFixRoundFinished, Stage: TaskEventStageReview,
		FixRound: &TaskEventFixRoundV1{Addressed: many(25), ChangesSummary: strings.Repeat("x", 9000), SelfCheck: &TaskEventSelfCheckV1{Commands: many(12)}}}
	fix.Bound()
	if len(fix.FixRound.Addressed) != 20 || len(fix.FixRound.SelfCheck.Commands) != 10 || len(fix.FixRound.ChangesSummary) > 8*1024 {
		t.Errorf("fix caps: addressed=%d commands=%d summary=%d", len(fix.FixRound.Addressed), len(fix.FixRound.SelfCheck.Commands), len(fix.FixRound.ChangesSummary))
	}

	pd := TaskEventV1{Type: TaskEventPreviewDeployFinished, Stage: TaskEventStageQA,
		PreviewDeploy: &TaskEventPreviewDeployV1{Kind: long, Error: long, DatabaseEngine: long, Roles: many(25),
			Migrate: &TaskEventCommandV1{Ran: true, ExitCode: -1, Summary: long}, Seed: &TaskEventCommandV1{Summary: long}}}
	pd.Bound()
	p := pd.PreviewDeploy
	checkRunes(t, "Kind", p.Kind, 200)
	checkRunes(t, "Error", p.Error, 200)
	checkRunes(t, "DatabaseEngine", p.DatabaseEngine, 200)
	checkRunes(t, "Migrate.Summary", p.Migrate.Summary, 200)
	checkRunes(t, "Seed.Summary", p.Seed.Summary, 200)
	if len(p.Roles) != 20 || p.Migrate.ExitCode != 0 {
		t.Errorf("preview caps: roles=%d exit=%d", len(p.Roles), p.Migrate.ExitCode)
	}

	qc := TaskEventV1{Type: TaskEventQACheckFinished, Stage: TaskEventStageQA,
		QACheck: &TaskEventQACheckV1{Criterion: -1, Tool: long, Method: long, Path: "/me?token=secret#frag", Role: long, Status: -1, DurationMs: -1, Result: long}}
	qc.Bound()
	if qc.QACheck.Path != "/me" {
		t.Errorf("Path = %q; want query stripped", qc.QACheck.Path)
	}
	if qc.QACheck.Criterion != 0 || qc.QACheck.Status != 0 || qc.QACheck.DurationMs != 0 {
		t.Errorf("negative QA check ints kept: %+v", qc.QACheck)
	}
	checkRunes(t, "Tool", qc.QACheck.Tool, 200)
	checkRunes(t, "Method", qc.QACheck.Method, 200)
	checkRunes(t, "Role", qc.QACheck.Role, 200)
	checkRunes(t, "Result", qc.QACheck.Result, 200)
	qc.QACheck.Path = "/" + long
	qc.Bound()
	checkRunes(t, "Path", qc.QACheck.Path, 500)

	cr := TaskEventV1{Type: TaskEventQACriterionFinished, Stage: TaskEventStageQA, QACriterion: &TaskEventQACriterionV1{Text: long, State: long, Reason: long}}
	cr.Bound()
	checkRunes(t, "Text", cr.QACriterion.Text, 500)
	checkRunes(t, "Reason", cr.QACriterion.Reason, 1000)

	qr := TaskEventV1{Type: TaskEventQARoundFinished, Stage: TaskEventStageQA, QARound: &TaskEventQARoundV1{Passed: -1, Failed: -1, NotChecked: -1, Summary: long}}
	qr.Bound()
	if qr.QARound.Passed != 0 || qr.QARound.Failed != 0 || qr.QARound.NotChecked != 0 {
		t.Errorf("negative QA round ints kept: %+v", qr.QARound)
	}
	checkRunes(t, "QARound.Summary", qr.QARound.Summary, 200)

	pr := TaskEventV1{Type: TaskEventPROpened, Stage: TaskEventStagePR, PullRequest: &TaskEventPullRequestV1{Number: -1, Title: long, InlineComments: -3}}
	pr.Bound()
	checkRunes(t, "Title", pr.PullRequest.Title, 200)
	if pr.PullRequest.Number != 0 || pr.PullRequest.InlineComments != 0 {
		t.Errorf("negative PR ints kept: %+v", pr.PullRequest)
	}

	st := TaskEventV1{Type: TaskEventStageFinished, Stage: TaskEventStageQA, StageDetail: &TaskEventStageV1{Outcome: long}}
	st.Bound()
	checkRunes(t, "Outcome", st.StageDetail.Outcome, 200)

	short := "ok é"
	if got := capRunes(short, 200); got != short {
		t.Errorf("short string changed: %q", got)
	}
}

func TestTaskEventBoundURLs(t *testing.T) {
	longURL := "https://p.example.com/" + strings.Repeat("a", 2048) // cut to 2048 bytes, not cleared
	cases := map[string]string{
		"https://p.example.com/x":         "https://p.example.com/x",
		"http://p.example.com/x":          "",
		"https://user:pw@p.example.com/x": "",
		"https://user@p.example.com/x":    "",
		"javascript:alert(1)":             "",
		"not a url":                       "",
		"https://":                        "",
		longURL: longURL[:TaskEventURLBytes-len("…")] + "…",
	}
	for in, want := range cases {
		pr := TaskEventV1{Type: TaskEventPROpened, Stage: TaskEventStagePR, PullRequest: &TaskEventPullRequestV1{URL: in}}
		pr.Bound()
		if pr.PullRequest.URL != want {
			t.Errorf("PR URL %q -> %q; want %q", in, pr.PullRequest.URL, want)
		}
		qc := TaskEventV1{Type: TaskEventQACheckFinished, Stage: TaskEventStageQA, QACheck: &TaskEventQACheckV1{EvidenceURL: in}}
		qc.Bound()
		if qc.QACheck.EvidenceURL != want {
			t.Errorf("EvidenceURL %q -> %q; want %q", in, qc.QACheck.EvidenceURL, want)
		}
		pd := TaskEventV1{Type: TaskEventPreviewDeployFinished, Stage: TaskEventStageQA, PreviewDeploy: &TaskEventPreviewDeployV1{URL: in}}
		pd.Bound()
		if pd.PreviewDeploy.URL != want {
			t.Errorf("preview URL %q -> %q; want %q", in, pd.PreviewDeploy.URL, want)
		}
	}
}
