package model

import "fmt"

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

// String implements fmt.Stringer interface.
func (s MetricType) String() string {
	return string(s)
}

// Validate performs enum validation.
func (s MetricType) Validate() error {
	switch s {
	case MetricTypeGauge, MetricTypeCounter:
		return nil
	default:
		return fmt.Errorf("unkown ProductStatus: %s", s)
	}
}

type Metric struct {
	Name       string
	MetricType MetricType
	Value      float64
}
