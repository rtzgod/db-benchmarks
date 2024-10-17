package main

import (
	"github.com/rtzgod/db-benchmarks/internal/config"
	"github.com/rtzgod/db-benchmarks/internal/db/influx"
	"github.com/rtzgod/db-benchmarks/internal/db/timescale"
	"log"
)

const (
	dataPoints    = 10000
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

	timescaledb, err := timescale.New(
		cfg.Timescale.URL,
		cfg.Timescale.Timeout)
	if err != nil {
		log.Fatalf("timescale init failed")
	}

	defer influxdb.Client.Close()
	defer timescaledb.Close()

	//influxdb.BenchmarkWrite(dataPoints, sensorsAmount)
	//influxdb.BenchmarkRead(queriesAmount)

	//timescaledb.BenchmarkWrite(dataPoints, sensorsAmount)
	timescaledb.BenchmarkRead(queriesAmount)
}
