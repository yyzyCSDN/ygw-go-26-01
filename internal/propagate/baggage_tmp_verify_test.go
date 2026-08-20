package propagate

import (
	"strings"
	"testing"

	"example.com/tracelink/internal/model"
)

func mkHeader(n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = "k" + string(rune('a'+i%26)) + string(rune('a'+i/26)) + "=v"
	}
	return strings.Join(parts, ",")
}

// headerN gives n distinct keys a, b, c, ... (wraps names for n>26).
func distinctHeader(n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		// produce distinct-ish keys
		parts[i] = "key" + string(rune('a'+i)) + "=v"
	}
	return strings.Join(parts, ",")
}

func TestMergeBaggageOverLimitZeroPollution(t *testing.T) {
	dst := model.Baggage{}
	// 17 distinct entries exceeds maxBaggageEntries(16).
	err := MergeBaggage(dst, distinctHeader(17))
	if err != model.ErrInvalidBaggage {
		t.Fatalf("err = %v, want ErrInvalidBaggage", err)
	}
	if len(dst) != 0 {
		t.Fatalf("dst polluted on failure: got %d entries %v", len(dst), dst)
	}
}

func TestMergeBaggageMalformedMidwayZeroPollution(t *testing.T) {
	dst := model.Baggage{}
	// valid entry then a malformed (no '=') entry.
	err := MergeBaggage(dst, "good=v,badnoequals")
	if err != model.ErrInvalidBaggage {
		t.Fatalf("err = %v, want ErrInvalidBaggage", err)
	}
	if _, ok := dst["good"]; ok {
		t.Fatalf("dst polluted on mid-parse failure: %v", dst)
	}
}

func TestMergeBaggageValidCommits(t *testing.T) {
	dst := model.Baggage{}
	err := MergeBaggage(dst, distinctHeader(16))
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(dst) != 16 {
		t.Fatalf("dst len = %d, want 16", len(dst))
	}
}

func TestMergeBaggagePreexistingCount(t *testing.T) {
	dst := model.Baggage{"pre": "1"}
	// dst already has 1; adding 16 more distinct => 17 total => over limit.
	err := MergeBaggage(dst, distinctHeader(16))
	if err != model.ErrInvalidBaggage {
		t.Fatalf("err = %v, want ErrInvalidBaggage", err)
	}
	// dst must be untouched: only "pre" remains.
	if len(dst) != 1 || dst["pre"] != "1" {
		t.Fatalf("dst polluted: %v", dst)
	}
}

func TestMergeBaggageDuplicateKeyNearLimit(t *testing.T) {
	// 16 distinct valid entries; one of them repeats a key already in dst.
	dst := model.Baggage{"a": "0"}
	// header: a=1 + 15 other distinct => merged unique = 16 (within limit).
	hdr := "a=1," + distinctHeader(15)
	err := MergeBaggage(dst, hdr)
	if err != nil {
		t.Fatalf("err = %v, want nil (duplicate should not inflate count)", err)
	}
	if dst["a"] != "1" {
		t.Fatalf("duplicate key not overwritten: %v", dst)
	}
	// merged unique = 1 (a) + 15 = 16
	if len(dst) != 16 {
		t.Fatalf("dst len = %d, want 16", len(dst))
	}
}

func TestMergeBaggageSizeLimit(t *testing.T) {
	dst := model.Baggage{}
	err := MergeBaggage(dst, strings.Repeat("a", maxBaggageBytes+1))
	if err != model.ErrInvalidBaggage {
		t.Fatalf("err = %v, want ErrInvalidBaggage", err)
	}
	if len(dst) != 0 {
		t.Fatalf("dst polluted on size failure: %v", dst)
	}
}

func TestMergeBaggageEmpty(t *testing.T) {
	dst := model.Baggage{"x": "1"}
	if err := MergeBaggage(dst, ""); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(dst) != 1 || dst["x"] != "1" {
		t.Fatalf("empty header disturbed dst: %v", dst)
	}
}
