// Package geo .

package geo

import (
	"gs/pb/cspb"
	"math"
)

// Coord 坐标
type Coord struct {
	X int32 `bson:"x"`
	Z int32 `bson:"z"`
}

// NewCoord 新建
func NewCoord(x, z int32) Coord {
	return Coord{
		X: x,
		Z: z,
	}
}

// NewCoord2 New from int64
func NewCoord2(v int64) Coord {
	z := int32(v & 0x00000000FFFFFFFF)
	x := int32(v >> 32)
	return Coord{X: x, Z: z}
}

// IsEqual 判断是否相等
func (c Coord) IsEqual(target Coord) bool {
	return c == target
}

// Int64 转化为int64
func (c *Coord) Int64() int64 {
	return (int64(c.X) << 32) | int64(c.Z)
}

// Add 相加
func (c *Coord) Add(coord Coord) Coord {
	return Coord{
		X: c.X + coord.X,
		Z: c.Z + coord.Z,
	}
}

// Minus 相减
func (c *Coord) Minus(coord Coord) Coord {
	return Coord{
		X: c.X - coord.X,
		Z: c.Z - coord.Z,
	}
}

// GetLocationToBorder 获取和边界的位置关系
func (c *Coord) GetLocationToBorder(b *Border) LocationState {
	return b.CoordLocation(*c)
}

// CalDstCoordToCoord 计算点和点之间的距离
func CalDstCoordToCoord(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return math.Sqrt(float64(squared))
}

// CalDstCoordToCoordWithoutSqrt 计算点和点之间的距离,不开根号
func CalDstCoordToCoordWithoutSqrt(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return float64(squared)
}

//通知客户端时，大地图坐标系转内城坐标系
func (c *Coord) ToInnerMapCoord() cspb.Coord {
	res := cspb.Coord{}
	res.X = c.X%4000000 + 4000000
	res.Z = c.Z%4000000 + 4000000
	return res
}
