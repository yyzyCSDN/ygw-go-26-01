package exporter

// Observer receives batch lifecycle callbacks.
type Observer interface {
	OnBatchSuccess(count int)
	OnBatchFailure(err error)
}

// SetObserver installs a batch lifecycle observer.
func (b *BatchExporter) SetObserver(observer Observer) {
	b.observer = observer
}

func (b *BatchExporter) notifySuccess(count int) {
	if b.observer != nil {
		b.observer.OnBatchSuccess(count)
	}
}

func (b *BatchExporter) notifyFailure(err error) {
	if b.observer != nil {
		b.observer.OnBatchFailure(err)
	}
}

// MetricsObserver adapts a Metrics set to the Observer interface.
func MetricsObserver(m *Metrics) Observer {
	return metricsObserver{m: m}
}
