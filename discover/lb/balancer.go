package lb

import (
	"errors"
)

// Balancer yields endpoints according to some heuristic.
type Balancer interface {
	Endpoint() (string, error)
}

// ErrNoEndpoints is returned when no qualifying endpoints are available.
var ErrNoEndpoints = errors.New("no endpoints available")
