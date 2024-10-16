package influx

import (
	"context"
	"fmt"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rtzgod/db-benchmarks/internal/utils"
	"log"
	"sync"
	"time"
)

type Influx struct {
	Client     influxdb.Client
	WriteAPI   api.WriteAPI
	BucketName string
	Timeout    time.Duration
}

func New(URL, token, org, bucket string, timeout time.Duration) *Influx {
	client := influxdb.NewClient(URL, token)
	writeAPI := client.WriteAPI(org, bucket)

	return &Influx{
		Client:     client,
		WriteAPI:   writeAPI,
		BucketName: bucket,
		Timeout:    timeout,
	}
}

func (i *Influx) BenchmarkInfluxWrite(dataPoints, sensorsAmount float32) {
	var wg sync.WaitGroup
	startTime := time.Now()

	// Launch a goroutine for each sensor
	for sensor := 0; sensor < int(sensorsAmount); sensor++ {
		wg.Add(1)

		go func(sensorID int) {
			defer wg.Done()
			i.WritePoint(int(dataPoints), sensorID)
		}(sensor)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	elapsed := time.Since(startTime)

	fmt.Printf("Inserted %v points in %v, average writing speed: %v rows/s\n", dataPoints*sensorsAmount, elapsed, (dataPoints*sensorsAmount)/(float32(elapsed)/1000000000.0))
}

func (i *Influx) WritePoint(dataPoints int, sensor int) {
	for j := 0; j < dataPoints; j++ {
		point := influxdb.NewPoint(
			"sensor_data",
			map[string]string{"sensor_id": fmt.Sprintf("sensor-%d", sensor)},
			map[string]interface{}{"temperature": utils.GenerateTemperature(time.Now().Unix(), float64(sensor))},
			time.Now(),
		)

		i.WriteAPI.WritePoint(point)

		if j%1000 == 0 {
			fmt.Printf("Inserted %d points\n", j)
		}

		time.Sleep(i.Timeout)
	}
}

func (i *Influx) BenchmarkInfluxRead(queries int) {
	query := fmt.Sprintf("from(bucket:\"%s\") |> range(start: -1h)", i.BucketName)
	queryAPI := i.Client.QueryAPI("my-org")

	startTime := time.Now()

	for i := 0; i < queries; i++ {
		_, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		if i%10 == 0 {
			fmt.Printf("Executed %d queries\n", i)
		}
	}
	fmt.Printf("Total time for %d queries: %v, average reading speed: 1 query in %v\n", queries, time.Since(startTime), time.Since(startTime)/time.Duration(queries))
}
