package utils

import "math"

var SinFloat64 = sinFloat64()

// 正弦函数
func sinFloat64() []float64 {
	n := 400
	data := make([]float64, n)
	for i := range data {
		data[i] = 1 + math.Sin(float64(i)/5)
	}
	return data
}
