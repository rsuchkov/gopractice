package serverstats

import "github.com/rsuchkov/gopractice/model"

type StorageExpected interface {
	SaveMetric(string, model.MetricType, float64)
	GetMetrics() map[string]model.Metric
	GetMetric(string, model.MetricType) (model.Metric, bool)
	IncMetric(string, model.MetricType, float64) error
}
