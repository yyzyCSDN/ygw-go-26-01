package service

import (
	"context"
	"errors"
	"sync"

	"example.com/tracelink/internal/model"
)

var errWorkerFull = errors.New("ingest worker queue is full")

// IngestWorker drains spans from a channel into the service.
type IngestWorker struct {
	svc      *Service
	in       chan model.Span
	done     chan struct{}
	mu       sync.Mutex
	ingested int64
	dropped  int64
}

// StartIngestWorker starts a background worker with the given channel capacity.
func (s *Service) StartIngestWorker(capacity int) *IngestWorker {
	w := &IngestWorker{
		svc:  s,
		in:   make(chan model.Span, capacity),
		done: make(chan struct{}),
	}
	go w.loop()
	return w
}

func (w *IngestWorker) loop() {
	for span := range w.in {
		if err := w.svc.Ingest(context.Background(), span); err != nil {
			w.mu.Lock()
			w.dropped++
			w.mu.Unlock()
			continue
		}
		w.mu.Lock()
		w.ingested++
		w.mu.Unlock()
	}
	close(w.done)
}

// Submit queues a span for ingestion.
func (w *IngestWorker) Submit(span model.Span) error {
	select {
	case w.in <- span:
		return nil
	default:
		return errWorkerFull
	}
}

// Stop closes the input channel and waits for the worker to drain.
func (w *IngestWorker) Stop() {
	close(w.in)
	<-w.done
}

// Counts returns the ingested and dropped counters.
func (w *IngestWorker) Counts() (ingested int64, dropped int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.ingested, w.dropped
}
