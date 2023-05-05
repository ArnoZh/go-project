// Package geo .
 
package geo

import "math/rand"

// Rectangle 矩形对象
type Rectangle struct {
	Coord  // 左下角点，即X Z均最小的点
	Width  int32
	Height int32
}

// NewRectangle 创建矩形
func NewRectangle(x, z, width, height int32) Rectangle {
	return Rectangle{
		Coord: Coord{
			X: x,
			Z: z,
		},
		Width:  width,
		Height: height,
	}
}

// RandCoord 随机一个坐标
func (rec *Rectangle) RandCoord() Coord {
	return Coord{
		X: rec.X + rand.Int31n(rec.Width),
		Z: rec.Z + rand.Int31n(rec.Height),
	}
}

// GetVerticeCoords 获取4个顶点坐标,边界点按照逆时针排列
func (rec *Rectangle) GetVerticeCoords() [4]Coord {
	var p [4]Coord
	p[0] = Coord{X: rec.Coord.X, Z: rec.Coord.Z}
	p[1] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z}
	p[2] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z + rec.Height}
	p[3] = Coord{X: rec.Coord.X, Z: rec.Coord.Z + rec.Height}
	return p
}

// GetVectors 获取4条线段，按照逆时针排列
func (rec *Rectangle) GetVectors() [4]Vector {
	coords := rec.GetVerticeCoords()
	return [4]Vector{
		NewVectorByCoord(coords[0]),
		NewVectorByCoord(coords[1]),
		NewVectorByCoord(coords[2]),
		NewVectorByCoord(coords[3]),
	}
}

// GetLocationToBorder 获得矩形和给定边界的位置关系
func (rec *Rectangle) GetLocationToBorder(b *Border) LocationState {
	minX := rec.X
	maxX := rec.X + rec.Width
	minZ := rec.Z
	maxZ := rec.Z + rec.Height
	return b.RectLocation(minX, maxX, minZ, maxZ)
}

// IsCoordInside 点是否在矩形内
func (rec *Rectangle) IsCoordInside(p Coord) bool {
	return p.X >= rec.Coord.X && p.X <= rec.Coord.X+rec.Width &&
		p.Z >= rec.Coord.Z && p.Z <= rec.Coord.Z+rec.Height
}

// IsRectIntersect 判断两个矩形是否相交（重合面积大于0就算相交）
func IsRectIntersect(recA *Rectangle, recB *Rectangle) bool {
	if recB.X+recB.Width <= recA.X ||
		recB.X >= recA.X+recA.Width ||
		recB.Z >= recA.Z+recA.Height ||
		recB.Z+recB.Height <= recA.Z {
		return false
	}
	return true
}
