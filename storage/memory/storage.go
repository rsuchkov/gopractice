package memory

import (
	"fmt"

	"github.com/rsuchkov/gopractice/model"
)

type (
	Storage struct {
		metrics model.AgentStats
	}

	StorageOption func(st *Storage) error
)

func New(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		metrics: model.AgentStats{
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

func (st *Storage) SaveMetric(name string, mtype model.MetricType, value float64) {
	st.metrics.Metrics[name] = model.Metric{
		Name:       name,
		MetricType: mtype,
		Value:      value,
	}
}

func (st *Storage) GetMetrics() map[string]model.Metric {
	return st.metrics.Metrics
}
