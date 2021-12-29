package agentstats

import (
	"math/rand"
	"runtime"

	"github.com/rsuchkov/gopractice/model"
)

func (svc *Processor) CollectMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	svc.statsStorage.SaveMetric("Alloc", model.MetricTypeGauge, float64(rtm.Alloc))
	svc.statsStorage.SaveMetric("BuckHashSys", model.MetricTypeGauge, float64(rtm.BuckHashSys))
	svc.statsStorage.SaveMetric("Frees", model.MetricTypeGauge, float64(rtm.Frees))
	svc.statsStorage.SaveMetric("GCCPUFraction", model.MetricTypeGauge, float64(rtm.GCCPUFraction))
	svc.statsStorage.SaveMetric("GCSys", model.MetricTypeGauge, float64(rtm.GCSys))
	svc.statsStorage.SaveMetric("HeapAlloc", model.MetricTypeGauge, float64(rtm.HeapAlloc))
	svc.statsStorage.SaveMetric("HeapIdle", model.MetricTypeGauge, float64(rtm.HeapIdle))
	svc.statsStorage.SaveMetric("HeapInuse", model.MetricTypeGauge, float64(rtm.HeapInuse))
	svc.statsStorage.SaveMetric("HeapObjects", model.MetricTypeGauge, float64(rtm.HeapObjects))
	svc.statsStorage.SaveMetric("HeapReleased", model.MetricTypeGauge, float64(rtm.HeapReleased))
	svc.statsStorage.SaveMetric("HeapSys", model.MetricTypeGauge, float64(rtm.HeapSys))
	svc.statsStorage.SaveMetric("LastGC", model.MetricTypeGauge, float64(rtm.LastGC))
	svc.statsStorage.SaveMetric("Lookups", model.MetricTypeGauge, float64(rtm.Lookups))
	svc.statsStorage.SaveMetric("MCacheInuse", model.MetricTypeGauge, float64(rtm.MCacheInuse))
	svc.statsStorage.SaveMetric("MCacheSys", model.MetricTypeGauge, float64(rtm.MCacheSys))
	svc.statsStorage.SaveMetric("MSpanInuse", model.MetricTypeGauge, float64(rtm.MSpanInuse))
	svc.statsStorage.SaveMetric("MSpanSys", model.MetricTypeGauge, float64(rtm.MSpanSys))
	svc.statsStorage.SaveMetric("NextGC", model.MetricTypeGauge, float64(rtm.NextGC))
	svc.statsStorage.SaveMetric("NumForcedGC", model.MetricTypeGauge, float64(rtm.NumForcedGC))
	svc.statsStorage.SaveMetric("NumGC", model.MetricTypeGauge, float64(rtm.NumGC))
	svc.statsStorage.SaveMetric("OtherSys", model.MetricTypeGauge, float64(rtm.OtherSys))
	svc.statsStorage.SaveMetric("PauseTotalNs", model.MetricTypeGauge, float64(rtm.PauseTotalNs))
	svc.statsStorage.SaveMetric("StackInuse", model.MetricTypeGauge, float64(rtm.StackInuse))
	svc.statsStorage.SaveMetric("StackSys", model.MetricTypeGauge, float64(rtm.StackSys))
	svc.statsStorage.SaveMetric("Sys", model.MetricTypeGauge, float64(rtm.Sys))
	svc.statsStorage.SaveMetric("RandomValue", model.MetricTypeGauge, float64(rand.Intn(100)))

	metrics := svc.statsStorage.GetMetrics()
	i, ok := metrics["PollCount"]
	if ok {
		svc.statsStorage.SaveMetric("PollCount", model.MetricTypeCounter, i.Value+1)
	} else {
		svc.statsStorage.SaveMetric("PollCount", model.MetricTypeCounter, 1)
	}
}
