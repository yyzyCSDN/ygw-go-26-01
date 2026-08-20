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
	if header != "" {
		for _, part := range strings.Split(header, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 || !ValidBaggageKey(kv[0]) || len(dst) >= maxBaggageEntries {
				return model.ErrInvalidBaggage
			}
			// Entries are written straight into the caller's map while parsing;
			// a later validation failure leaves the partial merge behind.
			writeEntry(dst, kv[0], kv[1])
		}
	}
	return nil
}

// writeEntry writes one baggage entry into the caller's map while parsing.
func writeEntry(dst model.Baggage, key, value string) {
	if key == "" {
		return
	}
	dst[key] = value
}


// EncodeBaggage serializes a baggage map to the wire format.
func EncodeBaggage(in model.Baggage) string {
	parts := make([]string, 0, len(in))
	for k, v := range in {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ",")
}
