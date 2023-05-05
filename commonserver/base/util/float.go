// Package util .
 
package util

const (
	esp = 1e-4
)

// MaxFloat32 浮点数比较
func MaxFloat32(a, b float32) float32 {
	if Float32HI(a, b) {
		return a
	}
	return b
}

// MaxFloat64 浮点数比较
func MaxFloat64(a, b float64) float64 {
	if Float64HI(a, b) {
		return a
	}
	return b
}

// AbsFloat32 浮点数的绝对值
func AbsFloat32(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}

// AbsInt64 int64的绝对值
func AbsInt64(x int64) int64 {
	if x > 0 {
		return x
	}
	return -x
}

// Float32EQ ==
func Float32EQ(a, b float32) bool {
	return AbsFloat32(a-b) < esp
}

// Float32NE !=
func Float32NE(a, b float32) bool {
	return !Float32EQ(a, b)
}

// Float32HI >
func Float32HI(a, b float32) bool {
	return a-b > esp
}

// Float64HI >
func Float64HI(a, b float64) bool {
	return a-b > esp
}

// Float32HS >=
func Float32HS(a, b float32) bool {
	if a-b > esp {
		return true
	}
	return Float32EQ(a, b)
}

// Float32LO <
func Float32LO(a, b float32) bool {
	return Float32HI(b, a)
}

// Float32LS <=
func Float32LS(a, b float32) bool {
	return Float32HS(b, a)
}
