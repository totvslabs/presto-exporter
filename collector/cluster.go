package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/totvslabs/presto-exporter/client"
)

type clusterCollector struct {
	mutex  sync.Mutex
	client client.Client

	up               *prometheus.Desc
	scrapeDuration   *prometheus.Desc
	runningQueries   *prometheus.Desc
	blockedQueries   *prometheus.Desc
	queuedQueries    *prometheus.Desc
	activeWorkers    *prometheus.Desc
	runningDrivers   *prometheus.Desc
	reservedMemory   *prometheus.Desc
	totalInputRows   *prometheus.Desc
	totalInputBytes  *prometheus.Desc
	totalCPUTimeSecs *prometheus.Desc
}

// NewCluster presto collector
func NewCluster(client client.Client) prometheus.Collector {
	const subsystem = "cluster"
	// nolint: lll
	return &clusterCollector{
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
		runningQueries: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "running_queries"),
			"Running requests of the presto cluster.",
			nil, nil,
		),
		blockedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "blocked_queries"),
			"Blocked queries of the presto cluster.",
			nil, nil,
		),
		queuedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "queued_queries"),
			"Queued queries of the presto cluster.",
			nil, nil,
		),
		activeWorkers: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "active_workers"),
			"Active workers of the presto cluster.",
			nil, nil,
		),
		runningDrivers: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "running_drivers"),
			"Running drivers of the presto cluster.",
			nil, nil,
		),
		reservedMemory: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "reserved_memory_bytes"),
			"Reserved memory of the presto cluster.",
			nil, nil,
		),
		totalInputRows: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "input_rows_total"),
			"Total input rows of the presto cluster.",
			nil, nil,
		),
		totalInputBytes: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "input_bytes_total"),
			"Total input bytes of the presto cluster.",
			nil, nil,
		),
		totalCPUTimeSecs: prometheus.NewDesc(
			prometheus.BuildFQName(ns, subsystem, "cpu_seconds_total"),
			"Total CPU time of the presto cluster.",
			nil, nil,
		),
	}
}

// Describe all metrics
func (c *clusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.runningQueries
	ch <- c.blockedQueries
	ch <- c.queuedQueries
	ch <- c.activeWorkers
	ch <- c.runningDrivers
	ch <- c.reservedMemory
	ch <- c.totalInputRows
	ch <- c.totalInputBytes
	ch <- c.totalCPUTimeSecs
}

// Collect all metrics
func (c *clusterCollector) Collect(ch chan<- prometheus.Metric) {
	var start = time.Now()
	defer func() {
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
	}()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Info("Collecting cluster metrics...")

	metrics, err := c.client.Cluster()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape tasks")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.runningQueries, prometheus.GaugeValue, metrics.RunningQueries)
	ch <- prometheus.MustNewConstMetric(c.blockedQueries, prometheus.GaugeValue, metrics.BlockedQueries)
	ch <- prometheus.MustNewConstMetric(c.queuedQueries, prometheus.GaugeValue, metrics.QueuedQueries)
	ch <- prometheus.MustNewConstMetric(c.activeWorkers, prometheus.GaugeValue, metrics.ActiveWorkers)
	ch <- prometheus.MustNewConstMetric(c.runningDrivers, prometheus.GaugeValue, metrics.RunningDrivers)
	ch <- prometheus.MustNewConstMetric(c.reservedMemory, prometheus.GaugeValue, metrics.ReservedMemory)
	ch <- prometheus.MustNewConstMetric(c.totalInputRows, prometheus.CounterValue, metrics.TotalInputRows)
	ch <- prometheus.MustNewConstMetric(c.totalInputBytes, prometheus.CounterValue, metrics.TotalInputBytes)
	ch <- prometheus.MustNewConstMetric(c.totalCPUTimeSecs, prometheus.CounterValue, metrics.TotalCPUTimeSecs)
}
