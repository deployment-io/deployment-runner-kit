package llm_provider_enums

import (
	"math"
	"testing"
)

// The published rates, pinned. Not busywork: these numbers turn into figures a
// customer reads, and a fat-fingered decimal is invisible in review — 0.175 and
// 1.75 look alike and differ tenfold.
func TestRateFor_PublishedRatesArePinned(t *testing.T) {
	cases := map[Model]Rate{
		Gpt55:      {InputPerM: 5.00, CachedPerM: 0.50, OutputPerM: 30.00},
		Gpt54:      {InputPerM: 2.50, CachedPerM: 0.25, OutputPerM: 15.00},
		Gpt53Codex: {InputPerM: 1.75, CachedPerM: 0.175, OutputPerM: 14.00},
		// The 5.6 family at STANDARD list price, not the launch promo — pinning
		// the promo would bake an expiry date into a table nothing re-checks.
		// These also do not price the long-context surcharge; see the comment on
		// modelProviderRate for why that understatement is deliberate.
		Gpt56Sol:   {InputPerM: 5.00, CachedPerM: 0.50, OutputPerM: 30.00},
		Gpt56Terra: {InputPerM: 2.00, CachedPerM: 0.20, OutputPerM: 12.00},
		Gpt56Luna:  {InputPerM: 0.20, CachedPerM: 0.02, OutputPerM: 1.20},
	}
	for m, want := range cases {
		got, ok := m.RateFor(OpenAIDirect)
		if !ok {
			t.Errorf("%s has no rate at OpenAI Direct", m)
			continue
		}
		if got != want {
			t.Errorf("%s rate = %+v, want %+v", m, got, want)
		}
	}
	// Cached input must never cost more than fresh input — the discount is the
	// reason the field exists, and inverting it would overstate every long run.
	for m := ClaudeHaiku45; m < MaxModel; m++ {
		for _, p := range AllProviders() {
			r, ok := m.RateFor(p)
			if !ok {
				continue
			}
			if r.CachedPerM > r.InputPerM {
				t.Errorf("%s at %s: cached %.3f > fresh %.3f", m, p, r.CachedPerM, r.InputPerM)
			}
			if r.InputPerM < 0 || r.OutputPerM < 0 || r.CachedPerM < 0 {
				t.Errorf("%s at %s has a negative rate: %+v", m, p, r)
			}
		}
	}
}

// Absent must stay distinguishable from zero. A missing rate has to render as
// an unknown cost; substituting 0.0 tells a customer their run was free.
func TestRateFor_MissingPairIsNotZero(t *testing.T) {
	// codex talks only to OpenAI, so this pair is not merely unpriced — it is
	// unreachable, and must never acquire a rate by accident.
	if _, ok := Gpt55.RateFor(AWSBedrock); ok {
		t.Error("Gpt55 reports a Bedrock rate; codex cannot reach Bedrock at all")
	}
	// A subscription is a flat fee with no per-token price. This pair is
	// legitimately absent forever, not pending.
	if _, ok := ClaudeOpus5.RateFor(AnthropicSubscription); ok {
		t.Error("a subscription has no per-token rate; pricing one would invent a number")
	}
	if r, ok := ClaudeOpus5.RateFor(AnthropicSubscription); r != (Rate{}) || ok {
		t.Errorf("unpriced pair returned %+v, %v — the zero Rate must never be usable as a price", r, ok)
	}
}

func TestRate_CostUSD(t *testing.T) {
	r := Rate{InputPerM: 5.00, CachedPerM: 0.50, OutputPerM: 30.00}

	// The arithmetic that matters: input ALREADY INCLUDES the cached subset, so
	// 1M input with 600k cached bills 400k fresh + 600k cached — not 1M + 600k.
	// Double-counting here is the obvious bug, and it inflates a cache-heavy
	// run by roughly its whole context.
	got := r.CostUSD(1_000_000, 600_000, 0, 100_000)
	want := 0.4*5.00 + 0.6*0.50 + 0.1*30.00 // 2.00 + 0.30 + 3.00
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("CostUSD = %.6f, want %.6f", got, want)
	}

	// Cache WRITES are extra input, not part of the total, so they add on top.
	got = r.CostUSD(1_000_000, 600_000, 200_000, 100_000)
	want += 0.2 * 5.00
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("CostUSD with cache writes = %.6f, want %.6f", got, want)
	}

	// A provider reporting cacheRead > input must not produce a refund.
	if got := r.CostUSD(100, 500, 0, 0); got < 0 {
		t.Errorf("CostUSD = %.6f; a negative cost reads as a refund", got)
	}
	if got := r.CostUSD(0, 0, 0, 0); got != 0 {
		t.Errorf("CostUSD of an empty run = %.6f, want 0", got)
	}
}
