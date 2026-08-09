package parameters_enums

import (
	"strconv"
	"testing"
)

// Job parameters are PERSISTED as a map keyed by these strings —
// job.Parameters is a map[string]interface{} in Mongo, and the runner reads
// values back out by the same key. Changing one orphans that parameter on every
// existing job: the value survives under the old key and is simply never read,
// so the job runs with the field silently absent.
//
// This package had no tests at all, while carrying the most exposed strings in
// the catalogue work.
//
// The contract is that a key's string IS its decimal value, which makes this
// exhaustive BY CONSTRUCTION — it walks keyMap, so a key added tomorrow is
// covered without anyone remembering to list it. That is deliberately stronger
// than the hand-maintained pinning the model and provider tests use, where no
// such relationship exists to lean on.
func TestKeyStringsAreTheEnumValue(t *testing.T) {
	if len(keyMap) == 0 {
		t.Fatal("keyMap is empty; this test would pass vacuously")
	}
	for k, want := range keyMap {
		if got := strconv.Itoa(int(k)); got != want {
			t.Errorf("%s: keyMap says %q but the enum value is %s — the persisted key must be the value, or a typo silently orphans the parameter",
				k.String(), want, got)
		}
		got, err := k.Key()
		if err != nil || got != want {
			t.Errorf("%s.Key() = %q, %v; want %q", k.String(), got, err, want)
		}
	}
}

// Two keys sharing a string would make one overwrite the other on every job —
// last writer wins, and the loser reads back the wrong value rather than
// nothing.
func TestKeyStringsAreUnique(t *testing.T) {
	seen := map[string]Key{}
	for k, s := range keyMap {
		if prev, dup := seen[s]; dup {
			t.Errorf("key %q is shared by %s and %s; one would overwrite the other in job.Parameters", s, prev.String(), k.String())
		}
		seen[s] = k
	}
}

// An unmapped key must ERROR rather than return "", which would write every
// such parameter under the empty string and collide them all together.
func TestUnmappedKeyErrors(t *testing.T) {
	if _, err := Key(99999).Key(); err == nil {
		t.Error("an unmapped key must error, not return an empty key")
	}
}
