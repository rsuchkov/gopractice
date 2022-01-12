package memory

import (
	"sync"

	"github.com/rsuchkov/gopractice/model"
)

type AgentStats struct {
	mu      sync.Mutex
	Metrics map[string]model.Metric
}
