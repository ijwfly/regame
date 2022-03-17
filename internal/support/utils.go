package support

import (
	"math"
	"time"
)

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Round(input float32, num int) float32 {
	return float32(float64(int(float64(input)*math.Pow10(num))) / math.Pow10(num))
}

func Round2(input float32) float32 {
	return float32(int(input*100)) / 100
}
