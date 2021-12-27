package serverstats

import "github.com/rsuchkov/gopractice/model"

func (svc *Processor) SaveMetric(name string, mtype model.MetricType, value float64) error {
	err := mtype.Validate()
	if err != nil {
		return err
	}
	svc.statsStorage.SaveMetric(name, mtype, value)
	return nil
}
