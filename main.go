package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/totvslabs/presto-exporter/client"
	"github.com/totvslabs/presto-exporter/collector"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// nolint: gochecknoglobals,lll
var (
	version       = "dev"
	app           = kingpin.New("presto-exporter", "exports presto metrics in the prometheus format")
	listenAddress = app.Flag("web.listen-address", "Address to listen on for web interface and telemetry").Default("127.0.0.1:9430").String()
	metricsPath   = app.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	prestoURL     = app.Flag("presto.url", "Presto URL to scrape").Default("http://localhost:8080/v1/cluster").String()
)

func main() {
	log.AddFlags(app)
	app.Version(version)
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Infof("starting presto-exporter %s...", version)

	var client = client.New(*prestoURL)
	prometheus.MustRegister(collector.New(client))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
			<head><title>Presto Exporter</title></head>
			<body>
				<h1>Presto Exporter</h1>
				<p><a href="`+*metricsPath+`">Metrics</a></p>
			</body>
			</html>
		`)
	})

	log.Infof("server listening on %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
