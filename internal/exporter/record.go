package exporter

import (
	"time"

	"example.com/tracelink/internal/model"
)

// BuildRecord creates an export record for a completed trace.
func BuildRecord(completed *model.CompletedTrace, sampled bool) model.ExportRecord {
	return model.ExportRecord{
		TraceID:      completed.TraceID,
		SpanCount:    completed.SpanCount,
		Sampled:      sampled,
		ExportedAtNs: time.Now().UnixNano(),
	}
}
