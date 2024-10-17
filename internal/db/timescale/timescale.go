package timescale

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"sync"
	"time"

	"github.com/rtzgod/db-benchmarks/internal/utils"
)

// Timescale struct to manage the database connection and settings
type Timescale struct {
	Conn    *pgx.Conn
	Timeout time.Duration
}

// New initializes a new TimescaleDB connection
func New(url string, timeout time.Duration) (*Timescale, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to TimescaleDB: %v", err)
	}
	return &Timescale{Conn: conn, Timeout: timeout}, nil
}

// BenchmarkWrite inserts data into the TimescaleDB using multiple goroutines.
func (t *Timescale) BenchmarkWrite(dataPoints, sensorsAmount float32) {
	var wg sync.WaitGroup
	startTime := time.Now()

	for sensor := 0; sensor < int(sensorsAmount); sensor++ {
		wg.Add(1)
		go func(sensorID int) {
			defer wg.Done()
			t.WritePoints(int(dataPoints), sensorID)
		}(sensor)
	}

	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %v points in %v, average writing speed: %v rows/s\n",
		dataPoints*sensorsAmount, elapsed, (dataPoints*sensorsAmount)/(float32(elapsed)/1000000000.0))
}

// WritePoints inserts multiple data points for a given sensor.
func (t *Timescale) WritePoints(dataPoints int, sensorID int) {
	for i := 0; i < dataPoints; i++ {
		_, err := t.Conn.Exec(
			context.Background(),
			"INSERT INTO sensor_data (sensor_id, temperature, time) VALUES ($1, $2, $3)",
			fmt.Sprintf("sensor-%d", sensorID),
			utils.GenerateTemperature(time.Now().Unix(), float64(sensorID)),
			time.Now(),
		)
		if err != nil {
			log.Fatalf("Failed to insert data point: %v", err)
		}

		if i%1000 == 0 {
			fmt.Printf("Inserted %d points for sensor %d\n", i, sensorID)
		}

		time.Sleep(t.Timeout)
	}
}

// BenchmarkRead executes multiple queries to benchmark the read speed.
func (t *Timescale) BenchmarkRead(queries int) {
	query := "SELECT * FROM sensor_data WHERE time >= NOW() - INTERVAL '1 hour'"

	startTime := time.Now()

	for i := 0; i < queries; i++ {
		rows, err := t.Conn.Query(context.Background(), query)
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		rows.Close()

		if i%10 == 0 {
			fmt.Printf("Executed %d queries\n", i)
		}
	}

	fmt.Printf("Total time for %d queries: %v, average reading speed: 1 query in %v\n", queries, time.Since(startTime), time.Since(startTime)/time.Duration(queries))
}

// Close closes the TimescaleDB connection.
func (t *Timescale) Close() {
	if err := t.Conn.Close(context.Background()); err != nil {
		log.Fatalf("Failed to close TimescaleDB connection: %v", err)
	}
}
