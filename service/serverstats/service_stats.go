package serverstats

import (
	"fmt"

	"github.com/rsuchkov/gopractice/model"
)

func (svc *Processor) SaveMetric(name string, mtype model.MType, value float64) error {
	err := mtype.Validate()
	if err != nil {
		return err
	}

	switch mtype {
	case model.MetricTypeGauge:
		svc.statsStorage.SaveMetric(name, mtype, value)
	case model.MetricTypeCounter:
		svc.statsStorage.IncMetric(name, mtype, value)
	default:
		return fmt.Errorf("unknown metric type: %s", mtype)
	}

	return nil
}

func (svc *Processor) SaveMetric_v2(m model.Metric) (model.Metric, error) {

	switch m.MType {
	case model.MetricTypeGauge:
		return svc.statsStorage.SaveMetric(m.ID, m.MType, *m.Value)
	case model.MetricTypeCounter:
		return svc.statsStorage.IncMetric(m.ID, m.MType, *m.Value)
	default:
		return model.Metric{}, fmt.Errorf("unknown metric type: %s", m.MType)
	}
}

func (svc *Processor) GetMetric(name string, mtype model.MType) (float64, error) {
	if err := mtype.Validate(); err != nil {
		return 0, err
	}

	m, ok := svc.statsStorage.GetMetric(name, mtype)
	if !ok {
		return 0, fmt.Errorf("metric %s with type %s doesn't exist", name, mtype)
	}

	return *m.Value, nil
}

func (svc *Processor) GetMetrics() ([]model.Metric, error) {
	metrics := []model.Metric{}
	for _, value := range svc.statsStorage.GetMetrics() {
		metrics = append(metrics, value)
	}

	return metrics, nil
}
