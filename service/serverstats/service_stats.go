package serverstats

import (
	"github.com/rsuchkov/gopractice/model"
)

func (svc *Processor) SaveMetric(name string, mtype model.MetricType, value float64) error {
	err := mtype.Validate()
	if err != nil {
		return err
	}

	if mtype == model.MetricTypeGauge {
		svc.statsStorage.SaveMetric(name, mtype, value)
	} else { // model.MetricTypeCounter
		m, err := svc.statsStorage.GetMetric(name, mtype)
		if err != nil {
			svc.statsStorage.SaveMetric(name, mtype, value)
		} else {
			svc.statsStorage.SaveMetric(name, mtype, value+m.Value)
		}
	}

	return nil
}

func (svc *Processor) GetMetric(name string, mtype model.MetricType) (float64, error) {
	err := mtype.Validate()
	if err != nil {
		return 0, err
	}
	m, err := svc.statsStorage.GetMetric(name, mtype)
	if err != nil {
		return 0, err
	}
	return m.Value, nil
}

func (svc *Processor) GetMetrics() ([]model.Metric, error) {
	metrics := []model.Metric{}
	for _, value := range svc.statsStorage.GetMetrics() {
		metrics = append(metrics, value)
	}

	return metrics, nil
}
