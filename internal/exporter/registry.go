package exporter

import "sync"

// Registry maps sink names to sinks.
type Registry struct {
	mu    sync.Mutex
	sinks map[string]Sink
}

// NewRegistry creates an empty sink registry.
func NewRegistry() *Registry {
	return &Registry{sinks: make(map[string]Sink)}
}

// Register installs a named sink, refusing to overwrite an existing one.
func (r *Registry) Register(name string, sink Sink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.sinks[name]; ok {
		return errSinkExists
	}
	r.sinks[name] = sink
	return nil
}

// Sink returns the named sink.
func (r *Registry) Sink(name string) (Sink, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sinks[name]
	return s, ok
}

// Names returns the registered sink names.
func (r *Registry) Names() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, 0, len(r.sinks))
	for name := range r.sinks {
		out = append(out, name)
	}
	return out
}
