package propagate

import (
	"strings"

	"example.com/tracelink/internal/model"
)

const (
	maxBaggageEntries = 16
	maxBaggageBytes   = 4096
)

// MergeBaggage parses a wire-format baggage header and merges the entries into
// dst. The header is fully validated into a private buffer before dst is
// mutated, so an invalid header never leaves a partial merge behind.
func MergeBaggage(dst model.Baggage, header string) error {
	if len(header) > maxBaggageBytes {
		return model.ErrInvalidBaggage
	}
	// Parse and validate the whole header into a private buffer first; only
	// once every entry has passed do we commit into the caller's map. A
	// failure path therefore leaves dst untouched.
	tmp := make(model.Baggage)
	// seen tracks the keys that would exist in dst after merging the entries
	// parsed so far (including any already present in dst), so the entry
	// count limit is enforced against the merged result, not against dst.
	seen := make(map[string]struct{}, len(dst))
	for k := range dst {
		seen[k] = struct{}{}
	}
	if header != "" {
		for _, part := range strings.Split(header, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 || !ValidBaggageKey(kv[0]) || len(seen) >= maxBaggageEntries {
				return model.ErrInvalidBaggage
			}
			seen[kv[0]] = struct{}{}
			tmp[kv[0]] = kv[1]
		}
	}
	for k, v := range tmp {
		dst[k] = v
	}
	return nil
}


// EncodeBaggage serializes a baggage map to the wire format.
func EncodeBaggage(in model.Baggage) string {
	parts := make([]string, 0, len(in))
	for k, v := range in {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ",")
}
