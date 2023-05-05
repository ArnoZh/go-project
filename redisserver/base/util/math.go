// Package util .
 
package util

// MaxInt max int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxInt32 max int32
func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// MinInt32 min int32
func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// AbsInt32 abs int32
func AbsInt32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

// MaxInt64 max int64
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinInt64 get min value
func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// AbsFloat64 浮点绝对值
func AbsFloat64(x float64) float64 {
	if x > 0 {
		return x
	}
	return -x
}

// MinFloat32 小值
func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// MinFloat64 小值
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// FindInt64 查找切片中是否存在某个数字并返回所在的索引值，若不存在返回-1
func FindInt64(aim []int64, value int64) int {
	for i, v := range aim {
		if v == value {
			return i
		}
	}
	return -1
}

// FindInt32 查找切片中是否存在某个数字并返回所在的索引值，若不存在返回-1
func FindInt32(aim []int32, value int32) int {
	for i, v := range aim {
		if v == value {
			return i
		}
	}
	return -1
}

// RemoveInt32 删除切片中指定Value，返回删除后的切片
func RemoveInt32(aim []int32, value int32) []int32 {
	for i, v := range aim {
		if v == value {
			return append(aim[:i], aim[i+1:]...)
		}
	}
	return aim
}

// InsertInt64 唯一的在数组中添加v
func InsertInt64(a *[]int64, v int64) {
	for _, v0 := range *a {
		if v0 == v {
			return
		}
	}

	*a = append(*a, v)
}
