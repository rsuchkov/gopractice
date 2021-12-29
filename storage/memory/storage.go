package memory

import (
	"fmt"

	"github.com/rsuchkov/gopractice/model"
)

type (
	Storage struct {
		metrics AgentStats
	}

	StorageOption func(st *Storage) error
)

func New(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		metrics: AgentStats{
			Metrics: make(map[string]model.Metric),
		},
	}
	for optIdx, opt := range opts {
		if err := opt(st); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	return st, nil
}
func genKey(name string, mtype model.MetricType) string {
	return fmt.Sprintf("%s-%s", name, mtype)
}

// A private function for storing metrics in memory.
// Important: the function is not thread safe!
func (st *Storage) saveMetric(name string, mtype model.MetricType, value float64) {
	st.metrics.Metrics[genKey(name, mtype)] = model.Metric{
		Name:       name,
		MetricType: mtype,
		Value:      value,
	}
}

func (st *Storage) SaveMetric(name string, mtype model.MetricType, value float64) {
	st.metrics.mu.Lock()
	defer st.metrics.mu.Unlock()
	st.saveMetric(name, mtype, value)
}

func (st *Storage) GetMetrics() map[string]model.Metric {
	return st.metrics.Metrics
}

func (st *Storage) GetMetric(name string, mtype model.MetricType) (model.Metric, bool) {
	i, ok := st.metrics.Metrics[genKey(name, mtype)]
	if !ok {
		return model.Metric{}, false
	}
	return i, true
}

func (st *Storage) IncMetric(name string, mtype model.MetricType, value float64) error {

	st.metrics.mu.Lock()
	defer st.metrics.mu.Unlock()

	m, ok := st.GetMetric(name, mtype)
	if !ok {
		st.saveMetric(name, mtype, value)
	} else {
		st.saveMetric(name, mtype, value+m.Value)

	}
	return nil
}
