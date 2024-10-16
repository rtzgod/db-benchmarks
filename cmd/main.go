package main

import (
	"github.com/rtzgod/db-benchmarks/internal/config"
	"github.com/rtzgod/db-benchmarks/internal/db/influx"
)

const (
	dataPoints    = 1000000
	queriesAmount = 100
	sensorsAmount = 5
)

func main() {
	cfg := config.MustLoad()

	influxdb := influx.New(
		cfg.Influx.URL,
		cfg.Influx.Token,
		cfg.Influx.Org,
		cfg.Influx.Bucket,
		cfg.Influx.Timeout)

	defer influxdb.Client.Close()

	// influxdb.BenchmarkInfluxWrite(dataPoints, sensorsAmount)
	influxdb.BenchmarkInfluxRead(queriesAmount)
}
