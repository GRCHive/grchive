package core

import (
	"math"
)

func Kilobytes(in float64) int64 {
	return int64(math.Round(in * 1024))
}

func Megabytes(in float64) int64 {
	return Kilobytes(in) * 1024
}

func Gigabytes(in float64) int64 {
	return Megabytes(in) * 1024
}
