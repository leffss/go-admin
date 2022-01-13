package common

import (
	"fmt"
	"math"
)

func logN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logN(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+" %s", val, suffix)
}

func FileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB"}
	return humanBytes(s, 1024, sizes)
}
