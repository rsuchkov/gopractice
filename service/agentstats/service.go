package agentstats

import "fmt"

type (
	Processor struct {
		statsStorage StorageExpected
	}

	ProcessorOption func(svc *Processor)
)

func WithStatsStorage(st StorageExpected) ProcessorOption {
	return func(svc *Processor) {
		svc.statsStorage = st
	}
}

func New(opts ...ProcessorOption) (*Processor, error) {
	svc := &Processor{}
	for _, opt := range opts {
		opt(svc)
	}

	if svc.statsStorage == nil {
		return nil, fmt.Errorf("orderStorage: nil")
	}

	return svc, nil
}
