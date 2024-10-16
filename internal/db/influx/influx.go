package influx

import (
	"fmt"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rtzgod/db-benchmarks/internal/utils"
	"sync"
	"time"
)

type Influx struct {
	Client   influxdb.Client
	WriteAPI api.WriteAPI
	Timeout  time.Duration
}

func New(URL, token, org, bucket string, timeout time.Duration) *Influx {
	client := influxdb.NewClient(URL, token)
	writeAPI := client.WriteAPI(org, bucket)

	return &Influx{
		Client:   client,
		WriteAPI: writeAPI,
		Timeout:  timeout,
	}
}

func (i *Influx) BenchmarkInfluxWrite(dataPoints int, sensorsAmount int) {
	var wg sync.WaitGroup

	// Launch a goroutine for each sensor
	for sensor := 0; sensor < sensorsAmount; sensor++ {
		wg.Add(1)

		go func(sensorID int) {
			defer wg.Done()
			i.WritePoint(dataPoints, sensorID)
		}(sensor)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All data points inserted.")
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
