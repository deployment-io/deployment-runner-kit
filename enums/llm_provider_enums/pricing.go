package llm_provider_enums

// Token pricing, keyed by (model, provider).
//
// ⚠️ THESE RATES DESCRIBE THE PRESENT, NEVER THE PAST.
//
// A rate here is an input to pricing a run AT THE MOMENT IT COMPLETES. The
// result is then recorded with the run and never recomputed, because vendors
// change prices and a Task's cost is a fact about when it happened. Deriving a
// historical cost from this table would silently rewrite every past Task the
// next time someone edited a number — no migration, no audit trail, and the
// dashboard quietly disagreeing with the invoice.
//
// That is not hypothetical: it is what app-server did before this existed. Its
// codex estimate ran during response mapping, so the cost of a June run was
// whatever the table said today.
//
// KEYED BY PROVIDER, not by model alone, because the same model costs
// different amounts depending on the route. Bedrock, Vertex and the gateway
// providers (Vercel, Cloudflare) each publish their own price for a model
// whose vendor is someone else entirely — and a gateway's margin is exactly
// the difference a model-only table cannot express.

// Rate is a model's token pricing at one provider, in USD per million tokens.
//
// Three rates rather than two because cached input bills lower than fresh
// input at every provider that offers it, and an agent re-reading a large
// repo context is mostly cache hits. Treating them as one rate overstates a
// long run substantially.
type Rate struct {
	// InputPerM prices fresh (non-cached) input tokens.
	InputPerM float64
	// CachedPerM prices input served from cache.
	CachedPerM float64
	// OutputPerM prices output tokens, reasoning tokens included.
	OutputPerM float64
}

// modelProviderRate holds published rates per (model, provider).
//
// SPARSE ON PURPOSE. An absent entry means "we cannot price this", which
// callers must render as an unknown cost rather than as zero — a confidently
// wrong number is worse than an honest blank.
//
// Only the pairs we actually need to price are here. Claude Code and opencode
// report their own cost, so nothing needs a rate to price them; codex reports
// token counts alone, which is why its models are the entries below. Add a
// pair when something has to price a run we cannot otherwise measure.
//
// Source: OpenAI's published API pricing, cross-checked against the Codex
// credit rate card (1 credit = $0.04). Nothing calls a vendor at runtime — a
// maintainer consults the pricing page by hand when updating this.
var modelProviderRate = map[Model]map[Provider]Rate{
	Gpt55:      {OpenAIDirect: {InputPerM: 5.00, CachedPerM: 0.50, OutputPerM: 30.00}},
	Gpt54:      {OpenAIDirect: {InputPerM: 2.50, CachedPerM: 0.25, OutputPerM: 15.00}},
	Gpt53Codex: {OpenAIDirect: {InputPerM: 1.75, CachedPerM: 0.175, OutputPerM: 14.00}},
	// The 5.6 family, at OpenAI's STANDARD list price — not the launch promo,
	// which is time-limited and would expire into a table nothing re-checks.
	// Understating a rate for a while is a wrong number either way; the promo
	// is the one that goes wrong silently and on a date nobody wrote down.
	//
	// ⚠️ THE LONG-CONTEXT SURCHARGE IS NOT PRICED. OpenAI bills 5.6 at a higher
	// input rate above a request-size threshold, and Rate has no tiered field —
	// three flat numbers per pair is the whole model. So a run whose prompts
	// cross that threshold is priced LOW here, and knowingly: an agent re-reading
	// a large repo does cross it. The alternative was a fourth field and a
	// threshold nothing else in the catalogue needs, for a cost line that is
	// already an estimate off token counts codex reports after the fact. Give
	// Rate a tiered input before anything bills a customer from this number.
	//
	// Cached input is 10% of fresh input at each, which is OpenAI's published
	// ratio for the family rather than an arithmetic convenience.
	Gpt56Sol:   {OpenAIDirect: {InputPerM: 5.00, CachedPerM: 0.50, OutputPerM: 30.00}},
	Gpt56Terra: {OpenAIDirect: {InputPerM: 2.00, CachedPerM: 0.20, OutputPerM: 12.00}},
	Gpt56Luna:  {OpenAIDirect: {InputPerM: 0.20, CachedPerM: 0.02, OutputPerM: 1.20}},
}

// RateFor returns the published rate for this model at this provider.
//
// The bool is the whole point: false means "no published rate for this pair",
// and the caller must not substitute zero. A subscription has no per-token
// price at all, so that pair is legitimately absent forever rather than
// pending.
func (m Model) RateFor(p Provider) (Rate, bool) {
	byProvider, ok := modelProviderRate[m]
	if !ok {
		return Rate{}, false
	}
	r, ok := byProvider[p]
	return r, ok
}

// CostUSD prices one run from its token counts.
//
// inputTokens is the provider's total input count, which ALREADY INCLUDES the
// cached subset — that is how both OpenAI and Anthropic report it. The
// fresh-input portion is therefore (input - cacheRead), and double-counting
// the cached tokens at the full rate is the obvious way to get this wrong.
//
// cacheWriteTokens bills at the fresh-input rate here. Anthropic prices cache
// writes at a premium over fresh input, so this UNDERSTATES a cache-heavy
// Claude run — acceptable only because nothing prices Claude from this table
// today (Claude Code reports its own cost). Give Rate a fourth field before
// that changes.
//
// Negative arithmetic is clamped rather than trusted: a provider reporting
// cacheRead > input would otherwise produce a negative cost, which reads as a
// refund.
func (r Rate) CostUSD(inputTokens, cacheReadTokens, cacheWriteTokens, outputTokens int) float64 {
	const perMillion = 1_000_000.0
	fresh := inputTokens - cacheReadTokens
	if fresh < 0 {
		fresh = 0
	}
	if cacheReadTokens < 0 {
		cacheReadTokens = 0
	}
	if cacheWriteTokens < 0 {
		cacheWriteTokens = 0
	}
	if outputTokens < 0 {
		outputTokens = 0
	}
	return float64(fresh+cacheWriteTokens)/perMillion*r.InputPerM +
		float64(cacheReadTokens)/perMillion*r.CachedPerM +
		float64(outputTokens)/perMillion*r.OutputPerM
}
