package propagate

import (
	"context"

	"example.com/tracelink/internal/model"
)

type traceKey struct{}

// WithTrace stores the trace context in ctx.
func WithTrace(ctx context.Context, tc model.TraceContext) context.Context {
	return context.WithValue(ctx, traceKey{}, tc)
}

// TraceFrom extracts the trace context from ctx.
func TraceFrom(ctx context.Context) (model.TraceContext, bool) {
	tc, ok := ctx.Value(traceKey{}).(model.TraceContext)
	return tc, ok
}
