// Package util .
 
package util

import "math"

const (
	// Epsilon 浮点正常精度
	Epsilon = 0.00000001
	// LowEpsilon 浮点低精度
	LowEpsilon = 0.01
)

// Ac 浮点精度函数
var Ac Accuracy = func() float64 { return LowEpsilon }

// Accuracy 精度
type Accuracy func() float64

// Equal 相等比较
func (ac Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < LowEpsilon
}

// Greater 大于
func (ac Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > ac()
}

// Smaller 小于
func (ac Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > ac()
}

// GreaterOrEqual 大于等于
func (ac Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < ac()
}

// SmallerOrEqual 小于等于
func (ac Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < ac()
}
