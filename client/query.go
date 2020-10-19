package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Query gets presto query metrics
func (c Client) Query() ([]QueryMetric, error) {
	var metrics []QueryMetric
	url, err := c.withPath("query")
	if err != nil {
		return metrics, errors.Wrap(err, "failed to get query path")
	}
	resp, err := http.Get(url)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to get query metrics")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metrics, errors.Wrap(err, "failed to read query response body")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return metrics, fmt.Errorf("failed to get metrics: %s %d", string(bts), resp.StatusCode)
	}

	if err := json.Unmarshal(bts, &metrics); err != nil {
		return metrics, errors.Wrapf(err, "failed to unmarshall query metrics output: %s", string(bts))
	}
	return metrics, nil
}

// QueryMetric is a single query information
type QueryMetric struct {
	ResourceGroupID []string   `json:"resourceGroupId"`
	State           string     `json:"state"`
	QueryStats      QueryStats `json:"queryStats"`
}

// QueryStats is the stats of a single query execution
type QueryStats struct {
	QueuedTime         Duration `json:"queuedTime"`
	ElapsedTime        Duration `json:"elapsedTime"`
	ExecutionTime      Duration `json:"executionTime"`
	TotalCPUTime       Duration `json:"totalCpuTime"`
	TotalScheduledTime Duration `json:"totalScheduledTime"`
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
