package collector

import (
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/totvslabs/presto-exporter/client"
)

type queryCollector struct {
	mutex          sync.Mutex
	client         client.Client
	up             *prometheus.Desc
	scrapeDuration *prometheus.Desc
	queueTime      *prometheus.Desc
	elapsedTime    *prometheus.Desc
	executionTime  *prometheus.Desc
	totalQueries   *prometheus.Desc
}

// NewQuery presto collector
func NewQuery(client client.Client) prometheus.Collector {
	const subsystem = "queries"
	// nolint: lll
	return &queryCollector{
		client: client,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "up"),
			"Presto API is responding",
			nil, nil,
		),
		scrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "scrape_duration_seconds"),
			"Scrape duration in seconds",
			nil, nil,
		),
		queueTime: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "queue_time_seconds"),
			"Queue time in seconds",
			[]string{"resource_group"}, nil,
		),
		elapsedTime: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "elapsed_time_seconds"),
			"Elapsed time in seconds",
			[]string{"resource_group"}, nil,
		),
		executionTime: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "execution_time_seconds"),
			"Execution time in seconds",
			[]string{"resource_group"}, nil,
		),
		totalQueries: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "total"),
			"Query count",
			[]string{"state"}, nil,
		),
	}
}

// Describe all metrics
func (c *queryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.queueTime
	ch <- c.elapsedTime
	ch <- c.executionTime
	ch <- c.totalQueries
}

// Collect all metrics
func (c *queryCollector) Collect(ch chan<- prometheus.Metric) {
	var start = time.Now()
	defer func() {
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
	}()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Info("Collecting query metrics...")

	queries, err := c.client.Query()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape tasks")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	var durations = map[string]queryMetrics{}
	var states = map[string]int64{}
	for _, query := range queries {
		group := strings.Join(query.ResourceGroupID, ".")
		m, ok := durations[group]
		if !ok {
			m = queryMetrics{}
		}
		m.queue += time.Duration(query.QueryStats.QueuedTime)
		m.elapsed += time.Duration(query.QueryStats.ElapsedTime)
		m.execution += time.Duration(query.QueryStats.ExecutionTime)
		durations[group] = m
		if query.ErrorCode.Type == "USER_ERROR" {
			states["USER_ERROR"]++
		} else {
			states[query.State]++
		}
	}

	for k, v := range durations {
		ch <- prometheus.MustNewConstMetric(c.queueTime, prometheus.GaugeValue, v.queue.Seconds(), k)
		ch <- prometheus.MustNewConstMetric(c.elapsedTime, prometheus.GaugeValue, v.elapsed.Seconds(), k)
		ch <- prometheus.MustNewConstMetric(c.executionTime, prometheus.GaugeValue, v.execution.Seconds(), k)
	}

	for k, v := range states {
		ch <- prometheus.MustNewConstMetric(c.totalQueries, prometheus.GaugeValue, float64(v), k)
	}
}

type queryMetrics struct {
	queue, elapsed, execution time.Duration
}
