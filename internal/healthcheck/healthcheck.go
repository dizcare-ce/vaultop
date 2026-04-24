// Package healthcheck provides a simple liveness and readiness probe
// mechanism for vaultop. Each registered probe returns a Status indicating
// whether the component is healthy, degraded, or unhealthy.
package healthcheck

import (
	"fmt"
	"sync"
	"time"
)

// Level describes the severity of a health status.
type Level string

const (
	LevelOK       Level = "ok"
	LevelDegraded Level = "degraded"
	LevelUnhealthy Level = "unhealthy"
)

// Status is the result returned by a single probe.
type Status struct {
	Name    string        `json:"name"`
	Level   Level         `json:"level"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency_ms"`
}

// Probe is a function that checks a single component.
type Probe func() Status

// Checker runs all registered probes and aggregates results.
type Checker struct {
	mu     sync.RWMutex
	probes map[string]Probe
}

// New returns an empty Checker.
func New() *Checker {
	return &Checker{probes: make(map[string]Probe)}
}

// Register adds a named probe. Registering the same name twice overwrites the
// previous probe.
func (c *Checker) Register(name string, p Probe) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.probes[name] = p
}

// RunAll executes every registered probe and returns the collected statuses
// together with an overall Level (worst of all individual levels).
func (c *Checker) RunAll() ([]Status, Level) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	overall := LevelOK
	results := make([]Status, 0, len(c.probes))

	for name, probe := range c.probes {
		start := time.Now()
		s := probe()
		s.Latency = time.Since(start)
		if s.Name == "" {
			s.Name = name
		}
		results = append(results, s)
		overall = worstLevel(overall, s.Level)
	}
	return results, overall
}

// OK is a convenience constructor for a healthy Status.
func OK(name string) Status { return Status{Name: name, Level: LevelOK} }

// Unhealthy is a convenience constructor for an unhealthy Status.
func Unhealthy(name string, err error) Status {
	return Status{Name: name, Level: LevelUnhealthy, Message: fmt.Sprintf("%v", err)}
}

func worstLevel(a, b Level) Level {
	rank := map[Level]int{LevelOK: 0, LevelDegraded: 1, LevelUnhealthy: 2}
	if rank[b] > rank[a] {
		return b
	}
	return a
}
