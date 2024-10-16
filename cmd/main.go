package main

import (
	"fmt"
	"github.com/rtzgod/db-benchmarks/internal/config"
	"github.com/rtzgod/db-benchmarks/internal/db/influx"
	"time"
)

const (
	dataPoints    = 2000000
	divider       = 1000000000
	sensorsAmount = 5
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(500000 / (float32(time.Second+time.Second/2) / divider))
	influxdb := influx.New(
		cfg.Influx.URL,
		cfg.Influx.Token,
		cfg.Influx.Org,
		cfg.Influx.Bucket,
		cfg.Influx.Timeout)

	defer influxdb.Client.Close()

	startTime := time.Now()

	influxdb.BenchmarkInfluxWrite(dataPoints, sensorsAmount)

	elapsed := time.Since(startTime)

	fmt.Printf("Inserted %d points in %v, average writing speed: %v rows/s\n", dataPoints*sensorsAmount, elapsed, (dataPoints*sensorsAmount)/(float32(elapsed)/divider))
}
