// Package service wires the trace pipeline components.
package service

import (
	"sync"
	"time"

	"example.com/tracelink/internal/audit"
	"example.com/tracelink/internal/collector"
	"example.com/tracelink/internal/exporter"
	"example.com/tracelink/internal/limiter"
	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/trace"
)

// Service is the public facade of the trace pipeline.
type Service struct {
	mu        sync.Mutex
	collector *collector.Collector
	head      *sampler.HeadSampler
	tail      *sampler.TailSampler
	trees     map[string]*trace.Tree
	sampled   map[string]bool
	store     *trace.Store
	exporter  *exporter.BatchExporter
	sink      *exporter.MemorySink
	watermark *exporter.Watermark
	registry  *limiter.Registry
	ledger    *audit.Ledger
	window    *sampler.Window
	metrics   *exporter.Metrics
}

// Options configures the service.
type Options struct {
	CollectorSize int
	TailWindow    int
	StoreLimit    int
	BatchSize     int
	Retries       int
	Rate          float64
	Burst         float64
	Policies      map[string]sampler.Policy
	TailPolicy    sampler.Policy
	Sink          exporter.Sink
}

// New builds a pipeline service.
func New(opts Options) *Service {
	policies := opts.Policies
	if policies == nil {
		policies = map[string]sampler.Policy{}
	}
	sink := exporter.NewMemorySink()
	var sinkIf exporter.Sink = sink
	if opts.Sink != nil {
		sinkIf = opts.Sink
	}
	metrics := exporter.NewMetrics()
	svc := &Service{
		collector: collector.NewCollector(opts.CollectorSize),
		head:      sampler.NewHead(policies),
		tail:      sampler.NewTail(opts.TailWindow, opts.TailPolicy),
		trees:     make(map[string]*trace.Tree),
		sampled:   make(map[string]bool),
		store:     trace.NewStore(opts.StoreLimit),
		exporter:  exporter.NewBatch(sinkIf, opts.BatchSize, opts.Retries),
		sink:      sink,
		watermark: exporter.NewWatermark(),
		registry:  limiter.NewRegistry(opts.Rate, opts.Burst),
		ledger:    audit.NewLedger(),
		window:    sampler.NewWindow(30 * time.Second),
		metrics:   metrics,
	}
	svc.exporter.SetObserver(exporter.MetricsObserver(metrics))
	return svc
}

// SampleWithQuota records a sampling decision only after quota is reserved.
func (s *Service) SampleWithQuota(decision model.SamplingDecision) error {
	lim := s.registry.LimiterFor(decision.Tenant)
	return s.ledger.RecordWithQuota(decision, lim)
}
