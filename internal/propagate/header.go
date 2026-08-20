package propagate

import (
	"fmt"
	"strings"

	"example.com/tracelink/internal/model"
)

// EncodeTraceHeader serializes a trace context to the wire format.
func EncodeTraceHeader(tc model.TraceContext) string {
	sampled := "00"
	if tc.Sampled {
		sampled = "01"
	}
	return fmt.Sprintf("00-%s-%s-%s", tc.TraceID, tc.SpanID, sampled)
}

// DecodeTraceHeader parses a trace header into a trace context.
func DecodeTraceHeader(header string) (model.TraceContext, error) {
	parts := strings.Split(header, "-")
	if len(parts) != 4 || parts[0] != "00" {
		return model.TraceContext{}, model.ErrInvalidSpan
	}
	if err := model.ValidateSpan(model.Span{TraceID: parts[1], SpanID: parts[2], Name: "header"}); err != nil {
		return model.TraceContext{}, err
	}
	return model.TraceContext{TraceID: parts[1], SpanID: parts[2], Sampled: parts[3] == "01"}, nil
}
