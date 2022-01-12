package model

import "fmt"

type MType string

const (
	MetricTypeGauge   MType = "gauge"
	MetricTypeCounter MType = "counter"
)

// String implements fmt.Stringer interface.
func (s MType) String() string {
	return string(s)
}

// Validate performs enum validation.
func (s MType) Validate() error {
	switch s {
	case MetricTypeGauge, MetricTypeCounter:
		return nil
	default:
		return fmt.Errorf("unkown ProductStatus: %s", s)
	}
}

type Metric struct {
	ID    string   `json:"id"`
	MType MType    `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
