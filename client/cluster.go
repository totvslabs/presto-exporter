package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Cluster gets presto cluster metrics
func (c Client) Cluster() (ClusterMetrics, error) {
	var metrics ClusterMetrics
	url, err := c.withPath("cluster")
	if err != nil {
		return metrics, errors.Wrap(err, "failed to get cluster path")
	}
	resp, err := http.Get(url)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to get cluster metrics")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to read cluster response body")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return metrics, fmt.Errorf("failed to get metrics: %s %d", string(bts), resp.StatusCode)
	}

	if err := json.Unmarshal(bts, &metrics); err != nil {
		return metrics, errors.Wrapf(err, "failed to unmarshall cluster metrics output: %s", string(bts))
	}
	return metrics, nil
}

// ClusterMetrics is the presto v1/cluster output representation.
type ClusterMetrics struct {
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
