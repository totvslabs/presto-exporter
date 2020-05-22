package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Client is the presto client
type Client struct {
	url string
}

// New creates a new presto client
func New(url string) Client {
	return Client{
		url: url,
	}
}

// Get presto metrics
func (c Client) Get() (Metrics, error) {
	var metrics Metrics
	resp, err := http.Get(c.url)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to get metrics")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return metrics, fmt.Errorf("failed to get metrics: %s %d", string(bts), resp.StatusCode)
	}

	if err := json.Unmarshal(bts, &metrics); err != nil {
		return metrics, errors.Wrapf(err, "failed to unmarshall metrics output: %s", string(bts))
	}
	return metrics, nil
}

// Metrics is the presto api output representation.
type Metrics struct {
	RunningQueries   float64 `json:"runningQueries"`
	BlockedQueries   float64 `json:"blockedQueries"`
	QueuedQueries    float64 `json:"queuedQueries"`
	ActiveWorkers    float64 `json:"activeWorkers"`
	RunningDrivers   float64 `json:"runningDrivers"`
	ReservedMemory   float64 `json:"reservedMemory"`
	TotalInputRows   float64 `json:"totalInputRows"`
	TotalInputBytes  float64 `json:"totalInputBytes"`
	TotalCPUTimeSecs float64 `json:"totalCpuTimeSecs"`
}
