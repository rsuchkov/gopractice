package agentstats

import "github.com/rsuchkov/gopractice/model"

type StorageExpected interface {
	SaveMetric(string, model.MetricType, float64)
	GetMetrics() map[string]model.Metric
}
