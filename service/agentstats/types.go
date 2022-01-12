package agentstats

import "github.com/rsuchkov/gopractice/model"

type StorageExpected interface {
	SaveMetric(string, model.MType, float64) (model.Metric, error)
	GetMetrics() map[string]model.Metric
}
