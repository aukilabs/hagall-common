package smoketest

import "time"

type Status int

const (
	StatusUnknown Status = 0
	StatusSuccess Status = 1
	StatusFailed  Status = 2
	StatusTimeout Status = 3
	StatusError   Status = 4
)

type SmokeTestRequest struct {
	Endpoint           string        `json:"endpoint"`
	Token              string        `json:"token"`
	MaxSessionIDLength int           `json:"max_session_id_length"`
	Timeout            time.Duration `json:"timeout"`
}

type SmokeTestResults struct {
	FromEndpoint    string  `json:"from_endpoint"`
	ToEndpoint      string  `json:"to_endpoint"`
	LatencyMilliSec float64 `json:"latency"`
	Status          Status  `json:"status"`
}
