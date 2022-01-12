package serverstats

import "github.com/rsuchkov/gopractice/model"

type StorageExpected interface {
	SaveMetric(string, model.MType, float64) (model.Metric, error)
	GetMetrics() map[string]model.Metric
	GetMetric(string, model.MType) (model.Metric, bool)
	IncMetric(string, model.MType, float64) (model.Metric, error)
}
