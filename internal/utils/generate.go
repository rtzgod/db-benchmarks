package utils

import (
	"math"
	"math/rand"
)

func GenerateTemperature(timestamp int64, sensorId float64) float64 {
	baseTemp := 20.0 + sensorId
	amplitude := 5.0
	period := 24.0

	timeInHours := float64(timestamp) / 3600.0
	sineValue := amplitude * math.Sin(2*math.Pi*timeInHours/period)

	noise := rand.Float64()*2.0 - 0.5

	return baseTemp + sineValue + noise
}
