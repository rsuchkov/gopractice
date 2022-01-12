package serverstats

import (
	"fmt"

	"github.com/rsuchkov/gopractice/model"
)

func (svc *Processor) SaveMetric(m model.Metric) (model.Metric, error) {

	switch m.MType {
	case model.MetricTypeGauge:
		return svc.statsStorage.SaveMetric(m.ID, m.MType, *m.Value)
	case model.MetricTypeCounter:
		return svc.statsStorage.IncMetric(m.ID, m.MType, *m.Value)
	default:
		return model.Metric{}, fmt.Errorf("unknown metric type: %s", m.MType)
	}
}

func (svc *Processor) GetMetric(name string, mtype model.MType) (model.Metric, error) {
	if err := mtype.Validate(); err != nil {
		return model.Metric{}, err
	}

	m, ok := svc.statsStorage.GetMetric(name, mtype)
	if !ok {
		return model.Metric{}, fmt.Errorf("metric %s with type %s doesn't exist", name, mtype)
	}

	return m, nil
}

func (svc *Processor) GetMetrics() ([]model.Metric, error) {
	metrics := []model.Metric{}
	for _, value := range svc.statsStorage.GetMetrics() {
		metrics = append(metrics, value)
	}

	return metrics, nil
}
